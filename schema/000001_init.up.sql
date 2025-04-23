-- Пользователи
CREATE TABLE users
(
    id            serial      not null unique,
    first_name    varchar(255) not null,
    second_name   varchar(255),
    email         varchar(255) not null unique,
    password_hash varchar(255) not null
);
-- цель
CREATE TABLE todo_goals
(
    id              serial       PRIMARY KEY NOT NULL UNIQUE,
    title           varchar(255) NOT NULL,
    description     varchar(255),
    colour          INT,
    completed_tasks INT         default 0,
    total_tasks     INT         default 0
);
-- цели пользователя
CREATE TABLE users_goals
(
    id      serial not null unique,
    user_id int references users (id) on delete cascade not null,
    goal_id int references todo_goals (id) on delete cascade not null
);
-- задача
CREATE TABLE todo_tasks
(
    id          serial       PRIMARY KEY NOT NULL UNIQUE,
    user_id     int references users (id) on delete cascade not null,
    title       varchar(255) NOT NULL,
    description varchar(255),
    goal_id     int,
    end_date    date         NOT NULL,
    start_time  time,
    end_time    time,
    colour      int,
    done        boolean      DEFAULT false
);
-- задачи цели
CREATE TABLE goal_tasks
(
    id          serial       PRIMARY KEY NOT NULL UNIQUE,
    task_id     int references todo_tasks (id) on delete cascade not null,
    goal_id     int references todo_goals (id) on delete cascade not null
);