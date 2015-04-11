-- +migrate Up
CREATE TABLE `stats` (
  `owner_id` int(11) unsigned NOT NULL,
  `owner_type` int(11) NOT NULL,
  `date` date NOT NULL,
  `stat_type` int(11) NOT NULL,
  `stat_value` int(11) NOT NULL,
  UNIQUE KEY `stats_unique` (`owner_id`,`owner_type`,`date`,`stat_type`),
  KEY `owner_id` (`owner_id`),
  KEY `owner_type` (`owner_type`),
  KEY `date` (`date`),
  KEY `stat_type` (`stat_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- +migrate Down
DROP TABLE `stats`;
