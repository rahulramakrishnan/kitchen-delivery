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

CREATE TABLE `shelf_orders` (
  `uuid`                            char(36)           NOT NULL,
  `order_uuid`                      char(36)           NOT NULL,
  `shelf_type`                      varchar(191)       NOT NULL,
  `order_status`                    varchar(191)       NOT NULL,
  `version`                         INTEGER            NOT NULL,
  `expires_at`                      DATETIME           NOT NULL,
  `created_at`                      DATETIME           NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`                      DATETIME           DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uuid`),
  FOREIGN KEY (`order_uuid`) REFERENCES orders(`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `shelf_orders` ADD INDEX (`order_uuid`);
ALTER TABLE `shelf_orders` ADD INDEX (`shelf_type`, `order_status`);
ALTER TABLE `shelf_orders` ADD INDEX (`expires_at`);