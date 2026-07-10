package repositories

import (
	"context"
	"os"
	"social-engine/common/apiErrors"
	lgg "social-engine/common/logger"
	"social-engine/common/repositories/interfaces"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Connect(ctx context.Context) error {
	if db != nil {
		return nil
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		lgg.L(ctx).Error("Missing DATABASE_URL env var")
		return apiErrors.ErrDatabaseURL
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		lgg.L(ctx).Error("Failed to connect to db", zap.Error(err))
		return apiErrors.ErrConnectDB
	}

	return nil
}

func InitializeRepoInstance(ctx context.Context) interfaces.IRepository {
	testControl := ctx.Value(GormTestContext)
	if testControl == nil {
		return NewRepository(ctx)
	}

	mockData, ok := testControl.([]MockPayload)
	if !ok {
		panic("test control context value is not of type []MockPayload")
	}

	return NewMockRepository(mockData)
}
