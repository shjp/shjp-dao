-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- select users
CREATE VIEW users_full AS
  WITH user_groups AS (
    SELECT
      groups.*,
      users.id AS user_id,
      roles.name AS role_name,
      roles.privilege,
      gu.status
    FROM groups_users AS gu
    INNER JOIN users ON users.id = gu.user_id
    INNER JOIN groups ON groups.id = gu.group_id
    INNER JOIN roles ON roles.id = gu.role_id
  )
  SELECT
    users.*,
    COALESCE(json_agg(user_groups) FILTER (WHERE user_groups IS NOT NULL)) AS groups
  FROM users
  LEFT JOIN user_groups ON users.id = user_groups.user_id
  GROUP BY users.id;

-- select announcements
CREATE VIEW announcements_full AS
  SELECT
    announcements.*,
    row_to_json(users) AS creator
  FROM announcements
  INNER JOIN users ON announcements.author_id = users.id;

-- select events
CREATE VIEW events_full AS
  SELECT
    events.*,
    row_to_json(users) AS author
  FROM events
  INNER JOIN users ON events.creator = users.id;

-- select groups
CREATE VIEW groups_full AS
  SELECT
    groups.*,
    COALESCE(json_agg(users) FILTER (WHERE users.id IS NOT NULL), '[]') AS members
  FROM groups
  LEFT JOIN groups_users AS gu ON gu.group_id = groups.id
  LEFT JOIN users ON users.id = gu.user_id
  GROUP BY groups.id;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP VIEW groups_full;
DROP VIEW events_full;
DROP VIEW announcements_full;
DROP VIEW users_full;
