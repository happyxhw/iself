DROP TABLE IF EXISTS strava_goal;
CREATE TABLE strava_goal
(
    id         bigserial   NOT NULL PRIMARY KEY,
    athlete_id bigint      NOT NULL,
    "type"     varchar(10) NOT NULL,
    field      varchar(10) NOT NULL,
    freq       varchar(10) NOT NULL,
    "value"    float       NOT NULL,
    created_at timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at bigint      NOT NULL DEFAULT 0,

    UNIQUE (athlete_id, "type", field, freq, deleted_at)
);

COMMENT ON TABLE strava_goal IS '目标表： 每周，每月，年度';

COMMENT ON COLUMN strava_goal.type IS '运动类型';
COMMENT ON COLUMN strava_goal.field IS '目标类型：距离, 时间, 卡路里等';
COMMENT ON COLUMN strava_goal.freq IS '目标类型： weekly, monthly, yearly';
COMMENT ON COLUMN strava_goal.value IS '目标值';