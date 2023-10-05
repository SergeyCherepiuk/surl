create table users(
    username varchar(30) primary key,
    email varchar(320) not null unique,
    password varchar(256) not null,
    is_verified boolean not null default false
);

create table verification_requests(
    id varchar(36) not null,
    username varchar(30) not null references users(username) on delete cascade on update cascade,
    expires_at timestamp not null,
    primary key(id, username)
);

create table password_reset_requests(
    id varchar(36) not null,
    username varchar(30) not null references users(username) on delete cascade on update cascade,
    expires_at timestamp not null,
    primary key(id, username)
);

create table urls(
    username varchar(30) not null references users(username) on delete cascade on update cascade,
    hash varchar(8) not null,
    origin text not null,
    expires_at timestamp not null,
    created_at timestamp not null default now(),
    last_used_at timestamp not null default 'epoch',
    primary key(username, hash)
);
