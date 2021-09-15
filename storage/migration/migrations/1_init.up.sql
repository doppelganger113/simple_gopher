CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- USERS
CREATE TABLE IF NOT EXISTS users
(
    id           UUID PRIMARY KEY    NOT NULL DEFAULT uuid_generate_v4(),
    email        VARCHAR(255) UNIQUE NOT NULL,
    role         VARCHAR(100)        NOT NULL,
    cog_username VARCHAR(255) UNIQUE NOT NULL,
    cog_sub      VARCHAR(255) UNIQUE NOT NULL,
    cog_name     VARCHAR(255) UNIQUE NOT NULL,
    created_at   timestamp           NOT NULL DEFAULT now(),
    updated_at   timestamp,
    disabled     BOOLEAN                      DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_users_created_at ON users (created_at);
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);

-- IMAGES
CREATE TABLE IF NOT EXISTS images
(
    id         UUID PRIMARY KEY    NOT NULL DEFAULT uuid_generate_v4(),
    name       VARCHAR(255) UNIQUE NOT NULL,
    created_at timestamp           NOT NULL DEFAULT now(),
    updated_at timestamp,
    format     VARCHAR(30)         NOT NULL,
    original   VARCHAR(255)        NOT NULL,
    domain     VARCHAR(255)        NOT NULL,
    path       VARCHAR(255)        NOT NULL,
    sizes      jsonb               NOT NULL,
    author_id  UUID                NOT NULL,

    CONSTRAINT author_fk
        FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_images_createdAt ON images (created_at);
CREATE INDEX IF NOT EXISTS idx_images_authorId ON images (author_id);

-- TAGS
CREATE TABLE IF NOT EXISTS tags
(
    id         UUID PRIMARY KEY    NOT NULL DEFAULT uuid_generate_v4(),
    value      VARCHAR(255) UNIQUE NOT NULL,
    created_at timestamp           NOT NULL DEFAULT now(),
    updated_at timestamp,
    author_id  UUID,

    CONSTRAINT author_fk
        FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_tags_createdAt ON tags (created_at);
CREATE INDEX IF NOT EXISTS idx_tags_authorId ON tags (author_id);

-- IMAGES_TAGS
CREATE TABLE IF NOT EXISTS images_tags
(
    id       UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    tag_id   UUID             NOT NULL,
    image_id UUID             NOT NULL,

    CONSTRAINT tag_fk
        FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE SET NULL,
    CONSTRAINT image_fk
        FOREIGN KEY (image_id) REFERENCES images (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_images_tags ON images_tags (tag_id, image_id);
