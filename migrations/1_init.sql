-- +migrate Up
CREATE TABLE "users"
(
    "user_id"    text primary key,
    "full_name"  text,
    "email"      text UNIQUE,
    "password"   text,
    "role"       text,
    "created_at" timestamptz not null,
    "updated_at" timestamptz not null
);

CREATE TABLE "repos"
(
    "name"        text primary key,
    "description" text,
    "url"         text UNIQUE,
    "color"       text,
    "lang"        text,
    "fork"        text,
    "stars"       text,
    "stars_today" text,
    "build_by"    text,
    "created_at"  timestamptz not null,
    "updated_at"  timestamptz not null
);

CREATE TABLE "bookmarks"
(
    "bid"        text primary key,
    "user_id"    text,
    "repo_name"  text UNIQUE,
    "created_at" timestamptz not null,
    "updated_at" timestamptz not null,
    unique (user_id, repo_name)
);

--ALTER TABLE "bookmarks" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
-- p0EE6bmdpH5cpLeu
--https://cloud.google.com/vpc/docs/using-firewalls#creating_firewall_rules
ALTER TABLE "bookmarks"
    ADD FOREIGN KEY ("repo_name") REFERENCES "repos" ("name");

-- +migrate Down
DROP TABLE bookmarks;
DROP TABLE users;
DROP TABLE repos;
