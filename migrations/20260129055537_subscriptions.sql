-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_name VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    price INTEGER CHECK (price >=0),
    date_start DATE NOT NULL,
    date_end DATE
);

CREATE INDEX idx_subscriptions_user_service ON subscriptions(user_id, service_name);
CREATE INDEX idx_subscriptions_service ON subscriptions(service_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_subscriptions_user_service;
DROP INDEX IF EXISTS idx_subscriptions_service;
DROP TABLE IF EXыISTS subscriptions;
-- +goose StatementEnd
