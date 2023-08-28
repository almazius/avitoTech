CREATE DATABASE avitoTest;

GRANT ALL PRIVILEGES ON DATABASE avitoTest to Postgres;
\c avitoTest

CREATE TABLE IF NOT exists segments(
    segment text primary key,
    userId bigint[] not null
);

CREATE TABLE IF NOT exists  tempSegment(
                            userid bigint,
                            segmentName text,
                            timeEnd time
);