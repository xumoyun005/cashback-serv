-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sources (
    id BIGSERIAL PRIMARY KEY,
    host_ip VARCHAR(50) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cashback (
    id SERIAL PRIMARY KEY,
    cashback_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    turon_user_id INTEGER NOT NULL,
    cinerama_user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS cashback_history (
    id BIGSERIAL PRIMARY KEY,
    cashback_id BIGINT NOT NULL,
    source_id BIGINT,
    cashback_amount DECIMAL(10,2) NOT NULL,
    host_ip VARCHAR(50) NOT NULL,
    type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (cashback_id) REFERENCES cashback(id),
    FOREIGN KEY (source_id) REFERENCES sources(id)
);

CREATE INDEX idx_cashback_turon_user_id ON cashback(turon_user_id);
CREATE INDEX idx_cashback_cinerama_user_id ON cashback(cinerama_user_id);
CREATE INDEX idx_cashback_history_cashback_id ON cashback_history(cashback_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cashback_history;
DROP TABLE IF EXISTS cashback;
DROP TABLE IF EXISTS sources;
-- +goose StatementEnd 