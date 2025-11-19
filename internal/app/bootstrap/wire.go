//go:build wireinject
// +build wireinject

package bootstrap

import (
	"erp/internal/interfaces/cleanhandler"
	wireSet "erp/internal/wire" // ← alias, чтобы не было конфликта с wire.Build
	"github.com/google/wire"
	"gorm.io/gorm"
)

// BuildAggregatorPayoutHandler создаёт хендлер через wire.
// Реальная реализация появится в wire_gen.go после генерации.
func BuildAggregatorPayoutHandler(db *gorm.DB) (*cleanhandler.AggregatorPayoutHandler, func(), error) {
	panic(wire.Build(wireSet.AggregatorPayoutSet))
}
