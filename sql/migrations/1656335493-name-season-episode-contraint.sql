-- Migration: name-season-episode-contraint
-- Created at: 2022-06-27 09:11:33
-- ====  UP  ====

BEGIN;

ALTER TABLE songs
  ADD CONSTRAINT name_season_ep_uniq UNIQUE (name, season, episode);

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE songs
  DROP CONSTRAINT name_season_ep_uniq;

COMMIT;
