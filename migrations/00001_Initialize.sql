-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE groups (
  id UUID PRIMARY KEY NOT NULL,
  name VARCHAR(60) NOT NULL,
  description TEXT,
  image_url TEXT
);

CREATE TABLE roles (
  id UUID PRIMARY KEY NOT NULL,
  group_id UUID NOT NULL REFERENCES groups (id),
  name TEXT NOT NULL,
  privilege INT NOT NULL
);

CREATE TABLE users (
  id UUID PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  email TEXT,
  baptismal_name TEXT,
  birthday TIMESTAMP,
  feastday TIMESTAMP,
  created TIMESTAMP NOT NULL DEFAULT now(),
  last_active TIMESTAMP,
  account_type TEXT,
  account_secret TEXT
);

CREATE TABLE events (
  id UUID PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  start TIMESTAMP,
  "end" TIMESTAMP,
  creator UUID NOT NULL REFERENCES users (id),
  deadline TIMESTAMP,
  allow_maybe BOOLEAN NOT NULL,
  description TEXT,
  location POINT,
  location_description TEXT,
  created TIMESTAMP NOT NULL DEFAULT now(),
  updated TIMESTAMP
);

CREATE TABLE announcements (
  id UUID PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT now(),
  updated TIMESTAMP,
  author_id UUID NOT NULL REFERENCES users (id),
  content TEXT
);

CREATE TABLE comments (
  id UUID PRIMARY KEY NOT NULL,
  author_id UUID NOT NULL REFERENCES users (id),
  created TIMESTAMP DEFAULT now(),
  updated TIMESTAMP,
  parent_id UUID NOT NULL,
  parent_type VARCHAR(20) NOT NULL
    CHECK(parent_type IN ('announcement', 'event')),
  content TEXT
);

CREATE TABLE groups_users (
  user_id UUID NOT NULL REFERENCES users (id),
  group_id UUID NOT NULL REFERENCES groups (id),
  role_id UUID REFERENCES roles (id),
  status VARCHAR(20) NOT NULL
    CHECK(status IN ('accepted', 'pending'))
);

CREATE UNIQUE INDEX group_user_index ON groups_users (user_id, group_id);

CREATE TABLE groups_events (
  group_id UUID NOT NULL REFERENCES groups (id),
  event_id UUID NOT NULL REFERENCES events (id)
);

CREATE TABLE groups_announcements (
  group_id UUID NOT NULL REFERENCES groups (id),
  announcement_id UUID NOT NULL REFERENCES announcements (id)
);

CREATE UNIQUE INDEX groups_announcements_index ON groups_announcements (group_id, announcement_id);

CREATE TABLE users_events (
  user_id UUID NOT NULL REFERENCES users (id),
  event_id UUID NOT NULL REFERENCES events (id),
  rsvp VARCHAR(10) NOT NULL
    CHECK(rsvp IN ('yes', 'no', 'maybe', 'unanswered'))
);

CREATE UNIQUE INDEX users_events_index ON users_events (user_id, event_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE users_events;
DROP TABLE groups_announcements;
DROP TABLE groups_events;
DROP TABLE groups_users;
DROP TABLE comments;
DROP TABLE events;
DROP TABLE announcements;
DROP TABLE users;
DROP TABLE roles;
DROP TABLE groups;