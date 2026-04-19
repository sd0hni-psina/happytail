CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    animal_id INTEGER REFERENCES animals(id),
    listing_type TEXT NOT NULL CHECK(listing_type IN ('sale', 'give')),
    price_amount BIGINT,
    price_currency TEXT,
    reason TEXT,
    photo_urls TEXT[],
    contact_info TEXT,
    status TEXT NOT NULL DEFAULT 'inactive' CHECK (status in ('active', 'inactive', 'deleted')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
)