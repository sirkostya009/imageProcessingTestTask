create table if not exists users (
    id serial primary key,
    username text not null unique,
    email text,
    password_hash text not null,
    constraint username_unique unique (username)
);

create table if not exists images (
    id serial primary key,
    user_id int not null references users(id),
    url text not null
);

create index if not exists idx_images_user_id on images(user_id);
