# MySQL 

# 基础篇

## 通用语法及分类
- DDL（Data Definition Language）：数据定义语言，用于定义数据库对象，（数据库、表、字段）
- DML（Data Manipulation Language）：数据操作语言，用于对数据库表中的数据进行操作，（增删改查）
- DQL（Data Query Language）：数据查询语言，用于从数据库表中查询数据，（SELECT）
- DCL（Data Control Language）：数据控制语言，用于控制数据库的访问权限，（权限管理）

## DDL（数据定义语言）

### 数据库操作

查询所有数据库：
`show databases;`

查询当前数据库：
`show database();`

创建数据库：
`create database [ if not exists ] 数据库名 [ default character set 字符集 ] [ collate 排序规则 ];`

删除数据库：
`drop database [ if exists ] 数据库名;`

使用数据库：
`use 数据库名;`





