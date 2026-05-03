ALTER TABLE shelters ADD COLUMN deleted_at TIMESTAMPTZ;
 
CREATE INDEX idx_shelters_deleted_at ON shelters(deleted_at)
    WHERE deleted_at IS NULL;