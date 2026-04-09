ALTER TABLE adoptions 
DROP COLUMN updated_at;

ALTER TABLE adoptions 
RENAME COLUMN created_at TO adopted_at;
