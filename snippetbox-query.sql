CREATE TABLE snippets(
    id varchar(50) NOT NULL PRIMARY KEY,
    title varchar(100) NOT NULL,
    content text NOT NULL,
    created_at timestamp NOT NULL,
    expires timestamp NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created_at);

CREATE TABLE sessions(
    token text PRIMARY KEY,
    data bytea NOT NULL,
    expiry timestamptz NOT NULL
);1

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

TRUNCATE
    snippets INSERT INTO snippets(id, title, content, created_at, expires)
        VALUES ('snippet-1', 'An old silent pond', 'Pond de ring', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1' YEAR);

INSERT INTO snippets(id, title, content, created_at, expires)
    VALUES ('snippet-2', 'Over the wintry forest', 'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + INTERVAL '1' YEAR);

CREATE TABLE users(
    id serial NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    hashed_password char(60) NOT NULL,
    created timestamptz NOT NULL
);

ALTER TABLE users
    ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO snippets(id, title, content, created_at, expires)
    VALUES ('snippet-3', < div > < label > DELETE in: < / label > < input type = 'radio' name = 'expires' value = '365' checked > One Year < input type = 'radio' name = 'expires' value = '7' > One Week < input type = 'radio' name = 'expires' value = '1' > One Day < / div > < div > < input type = 'submit' value = 'Publish snippet' > < / div > 'First autumn morning',
        'First autumn morning\nthe mirror I stare into\nshows my husband''s face.\n\n– Tsukatsuki Rio',
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP + INTERVAL '1' YEAR);

TRUNCATE
    snippets UPDATE
        snippets
    SET
        title = 'anjay',
        content = 'anjay' CREATE DATABASE test_snippetbox character
        SET
            utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE DATABASE test_snippetbox WITH ENCODING 'UTF8' LC_COLLATE 'en_US.UTF-8' LC_CTYPE 'en_US.UTF-8' TEMPLATE template0;

