CREATE DATABASE  if not exists tiktok_follow;
USE tiktok_follow;

create table follow
(
    id          bigint auto_increment      primary key      ,
    user_id     bigint                             not null,
    follow_id    bigint                             not null,
    cancel      bool        default false           not null,
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    index idx_follow_id (user_id, follow_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='关注表';


# create table follow_count
# (
#     id              bigint auto_increment         primary key                      ,
#     user_id         bigint                                      not null,
#     followed        int         default 0                       not null,
#     follower        int         default 0                       not null,
#     create_time     timestamp   default CURRENT_TIMESTAMP           null,
#     update_time     timestamp   default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
#     index idx_user_id (user_id)
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='关注计数表';
