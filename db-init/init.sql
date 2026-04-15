-- CREATE TABLE IF NOT EXISTS contacts (
--     id SERIAL PRIMARY KEY,
--     name TEXT UNIQUE,
--     phone TEXT,
--     image_url TEXT
-- );

-- 1. Create the table only if it doesn't exist
CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE,
    phone TEXT,
    image_url TEXT
);

-- 2. Add one dummy user (skips if 'John Doe' already exists)
INSERT INTO contacts (name, phone)
VALUES ('John Doe', '1234567890')
ON CONFLICT (name) DO NOTHING;