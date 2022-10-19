-- https://developers.strava.com/docs/reference/#api-models-StreamSet
DROP TABLE IF EXISTS strava_activity_stream;
CREATE TABLE strava_activity_stream
(
    id              bigint NOT NULL PRIMARY KEY,
    "time"            jsonb,
    distance        jsonb,
    latlng          jsonb,
    altitude        jsonb,
    velocity_smooth jsonb,
    heartrate       jsonb,
    cadence         jsonb,
    watts           jsonb,
    temp            jsonb,
    moving          jsonb,
    grade_smooth    jsonb,
    created_at      timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      bigint NOT NULL          DEFAULT 0
);

COMMENT ON TABLE strava_activity_stream IS '曲线信息';

COMMENT ON COLUMN strava_activity_stream.time IS '时间曲线';
COMMENT ON COLUMN strava_activity_stream.distance IS '距离曲线';
COMMENT ON COLUMN strava_activity_stream.latlng IS '坐标曲线';
COMMENT ON COLUMN strava_activity_stream.altitude IS '高度曲线';
COMMENT ON COLUMN strava_activity_stream.velocity_smooth IS '平滑速度曲线';
COMMENT ON COLUMN strava_activity_stream.heartrate IS '心率曲线';
COMMENT ON COLUMN strava_activity_stream.cadence IS '步频、踏频曲线';
COMMENT ON COLUMN strava_activity_stream.watts IS '功率曲线';
COMMENT ON COLUMN strava_activity_stream.temp IS '温度曲线';
COMMENT ON COLUMN strava_activity_stream.moving IS '是否移动，true or false';
COMMENT ON COLUMN strava_activity_stream.grade_smooth IS '等级曲线';
