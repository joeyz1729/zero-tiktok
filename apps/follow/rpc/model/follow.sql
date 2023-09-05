CREATE DATABASE tiktok_follow;
USE tiktok_follow;
create table follow
(
    id          bigint auto_increment            ,
    user_id     bigint                             not null,
    follower_id    bigint                             not null,
    cancel      bool        default false           not null,
    create_time timestamp default CURRENT_TIMESTAMP null,
    update_time timestamp default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
#     index idx_video_id (video_id),
    index idx_user_id (user_id),
    index idx_follower_id (follower_id),
#     index idx_title (title),
#     constraint uk_play_url  unique (play_url),
#     constraint fk_author_id  foreign key (author_id) references tiktok_user.user(user_id),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='关注表';

create table follow_count
(
    id              bigint auto_increment                               ,
    user_id         bigint                                      not null,
    following       int         default 0                       not null,
    whisper         int         default 0                       not null,
    black           int         default 0                       not null,
    follower        int         default 0                       not null,
    create_time     timestamp   default CURRENT_TIMESTAMP           null,
    update_time     timestamp   default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,

    index idx_user_id (user_id),

    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='关注计数表';
