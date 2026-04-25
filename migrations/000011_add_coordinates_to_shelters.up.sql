ALTER TABLE shelters 
ADD COLUMN latitude  DOUBLE PRECISION,
ADD COLUMN longitude DOUBLE PRECISION;


CREATE INDEX idx_shelters_coordinates ON shelters(latitude, longitude)
WHERE latitude IS NOT NULL AND longitude IS NOT NULL;