-- +migrate Up
CREATE TABLE `clubs_groups` (
  `club_id` smallint(5) unsigned DEFAULT NULL,
  `group_id` int(11) unsigned DEFAULT NULL,
  UNIQUE KEY `club_id` (`club_id`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- +migrate Down
DROP TABLE `clubs_groups`;
