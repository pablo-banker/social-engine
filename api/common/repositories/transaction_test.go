package repositories

import (
	"context"
	"errors"
	"testing"

	"social-engine/common/repositories/constants"
	"social-engine/common/repositories/interfaces"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newSQLiteRepo opens a real (non dry-run) in-memory SQLite database with a
// single connection so transaction commit/rollback visibility behaves like a
// normal RDBMS within the test.
func newSQLiteRepo(t *testing.T) *BaseRepository {
	t.Helper()

	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := gdb.Exec("CREATE TABLE widgets (id INTEGER PRIMARY KEY)").Error; err != nil {
		t.Fatalf("create table: %v", err)
	}

	return &BaseRepository{db: gdb}
}

func countWidgets(t *testing.T, repo *BaseRepository) int64 {
	t.Helper()
	var n int64
	if err := repo.Raw(context.Background(), &n, "SELECT count(*) FROM widgets"); err != nil {
		t.Fatalf("count widgets: %v", err)
	}
	return n
}

func insertWidget(ctx context.Context, tx interfaces.IRepository, id int) error {
	return tx.Raw(ctx, nil, "INSERT INTO widgets(id) VALUES (?)", id)
}

func TestWithTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("commits when fn returns nil", func(t *testing.T) {
		repo := newSQLiteRepo(t)

		err := repo.WithTransaction(ctx, func(tx interfaces.IRepository) error {
			return insertWidget(ctx, tx, 1)
		})

		assert.NoError(t, err)
		assert.Equal(t, int64(1), countWidgets(t, repo), "row should be committed")
	})

	t.Run("rolls back when fn returns an error", func(t *testing.T) {
		repo := newSQLiteRepo(t)
		wantErr := errors.New("boom")

		err := repo.WithTransaction(ctx, func(tx interfaces.IRepository) error {
			if e := insertWidget(ctx, tx, 2); e != nil {
				return e
			}
			return wantErr
		})

		assert.ErrorIs(t, err, wantErr)
		assert.Equal(t, int64(0), countWidgets(t, repo), "row must be rolled back on error")
	})

	t.Run("rolls back and re-panics on panic", func(t *testing.T) {
		repo := newSQLiteRepo(t)

		assert.PanicsWithValue(t, "kaboom", func() {
			_ = repo.WithTransaction(ctx, func(tx interfaces.IRepository) error {
				_ = insertWidget(ctx, tx, 3)
				panic("kaboom")
			})
		})

		assert.Equal(t, int64(0), countWidgets(t, repo), "row must be rolled back on panic")
	})

	t.Run("returns BeginTx error without invoking fn", func(t *testing.T) {
		repo := newSQLiteRepo(t)
		if err := repo.db.Exec("DROP TABLE widgets").Error; err != nil {
			t.Fatalf("drop: %v", err)
		}
		// Close the pool so Begin fails deterministically.
		sqlDB, _ := repo.db.DB()
		_ = sqlDB.Close()

		called := false
		err := repo.WithTransaction(ctx, func(tx interfaces.IRepository) error {
			called = true
			return nil
		})

		assert.Error(t, err)
		assert.False(t, called, "fn must not run when BeginTx fails")
	})
}

// TestRollbackNilActuallyRollsBack is the regression test for the leak: the old
// Rollback(nil) returned without rolling back, leaving the connection open as
// "idle in transaction". It must now roll back.
func TestRollbackNilActuallyRollsBack(t *testing.T) {
	ctx := context.Background()
	repo := newSQLiteRepo(t)

	tx, err := repo.BeginTx(ctx)
	assert.NoError(t, err)
	assert.NoError(t, insertWidget(ctx, tx, 10))

	assert.NoError(t, tx.Rollback(nil))
	assert.Equal(t, int64(0), countWidgets(t, repo), "Rollback(nil) must actually roll back")
}

func TestRollbackAfterCommitIsNoOp(t *testing.T) {
	ctx := context.Background()
	repo := newSQLiteRepo(t)

	tx, err := repo.BeginTx(ctx)
	assert.NoError(t, err)
	assert.NoError(t, insertWidget(ctx, tx, 11))
	assert.NoError(t, tx.Commit())

	// Calling Rollback after a successful Commit must be a harmless no-op.
	assert.NoError(t, tx.Rollback(nil))
	assert.Equal(t, int64(1), countWidgets(t, repo), "committed row must survive a post-commit rollback")
}

func TestRollbackPropagatesOriginalError(t *testing.T) {
	ctx := context.Background()
	repo := newSQLiteRepo(t)

	tx, err := repo.BeginTx(ctx)
	assert.NoError(t, err)
	assert.NoError(t, insertWidget(ctx, tx, 12))

	orig := errors.New("db failure")
	assert.ErrorIs(t, tx.Rollback(orig), orig, "Rollback(err) returns the original error after rolling back")
	assert.Equal(t, int64(0), countWidgets(t, repo))
}

func TestMockWithTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("commits on success", func(t *testing.T) {
		m := NewMockRepository([]MockPayload{
			{Name: "Begin Tx", Type: constants.RepositoryBeginTx},
			{Name: "Commit", Type: constants.RepositoryCommit},
		})

		err := m.WithTransaction(ctx, func(tx interfaces.IRepository) error { return nil })
		assert.NoError(t, err)
	})

	t.Run("rolls back on error", func(t *testing.T) {
		m := NewMockRepository([]MockPayload{
			{Name: "Begin Tx", Type: constants.RepositoryBeginTx},
			{Name: "Rollback", Type: constants.RepositoryRollback},
		})

		wantErr := errors.New("boom")
		err := m.WithTransaction(ctx, func(tx interfaces.IRepository) error { return wantErr })
		assert.ErrorIs(t, err, wantErr)
	})

	t.Run("returns BeginTx error without invoking fn", func(t *testing.T) {
		m := NewMockRepository([]MockPayload{
			{Name: "Begin Tx fails", Type: constants.RepositoryBeginTx, ExpectedError: errors.New("tx error")},
		})

		called := false
		err := m.WithTransaction(ctx, func(tx interfaces.IRepository) error {
			called = true
			return nil
		})
		assert.Error(t, err)
		assert.False(t, called)
	})
}
