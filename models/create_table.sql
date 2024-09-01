CREATE TABLE `user`
(
    `id`          bigint(20)                             NOT NULL AUTO_INCREMENT,
    `user_id`     bigint(20)                             NOT NULL,
    `name`    varchar(64) COLLATE utf8mb4_general_ci     NOT NULL,
    `password`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
    `email`       varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
    `gender`      tinyint(4)                             NOT NULL DEFAULT '0',
    `telephone`   varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
    `boss`        tinyint(4)                             NOT NULL DEFAULT '0',
    `company_id`     bigint(20)                          ,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_email` (`email`) USING BTREE,
    UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT COLLATE=utf8mb4_general_ci;

ALTER TABLE user
    ADD COLUMN department VARCHAR(255),
    ADD COLUMN `role` VARCHAR(20);

ALTER TABLE user MODIFY COLUMN company_id varchar(20) NULL;

CREATE TABLE `company`
(
    `id`           bigint(20)                             NOT NULL AUTO_INCREMENT,
    `owner_id`     bigint(20)                             NOT NULL,
    `name`         varchar(64) COLLATE utf8mb4_general_ci     NOT NULL,
    `company_id`   bigint(20)                             NOT NULL ,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_owner_id` (`owner_id`) USING BTREE
) ENGINE=InnoDB DEFAULT COLLATE=utf8mb4_general_ci;

CREATE TABLE `logger`
(
    `id`            bigint(20)                             NOT NULL AUTO_INCREMENT,
    `user_id`       bigint(20)                             NOT NULL,
    `user_name`     varchar(64) COLLATE utf8mb4_general_ci     NOT NULL,
    `company_id`    bigint(20)                          NOT NULL,
    `department`    varchar(64) COLLATE utf8mb4_general_ci     NOT NULL,
    `role`          varchar(64) COLLATE utf8mb4_general_ci     NOT NULL,
    `log_time`      timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT COLLATE=utf8mb4_general_ci;

CREATE TABLE `post`
(
    `id`            bigint(20)                              NOT NULL AUTO_INCREMENT,
    `company_id`    bigint(20)                              NOT NULL,
    `post_name`     varchar(64) COLLATE utf8mb4_general_ci  NOT NULL,
    `post_by`       bigint(20)                              NOT NULL,
    `description`   varchar(64) COLLATE utf8mb4_general_ci  NOT NULL,
    `field_id`      varchar(64) COLLATE utf8mb4_general_ci  NULL,
    `create_time`    timestamp                               NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`   timestamp                               NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE
        CURRENT_TIMESTAMP COMMENT '更新时间',
    `end_time`      timestamp                               NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结束时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_field_id` (`field_id`) USING BTREE
) ENGINE=InnoDB DEFAULT COLLATE=utf8mb4_general_ci;

CREATE TABLE `post_task`
(
    `post_id`       bigint(20)                              NOT NULL,
    `post_to`       bigint(20)                              NOT NULL,
    `status`        int(4)                                  NOT NULL DEFAULT '0',
    `complete_time` timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
    PRIMARY KEY (post_id,post_to),
    FOREIGN KEY (post_id) REFERENCES post(id),
    FOREIGN KEY (post_to) REFERENCES user(user_id)
) ENGINE=InnoDB DEFAULT COLLATE=utf8mb4_general_ci;
