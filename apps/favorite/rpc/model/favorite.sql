CREATE DATABASE if not exists tiktok_favorite;
USE tiktok_favorite;
create table favorite
(
    id          bigint auto_increment            ,
    user_id     bigint                             not null,
    video_id    bigint                             not null,
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    unique index idx_user_video_id (user_id, video_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='点赞表';

