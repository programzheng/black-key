-- +goose Up
-- +goose StatementBegin
CREATE TABLE `files` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `hash_id` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `reference` varchar(255) DEFAULT NULL,
  `system` varchar(255) DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `path` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `hash_id` (`hash_id`),
  KEY `idx_files_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `files`;
-- +goose StatementEnd
