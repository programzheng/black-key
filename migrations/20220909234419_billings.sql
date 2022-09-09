-- +goose Up
-- +goose StatementBegin
CREATE TABLE `billings` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL COMMENT '標題',
  `amount` int DEFAULT NULL COMMENT '總付款金額',
  `payer` varchar(255) DEFAULT NULL COMMENT '付款人',
  `note` varchar(255) DEFAULT NULL COMMENT '備註',
  PRIMARY KEY (`id`),
  KEY `idx_billings_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `billings`;
-- +goose StatementEnd
