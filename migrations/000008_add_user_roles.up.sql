CREATE TABLE user_roles (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER REFERENCES users(id),
    role        TEXT NOT NULL CHECK (role IN ('admin', 'shelter_admin', 'user')),
    shelter_id  INTEGER REFERENCES shelters(id), -- NULL если роль глобальная
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, role, shelter_id)
)