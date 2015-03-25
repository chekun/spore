-- +migrate Up
CREATE TABLE `users_groups` (
  `user_id` int(11) unsigned DEFAULT NULL,
  `group_id` int(11) unsigned DEFAULT NULL,
  UNIQUE KEY `user_id` (`user_id`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- +migrate Down
DROP TABLE `users_groups`;
