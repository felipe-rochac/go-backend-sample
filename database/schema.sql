CREATE SCHEMA internal;

USE internal;

CREATE DATABASE IF NOT EXISTS users
    CHARACTER SET = utf8
    COLLATE = utf8mb4
    ENCRYPTION = 'Y';

USE `users`;

CREATE TABLE IF NOT EXISTS `user`
(
    user_id BINARY(16) default (UUID_TO_BIN(UUID())) PRIMARY KEY,
    `name` VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(100) NOT NULL,
);