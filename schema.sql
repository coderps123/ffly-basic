-- 创建数据库（如果不存在）
create database if not exists `ffly_basic` 
default char set utf8mb4 -- 设置数据库字符集 utf8mb4
collate utf8mb4_unicode_ci;   -- 设置数据库排序规则 utf8mb4_unicode_ci  

-- 修改现有数据库的字符集和排序规则（针对已经存在的数据库）
alter database `ffly_basic` 
default char set utf8mb4  -- 设置数据库字符集 utf8mb4
collate utf8mb4_unicode_ci;    -- 设置数据库排序规则 utf8mb4_unicode_ci

-- 使用新创建或已存在的数据库
use `ffly_basic`;

-- 创建用户表
create table if not exists `users` (
  `id` bigint unsigned not null auto_increment comment '用户id',
  `username` varchar(50) not null comment '用户名',
  `password` varchar(255) not null comment '密码',
  `nickname` varchar(50) default null comment '昵称',
  `email` varchar(100) default null comment '邮箱',
  `phone` varchar(20) default null comment '手机号',
  `status` tinyint unsigned not null default '1' comment '状态 1: 启用 2: 禁用',
  `role_id` bigint unsigned not null comment '角色id',
  `created_at` timestamp not null default current_timestamp comment '创建时间',
  `updated_at` timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
  `deleted_at` timestamp null default null comment '删除时间',
  primary key (`id`), -- 主键
  unique key `uk_username` (`username`), -- 唯一索引 username
  unique key `uk_email` (`email`), -- 唯一索引 email
  unique key `uk_phone` (`phone`), -- 唯一索引 phone
  unique key `uk_username_email_phone` (`username`, `email`, `phone`), -- 联合唯一索引 username, email, phone
  key `idx_role_id` (`role_id`), -- 索引 role_id
  key `idx_deleted_at` (`deleted_at`) -- 索引 deleted_at
) engine=innodb auto_increment=1 comment='用户表';

-- 创建角色表
create table if not exists `roles` (
  `id` bigint unsigned not null auto_increment comment '角色id',
  `name` varchar(50) not null comment '角色名称',
  `code` varchar(50) not null comment '角色代码',
  `status` tinyint unsigned not null default '1' comment '状态 1: 启用 2: 禁用',
  `remark` varchar(255) default null comment '备注',
  `created_at` timestamp not null default current_timestamp comment '创建时间',
  `updated_at` timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
  `deleted_at` timestamp null default null comment '删除时间',
  primary key (`id`), -- 主键
  unique key `uk_name` (`name`), -- 唯一索引 name
  unique key `uk_code` (`code`), -- 唯一索引 code
  key `idx_deleted_at` (`deleted_at`) -- 索引 deleted_at
) engine=innodb auto_increment=1 comment='角色表';

-- 创建用户角色关联表
create table if not exists `user_roles` (
  `id` bigint unsigned not null auto_increment comment 'ID',
  `user_id` bigint unsigned not null comment '用户id',
  `role_id` bigint unsigned not null comment '角色id',
  `created_at` timestamp not null default current_timestamp comment '创建时间',
  `updated_at` timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
  `deleted_at` timestamp null default null comment '删除时间',
  primary key (`id`), -- 主键
  unique key `uk_user_role` (`user_id`, `role_id`), -- 联合唯一索引 user_id, role_id
  key `idx_user_id` (`user_id`), -- 索引 user_id
  key `idx_role_id` (`role_id`), -- 索引 role_id
  key `idx_deleted_at` (`deleted_at`), -- 索引 deleted_at
  constraint `fk_user_roles_user_id` foreign key (`user_id`) -- 外键 user_id
  references `users` (`id`) on delete cascade on update cascade, -- 引用 users.id 并设置级联删除和更新
  constraint `fk_user_roles_role_id` foreign key (`role_id`) -- 外键 role_id
  references `roles` (`id`) on delete cascade on update cascade -- 引用 roles.id 并设置级联删除和更新
) engine=innodb auto_increment=1 comment='用户角色关联表';

-- 创建权限表
create table if not exists `permissions` (
  `id` bigint unsigned not null auto_increment comment '权限id',
  `name` varchar(50) not null comment '权限名称',
  `type` enum('menu', 'button') not null comment '权限类型, menu: 菜单, button: 按钮',
  `path` varchar(255) default null comment '菜单路径',
  `code` varchar(50) default null comment '权限代码', -- 按钮权限的标识符
  `component` varchar(255) default null comment '组件路径', -- 菜单权限的组件
  `icon` varchar(255) default null comment '菜单图标',
  `sort` int not null default '0' comment '排序',
  `parent_id` bigint unsigned not null default '0' comment '父权限id',
  `status` tinyint unsigned not null default '1' comment '状态 1: 启用 2: 禁用',
  `remark` varchar(255) default null comment '备注',
  `created_at` timestamp not null default current_timestamp comment '创建时间',
  `updated_at` timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
  `deleted_at` timestamp null default null comment '删除时间',
  primary key (`id`), -- 主键
  unique key `uk_path` (`path`), -- 唯一索引 path
  unique key `uk_code` (`code`), -- 唯一索引 code
  key `idx_parent_id` (`parent_id`), -- 索引 parent_id
  key `idx_deleted_at` (`deleted_at`), -- 索引 deleted_at
  constraint `chk_type_path_code` check ( -- 约束 type, path, code 
    (type = 'menu' and path is not null and code is null) or 
    (type = 'button' and path is null and code is not null)
  )
) engine=innodb auto_increment=1 comment='权限表';

-- 创建角色权限关联表
create table if not exists `role_permissions` (
  `id` bigint unsigned not null auto_increment comment 'ID',
  `role_id`bigint unsigned not null comment '角色id',
  `permission_id` bigint unsigned not null comment '权限id',
  `created_at` timestamp not null default current_timestamp comment '创建时间',
  `updated_at` timestamp not null default current_timestamp on update current_timestamp comment '更新时间',
  `deleted_at` timestamp null default null comment '删除时间',
  primary key (`id`), -- 主键
  unique key `uk_role_permission` (`role_id`, `permission_id`), -- 联合唯一索引 role_id, permission_id
  key `idx_role_id` (`role_id`), -- 索引 role_id
  key `idx_permission_id` (`permission_id`), -- 索引 permission_id
  key `idx_deleted_at` (`deleted_at`), -- 索引 deleted_at
  constraint `fk_role_permissions_role_id` foreign key (`role_id`) -- 外键 role_id
  references `roles` (`id`) on delete cascade on update cascade, -- 引用 roles.id 并设置级联删除和更新
  constraint `fk_role_permissions_permission_id` foreign key (`permission_id`) -- 外键 permission_id
  references `permissions` (`id`) on delete cascade on update cascade -- 引用 permissions.id 并设置级联删除和更新
) engine=innodb auto_increment=1 comment='角色权限关联表';



-- 修改表
alter table `users` 
add unique key `uk_username` (`username`),
add unique key `uk_email` (`email`),
add unique key `uk_phone` (`phone`);

alter table `permissions` 
modify column `type` enum('menu', 'button') not null comment '权限类型, menu: 菜单, button: 按钮';