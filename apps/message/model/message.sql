CREATE DATABASE tiktok_message;
USE tiktok_message;
create table message
(
    id          bigint      auto_increment                  ,
    user_id     bigint                              not null,
    to_user_id  bigint                              not null,
    content     text                                not null,
    create_time timestamp   default CURRENT_TIMESTAMP null,
    update_time timestamp   default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    index uid (user_id),
    index to_uid (to_user_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='用户表';

