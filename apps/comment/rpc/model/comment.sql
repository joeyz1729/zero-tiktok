CREATE DATABASE tiktok_comment;
USE tiktok_comment;
create table comment
(
    id          bigint      auto_increment                  ,
    video_id    bigint                              not null,
    user_id     bigint                              not null,
    comment_id  bigint                              not null,
    content     text                                not null,
#     cancel      bool        default false           not null,
    create_time timestamp   default CURRENT_TIMESTAMP null,
    update_time timestamp   default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='用户表';

