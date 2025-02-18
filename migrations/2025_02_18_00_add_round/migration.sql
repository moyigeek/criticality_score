ALTER TABLE scores
ADD COLUMN round integer;

UPDATE scores
SET round = 1;

ALTER TABLE scores
ALTER COLUMN round SET NOT NULL;