DROP TABLE IF EXISTS comments;

CREATE TABLE IF NOT EXISTS comments
(
    id SERIAL PRIMARY KEY,
    author character varying NOT NULL,
    message character varying NOT NULL,
    upvote integer,
    downvote integer,
    deleted boolean,
    uuid character varying NOT NULL UNIQUE,
    content_uuid character varying NOT NULL REFERENCES contents(uuid),
    created timestamp with time zone NOT NULL,
    last_edited timestamp with time zone
);


