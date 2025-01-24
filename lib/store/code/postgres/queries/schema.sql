CREATE TABLE data (
                      id           SERIAL PRIMARY KEY,
                      code_block   TEXT,
                      file_path    TEXT,
                      line         INTEGER,
                      description  TEXT
);
