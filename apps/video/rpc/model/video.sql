CREATE DATABASE tiktok_video;
USE tiktok_video;
create table video
(
    id          bigint auto_increment            ,
    video_id     bigint                             not null,
    author_id    bigint                             not null,
    title       varchar(64)                         not null,
    play_url    varchar(64)                         not null,
    cover_url   varchar(64)                         not null,
    data        blob                                null,
    publish_time timestamp default CURRENT_TIMESTAMP null,
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint uk_video_id  unique (video_id),
#     index idx_video_id (video_id),
#     index idx_author_id (author_id),
#     index idx_title (title),
    constraint uk_play_url  unique (play_url),
#     constraint fk_author_id  foreign key (author_id) references tiktok_user.user(user_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='视频表';

