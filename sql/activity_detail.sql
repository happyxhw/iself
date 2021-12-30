DROP TABLE IF EXISTS activity_detail;
CREATE TABLE activity_detail
(
    id                   bigint                   NOT NULL PRIMARY KEY,
    athlete_id           bigint                   NOT NULL,
    name                 varchar(128)             NOT NULL,
    type                 varchar(32)              NOT NULL,
    distance             float                    NOT NULL DEFAULT 0.0,
    moving_time          integer                  NOT NULL DEFAULT 0,
    elapsed_time         integer                  NOT NULL DEFAULT 0,
    total_elevation_gain float                    NOT NULL DEFAULT 0.0,
    start_date_local     timestamp                NOT NULL,
    polyline             text                     NOT NULL DEFAULT '',
    summary_polyline     text                     NOT NULL DEFAULT '',
    average_speed        float                    NOT NULL DEFAULT 0.0,
    max_speed            float                    NOT NULL DEFAULT 0.0,
    average_heartrate    float                    NOT NULL DEFAULT 0.0,
    max_heartrate        float                    NOT NULL DEFAULT 0.0,
    elev_high            float                    NOT NULL DEFAULT 0.0,
    elev_low             float                    NOT NULL DEFAULT 0.0,
    calories             float                    NOT NULL DEFAULT 0.0,
    splits_metric        jsonb,
    best_efforts         jsonb,
    device_name          varchar(50),

    created_at           timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at           bigint                   NOT NULL DEFAULT 0
);

-- where id = athlete_id order by start_date_local
CREATE INDEX athlete_id_index ON activity_detail (athlete_id, type, start_date_local);


COMMENT ON TABLE activity_detail IS '活动详情表';

COMMENT ON COLUMN activity_detail.id IS '活动id, 是strava返回的活动id';
COMMENT ON COLUMN activity_detail.athlete_id IS 'strava用户id';
COMMENT ON COLUMN activity_detail.name IS '活动名称';
COMMENT ON COLUMN activity_detail.type IS '活动类型';
COMMENT ON COLUMN activity_detail.distance IS '距离';
COMMENT ON COLUMN activity_detail.moving_time IS '实际运动时间，单位秒';
COMMENT ON COLUMN activity_detail.elapsed_time IS '经过时间，单位秒';
COMMENT ON COLUMN activity_detail.total_elevation_gain IS '总高度增益';
COMMENT ON COLUMN activity_detail.start_date_local IS '开始时间';
COMMENT ON COLUMN activity_detail.polyline IS '编码后的地图';
COMMENT ON COLUMN activity_detail.summary_polyline IS '编码后的地图，压缩版';
COMMENT ON COLUMN activity_detail.average_speed IS '均速';
COMMENT ON COLUMN activity_detail.max_speed IS '最大速';
COMMENT ON COLUMN activity_detail.average_heartrate IS '平均心率';
COMMENT ON COLUMN activity_detail.max_heartrate IS '最大心率';
COMMENT ON COLUMN activity_detail.elev_high IS '最大高度';
COMMENT ON COLUMN activity_detail.elev_low IS '最低高度';
COMMENT ON COLUMN activity_detail.calories IS '卡路里';
COMMENT ON COLUMN activity_detail.splits_metric IS '每公里数据，json 列表';
COMMENT ON COLUMN activity_detail.best_efforts IS '最佳，json 列表';
COMMENT ON COLUMN activity_detail.device_name IS '设备名称';
