CREATE TABLE db_group (
    id BIGSERIAL PRIMARY KEY,
    row_status row_status NOT NULL DEFAULT 'NORMAL',
    creator_id INTEGER NOT NULL REFERENCES principal (id),
    created_ts BIGINT NOT NULL DEFAULT extract(epoch from now()),
    updater_id INTEGER NOT NULL REFERENCES principal (id),
    updated_ts BIGINT NOT NULL DEFAULT extract(epoch from now()),
    project_resource_id TEXT NOT NULL REFERENCES project (resource_id),
    resource_id TEXT NOT NULL,
    placeholder TEXT NOT NULL DEFAULT '',
    expression JSONB NOT NULL DEFAULT '{}'
);

CREATE UNIQUE INDEX idx_db_group_unique_resource_id ON db_group(resource_id);

CREATE INDEX idx_db_group_project_resource_id ON db_group(project_resource_id);

ALTER SEQUENCE db_group_id_seq RESTART WITH 101;

CREATE TRIGGER update_db_group_updated_ts
BEFORE
UPDATE
    ON db_group FOR EACH ROW
EXECUTE FUNCTION trigger_update_updated_ts();

CREATE TABLE schema_group (
    id BIGSERIAL PRIMARY KEY,
    row_status row_status NOT NULL DEFAULT 'NORMAL',
    creator_id INTEGER NOT NULL REFERENCES principal (id),
    created_ts BIGINT NOT NULL DEFAULT extract(epoch from now()),
    updater_id INTEGER NOT NULL REFERENCES principal (id),
    updated_ts BIGINT NOT NULL DEFAULT extract(epoch from now()),
    db_group_resource_id TEXT NOT NULL REFERENCES db_group (resource_id),
    resource_id TEXT NOT NULL,
    placeholder TEXT NOT NULL DEFAULT '',
    expression JSONB NOT NULL DEFAULT '{}'
);

CREATE UNIQUE INDEX idx_schema_group_unique_resource_id ON schema_group(resource_id);

CREATE INDEX idx_schema_group_db_group_resource_id ON schema_group(db_group_resource_id);

ALTER SEQUENCE schema_group_id_seq RESTART WITH 101;

CREATE TRIGGER update_schema_group_updated_ts
BEFORE
UPDATE
    ON schema_group FOR EACH ROW
EXECUTE FUNCTION trigger_update_updated_ts();