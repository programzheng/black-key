-- +goose Up
-- +goose StatementBegin
CREATE TABLE `line_bot_events` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `type` enum('text', 'postback', 'image', 'video', 'audio', 'file', 'location', 'sticker') NOT NULL,
  `key` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_line_bot_events_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `line_bot_events`;
-- +goose StatementEnd
