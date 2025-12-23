CREATE TABLE snippets(
    id varchar(50) NOT NULL PRIMARY KEY,
    title varchar(100) NOT NULL,
    content text NOT NULL,
    created_at timestamp NOT NULL,
    expires timestamp NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created_at);

CREATE TABLE users(
    id serial NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    hashed_password char(60) NOT NULL,
    created timestamptz NOT NULL
);

ALTER TABLE users
    ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO users(name, email, hashed_password, created)
    VALUES ('Alice Jones', 'alice@example.com', '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG', '2022-01-01 09:18:24');

