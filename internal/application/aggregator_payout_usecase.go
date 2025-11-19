package application

import "erp/internal/domain"

type AggregatorPayoutUseCase struct {
	repo domain.AggregatorPayoutRepository
}

func NewAggregatorPayoutUseCase(r domain.AggregatorPayoutRepository) *AggregatorPayoutUseCase {
	return &AggregatorPayoutUseCase{repo: r}
}

func (uc *AggregatorPayoutUseCase) GetDayPayouts(from, to string) ([]*domain.AggregatorDayPayout, error) {
	return uc.repo.GetDayPayouts(from, to)
}

func (uc *AggregatorPayoutUseCase) MarkPaid(channel, date string, amount float64) error {
	return uc.repo.MarkPaid(channel, date, amount)
}
