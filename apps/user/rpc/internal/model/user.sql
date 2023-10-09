CREATE DATABASE if not exists tiktok_user;
USE tiktok_user;
create table user
(
    id          bigint auto_increment
        primary key,
    user_id     bigint                              not null,
    username    varchar(64)                         not null,
    password    varchar(64)                         not null,
    email       varchar(64)                         null,
    gender      tinyint   default 0                 not null,
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint idx_user_id
        unique (user_id),
    constraint idx_username
        unique (username),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='用户表';


create table user_count
(
    id bigint auto_increment primary key,
    user_id bigint not null,
    followed_count bigint default 0  not null,
    follower_count bigint default 0 not null,
    total_favorited bigint default 0 not null,
    work_count bigint default 0 not null,
    favorite_count bigint default 0 not null,
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint idx_user_id unique (user_id)
) ENGINE = InnoDB Default CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='用户计数表';