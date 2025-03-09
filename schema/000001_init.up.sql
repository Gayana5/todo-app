CREATE TABLE users
(
    id            serial       not null unique,
    first_name          varchar(255) not null,
    second_name          varchar(255),
    email               varchar(255) not null unique,
    password_hash varchar(255) not null
);
CREATE TABLE todo_goals
(
    id          serial       PRIMARY KEY NOT NULL UNIQUE,
    title       varchar(255) NOT NULL,
    description varchar(255),
    date        date         NOT NULL,
    start_time  time,
    end_time    time,
    priority    boolean      NOT NULL DEFAULT false
);
CREATE TABLE users_goals
(
    id serial not null unique,
    user_id int references users (id) on delete cascade not null,
    list_id int references todo_lists (id) on delete cascade not null
);
CREATE TABLE todo_items
(
    id          serial       PRIMARY KEY NOT NULL UNIQUE,
    title       varchar(255) NOT NULL,
    description varchar(255),
    date        time         NOT NULL,
    start_time  time,
    end_time    time,
    priority    boolean      NOT NULL DEFAULT false
    done        boolean      DEFAULT false
);
CREATE TABLE goal_items
(
    id          serial       PRIMARY KEY NOT NULL UNIQUE,
    title       varchar(255) NOT NULL,
    description varchar(255),
    date        time         NOT NULL,
    start_time  time,
    end_time    time,
    priority    boolean      NOT NULL DEFAULT false
    done        boolean      DEFAULT false
);