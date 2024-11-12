CREATE TYPE STATUSCONTENT AS ENUM ('DRAFT', 'PUBLISH', 'PENDING');

CREATE TABLE IF NOT EXISTS contents (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    created_by_id INT REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(200) NOT NULL,
    excerpt TEXT NOT NULL,
    description TEXT NOT NULL,
    image TEXT NULL,
    status STATUSCONTENT NOT NULL DEFAULT 'PUBLISH',
    tags TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_contents_category_id ON contents(category_id);
CREATE INDEX idx_contents_created_by_id ON contents(created_by_id);