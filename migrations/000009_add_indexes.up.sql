CREATE INDEX idx_animals_status ON animals(status);
CREATE INDEX idx_animals_type ON animals(animal_type);
CREATE INDEX idx_animals_shelter_id ON animals(shelter_id);

CREATE INDEX idx_animal_photos_animal_id ON animal_photos(animal_id);

CREATE INDEX idx_adoptions_user_id ON adoptions(user_id);
CREATE INDEX idx_adoptions_animal_id ON adoptions(animal_id);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_animal_id ON posts(animal_id);
CREATE INDEX idx_posts_status ON posts(status);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);

CREATE INDEX idx_users_email ON users(email);