-- +goose Up
CREATE TABLE departments (
    id SERIAL PRIMARY KEY,

    name VARCHAR(200) NOT NULL,

    parent_id INT REFERENCES departments(id)
    ON DELETE CASCADE,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX departments_root_unique
ON departments(name)
WHERE parent_id IS NULL;

CREATE UNIQUE INDEX departments_child_unique
ON departments(parent_id, name)
WHERE parent_id IS NOT NULL;


-- +goose Down
DROP TABLE departments;