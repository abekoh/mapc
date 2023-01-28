-- +migrate Up
CREATE TABLE users
(
    id         TEXT PRIMARY KEY,
    name       TEXT     NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE TABLE tasks
(
    id            TEXT PRIMARY KEY,
    user_id       TEXT     NOT NULL,
    title         TEXT     NOT NULL,
    description   TEXT     NOT NULL,
    story_point   INTEGER,
    registered_at DATETIME NOT NULL,
    created_at    DATETIME NOT NULL,
    updated_at    DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE sub_tasks
(
    id            TEXT PRIMARY KEY,
    task_id       TEXT     NOT NULL,
    user_id       TEXT     NOT NULL,
    title         TEXT     NOT NULL,
    description   TEXT     NOT NULL,
    registered_at DATETIME NOT NULL,
    created_at    DATETIME NOT NULL,
    updated_at    DATETIME NOT NULL,
    FOREIGN KEY (task_id) REFERENCES tasks (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +migrate Down
DROP TABLE sub_tasks;
DROP TABLE tasks;
DROP TABLE users;