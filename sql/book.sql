CREATE TABLE `book` (
  `id`            INT(11)      NOT NULL      AUTO_INCREMENT PRIMARY KEY,
  `name`          VARCHAR(255) NOT NULL      DEFAULT '',
  `writer`        TEXT,

  `operator_uid`  VARCHAR(50)  NOT NULL      DEFAULT '',
  `operator_name` VARCHAR(50)  NOT NULL      DEFAULT '',
  `created_at`    TIMESTAMP    NOT NULL      DEFAULT CURRENT_TIMESTAMP,
  `updated_at`    TIMESTAMP    NOT NULL      DEFAULT CURRENT_TIMESTAMP
  ON UPDATE CURRENT_TIMESTAMP
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;