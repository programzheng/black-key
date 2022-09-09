-- +goose Up
-- +goose StatementBegin
CREATE TABLE `line_billings` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `billing_id` int unsigned NOT NULL,
  `group_id` varchar(255) DEFAULT NULL,
  `room_id` varchar(255) DEFAULT NULL,
  `user_id` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `billing_id` (`billing_id`),
  KEY `idx_line_billings_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `line_billings`;
-- +goose StatementEnd
