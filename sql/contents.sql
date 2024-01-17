DROP TABLE IF EXISTS contents;

CREATE TABLE IF NOT EXISTS contents (
    id SERIAL PRIMARY KEY,
    category character varying NOT NULL,
    title character varying NOT NULL,
    author character varying NOT NULL,
    message character varying NOT NULL,
    uuid character varying NOT NULL,
    link character varying,
    image_link character varying,
    upvote integer,
    downvote integer,
    comment_count integer,
    deleted boolean,
    created timestamp with time zone NOT NULL,
    last_edited timestamp with time zone
);