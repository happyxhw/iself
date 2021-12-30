DROP TABLE IF EXISTS user_oauth2;
CREATE table user_oauth2
(
    id         bigserial     NOT NULL PRIMARY KEY,
    name       varchar(64)   NOT NULL,
    source_id  bigint        NOT NULL,
    source     oauth2_source NOT NULL,
    avatar_url varchar(255)  NOT NULL   DEFAULT '',

    created_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at bigint        NOT NULL   DEFAULT 0
);

CREATE UNIQUE INDEX uni_source ON user_oauth2 (source_id, source, deleted_at);

COMMENT ON TABLE user_oauth2 IS 'oauth2 登录表';

COMMENT ON COLUMN user_oauth2.source_id IS 'oauth2用户id';
COMMENT ON COLUMN user_oauth2.source IS 'oauth2源';
COMMENT ON COLUMN user_oauth2.avatar_url IS '头像链接';
