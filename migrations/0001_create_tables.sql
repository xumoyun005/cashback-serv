-- +goose Up
-- +goose StatementBegin
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
    id SERIAL PRIMARY KEY,
    cashback_id INTEGER NOT NULL REFERENCES cashback(id),
    cashback_amount DECIMAL(10,2) NOT NULL,
    host_ip VARCHAR(45),
    device VARCHAR(255),
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_cashback_turon_user_id ON cashback(turon_user_id);
CREATE INDEX idx_cashback_cinerama_user_id ON cashback(cinerama_user_id);
CREATE INDEX idx_cashback_history_cashback_id ON cashback_history(cashback_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cashback_history;
DROP TABLE IF EXISTS cashback;
-- +goose StatementEnd 