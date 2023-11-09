CREATE DATABASE if not exists tiktok_user;
USE tiktok_user;
create table if not exists user
(
#     id          bigint auto_increment
#         primary key,
    user_id     bigint                    not null primary key,
    username    varchar(64)                         not null,
    password    varchar(64)                         not null,
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    constraint idx_username_userid
        unique (username, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='用户表';


create table if not exists user_count
(
#     id bigint auto_increment primary key,
    user_id bigint not null primary key,
    followed_count bigint default 0  not null,
    follower_count bigint default 0 not null,
    total_favorited bigint default 0 not null,
    work_count bigint default 0 not null,
    favorite_count bigint default 0 not null,
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP
) ENGINE = InnoDB Default CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='用户计数表';

insert into tiktok_user.user(user_id, username, password)
values
                                                (1, "auth1", "123"),
                                                (2, "auth2", "123");
insert into tiktok_user.user(user_id, username, password)
values
    (3, "auth3", "123"),
    (4, "auth4", "123"),
    (5, "auth5", "123"),
    (6, "auth6", "123"),
    (7, "auth7", "123"),
    (8, "auth8", "123"),
    (9, "auth9", "123"),
    (10, "auth10", "123");


insert into tiktok_user.user_count(user_id)
values
    (1),
    (2),
    (3),
    (4),
    (5),
    (6),
    (7),
    (8),
    (9),
    (10);
