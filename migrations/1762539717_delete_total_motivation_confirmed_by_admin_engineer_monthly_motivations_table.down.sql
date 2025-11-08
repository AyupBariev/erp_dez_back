-- +migrate Down

ALTER TABLE engineer_monthly_motivations
    RENAME COLUMN total_motivation_amount TO total_motivation,
    ADD COLUMN confirmed_by_admin BOOLEAN DEFAULT FALSE