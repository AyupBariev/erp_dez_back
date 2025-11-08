-- +migrate Up

ALTER TABLE engineer_monthly_motivations
    DROP COLUMN confirmed_by_admin,
    RENAME COLUMN total_motivation TO total_motivation_amount
