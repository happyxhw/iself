DROP TABLE IF EXISTS goal;
CREATE TABLE goal
(
    id         bigserial   NOT NULL PRIMARY KEY,
    athlete_id bigint      NOT NULL,
    type       varchar(10) NOT NULL,
    field      varchar(10) NOT NULL,
    freq       varchar(10) NOT NULL,
    value      float       NOT NULL,
    created_at timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at bigint      NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX uni_athlete_type_freq ON goal (athlete_id, type, field, freq, deleted_at);

COMMENT ON TABLE goal IS '目标表： 每周，每月，年度';

COMMENT ON COLUMN goal.type IS '运动类型';
COMMENT ON COLUMN goal.field IS '目标类型：距离, 时间, 卡路里等';
COMMENT ON COLUMN goal.freq IS '目标类型： weekly, monthly, yearly';
COMMENT ON COLUMN goal.value IS '目标值';