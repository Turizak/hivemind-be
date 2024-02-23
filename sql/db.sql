DROP TABLE IF EXISTS comments;

DROP TABLE IF EXISTS contents;

DROP TABLE IF EXISTS hives;

DROP TABLE IF EXISTS accounts;

CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    username character varying NOT NULL UNIQUE,
    email character varying NOT NULL UNIQUE,
    password character varying NOT NULL,
    uuid character varying NOT NULL UNIQUE,
    deleted boolean,
    banned boolean,
    created timestamp with time zone NOT NULL
);

CREATE TABLE IF NOT EXISTS hives (
    id SERIAL PRIMARY KEY,
    name character varying NOT NULL UNIQUE,
    creator character varying NOT NULL,
    description character varying NOT NULL,
    uuid character varying NOT NULL UNIQUE,
    account_uuid character varying NOT NULL REFERENCES accounts(uuid),
    member_count integer,
    total_upvotes integer,
    total_downvotes integer,
    total_comments integer,
    total_content integer,
    archived boolean,
    banned boolean,
    created timestamp with time zone NOT NULL,
    last_edited timestamp with time zone
);

CREATE TABLE IF NOT EXISTS contents (
    id SERIAL PRIMARY KEY,
    hive character varying NOT NULL,
    title character varying NOT NULL,
    author character varying NOT NULL,
    message character varying NOT NULL,
    uuid character varying NOT NULL UNIQUE,
    hive_uuid character varying NOT NULL REFERENCES hives(uuid),
    account_uuid character varying NOT NULL REFERENCES accounts(uuid),
    link character varying,
    image_link character varying,
    upvote integer,
    downvote integer,
    comment_count integer,
    deleted boolean,
    created timestamp with time zone NOT NULL,
    last_edited timestamp with time zone
);

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    author character varying NOT NULL,
    message character varying NOT NULL,
    upvote integer,
    downvote integer,
    deleted boolean,
    uuid character varying NOT NULL UNIQUE,
    parent_uuid character varying,
    content_uuid character varying NOT NULL REFERENCES contents(uuid),
    account_uuid character varying NOT NULL REFERENCES accounts(uuid),
    created timestamp with time zone NOT NULL,
    last_edited timestamp with time zone
);