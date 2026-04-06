CREATE TABLE shelters (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  address TEXT NOT NULL,
  phone_number VARCHAR(15) CHECK(
    phone_number ~ '^\+[0-9]+$'
    AND
    char_length(phone_number) BETWEEN 10 AND 15
  ),
  created_at TIMESTAMPTZ DEFAULT NOW()  
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    points INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE partners (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    discount_percent NUMERIC(5,2),
    created_at TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE animals (
    id SERIAL PRIMARY KEY,
    animal_type TEXT NOT NULL,
    name TEXT NOT NULL,
    age int,
    breed TEXT ,
    color TEXT ,
    is_vaccinated BOOLEAN DEFAULT false,
    has_vet_passport BOOLEAN DEFAULT false,
    description TEXT,
    shelter_id INTEGER REFERENCES shelters(id),
    status           TEXT         NOT NULL DEFAULT 'available'
                                  CHECK (status IN ('available', 'adopted', 'reserved')),
    share_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_badges (
    id         SERIAL       PRIMARY KEY,
    user_id    INTEGER      REFERENCES users(id),
    badge      TEXT         NOT NULL
                            CHECK (badge IN ('Друг Животных', 'Спасатель Хвостиков', 'Волонтер')),
    created_at TIMESTAMPTZ  DEFAULT NOW(),
    UNIQUE (user_id, badge)
);

CREATE TABLE animal_photos (
    id         SERIAL       PRIMARY KEY,
    animal_id  INTEGER      REFERENCES animals(id),
    url        TEXT         NOT NULL,
    is_main    BOOLEAN      DEFAULT false,
    created_at TIMESTAMPTZ  DEFAULT NOW()
);

CREATE TABLE adoptions (
    id         SERIAL       PRIMARY KEY,
    user_id    INTEGER      REFERENCES users(id),
    animal_id  INTEGER      REFERENCES animals(id),
    adopted_at TIMESTAMPTZ  DEFAULT NOW()
);