package interfaces

import (
	"context"
	"social-engine/common/repositories/entities"
)

type IRepository interface {
	WithTableName(tableName string) IRepository
	Ping(ctx context.Context) error
	Find(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (*entities.PaginatedResult, error)
	FindByID(ctx context.Context, entityInstance entities.IEntity, id any, params *entities.QueryParams) error
	FindOne(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error
	FindAll(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (any, error)
	Save(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error
	Update(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error
	Delete(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error
	WithTransaction(ctx context.Context, fn func(tx IRepository) error) error
	BeginTx(ctx context.Context) (IRepository, error)
	Rollback(err error) error
	Commit() error
	Raw(ctx context.Context, entityInstance any, query string, values ...any) error
	Verify(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (bool, error)
	Count(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) (int64, error)
	CountDistinct(ctx context.Context, entityInstance entities.IEntity, field string, params *entities.QueryParams) (int64, error)
	SaveOrUpdate(ctx context.Context, entityInstance entities.IEntity, params *entities.QueryParams) error
	BulkSave(ctx context.Context, entities []entities.IEntity) error
	BulkUpdate(ctx context.Context, entity entities.IEntity, params *entities.QueryParams) error
	BulkDelete(ctx context.Context, entity entities.IEntity, ids []any, params *entities.QueryParams) error
}
