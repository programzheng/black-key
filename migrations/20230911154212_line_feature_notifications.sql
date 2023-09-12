-- +goose Up
-- +goose StatementBegin
CREATE TABLE `line_feature_notifications` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `feature` ENUM('new_rent_homes') NOT NULL COMMENT '功能',
  `push_cycle` ENUM('specify') NOT NULL COMMENT '發送週期',
  `push_date_time` datetime DEFAULT NULL COMMENT '發送時間',
  `limit` int DEFAULT -1 COMMENT '限制次數',
  `group_id` varchar(255) NOT NULL COMMENT '群組ID',
  `room_id` varchar(255) NOT NULL COMMENT '多人群組聊天ID',
  `user_id` varchar(255) NOT NULL COMMENT '使用者ID',
  `request` json COMMENT 'request',
  PRIMARY KEY (`id`),
  KEY `idx_line_feature_notifications_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `line_feature_notifications`;
-- +goose StatementEnd
