DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS contents;

CREATE TABLE IF NOT EXISTS contents (
    id SERIAL PRIMARY KEY,
    category character varying NOT NULL,
    title character varying NOT NULL,
    author character varying NOT NULL,
    message character varying NOT NULL,
    uuid character varying NOT NULL UNIQUE,
    link character varying,
    image_link character varying,
    upvote integer,
    downvote integer,
    comment_count integer,
    deleted boolean,
    created timestamp with time zone NOT NULL,
    last_edited timestamp with time zone
);

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


