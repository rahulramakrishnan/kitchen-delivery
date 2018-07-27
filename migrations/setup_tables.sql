CREATE TABLE `orders` (
  `uuid`                            char(36)           NOT NULL,
  `name`                            varchar(255)       NOT NULL,
  `temp`                            varchar(191)       NOT NULL,
  `shelf_life`                      INTEGER            NOT NULL,
  `decay_rate`                      FLOAT              NOT NULL,
  `created_at`                      DATETIME           NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `orders` ADD INDEX (`temp`);
ALTER TABLE `orders` ADD INDEX (`shelf_life`, `decay_rate`);
ALTER TABLE `orders` ADD INDEX (`created_at`);


CREATE TABLE `order_logs` (
  `uuid`                            char(36)           NOT NULL,
  `order_uuid`                      char(36)           NOT NULL,
  `order_status`                    varchar(191)       NOT NULL,
  `description`                     TEXT               NOT NULL,
  `created_at`                      DATETIME           NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`uuid`),
  FOREIGN KEY (`order_uuid`) REFERENCES orders(`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `order_logs` ADD INDEX (`order_uuid`);
ALTER TABLE `order_logs` ADD INDEX (`order_status`);
ALTER TABLE `order_logs` ADD INDEX (`created_at`);