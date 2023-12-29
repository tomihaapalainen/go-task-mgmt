-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT,
    password_hash TEXT,
    role_id INTEGER,
    FOREIGN KEY (role_id) REFERENCES role(id),
    UNIQUE (email)
);

CREATE INDEX IF NOT EXISTS user_email_ix ON user (email);

CREATE TABLE IF NOT EXISTS role (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
);

INSERT INTO role (name) values('admin');
INSERT INTO role (name) values('project manager');
INSERT INTO role (name) values('user');

CREATE TABLE IF NOT EXISTS permission (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
);

INSERT INTO permission (name) values('all');
INSERT INTO permission (name) values('create project');
INSERT INTO permission (name) values('read project');
INSERT INTO permission (name) values('update project');
INSERT INTO permission (name) values('delete project');
INSERT INTO permission (name) values('create task');
INSERT INTO permission (name) values('read task');
INSERT INTO permission (name) values('update task');
INSERT INTO permission (name) values('delete task');
INSERT INTO permission (name) values('manage users');
INSERT INTO permission (name) values('manage roles');


CREATE TABLE IF NOT EXISTS role_permission (
    role_id INTEGER,
    permission_id INTEGER,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES role(id),
    FOREIGN KEY (permission_id) REFERENCES permission(id)
);

INSERT INTO role_permission (role_id, permission_id) values(1, 1);
INSERT INTO role_permission (role_id, permission_id) values(2, 2);
INSERT INTO role_permission (role_id, permission_id) values(2, 3);
INSERT INTO role_permission (role_id, permission_id) values(2, 4);
INSERT INTO role_permission (role_id, permission_id) values(2, 5);
INSERT INTO role_permission (role_id, permission_id) values(2, 6);
INSERT INTO role_permission (role_id, permission_id) values(2, 7);
INSERT INTO role_permission (role_id, permission_id) values(2, 8);
INSERT INTO role_permission (role_id, permission_id) values(2, 9);
INSERT INTO role_permission (role_id, permission_id) values(3, 6);
INSERT INTO role_permission (role_id, permission_id) values(3, 7);
INSERT INTO role_permission (role_id, permission_id) values(3, 8);
INSERT INTO role_permission (role_id, permission_id) values(3, 9);

CREATE TABLE IF NOT EXISTS project (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    name TEXT,
    description TEXT,
    FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
    UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS project_name_ix ON project (name);

CREATE TABLE IF NOT EXISTS task (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER,
    assignee_id INTEGER,
    creator_id INTEGER,
    title TEXT,
    content TEXT,
    status TEXT,
    FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE CASCADE,
    FOREIGN KEY (assignee_id) REFERENCES user(id),
    FOREIGN KEY (creator_id) REFERENCES user(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE task;
DROP TABLE project;
DROP TABLE role_permission;
DROP TABLE permission;
DROP TABLE role;
DROP TABLE user;
-- +goose StatementEnd
