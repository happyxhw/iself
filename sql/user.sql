DROP TYPE IF EXISTS oauth2_source CASCADE;
CREATE TYPE oauth2_source AS ENUM ('github', 'google', 'strava', '');

DROP TABLE IF EXISTS "user";
CREATE table "user"
(
    id         bigserial     NOT NULL   PRIMARY KEY,
    name       varchar(64)   NOT NULL,
    email      varchar(255)  NOT NULL,
    password   varchar(64)   NOT NULL,
    avatar_url varchar(255)  NOT NULL   DEFAULT '',
    role       int2          NOT NULL   DEFAULT 0,
    source     oauth2_source NOT NULL   DEFAULT '',
    source_id  bigint        NOT NULL   DEFAULT 0,
    status     int2          NOT NULL   DEFAULT 1,

    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at bigint        NOT NULL   DEFAULT 0,

    UNIQUE (email, deleted_at)
);1991&+1shaha

CREATE INDEX user_source ON "user" (source, source_id);

COMMENT ON TABLE "user" IS '用户表';

COMMENT ON COLUMN "user".name IS '用户名';
COMMENT ON COLUMN "user".email IS '邮箱, 如果是 oauth2 用户则是: source@source_id';
COMMENT ON COLUMN "user".password IS '密码';
COMMENT ON COLUMN "user".avatar_url IS '头像链接';
COMMENT ON COLUMN "user".role IS '身份';
COMMENT ON COLUMN "user".source IS '来源，oauth2: github, google ...';
COMMENT ON COLUMN "user".source_id IS '来源，oauth2 user id';
COMMENT ON COLUMN "user".status IS '用户状态: 0 未激活, 1 已激活';
