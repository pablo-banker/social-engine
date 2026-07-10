package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"social-engine/common/repositories/entities"
	"social-engine/common/repositories/interfaces"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ interfaces.IRepository = (*BaseRepository)(nil)

type BaseRepository struct {
	db *gorm.DB
}

func NewRepository(ctx context.Context) *BaseRepository {
	return &BaseRepository{db: db.WithContext(ctx)}
}

func (r *BaseRepository) WithTableName(tableName string) interfaces.IRepository {
	return &BaseRepository{db: r.db.Table(tableName)}
}

func (r *BaseRepository) Ping(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	return sqlDB.PingContext(ctx)
}

func (r *BaseRepository) BeginTx(ctx context.Context) (interfaces.IRepository, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &BaseRepository{db: tx}, nil
}

func (r *BaseRepository) Rollback(err error) error {
	if p := recover(); p != nil {
		_ = r.db.Rollback()
		panic(p)
	}

	// Always issue the rollback. A previous version skipped it when err was nil,
	// which turned the common `defer repoTx.Rollback(nil)` guard into a no-op and
	// leaked the connection as "idle in transaction" on error/early-return paths.
	// Rolling back a transaction that was already committed (or rolled back) is a
	// harmless no-op, so the after-commit case is ignored explicitly.
	if rbErr := r.db.Rollback().Error; rbErr != nil && !isTxAlreadyClosed(rbErr) {
		if err != nil {
			return fmt.Errorf("rollback error: %w; original error: %v", rbErr, err)
		}
		return rbErr
	}

	return err
}

func (r *BaseRepository) Commit() error {
	return r.db.Commit().Error
}

func (r *BaseRepository) WithTransaction(ctx context.Context, fn func(tx interfaces.IRepository) error) (err error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}

		if err != nil {
			if rbErr := tx.Rollback().Error; rbErr != nil && !isTxAlreadyClosed(rbErr) {
				err = fmt.Errorf("rollback error: %w; original error: %v", rbErr, err)
			}
			return
		}

		if cErr := tx.Commit().Error; cErr != nil {
			err = cErr
		}
	}()

	return fn(&BaseRepository{db: tx})
}

// isTxAlreadyClosed reports whether err means the transaction was already
// committed or rolled back, in which case a second rollback is a safe no-op.
func isTxAlreadyClosed(err error) bool {
	return errors.Is(err, sql.ErrTxDone) || errors.Is(err, gorm.ErrInvalidTransaction)
}

func (r *BaseRepository) FindByID(ctx context.Context, entityInstance entities.IEntity, id any, params *entities.QueryParams) error {
	query := buildBaseQuery(r.db.WithContext(ctx), entityInstance, params)
	result := query.Where("id = ?", id).First(entityInstance)
	return result.Error
}

func (r *BaseRepository) FindOne(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error {
	query := buildBaseQuery(r.db.WithContext(ctx).Model(entityInstance), entityInstance, params)
	dbResult := query.First(entityInstance)
	return dbResult.Error
}

func (r *BaseRepository) Find(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (*entities.PaginatedResult, error) {
	if params == nil {
		params = &entities.QueryParams{}
	}
	sortField := params.Sort
	if sortField == "" {
		sortField = "id desc"
	}

	query := buildBaseQuery(r.db.WithContext(ctx).Model(entityInstance), entityInstance, params)

	totalRecords := int64(0)
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, err
	}

	result, offset, limit := handlePagination(totalRecords, params.Page, params.Limit)
	query = query.Offset(offset).Limit(limit).Order(sortField)

	entityList := createEntitySlice(entityInstance)
	if err := query.Find(entityList).Error; err != nil {
		return nil, err
	}

	result.Data = entityList
	result.TotalItems = totalRecords
	return result, nil
}

func (r *BaseRepository) FindAll(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (any, error) {
	if params == nil {
		params = &entities.QueryParams{}
	}
	sortField := params.Sort
	if sortField == "" {
		sortField = "id desc"
	}
	query := buildBaseQuery(r.db.WithContext(ctx).Model(entityInstance), entityInstance, params)
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	entityList := createEntitySlice(entityInstance)
	if err := query.Order(sortField).Find(entityList).Error; err != nil {
		return nil, err
	}

	return reflect.ValueOf(entityList).Elem().Interface(), nil
}

func (r *BaseRepository) Save(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error {
	query := r.db.WithContext(ctx)
	if params != nil && params.SelectFields != nil && len(params.SelectFields) > 0 {
		query = query.Select(params.SelectFields)
	}

	return query.Create(entityInstance).Error
}

func (r *BaseRepository) Update(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error {
	query := buildBaseQuery(r.db.WithContext(ctx).Model(entityInstance), entityInstance, params)

	id := entityInstance.GetID()
	if !isEmptyID(id) {
		query = query.Where("id = ?", id)
	}

	if params != nil && params.Query.From != "" {
		name, alias := splitNameAlias(params.Query.From)
		query = query.Clauses(clause.From{Tables: []clause.Table{{Name: name, Alias: alias}}})
	}

	dbResult := query.Updates(buildUpdateMap(entityInstance, params.UpdateFields, params.IncrementFields))
	if dbResult.Error != nil {
		return dbResult.Error
	}

	if dbResult.RowsAffected == 0 {
		return fmt.Errorf("no rows affected for ID %v", id)
	}

	if !isEmptyID(id) {
		if err := r.FindByID(ctx, entityInstance, id, &entities.QueryParams{
			SelectFields: params.SelectFields,
		}); err != nil {
			return fmt.Errorf("failed to find updated entity by ID %v: %w", id, err)
		}
	}

	return nil
}

func (r *BaseRepository) Delete(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error {
	query := buildBaseQuery(r.db.WithContext(ctx).Model(entityInstance), entityInstance, params)
	id := entityInstance.GetID()
	if !isEmptyID(id) {
		query = query.Where("id = ?", id)
	}
	return query.Delete(entityInstance).Error
}

func (r *BaseRepository) Raw(ctx context.Context, entityInstance any, query string, values ...any) error {
	return r.db.WithContext(ctx).Raw(query, values...).Scan(entityInstance).Error
}

func (r *BaseRepository) Verify(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (bool, error) {
	rawSQL := "SELECT EXISTS (SELECT 1 FROM " + entityInstance.TableName()
	if params.Query.Joins != "" {
		rawSQL += " " + params.Query.Joins
	}
	if params.Query.Filters != "" {
		rawSQL += " WHERE " + params.Query.Filters
	}
	rawSQL += ")"
	var exists bool
	err := r.db.WithContext(ctx).Raw(rawSQL, params.Query.Values...).Scan(&exists).Error
	return exists, err
}

func (r *BaseRepository) Count(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (int64, error) {
	query := buildBaseQuery(r.db.WithContext(ctx).Model(entityInstance), entityInstance, params)
	var count int64
	return count, query.Count(&count).Error
}

func (r *BaseRepository) CountDistinct(ctx context.Context, entityInstance entities.IEntity, field string, params *entities.QueryParams) (int64, error) {
	var results []any

	query := buildBaseQuery(r.db.WithContext(ctx).Model(entityInstance), entityInstance, params)

	err := query.Distinct().Pluck(field, &results).Error
	if err != nil {
		return 0, err
	}

	return int64(len(results)), nil
}

func (r *BaseRepository) SaveOrUpdate(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error {
	if len(params.ConflictColumns) == 0 {
		return errors.New("conflict columns must be provided")
	}

	table := entityInstance.TableName()
	exprUpdateMap := buildUpdateMapExprOnlyWithTable(table, params.IncrementFields)

	for _, field := range params.UpdateFields {
		if _, exists := exprUpdateMap[field]; !exists {
			exprUpdateMap[field] = gorm.Expr(fmt.Sprintf("excluded.%s", field))
		}
	}

	return r.db.WithContext(ctx).
		Model(entityInstance).
		Clauses(clause.OnConflict{
			Columns:   toClauseColumns(params.ConflictColumns),
			DoUpdates: clause.Assignments(exprUpdateMap),
		}).
		Create(entityInstance).Error
}

func (r *BaseRepository) BulkSave(ctx context.Context, entityList []entities.IEntity) error {
	if len(entityList) == 0 {
		return nil
	}
	modelSlice := toModelSlice(entityList)
	return r.db.WithContext(ctx).Create(modelSlice.Interface()).Error
}

func (r *BaseRepository) BulkUpdate(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error {
	query := buildBaseQuery(
		r.db.WithContext(ctx).Model(entityInstance),
		entityInstance,
		params,
	)

	if params != nil && len(params.UpdateFields) > 0 {
		query = query.Select(params.UpdateFields)
	}

	return query.Updates(buildUpdateMap(entityInstance, params.UpdateFields, params.IncrementFields)).Error
}

func (r *BaseRepository) BulkDelete(ctx context.Context, entity entities.IEntity, ids []any, params *entities.QueryParams) error {
	if len(ids) == 0 {
		return errors.New("no IDs provided for bulk delete")
	}
	query := buildBaseQuery(r.db.WithContext(ctx).Model(entity), entity, params)
	return query.Where("id IN ?", ids).Delete(nil).Error
}

func ToIEntitySlice(list any) []entities.IEntity {
	s := reflect.ValueOf(list)
	if s.Kind() != reflect.Slice {
		panic("ToIEntitySlice expects a slice")
	}
	result := make([]entities.IEntity, s.Len())
	for i := 0; i < s.Len(); i++ {
		item := s.Index(i).Interface()
		entity, ok := item.(entities.IEntity)
		if !ok {
			panic(fmt.Sprintf("item at index %d does not implement IEntity", i))
		}
		result[i] = entity
	}
	return result
}
