-- name: GetUser :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY name;

-- name: GetTask :one
SELECT *
FROM tasks
WHERE id = ?
LIMIT 1;

-- name: ListTasks :many
SELECT *
FROM tasks
WHERE id = ?
ORDER BY title;

-- name: SubTask :one
SELECT *
FROM sub_tasks
WHERE id = ?
LIMIT 1;

-- name: ListSubTasks :many
SELECT *
FROM sub_tasks
WHERE id = ?
ORDER BY title;
