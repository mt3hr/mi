CREATE TABLE IF NOT EXISTS Task (
    TaskID TEXT PRIMARY KEY,
    CreatedTime TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS TaskTitleInfo (
    TaskTitleID TEXT PRIMARY KEY,
    TaskID TEXT NOT NULL,
    UpdatedTime TEXT NOT NULL,
    Title TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS CheckStateInfo (
    CheckStateID TEXT PRIMARY KEY,
    TaskID TEXT NOT NULL,
    UpdatedTime TEXT NOT NULL,
    IsChecked TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS LimitInfo (
    LimitID TEXT PRIMARY KEY,
    TaskID TEXT NOT NULL,
    UpdatedTime TEXT NOT NULL,
    LimitTime Text
);

CREATE TABLE IF NOT EXISTS StartInfo (
    StartID TEXT PRIMARY KEY,
    TaskID TEXT NOT NULL,
    UpdatedTime TEXT NOT NULL,
    StartTime Text
);

CREATE TABLE IF NOT EXISTS EndInfo (
    EndID TEXT PRIMARY KEY,
    TaskID TEXT NOT NULL,
    UpdatedTime TEXT NOT NULL,
    EndTime Text
);

CREATE TABLE IF NOT EXISTS BoardInfo (
    BoardInfoID TEXT PRIMARY KEY,
    TaskID TEXT NOT NULL,
    UpdatedTime TEXT NOT NULL,
    BoardName TEXT NOT NULL
);