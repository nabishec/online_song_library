CREATE TABLE songs (
    id  SERIAL PRIMARY KEY,
    song_name TEXT NOT NULL,
    group_name TEXT NOT NULL
);

CREATE TABLE songs_detail (
    id  SERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id) ON DELETE CASCADE,
    release_date TEXT NOT NULL,
    link  TEXT NOT NULL,
    text TEXT NOT NULL
);
