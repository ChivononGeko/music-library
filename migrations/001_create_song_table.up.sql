CREATE TABLE IF NOT EXISTS songs (
    id VARCHAR(255) PRIMARY KEY, 
    group_name VARCHAR(255) NOT NULL,
    song_name VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    text TEXT NOT NULL,
    link VARCHAR(255) NOT NULL,
    CONSTRAINT unique_song UNIQUE (group_name, song_name)
);

CREATE INDEX IF NOT EXISTS idx_group_name ON songs(group_name);
