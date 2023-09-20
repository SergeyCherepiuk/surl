create table users(
    username varchar(30) primary key,
    password varchar(256) not null
);

create table urls(
    username varchar(30) not null references users(username) on delete cascade on update cascade,
    hash varchar(8) not null,
    origin text not null,
    created_at timestamp not null default now(),
    last_used_at timestamp not null default 'epoch',
    primary key(username, hash)
);
