-- +migrate Down

ALTER TABLE repeat_requests
    DROP COLUMN confirmed_by,
    DROP COLUMN status,
    DROP COLUMN updated_at,
    DROP COLUMN deleted_at;
