create table categories
(
    id          integer
        constraint categories_pk
            primary key autoincrement,
    name        TEXT not null,
    description TEXT
);

create table courses
(
    id          integer
        constraint courses_pk
            primary key autoincrement,
    name        TEXT not null,
    description TEXT,
    category_id integer
        constraint courses_categories_id_fk
            references categories
)