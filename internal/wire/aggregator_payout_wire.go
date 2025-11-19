package wire

import (
	"erp/internal/application"
	"erp/internal/domain"
	"erp/internal/infrastructure/mysql"
	"erp/internal/interfaces/cleanhandler"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func ProvideAggPayoutRepo(db *gorm.DB) domain.AggregatorPayoutRepository {
	return mysql.NewMysqlAggregatorPayoutRepository(db)
}

func ProvideAggPayoutUC(repo domain.AggregatorPayoutRepository) *application.AggregatorPayoutUseCase {
	return application.NewAggregatorPayoutUseCase(repo)
}

var AggregatorPayoutSet = wire.NewSet(
	ProvideAggPayoutRepo,
	ProvideAggPayoutUC,
	cleanhandler.NewAggregatorPayoutHandler,
)
