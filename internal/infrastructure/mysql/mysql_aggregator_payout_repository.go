package mysql

import (
	"erp/internal/domain"
	"gorm.io/gorm"
	"time"
)

type mysqlAggregatorPayoutRepository struct{ db *gorm.DB }

func NewMysqlAggregatorPayoutRepository(db *gorm.DB) domain.AggregatorPayoutRepository {
	return &mysqlAggregatorPayoutRepository{db: db}
}

func (r *mysqlAggregatorPayoutRepository) GetDayPayouts(from, to string) ([]*domain.AggregatorDayPayout, error) {
	var rows []*domain.AggregatorDayPayout
	fromTime, _ := time.Parse("2006-01-02", from)
	toTime, _ := time.Parse("2006-01-02", to)
	toTime = toTime.AddDate(0, 0, 1)
	err := r.db.Raw(`
SELECT a.name                                     AS aggregator,
       COUNT(*)                                   AS order_count,
       SUM(CAST(o.finish_price AS DECIMAL(10,2))) AS orders_sum,
       AVG(CAST(o.finish_price AS DECIMAL(10,2))) AS avg_check,
       SUM(o.aggregator_payout)                   AS lead_cost
FROM aggregators a
JOIN orders o ON a.id = o.aggregator_id
WHERE o.status = 'closed_finally'
  AND o.scheduled_at BETWEEN ? AND ?
GROUP BY aggregator
ORDER BY lead_cost DESC, aggregator`, fromTime, toTime).Scan(&rows).Error
	return rows, err
}

// MarkPaid TODO возможно потом потребуется отмечать что оплатили агрегатору
func (r *mysqlAggregatorPayoutRepository) MarkPaid(channel string, date string, amount float64) error {
	return r.db.Exec(`
INSERT INTO aggregator_leads (aggregator_id, order_id, lead_cost, paid, paid_at)
SELECT a.id, o.id, CAST(o.finish_price AS DECIMAL(10,2)) * a.lead_percent/100, 1, NOW()
FROM orders o
JOIN aggregators a ON a.id = o.aggregator_id
WHERE DATE(o.created_at) = ? AND a.name = ?
ON DUPLICATE KEY UPDATE paid = VALUES(paid), paid_at = VALUES(paid_at)`, date, channel).Error
}
