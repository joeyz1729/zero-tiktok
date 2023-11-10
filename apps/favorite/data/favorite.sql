CREATE DATABASE if not exists tiktok_favorite;
USE tiktok_favorite;
# 查询user，video关系
# 查询user点赞的视频列表
create table favorite
(
    id          bigint auto_increment   primary key         ,
    user_id     bigint                             not null,
    video_id    bigint                             not null,
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    unique index idx_user_video_pair (user_id, video_id) # 根据需求选择联合索引
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='点赞表';

