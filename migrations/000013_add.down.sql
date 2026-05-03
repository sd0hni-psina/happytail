DROP INDEX IF EXISTS idx_shelters_deleted_at;
ALTER TABLE shelters DROP COLUMN IF EXISTS deleted_at;