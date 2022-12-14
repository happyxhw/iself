DROP TABLE IF EXISTS strava_activity_raw;
CREATE TABLE strava_activity_raw
(
    id         bigint                   NOT NULL PRIMARY KEY,
    data       jsonb,

    created_at timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at bigint                   NOT NULL DEFAULT 0
);

COMMENT ON TABLE strava_activity_raw IS '活动原始数据';

COMMENT ON COLUMN strava_activity_raw.data IS '原始数据';