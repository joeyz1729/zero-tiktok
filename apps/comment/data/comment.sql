CREATE DATABASE if not exists tiktok_comment;
USE tiktok_comment;
create table comment
(
    comment_id  bigint              not null primary key,
    video_id    bigint                              not null,
    user_id     bigint                              not null,
    content     text                                not null,
    create_time timestamp   default CURRENT_TIMESTAMP null,
    update_time timestamp   default CURRENT_TIMESTAMP null on update CURRENT_TIMESTAMP,
    index unified_index (video_id, comment_id, user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 collate = utf8mb4_general_ci COMMENT='用户表';

# 写入：根据userId，videoId，生成的commentId，插入数据
# commentId是一定唯一的，user_id和video_id可以重复
# 删除： 根据userId，videoId，指定的commentId，删除数据
# 添加三字段的联合索引
# 所以使用三个联合索引，方便查询，如果只用comment查还需要验证user和video
# 查询：根据videoId查询所有comment
# 把video_id放前面，方便查询所有comment列表