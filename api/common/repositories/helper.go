package repositories

import (
	"fmt"
	"reflect"
	"slices"
	"social-engine/common/repositories/entities"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func extractGormColumn(tag string) string {
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}
	return ""
}

func createEntitySlice(entityInstance interface{}) interface{} {
	elemType := reflect.TypeOf(entityInstance).Elem()
	sliceType := reflect.SliceOf(reflect.PointerTo(elemType))
	return reflect.New(sliceType).Interface()
}

func handlePagination(totalRecords int64, page, limit int) (*entities.PaginatedResult, int, int) {
	// checks the limit range and set a default when not provided or out of range
	if limit <= 0 {
		limit = 10
	}

	// the ceiling is 1000
	if limit > 1000 {
		limit = 1000
	}

	// calculates the totalPages
	totalPages := totalRecords / int64(limit)
	if totalRecords%int64(limit) != 0 {
		totalPages++
	}

	pageControl := 0
	if page > 0 && page <= int(totalPages) {
		pageControl = page
	} else {
		if page <= 0 {
			pageControl = 1
		}

		if page > int(totalPages) {
			pageControl = int(totalPages)
		}
	}

	nextPageControl := 0
	if pageControl+1 <= int(totalPages) {
		nextPageControl = pageControl + 1
	} else {
		nextPageControl = pageControl
	}

	previousPageControl := 0
	if pageControl-1 > 0 {
		previousPageControl = pageControl - 1
	} else {
		previousPageControl = pageControl
	}

	// calcutale the offset
	offset := (pageControl - 1) * limit

	result := &entities.PaginatedResult{
		TotalPages:   int64(totalPages),
		CurrentPage:  int64(pageControl),
		PreviousPage: int64(previousPageControl),
		NextPage:     int64(nextPageControl),
	}

	return result, offset, limit
}

func toClauseColumns(cols []string) []clause.Column {
	result := make([]clause.Column, len(cols))
	for i, col := range cols {
		result[i] = clause.Column{Name: col}
	}
	return result
}

func buildUpdateMap(entity any, updateFields []string, incrementFields []string) map[string]interface{} {
	updateMap := make(map[string]interface{})
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.IsValid() || !fieldVal.CanInterface() {
			continue
		}

		value := fieldVal.Interface()
		column := field.Name

		if tag := field.Tag.Get("gorm"); strings.Contains(tag, "column:") {
			for _, part := range strings.Split(tag, ";") {
				if strings.HasPrefix(part, "column:") {
					column = strings.TrimPrefix(part, "column:")
					break
				}
			}
		}

		// Incremento tem prioridade sobre update direto
		if slices.Contains(incrementFields, column) {
			updateMap[column] = gorm.Expr(fmt.Sprintf("%s + ?", column), value)
		} else if slices.Contains(updateFields, column) {
			updateMap[column] = value
		}
	}

	return updateMap
}

func buildBaseQuery(query *gorm.DB, entity entities.IEntity, params *entities.QueryParams) *gorm.DB {
	if params == nil {
		params = &entities.QueryParams{}
	}
	if params.Query.Joins != "" {
		query = query.Joins(params.Query.Joins, params.Query.JoinValues...)
	}
	if params.Query.Filters != "" {
		query = query.Where(params.Query.Filters, params.Query.Values...)
	}

	if len(params.SelectFields) > 0 {
		query = query.Select(params.SelectFields)
	}

	for _, field := range params.Search {
		if field.Value == "" {
			continue
		}
		query = query.Where(fmt.Sprintf("%s ILIKE ?", field.Field), "%"+field.Value+"%")
	}

	if params.Sort != "" {
		query = query.Order(params.Sort)
	}

	for _, association := range entity.LoadAssociations() {
		query = query.Preload(association)
	}
	return query
}

func toModelSlice(list []entities.IEntity) reflect.Value {
	if len(list) == 0 {
		panic("toModelSlice expects a non-empty IEntity slice")
	}
	typeOfEntity := reflect.TypeOf(list[0])
	slice := reflect.MakeSlice(reflect.SliceOf(typeOfEntity), len(list), len(list))
	for i, item := range list {
		slice.Index(i).Set(reflect.ValueOf(item))
	}
	return slice
}

func buildUpdateMapExprOnlyWithTable(tableName string, incrementFields []string) map[string]interface{} {
	exprs := make(map[string]interface{})
	for _, column := range incrementFields {
		exprs[column] = gorm.Expr(fmt.Sprintf("%s.%s + excluded.%s", tableName, column, column))
	}
	return exprs
}

func splitNameAlias(from string) (name, alias string) {
	parts := strings.Fields(from)
	name = parts[0]
	if len(parts) > 1 {
		alias = parts[1]
	}
	return
}
