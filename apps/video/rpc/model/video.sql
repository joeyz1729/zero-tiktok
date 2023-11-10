CREATE DATABASE if not exists tiktok_video;
USE tiktok_video;
create table video
(
    video_id     bigint                             not null primary key,
    author_id    bigint                             not null,
    title       varchar(255)                         not null default ' ',
    play_url    varchar(255)                         not null default ' ',
    cover_url   varchar(255)                         not null default ' ',
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    index idx_user_video_pair (author_id, video_id),    # 用于作者查找
    index idx_publish_time (create_time)   # 用于排序查找
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='视频表';

create table video_count
(
    video_id bigint not null primary key,   # 非空唯一，主键
    favorite_count bigint not null default 0 , # 计数
    comment_count bigint not null default 0 ,
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='视频计数表';