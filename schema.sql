DROP TABLE IF EXISTS `contribs`;
CREATE TABLE `contribs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `week` int(11) unsigned NOT NULL,
  `org_id` int(11) unsigned NOT NULL,
  `repo_id` int(11) unsigned NOT NULL,
  `user_id` int(11) unsigned NOT NULL,
  `commits` int(11) unsigned NOT NULL,
  `additions` int(11) unsigned NOT NULL,
  `deletions` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `record` (`week`,`org_id`,`repo_id`,`user_id`),
  KEY `search` (`org_id`,`week`,`repo_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `org_members`;
CREATE TABLE `org_members` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `org_id` int(11) unsigned NOT NULL,
  `user_id` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `membership` (`org_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `orgs`;
CREATE TABLE `orgs` (
  `id` int(11) unsigned NOT NULL,
  `login` varchar(255) NOT NULL DEFAULT '',
  `company` varchar(255) NOT NULL DEFAULT '',
  `avatar_url` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `login` (`login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `repos`;
CREATE TABLE `repos` (
  `id` int(11) unsigned NOT NULL,
  `org_id` int(11) unsigned NOT NULL,
  `name` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(255) NOT NULL,
  `is_private` tinyint(1) NOT NULL,
  `is_fork` tinyint(1) NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `owner_repo` (`org_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `team_members`;
CREATE TABLE `team_members` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `org_id` int(11) unsigned NOT NULL,
  `team_id` int(11) unsigned NOT NULL,
  `user_id` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `membership` (`team_id`,`user_id`),
  KEY `search` (`org_id`,`user_id`),
  KEY `team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `team_repos`;
CREATE TABLE `team_repos` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `org_id` int(11) unsigned NOT NULL,
  `team_id` int(11) unsigned NOT NULL,
  `repo_id` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `access` (`team_id`,`repo_id`),
  KEY `search` (`org_id`,`team_id`),
  KEY `team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `teams`;
CREATE TABLE `teams` (
  `id` int(11) unsigned NOT NULL,
  `org_id` int(11) NOT NULL,
  `slug` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `permission` varchar(20) NOT NULL DEFAULT '',
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `owner` (`org_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `tokens`;
CREATE TABLE `tokens` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user` varchar(255) NOT NULL DEFAULT '',
  `token` varchar(40) NOT NULL DEFAULT '',
  `quota` int(11) unsigned NOT NULL,
  `remaining` int(11) unsigned NOT NULL,
  `reset_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `owner` (`user`),
  UNIQUE KEY `token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL,
  `login` varchar(255) NOT NULL DEFAULT '',
  `name` varchar(255) NOT NULL DEFAULT '',
  `avatar_url` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `login` (`login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
