-- +goose Up
-- +goose StatementBegin
CREATE TABLE `line_bot_requests` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `type` varchar(255) NOT NULL,
  `group_id` varchar(255) DEFAULT NULL,
  `room_id` varchar(255) DEFAULT NULL,
  `user_id` varchar(255) NOT NULL,
  `reply_token` varchar(255) NOT NULL,
  `request` text NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_line_bot_requests_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `line_bot_requests`;
-- +goose StatementEnd
