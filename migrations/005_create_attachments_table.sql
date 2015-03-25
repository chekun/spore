-- +migrate Up
CREATE TABLE `attachments` (
  `owner_id` int(11) NOT NULL DEFAULT '0',
  `owner_type` int(1) unsigned DEFAULT NULL COMMENT '1=thread|2=post',
  `file_name` varchar(100) DEFAULT NULL,
  `width` int(11) unsigned DEFAULT NULL,
  `height` int(11) unsigned DEFAULT NULL,
  `file_type` varchar(10) DEFAULT NULL,
  `user_id` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`owner_id`),
  UNIQUE KEY `file_name` (`file_name`),
  KEY `owner_type` (`owner_type`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- +migrate Down
DROP TABLE `attachments`;
