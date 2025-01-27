CREATE TABLE data (
    project      TEXT,
    id           BIGINT,
    code_block   TEXT,
    file_path    TEXT,
    line         INTEGER,
    description  TEXT,

    PRIMARY KEY (project, id)
);
