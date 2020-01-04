-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE mass_files (
  id UUID PRIMARY KEY NOT NULL,
  type TEXT NOT NULL,
  name TEXT NOT NULL,
  date TEXT NOT NULL,
  url TEXT
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE mass_files;
