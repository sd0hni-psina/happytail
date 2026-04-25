DROP INDEX IF EXISTS idx_shelters_coordinates;

ALTER TABLE shelters 
DROP COLUMN IF EXISTS latitude,
DROP COLUMN IF EXISTS longitude;