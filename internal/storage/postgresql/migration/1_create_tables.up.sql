CREATE TABLE groups (
    id  UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE songs (
    id  UUID PRIMARY KEY,
    name TEXT NOT NULL,
    release_date DATE NOT NULL,
    link  TEXT NOT NULL,
    text TEXT NOT NULL,
    group_id UUID NOT NULL REFERENCES groups(id)
);
