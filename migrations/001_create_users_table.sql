-- +migrate Up
CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `screen_name` varchar(100) DEFAULT NULL,
  `avatar` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `screen_name` (`screen_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- +migrate Down
Drop TABLE `users`;
