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
    WHERE gu.status = 'accepted'
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
  WITH rsvps AS (
    SELECT
      row_to_json(users.*) AS user,
      ue.event_id AS event_id,
      ue.rsvp AS rsvp
    FROM users
    INNER JOIN users_events AS ue ON ue.user_id = users.id
  )
  SELECT
    events.*,
    row_to_json(users) AS author,
    COALESCE(json_agg(rsvps) FILTER (WHERE rsvps.event_id IS NOT NULL)) AS rsvps
  FROM events
  LEFT JOIN rsvps ON rsvps.event_id = events.id
  LEFT JOIN users ON events.creator = users.id
  GROUP BY events.id, users.id;

-- select groups
CREATE VIEW groups_full AS
  SELECT
    groups.*,
    COALESCE(json_agg(DISTINCT
      to_jsonb(users)
      || jsonb_build_object(
        'role_name', user_roles.name,
        'privilege', user_roles.privilege)
      || jsonb_build_object(
        'status', gu.status)
    ) FILTER (WHERE users.id IS NOT NULL), '[]') AS members,
    COALESCE(json_agg(DISTINCT roles) FILTER (WHERE roles IS NOT NULL), '[]') AS roles
  FROM groups
  LEFT JOIN groups_users AS gu ON gu.group_id = groups.id
  LEFT JOIN users ON users.id = gu.user_id
  LEFT JOIN roles AS user_roles ON user_roles.id = gu.role_id
  LEFT JOIN roles ON roles.group_id = groups.id
  GROUP BY groups.id;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP VIEW groups_full;
DROP VIEW events_full;
DROP VIEW announcements_full;
DROP VIEW users_full;
