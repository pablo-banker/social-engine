package repositories

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"social-engine/common/repositories/constants"
	"social-engine/common/repositories/entities"
	"social-engine/common/repositories/interfaces"
	"time"

	"github.com/shopspring/decimal"
)

var _ interfaces.IRepository = (*MockRepository)(nil)

type MockRepository struct {
	expectedResults []MockPayload
	indexHandler    int
}

type MockPayload struct {
	Name           string
	Type           constants.RepositoryType
	Params         *entities.QueryParams
	ExpectedResult any
	ExpectedError  error
}

func NewMockRepository(expectedResults []MockPayload) *MockRepository {
	return &MockRepository{
		expectedResults: expectedResults,
		indexHandler:    0,
	}
}

func (m *MockRepository) WithTableName(tableName string) interfaces.IRepository {
	return m
}

func (m *MockRepository) BeginTx(ctx context.Context) (interfaces.IRepository, error) {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryBeginTx); err != nil {
		return nil, m.formatError(err, constants.RepositoryBeginTx)
	}

	if data.ExpectedError != nil {
		return nil, data.ExpectedError
	}

	return m, nil
}

func (m *MockRepository) Rollback(err error) error {
	defer m.incrementIndex()
	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryRollback); err != nil {
		return m.formatError(err, constants.RepositoryRollback)
	}

	if err != nil {
		return err
	}

	return m.handleTypeError(data)
}

func (m *MockRepository) Commit() error {
	defer m.incrementIndex()
	data := m.expectedResults[m.indexHandler]

	if err := m.validateType(data.Type, constants.RepositoryCommit); err != nil {
		return m.formatError(err, constants.RepositoryCommit)
	}

	return m.handleTypeError(data)
}

// WithTransaction mirrors the real implementation against the ordered payload
// protocol: it consumes a BeginTx payload, runs fn (which consumes the inner
// operation payloads), then consumes a Commit payload on success or a Rollback
// payload when fn returns an error.
func (m *MockRepository) WithTransaction(ctx context.Context, fn func(tx interfaces.IRepository) error) error {
	txRepo, err := m.BeginTx(ctx)
	if err != nil {
		return err
	}

	if fnErr := fn(txRepo); fnErr != nil {
		_ = m.Rollback(fnErr)
		return fnErr
	}

	return m.Commit()
}

func (m *MockRepository) Ping(ctx context.Context) error {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryPing); err != nil {
		return m.formatError(err, constants.RepositoryPing)
	}

	if data.ExpectedError != nil {
		return data.ExpectedError
	}

	return nil
}

func (m *MockRepository) Find(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (*entities.PaginatedResult, error) {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryFind); err != nil {
		return nil, m.formatError(err, constants.RepositoryFind)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return nil, m.formatError(err, constants.RepositoryFind)
	}

	err := m.handleTypeError(m.expectedResults[m.indexHandler])
	if err != nil {
		return nil, err
	}

	paginatedResult := &entities.PaginatedResult{
		Data:         m.expectedResults[m.indexHandler],
		TotalPages:   1,
		CurrentPage:  1,
		PreviousPage: 1,
		NextPage:     1,
	}

	return paginatedResult, nil
}

func (m *MockRepository) FindByID(ctx context.Context, entityInstance entities.IEntity, id any, params *entities.QueryParams) error {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryFindByID); err != nil {
		return m.formatError(err, constants.RepositoryFindByID)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return m.formatError(err, constants.RepositoryFindByID)
	}

	return m.handleTypeTarget(entityInstance, data)
}

func (m *MockRepository) FindOne(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryFindOne); err != nil {
		return m.formatError(err, constants.RepositoryFindOne)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return m.formatError(err, constants.RepositoryFindOne)
	}

	return m.handleTypeTarget(entityInstance, m.expectedResults[m.indexHandler])
}

func (m *MockRepository) FindAll(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (any, error) {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryFindAll); err != nil {
		return nil, m.formatError(err, constants.RepositoryFindAll)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return nil, m.formatError(err, constants.RepositoryFindAll)
	}

	return m.handleTypeAny(data)
}

func (m *MockRepository) Save(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositorySave); err != nil {
		return m.formatError(err, constants.RepositorySave)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return m.formatError(err, constants.RepositorySave)
	}

	return m.handleTypeTarget(entityInstance, data)
}

func (m *MockRepository) Update(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryUpdate); err != nil {
		return m.formatError(err, constants.RepositoryUpdate)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return m.formatError(err, constants.RepositoryUpdate)
	}

	if err := m.validateUpdateFields(entityInstance, data.ExpectedResult, params); err != nil {
		return m.formatError(err, constants.RepositoryUpdate)
	}

	return m.handleTypeTarget(entityInstance, data)
}

func (m *MockRepository) Delete(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryDelete); err != nil {
		return m.formatError(err, constants.RepositoryDelete)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return m.formatError(err, constants.RepositoryDelete)
	}

	return data.ExpectedError
}

func (m *MockRepository) Raw(ctx context.Context, targetEntityInstance any, sql string, values ...any) error {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryRaw); err != nil {
		return m.formatError(err, constants.RepositoryRaw)
	}

	result, err := m.handleTypeAny(data)
	if err != nil {
		return m.formatError(err, constants.RepositoryRaw)
	}

	if targetEntityInstance == nil {
		return m.handleTypeError(data)
	}

	targetVal := reflect.ValueOf(targetEntityInstance)
	resultVal := reflect.ValueOf(result)

	// Ensure the target is a pointer
	if targetVal.Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer to the result target, got %s", targetVal.Kind())
	}

	elem := targetVal.Elem()

	// If the target is a slice pointer
	if elem.Kind() == reflect.Slice {
		if !resultVal.Type().AssignableTo(elem.Type()) {
			return fmt.Errorf("cannot assign result of type %s to target of type %s", resultVal.Type(), elem.Type())
		}
		elem.Set(resultVal)
		return nil
	}

	// If the target is a single struct or primitive
	if elem.CanSet() {
		if !resultVal.Type().AssignableTo(elem.Type()) {
			return fmt.Errorf("cannot assign result of type %s to target of type %s", resultVal.Type(), elem.Type())
		}
		elem.Set(resultVal)
		return nil
	}

	// Fallback: if somehow the pointer is invalid, try assigning to a slice directly
	if elem.Kind() == reflect.Invalid && targetVal.Kind() == reflect.Slice {
		if resultVal.Type().AssignableTo(targetVal.Type()) {
			targetVal.Set(resultVal)
			return nil
		}
		return m.formatError(fmt.Errorf("fallback failed, incompatible slice types: %s vs %s", resultVal.Type(), targetVal.Type()), constants.RepositoryRaw)
	}

	return m.formatError(fmt.Errorf("unsupported assignment from %s to %s", resultVal.Type(), targetVal.Type()), constants.RepositoryRaw)
}

func (m *MockRepository) Verify(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (bool, error) {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryVerify); err != nil {
		return false, m.formatError(err, constants.RepositoryVerify)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return false, m.formatError(err, constants.RepositoryVerify)
	}

	return m.handleBoolean(data)
}

func (m *MockRepository) Count(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (int64, error) {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryCount); err != nil {
		return 0, m.formatError(err, constants.RepositoryCount)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return 0, m.formatError(err, constants.RepositoryCount)
	}

	result, err := m.handleTypeAny(data)
	if err != nil {
		return 0, m.formatError(err, constants.RepositoryCount)
	}

	if count, ok := result.(int); ok {
		return int64(count), nil
	}

	return 0, fmt.Errorf("expected count to be of type int, got %T", result)
}

func (m *MockRepository) CountDistinct(ctx context.Context, entityInstance entities.IEntity, field string, params *entities.QueryParams) (int64, error) {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryCountDistinct); err != nil {
		return 0, m.formatError(err, constants.RepositoryCountDistinct)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return 0, m.formatError(err, constants.RepositoryCountDistinct)
	}

	result, err := m.handleTypeAny(data)
	if err != nil {
		return 0, m.formatError(err, constants.RepositoryCountDistinct)
	}

	if count, ok := result.(int); ok {
		return int64(count), nil
	}

	return 0, fmt.Errorf("expected count to be of type int, got %T", result)
}

func (m *MockRepository) SaveOrUpdate(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositorySaveOrUpdate); err != nil {
		return m.formatError(err, constants.RepositorySaveOrUpdate)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return m.formatError(err, constants.RepositorySaveOrUpdate)
	}

	return m.handleTypeError(data)
}

func (m *MockRepository) BulkSave(ctx context.Context, entities []entities.IEntity) error {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryBulkSave); err != nil {
		return m.formatError(err, constants.RepositoryBulkSave)
	}

	return m.handleTypeError(data)
}

func (m *MockRepository) BulkUpdate(ctx context.Context, entity entities.IEntity, params *entities.QueryParams) error {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryBulkUpdate); err != nil {
		return m.formatError(err, constants.RepositoryBulkUpdate)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return m.formatError(err, constants.RepositoryBulkUpdate)
	}

	return m.handleTypeError(data)
}

func (m *MockRepository) BulkDelete(ctx context.Context, entity entities.IEntity, ids []any, params *entities.QueryParams) error {
	defer m.incrementIndex()

	data := m.expectedResults[m.indexHandler]
	if err := m.validateType(data.Type, constants.RepositoryBulkDelete); err != nil {
		return m.formatError(err, constants.RepositoryBulkDelete)
	}

	if err := m.validateQueryParams(params, data.Params); err != nil {
		return m.formatError(err, constants.RepositoryBulkDelete)
	}

	return m.handleTypeError(data)
}

// ------- Secondary Methods -------

func (m *MockRepository) incrementIndex() {
	m.indexHandler++
}

func (m *MockRepository) handleTypeError(data MockPayload) error {
	return data.ExpectedError
}

func (m *MockRepository) handleTypeAny(data MockPayload) (any, error) {
	if data.ExpectedError != nil {
		return nil, data.ExpectedError
	}

	return data.ExpectedResult, nil
}

func (m *MockRepository) handleBoolean(data MockPayload) (bool, error) {
	if data.ExpectedError != nil {
		return false, data.ExpectedError
	}

	if val, ok := data.ExpectedResult.(bool); ok {
		return val, nil
	}

	return false, fmt.Errorf("expected boolean, got %T", data.ExpectedResult)
}

func (m *MockRepository) handleTypeTarget(target any, data MockPayload) error {
	if data.ExpectedError != nil {
		return data.ExpectedError
	}

	if data.ExpectedResult == nil {
		return fmt.Errorf("expected result is nil for type %s", data.Type)
	}

	targetVal := reflect.ValueOf(target)
	resultVal := reflect.ValueOf(data.ExpectedResult)

	if targetVal.Kind() != reflect.Ptr {
		return fmt.Errorf("handleTypeTarget: target must be a pointer, got %s", targetVal.Kind())
	}

	targetElem := targetVal.Elem()

	switch targetElem.Kind() {
	case reflect.Slice:
		if !resultVal.Type().AssignableTo(targetElem.Type()) {
			return fmt.Errorf("handleTypeTarget: cannot assign result of type %s to slice of type %s", resultVal.Type(), targetElem.Type())
		}
		targetElem.Set(resultVal)
		return nil

	case reflect.Struct:
		if resultVal.Kind() == reflect.Ptr && resultVal.Type().Elem() == targetElem.Type() {
			targetElem.Set(resultVal.Elem())
			return nil
		}

		if !resultVal.Type().AssignableTo(targetElem.Type()) {
			return fmt.Errorf("handleTypeTarget: cannot assign result of type %s to struct of type %s", resultVal.Type(), targetElem.Type())
		}
		targetElem.Set(resultVal)
		return nil

	default:
		if !resultVal.IsValid() || !targetElem.IsValid() {
			return fmt.Errorf("handleTypeTarget: one of the values is invalid (zero reflect.Value)")
		}

		if resultVal.Type().AssignableTo(targetElem.Type()) {
			targetElem.Set(resultVal)
			return nil
		}
		return fmt.Errorf(
			"handleTypeTarget: cannot assign result of type %s to target of type %s",
			resultVal.Type(), targetElem.Type(),
		)
	}
}

func (m *MockRepository) handleTypeEntitySlice(target any, data MockPayload) error {
	targetVal := reflect.ValueOf(target)
	resultVal := reflect.ValueOf(data.ExpectedResult)

	if targetVal.Kind() != reflect.Ptr || targetVal.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("target must be pointer to slice, got %T", target)
	}

	elem := targetVal.Elem()

	// Converte []struct → []interface
	newSlice := reflect.MakeSlice(elem.Type(), resultVal.Len(), resultVal.Len())

	for i := 0; i < resultVal.Len(); i++ {
		item := resultVal.Index(i).Interface()

		entity, ok := item.(entities.IEntity)
		if !ok {
			return fmt.Errorf("result[%d] does not implement IEntity", i)
		}

		newSlice.Index(i).Set(reflect.ValueOf(entity))
	}

	elem.Set(newSlice)
	return nil
}

func (m *MockRepository) validateQueryParams(got *entities.QueryParams, expected *entities.QueryParams) error {
	gotIsEmpty := got == nil || reflect.ValueOf(*got).IsZero()
	expectedIsEmpty := expected == nil || reflect.ValueOf(*expected).IsZero()

	if gotIsEmpty && expectedIsEmpty {
		return nil
	}

	if gotIsEmpty && !expectedIsEmpty {
		return errors.New("got QueryParams is nil but expected is not")
	}

	if !gotIsEmpty && expectedIsEmpty {
		return errors.New("got QueryParams is not nil but expected is")
	}

	// Query Filters
	if Normalize(got.Query.Filters) != Normalize(expected.Query.Filters) {
		return fmt.Errorf("filters mismatch: got '%s', expected '%s'", got.Query.Filters, expected.Query.Filters)
	}

	// Query Joins
	if Normalize(got.Query.Joins) != Normalize(expected.Query.Joins) {
		return fmt.Errorf("joins mismatch: got '%s', expected '%s'", got.Query.Joins, expected.Query.Joins)
	}

	// Query Joins Values
	if len(got.Query.JoinValues) != len(expected.Query.JoinValues) {
		return fmt.Errorf("joins values mismatch: got '%s', expected '%s'", got.Query.JoinValues, expected.Query.JoinValues)
	}

	// Query Values
	if len(got.Query.Values) != len(expected.Query.Values) {
		return fmt.Errorf("query values length mismatch: got %d, expected %d", len(got.Query.Values), len(expected.Query.Values))
	}

	for i := range got.Query.Values {
		gotVal := got.Query.Values[i]
		expectedVal := expected.Query.Values[i]

		switch gv := gotVal.(type) {
		case decimal.Decimal:
			ev, ok := expectedVal.(decimal.Decimal)
			if !ok || !gv.Equal(ev) {
				return fmt.Errorf("query value[%d] mismatch: got '%v', expected '%v'", i, gv, expectedVal)
			}
		default:
			if !reflect.DeepEqual(gotVal, expectedVal) {
				return fmt.Errorf("query value[%d] mismatch: got '%v', expected '%v'", i, gotVal, expectedVal)
			}
		}
	}

	// Sort
	if got.Sort != expected.Sort {
		return fmt.Errorf("sort mismatch: got '%s', expected '%s'", got.Sort, expected.Sort)
	}

	// Limit
	if got.Limit != expected.Limit {
		return fmt.Errorf("limit mismatch: got %d, expected %d", got.Limit, expected.Limit)
	}

	// Page
	if got.Page != expected.Page {
		return fmt.Errorf("page mismatch: got %d, expected %d", got.Page, expected.Page)
	}

	// UpdateFields
	if !slices.Equal(got.UpdateFields, expected.UpdateFields) {
		return fmt.Errorf("update fields mismatch: got %v, expected %v", got.UpdateFields, expected.UpdateFields)
	}

	// IncrementFields
	if !slices.Equal(got.IncrementFields, expected.IncrementFields) {
		return fmt.Errorf("increment fields mismatch: got %v, expected %v", got.IncrementFields, expected.IncrementFields)
	}

	// ConflictColumns
	if !slices.Equal(got.ConflictColumns, expected.ConflictColumns) {
		return fmt.Errorf("conflict columns mismatch: got %v, expected %v", got.ConflictColumns, expected.ConflictColumns)
	}

	// Search
	if len(got.Search) != len(expected.Search) {
		return fmt.Errorf("search fields length mismatch: got %d, expected %d", len(got.Search), len(expected.Search))
	}
	for i := range got.Search {
		if got.Search[i] != expected.Search[i] {
			return fmt.Errorf("search[%d] mismatch: got %+v, expected %+v", i, got.Search[i], expected.Search[i])
		}
	}

	return nil
}

func (m *MockRepository) validateType(got constants.RepositoryType, expected constants.RepositoryType) error {
	if got != expected {
		return fmt.Errorf("expected type %s, got %s", expected, got)
	}
	return nil
}

func (m *MockRepository) formatError(err error, funcType constants.RepositoryType) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("[%s] %w", funcType, err)
}

func (m *MockRepository) validateUpdateFields(got, expected any, params *entities.QueryParams) error {
	if got == nil || expected == nil || params == nil {
		return nil
	}

	gotVal := reflect.ValueOf(got).Elem()
	gotType := gotVal.Type()
	expectedVal := reflect.ValueOf(expected).Elem()

	var timeType = reflect.TypeOf(time.Time{})

	for _, field := range params.UpdateFields {
		if slices.Contains(params.IncrementFields, field) {
			continue
		}

		var matched bool
		for i := 0; i < gotVal.NumField(); i++ {
			structField := gotType.Field(i)
			gormTag := structField.Tag.Get("gorm")

			// Extrai "column:..." da tag GORM
			if columnName := extractGormColumn(gormTag); columnName == field {
				gotField := gotVal.Field(i)
				expectedField := expectedVal.Field(i)

				if gotField.Type() == timeType {
					matched = true
					break
				}

				if !reflect.DeepEqual(gotField.Interface(), expectedField.Interface()) {
					return fmt.Errorf("field '%s' mismatch: got '%v', expected '%v'", field, gotField.Interface(), expectedField.Interface())
				}

				matched = true
				break
			}
		}

		if !matched {
			return fmt.Errorf("field '%s' not found in entity (by gorm column tag)", field)
		}
	}

	return nil
}
