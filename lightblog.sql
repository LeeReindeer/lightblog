# lightblog.sql
# lightblog 数据库表设计
# author: LeeReindeer
# IMPORT: mysql -uroot -p < lightblog.sql
CREATE DATABASE IF NOT EXISTS lightblog;
USE lightblog;

# 用户表
DROP TABLE IF EXISTS user;
CREATE TABLE user (
  user_id INT AUTO_INCREMENT PRIMARY KEY ,
  user_name VARCHAR(20) UNIQUE, #用户名
  user_avatar VARCHAR(100), #用户头像链接
  user_password VARCHAR(40), #密码，服务器之存储密码的 hash
  user_register DATETIME, #用户注册时间
  user_bio VARCHAR(140), #用户简介
  user_followers INT, # 粉丝数量
  user_following INT #关注数量
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

# 轻博客表
DROP TABLE IF EXISTS blog;
CREATE TABLE blog (
  blog_id INT AUTO_INCREMENT PRIMARY KEY,
  blog_uid INT, # 作者 id
  blog_content VARCHAR(141), # lightblog 字数限制 141
  blog_time DATETIME, #微博发布时间

  blog_like INT, #喜欢数量
  blog_unlike INT, #unlike 数量
  blog_comment INT, #评论数量

  FOREIGN KEY (blog_uid)
      REFERENCES user(user_id)
      ON DELETE NO ACTION
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

# 点赞/反对/评论用户表
DROP TABLE IF EXISTS blog_counter;
CREATE TABLE blog_counter (
  blog_id INT,
  user_id INT,
  count_type TINYINT, #like 0 , unlike 1, comment 2

  primary key (blog_id, user_id),

  FOREIGN KEY (blog_id)
    REFERENCES blog(blog_id)
    ON DELETE CASCADE,
  FOREIGN KEY (user_id)
    REFERENCES user(user_id)
    ON DELETE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

# 评论表
DROP TABLE IF EXISTS comment;
CREATE TABLE comment (
  comm_id INT AUTO_INCREMENT PRIMARY KEY,
  comm_uid INT, #发表评论的用户 id
  comm_content VARCHAR(141), #评论内容
  comm_time DATETIME, #评论时间
  # fixme
  comm_host_id INT, #评论 blog 的评论，即为 blog_id；评论评论的评论，即为父评论 id
  comm_attach_blog BOOL, #评论类型： TRUE-> 评论 blog ，FALSE->多重评论
  comm_like INT,
  #comm_unlike INT,

  FOREIGN KEY (comm_uid)
    REFERENCES user(user_id)
    ON DELETE NO ACTION,
  FOREIGN KEY (comm_host_id)
    REFERENCES comment(comm_id)
    ON DELETE NO ACTION
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

# 关注关系表
DROP TABLE IF EXISTS fan_follow;
CREATE TABLE fan_follow (
  user_from INT,
  user_to INT,

  FOREIGN KEY (user_from)
    REFERENCES user(user_id)
    ON DELETE CASCADE, # delete followers and followings when you delete account
  FOREIGN KEY (user_to)
    REFERENCES user(user_id)
    ON DELETE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

# 标签表
DROP TABLE IF EXISTS tag;
CREATE TABLE tag (
  tag_id INT AUTO_INCREMENT PRIMARY KEY,
  tag_name VARCHAR(10), #标签名称，限制10个字符
  tag_time DATETIME #标签创建的时间
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

# 标签博客关系表
DROP TABLE IF EXISTS taged_blog;
CREATE TABLE taged_blog (
  blog_id INT,
  tag_id INT,

  FOREIGN KEY (blog_id)
    REFERENCES blog(blog_id)
    ON DELETE NO ACTION,
  FOREIGN KEY (tag_id)
    REFERENCES tag(tag_id)
    ON DELETE NO ACTION
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
