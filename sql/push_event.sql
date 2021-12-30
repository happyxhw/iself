DROP TYPE IF EXISTS "object" CASCADE;
CREATE TYPE "object" AS ENUM ('activity', 'athlete');

DROP TYPE IF EXISTS "aspect" CASCADE;
CREATE TYPE aspect AS ENUM ('create', 'update', 'delete');

DROP TABLE IF EXISTS push_event;
CREATE TABLE push_event
(
    id              bigserial NOT NULL PRIMARY KEY,
    owner_id        bigint    NOT NULL,
    aspect_type     aspect    NOT NULL,
    object_type     object    NOT NULL,
    object_id       bigint    NOT NULL,
    event_time      bigint    NOT NULL,
    updates         jsonb,
    status          integer   NOT NULL DEFAULT 0,

    created_at      timestamptz        DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamptz        DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uni_object_id ON push_event (object_id);

COMMENT ON TABLE push_event IS 'strava 推送记录表';

COMMENT ON COLUMN push_event.id IS '主键';
COMMENT ON COLUMN push_event.owner_id IS 'strava athlete id';
COMMENT ON COLUMN push_event.aspect_type IS 'event_type: create, update, delete';
COMMENT ON COLUMN push_event.object_type IS 'object_type: activity, athlete';
COMMENT ON COLUMN push_event.object_id IS 'athlete_id or activity_id';
COMMENT ON COLUMN push_event.updates IS '更新内容';
COMMENT ON COLUMN push_event.status IS '是否处理了信息, 1: 所有数据全部正常写入';