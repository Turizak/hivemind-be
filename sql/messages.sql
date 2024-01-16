DROP TABLE IF EXISTS messages;

CREATE TABLE IF NOT EXISTS messages
(
    id SERIAL PRIMARY KEY,
    category character varying NOT NULL,
    title character varying NOT NULL,
    author character varying NOT NULL,
    message character varying NOT NULL,
    type character varying NOT NULL,
    upvote integer,
    downvote integer,
    comment_count integer,
    deleted boolean,
    uuid character varying NOT NULL,
    created timestamp with time zone NOT NULL,
    last_edited timestamp with time zone
);