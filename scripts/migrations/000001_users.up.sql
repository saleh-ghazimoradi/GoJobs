CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_admin BOOLEAN DEFAULT false,
    profile_picture TEXT
);

INSERT INTO users (username, password, email, is_admin)
VALUES ('admin1', '$2a$10$K8vrm0oVPaDODXrSKNSnAOESnEW34TLZp33eJQs.WWiu9N.08fV/O', 'admin1@gmail.com', true)
    ON CONFLICT (username) DO NOTHING;

