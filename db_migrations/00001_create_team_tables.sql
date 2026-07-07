-- +goose Up
-- +goose StatementBegin
-- Legacy-compatible: create the teams / team_infos tables only if they do not already
-- exist (the legacy system may already own them). Columns match shared/db_models.Team /
-- TeamInfo. team_service reads/writes a subset; the extra legacy columns are created so a
-- fresh database still has the schema other services expect.
CREATE TABLE IF NOT EXISTS teams (
    id                  BIGSERIAL PRIMARY KEY,
    type                TEXT             NOT NULL DEFAULT '',
    name                TEXT             NOT NULL DEFAULT '',
    team_code           TEXT             NOT NULL DEFAULT '',
    description         TEXT             NOT NULL DEFAULT '',
    product_limit_count BIGINT           NOT NULL DEFAULT 0,
    product_count       BIGINT           NOT NULL DEFAULT 0,
    invoice_unpaid      DOUBLE PRECISION NOT NULL DEFAULT 0,
    invoice_not_final   DOUBLE PRECISION NOT NULL DEFAULT 0,
    deleted             BOOLEAN          NOT NULL DEFAULT false
);
CREATE UNIQUE INDEX IF NOT EXISTS team_code_unique ON teams (team_code);

CREATE TABLE IF NOT EXISTS team_infos (
    id                  BIGSERIAL PRIMARY KEY,
    team_id             BIGINT NOT NULL DEFAULT 0,
    return_warehouse_id BIGINT,
    return_user_id      BIGINT,
    contact_number      TEXT   NOT NULL DEFAULT '',
    bank_type           TEXT   NOT NULL DEFAULT '',
    bank_owner_name     TEXT   NOT NULL DEFAULT '',
    bank_account_number TEXT   NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_team_infos_team_id ON team_infos (team_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- No-op: never drop the (potentially legacy-owned) teams / team_infos tables on rollback.
SELECT 1;
-- +goose StatementEnd
