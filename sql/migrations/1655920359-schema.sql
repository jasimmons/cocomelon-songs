-- Migration: schema
-- Created at: 2022-06-22 13:52:39
-- ====  UP  ====

BEGIN;

CREATE TABLE `songs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `season` int NOT NULL,
  `episode` int NOT NULL,
  `start_time` varchar(255) DEFAULT '',
  PRIMARY KEY (`id`)
) CHARSET=utf8;

COMMIT;

-- ==== DOWN ====

BEGIN;

DROP TABLE `songs`;

COMMIT;
