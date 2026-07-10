package repositories

import (
	"reflect"
	"social-engine/common/repositories/entities"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type testEntity struct {
	ID         int    `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	Count      int    `gorm:"column:count"`
	NoTag      string
	unexported string
}

func (m *testEntity) TableName() string {
	return "mock_entities"
}

func (m *testEntity) LoadAssociations() []string {
	return []string{"Profile", "Cards"}
}

func (m *testEntity) GetID() any {
	return ""
}

func TestExtractGormColumn(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		want string
	}{
		{
			name: "should extract column from single tag",
			tag:  "column:user_id",
			want: "user_id",
		},
		{
			name: "should extract column from multiple tags",
			tag:  "primaryKey;column:full_name;not null",
			want: "full_name",
		},
		{
			name: "should return empty when column tag does not exist",
			tag:  "primaryKey;not null",
			want: "",
		},
		{
			name: "should return empty for empty tag",
			tag:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractGormColumn(tt.tag)
			if got != tt.want {
				t.Fatalf("extractGormColumn() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCreateEntitySlice(t *testing.T) {
	type sample struct {
		ID int
	}

	tests := []struct {
		name        string
		entity      interface{}
		wantType    reflect.Type
		wantIsPtr   bool
		wantIsSlice bool
	}{
		{
			name:        "should create pointer to slice of pointers",
			entity:      &sample{},
			wantType:    reflect.TypeOf(&[]*sample{}),
			wantIsPtr:   true,
			wantIsSlice: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createEntitySlice(tt.entity)

			if reflect.TypeOf(got) != tt.wantType {
				t.Fatalf("createEntitySlice() type = %v, want %v", reflect.TypeOf(got), tt.wantType)
			}

			gotVal := reflect.ValueOf(got)
			if tt.wantIsPtr && gotVal.Kind() != reflect.Ptr {
				t.Fatalf("createEntitySlice() kind = %v, want ptr", gotVal.Kind())
			}

			if tt.wantIsSlice && gotVal.Elem().Kind() != reflect.Slice {
				t.Fatalf("createEntitySlice() elem kind = %v, want slice", gotVal.Elem().Kind())
			}
		})
	}
}

func TestHandlePagination(t *testing.T) {
	tests := []struct {
		name         string
		totalRecords int64
		page         int
		limit        int
		wantResult   *entities.PaginatedResult
		wantOffset   int
		wantLimit    int
	}{
		{
			name:         "should use default limit when limit is zero",
			totalRecords: 95,
			page:         1,
			limit:        0,
			wantResult: &entities.PaginatedResult{
				TotalPages:   10,
				CurrentPage:  1,
				PreviousPage: 1,
				NextPage:     2,
			},
			wantOffset: 0,
			wantLimit:  10,
		},
		{
			name:         "should cap limit at 1000",
			totalRecords: 2500,
			page:         1,
			limit:        5000,
			wantResult: &entities.PaginatedResult{
				TotalPages:   3,
				CurrentPage:  1,
				PreviousPage: 1,
				NextPage:     2,
			},
			wantOffset: 0,
			wantLimit:  1000,
		},
		{
			name:         "should calculate middle page correctly",
			totalRecords: 100,
			page:         3,
			limit:        10,
			wantResult: &entities.PaginatedResult{
				TotalPages:   10,
				CurrentPage:  3,
				PreviousPage: 2,
				NextPage:     4,
			},
			wantOffset: 20,
			wantLimit:  10,
		},
		{
			name:         "should normalize page to 1 when page is zero",
			totalRecords: 100,
			page:         0,
			limit:        10,
			wantResult: &entities.PaginatedResult{
				TotalPages:   10,
				CurrentPage:  1,
				PreviousPage: 1,
				NextPage:     2,
			},
			wantOffset: 0,
			wantLimit:  10,
		},
		{
			name:         "should normalize page to last page when page exceeds total pages",
			totalRecords: 95,
			page:         99,
			limit:        10,
			wantResult: &entities.PaginatedResult{
				TotalPages:   10,
				CurrentPage:  10,
				PreviousPage: 9,
				NextPage:     10,
			},
			wantOffset: 90,
			wantLimit:  10,
		},
		{
			name:         "should handle zero records",
			totalRecords: 0,
			page:         1,
			limit:        10,
			wantResult: &entities.PaginatedResult{
				TotalPages:   0,
				CurrentPage:  0,
				PreviousPage: 0,
				NextPage:     0,
			},
			wantOffset: -10,
			wantLimit:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotOffset, gotLimit := handlePagination(tt.totalRecords, tt.page, tt.limit)

			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Fatalf("handlePagination() result = %+v, want %+v", gotResult, tt.wantResult)
			}

			if gotOffset != tt.wantOffset {
				t.Fatalf("handlePagination() offset = %d, want %d", gotOffset, tt.wantOffset)
			}

			if gotLimit != tt.wantLimit {
				t.Fatalf("handlePagination() limit = %d, want %d", gotLimit, tt.wantLimit)
			}
		})
	}
}

func TestToClauseColumns(t *testing.T) {
	tests := []struct {
		name string
		cols []string
		want []string
	}{
		{
			name: "should convert multiple columns",
			cols: []string{"id", "name", "created_at"},
			want: []string{"id", "name", "created_at"},
		},
		{
			name: "should handle empty slice",
			cols: []string{},
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toClauseColumns(tt.cols)

			if len(got) != len(tt.want) {
				t.Fatalf("toClauseColumns() len = %d, want %d", len(got), len(tt.want))
			}

			for i := range got {
				if got[i].Name != tt.want[i] {
					t.Fatalf("toClauseColumns()[%d].Name = %q, want %q", i, got[i].Name, tt.want[i])
				}
			}
		})
	}
}

func TestBuildUpdateMap(t *testing.T) {
	tests := []struct {
		name            string
		entity          any
		updateFields    []string
		incrementFields []string
		assert          func(t *testing.T, got map[string]interface{})
	}{
		{
			name: "should build update map with direct fields only",
			entity: &testEntity{
				ID:    1,
				Name:  "Pablo",
				Count: 7,
				NoTag: "value",
			},
			updateFields:    []string{"id", "name", "NoTag"},
			incrementFields: nil,
			assert: func(t *testing.T, got map[string]interface{}) {
				if len(got) != 3 {
					t.Fatalf("len(updateMap) = %d, want 3", len(got))
				}
				if got["id"] != 1 {
					t.Fatalf(`got["id"] = %v, want 1`, got["id"])
				}
				if got["name"] != "Pablo" {
					t.Fatalf(`got["name"] = %v, want Pablo`, got["name"])
				}
				if got["NoTag"] != "value" {
					t.Fatalf(`got["NoTag"] = %v, want value`, got["NoTag"])
				}
			},
		},
		{
			name: "should prioritize increment fields over update fields",
			entity: &testEntity{
				Count: 10,
			},
			updateFields:    []string{"count"},
			incrementFields: []string{"count"},
			assert: func(t *testing.T, got map[string]interface{}) {
				val, ok := got["count"]
				if !ok {
					t.Fatalf(`expected key "count" to exist`)
				}
				if _, ok := val.(clauseExprMatcher); ok {
					return
				}
				expr, ok := val.(clause.Expr)
				if !ok {
					t.Fatalf(`got["count"] type = %T, want clause.Expr`, val)
				}
				if expr.SQL != "count + ?" {
					t.Fatalf("expr.SQL = %q, want %q", expr.SQL, "count + ?")
				}
				if len(expr.Vars) != 1 || expr.Vars[0] != 10 {
					t.Fatalf("expr.Vars = %#v, want []interface{}{10}", expr.Vars)
				}
			},
		},
		{
			name: "should ignore fields not listed",
			entity: &testEntity{
				ID:    1,
				Name:  "Pablo",
				Count: 10,
			},
			updateFields:    []string{"name"},
			incrementFields: nil,
			assert: func(t *testing.T, got map[string]interface{}) {
				if len(got) != 1 {
					t.Fatalf("len(updateMap) = %d, want 1", len(got))
				}
				if got["name"] != "Pablo" {
					t.Fatalf(`got["name"] = %v, want Pablo`, got["name"])
				}
			},
		},
		{
			name: "should work with non-pointer entity",
			entity: testEntity{
				Name: "Direct",
			},
			updateFields:    []string{"name"},
			incrementFields: nil,
			assert: func(t *testing.T, got map[string]interface{}) {
				if got["name"] != "Direct" {
					t.Fatalf(`got["name"] = %v, want Direct`, got["name"])
				}
			},
		},
		{
			name: "should not include unexported field",
			entity: &testEntity{
				Name:       "Visible",
				unexported: "hidden",
			},
			updateFields:    []string{"name", "unexported"},
			incrementFields: nil,
			assert: func(t *testing.T, got map[string]interface{}) {
				if _, ok := got["unexported"]; ok {
					t.Fatalf(`did not expect "unexported" field in update map`)
				}
				if got["name"] != "Visible" {
					t.Fatalf(`got["name"] = %v, want Visible`, got["name"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildUpdateMap(tt.entity, tt.updateFields, tt.incrementFields)
			tt.assert(t, got)
		})
	}
}

func newDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("failed to open dry-run db: %v", err)
	}

	return db
}

func normalizeSQL(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func TestBuildBaseQuery(t *testing.T) {
	tests := []struct {
		name              string
		params            *entities.QueryParams
		wantSQLContains   []string
		wantSQLNotContain []string
		wantVars          []interface{}
		wantPreloads      []string
	}{
		{
			name:            "should handle nil params and only preload associations",
			params:          nil,
			wantSQLContains: []string{"SELECT * FROM `mock_entities`"},
			wantPreloads:    []string{"Profile", "Cards"},
		},
		{
			name: "should apply joins and filters",
			params: &entities.QueryParams{
				Query: entities.Query{
					Joins:      "LEFT JOIN profiles ON profiles.user_id = mock_entities.id",
					Filters:    "mock_entities.id = ?",
					JoinValues: []interface{}{},
					Values:     []interface{}{10},
				},
			},
			wantSQLContains: []string{
				"LEFT JOIN profiles ON profiles.user_id = mock_entities.id",
				"WHERE mock_entities.id = ?",
			},
			wantVars:     []interface{}{10},
			wantPreloads: []string{"Profile", "Cards"},
		},
		{
			name: "should apply select fields",
			params: &entities.QueryParams{
				SelectFields: []string{"id", "name"},
			},
			wantSQLContains: []string{
				"SELECT `id`,`name` FROM `mock_entities`",
			},
			wantPreloads: []string{"Profile", "Cards"},
		},
		{
			name: "should apply search filters and ignore empty values",
			params: &entities.QueryParams{
				Search: []entities.SearchField{
					{Field: "name", Value: "pablo"},
					{Field: "email", Value: ""},
					{Field: "nickname", Value: "banker"},
				},
			},
			wantSQLContains: []string{
				"WHERE name ILIKE ? AND nickname ILIKE ?",
			},
			wantSQLNotContain: []string{
				"email ILIKE ?",
			},
			wantVars: []interface{}{
				"%pablo%",
				"%banker%",
			},
			wantPreloads: []string{"Profile", "Cards"},
		},
		{
			name: "should apply sort",
			params: &entities.QueryParams{
				Sort: "name asc",
			},
			wantSQLContains: []string{
				"ORDER BY name asc",
			},
			wantPreloads: []string{"Profile", "Cards"},
		},
		{
			name: "should apply full query with joins filters select search and sort",
			params: &entities.QueryParams{
				Query: entities.Query{
					Joins:   "INNER JOIN guilds ON guilds.id = mock_entities.id",
					Filters: "mock_entities.active = ?",
					Values:  []interface{}{true},
				},
				SelectFields: []string{"mock_entities.id", "mock_entities.name"},
				Search: []entities.SearchField{
					{Field: "mock_entities.name", Value: "hero"},
				},
				Sort: "mock_entities.name desc",
			},
			wantSQLContains: []string{
				"INNER JOIN guilds ON guilds.id = mock_entities.id",
				"WHERE mock_entities.active = ? AND mock_entities.name ILIKE ?",
				"ORDER BY mock_entities.name desc",
			},
			wantVars: []interface{}{
				true,
				"%hero%",
			},
			wantPreloads: []string{"Profile", "Cards"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newDryRunDB(t)

			baseQuery := db.Model(&testEntity{})

			got := buildBaseQuery(baseQuery, &testEntity{}, tt.params)

			var result []testEntity
			executed := got.Find(&result)

			sql := normalizeSQL(executed.Statement.SQL.String())

			for _, expected := range tt.wantSQLContains {
				if !strings.Contains(sql, normalizeSQL(expected)) {
					t.Fatalf("expected SQL to contain %q\nSQL: %s", expected, sql)
				}
			}

			for _, unexpected := range tt.wantSQLNotContain {
				if strings.Contains(sql, normalizeSQL(unexpected)) {
					t.Fatalf("expected SQL to NOT contain %q\nSQL: %s", unexpected, sql)
				}
			}

			if len(executed.Statement.Vars) != len(tt.wantVars) {
				t.Fatalf("vars length = %d, want %d\nvars = %#v", len(executed.Statement.Vars), len(tt.wantVars), executed.Statement.Vars)
			}

			for i := range tt.wantVars {
				if executed.Statement.Vars[i] != tt.wantVars[i] {
					t.Fatalf("var[%d] = %#v, want %#v", i, executed.Statement.Vars[i], tt.wantVars[i])
				}
			}

			if len(got.Statement.Preloads) != len(tt.wantPreloads) {
				t.Fatalf("preloads length = %d, want %d\npreloads = %#v", len(got.Statement.Preloads), len(tt.wantPreloads), got.Statement.Preloads)
			}

			for _, preload := range tt.wantPreloads {
				if _, ok := got.Statement.Preloads[preload]; !ok {
					t.Fatalf("expected preload %q to exist, got %#v", preload, got.Statement.Preloads)
				}
			}
		})
	}
}

type clauseExprMatcher interface{}

func TestToModelSlice(t *testing.T) {
	tests := []struct {
		name      string
		list      []entities.IEntity
		wantLen   int
		wantType  string
		wantPanic bool
	}{
		{
			name: "should convert IEntity slice to reflected slice",
			list: []entities.IEntity{
				&testEntity{},
				&testEntity{},
			},
			wantLen:   2,
			wantType:  "[]entities.IEntity",
			wantPanic: false,
		},
		{
			name:      "should panic on empty slice",
			list:      []entities.IEntity{},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Fatalf("expected panic, got none")
					}
				}()
			}

			got := toModelSlice(tt.list)

			if tt.wantPanic {
				return
			}

			if got.Len() != tt.wantLen {
				t.Fatalf("toModelSlice().Len() = %d, want %d", got.Len(), tt.wantLen)
			}

			if got.Type().String() != "[]entities.IEntity" && got.Len() == 0 {
				t.Fatalf("unexpected slice type: %v", got.Type())
			}
		})
	}
}

func TestBuildUpdateMapExprOnlyWithTable(t *testing.T) {
	tests := []struct {
		name            string
		tableName       string
		incrementFields []string
		wantKeys        []string
	}{
		{
			name:            "should build expressions for all increment fields",
			tableName:       "users",
			incrementFields: []string{"xp", "coins"},
			wantKeys:        []string{"xp", "coins"},
		},
		{
			name:            "should return empty map for no fields",
			tableName:       "users",
			incrementFields: []string{},
			wantKeys:        []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildUpdateMapExprOnlyWithTable(tt.tableName, tt.incrementFields)

			if len(got) != len(tt.wantKeys) {
				t.Fatalf("len(got) = %d, want %d", len(got), len(tt.wantKeys))
			}

			for _, key := range tt.wantKeys {
				val, ok := got[key]
				if !ok {
					t.Fatalf("expected key %q to exist", key)
				}
				if reflect.TypeOf(val).String() != "clause.Expr" {
					t.Fatalf("got[%q] type = %T, want clause.Expr", key, val)
				}
			}
		})
	}
}

func TestSplitNameAlias(t *testing.T) {
	tests := []struct {
		name      string
		from      string
		wantName  string
		wantAlias string
	}{
		{
			name:      "should split name and alias",
			from:      "users u",
			wantName:  "users",
			wantAlias: "u",
		},
		{
			name:      "should return only name when alias does not exist",
			from:      "users",
			wantName:  "users",
			wantAlias: "",
		},
		{
			name:      "should ignore extra fields after alias",
			from:      "users u extra",
			wantName:  "users",
			wantAlias: "u",
		},
		{
			name:      "should handle multiple spaces",
			from:      "   users    u   ",
			wantName:  "users",
			wantAlias: "u",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotAlias := splitNameAlias(tt.from)

			if gotName != tt.wantName {
				t.Fatalf("splitNameAlias() name = %q, want %q", gotName, tt.wantName)
			}
			if gotAlias != tt.wantAlias {
				t.Fatalf("splitNameAlias() alias = %q, want %q", gotAlias, tt.wantAlias)
			}
		})
	}
}
