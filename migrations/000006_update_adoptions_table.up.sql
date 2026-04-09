ALTER TABLE adoptions 
RENAME COLUMN adopted_at TO created_at;

ALTER TABLE adoptions 
ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();