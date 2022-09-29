DROP TYPE IF EXISTS "strava_object" CASCADE;
CREATE TYPE "strava_object" AS ENUM ('activity', 'athlete');

DROP TYPE IF EXISTS "strava_aspect" CASCADE;
CREATE TYPE strava_aspect AS ENUM ('create', 'update', 'delete');

DROP TABLE IF EXISTS strava_push_event;
CREATE TABLE strava_push_event
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
    updated_at      timestamptz        DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (object_id)
);

COMMENT ON TABLE strava_push_event IS 'strava 推送记录表';

COMMENT ON COLUMN strava_push_event.id IS '主键';
COMMENT ON COLUMN strava_push_event.owner_id IS 'strava athlete id';
COMMENT ON COLUMN strava_push_event.aspect_type IS 'event_type: create, update, delete';
COMMENT ON COLUMN strava_push_event.object_type IS 'object_type: activity, athlete';
COMMENT ON COLUMN strava_push_event.object_id IS 'athlete_id or activity_id';
COMMENT ON COLUMN strava_push_event.updates IS '更新内容';
COMMENT ON COLUMN strava_push_event.status IS '是否处理了信息, 1: 所有数据全部正常写入';