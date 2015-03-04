DROP TABLE IF EXISTS `contributions`;
CREATE TABLE `contributions` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `week` int(11) unsigned NOT NULL,
  `author` varchar(255) NOT NULL DEFAULT '',
  `owner` varchar(255) NOT NULL DEFAULT '',
  `repo` varchar(255) NOT NULL DEFAULT '',
  `commits` int(11) unsigned NOT NULL,
  `additions` int(11) unsigned NOT NULL,
  `deletions` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `week_author_repo` (`week`,`author`,`owner`,`repo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `members`;
CREATE TABLE `members` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `team_id` int(10) unsigned NOT NULL,
  `login` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `team` (`team_id`),
  KEY `login` (`login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `owners`;
CREATE TABLE `owners` (
  `login` varchar(255) NOT NULL DEFAULT '',
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `repos`;
CREATE TABLE `repos` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `owner` varchar(255) NOT NULL DEFAULT '',
  `repo` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime DEFAULT NULL,
  `is_private` tinyint(1) NOT NULL,
  `is_fork` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `owner_repo` (`owner`,`repo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `teams`;
CREATE TABLE `teams` (
  `id` int(11) unsigned NOT NULL,
  `owner` varchar(255) NOT NULL DEFAULT '',
  `team` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `owner` (`owner`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `tokens`;
CREATE TABLE `tokens` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `owner` varchar(255) NOT NULL DEFAULT '',
  `token` varchar(40) NOT NULL DEFAULT '',
  `limit` int(11) unsigned NOT NULL,
  `remaining` int(11) unsigned NOT NULL,
  `reset_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `owner` (`owner`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
