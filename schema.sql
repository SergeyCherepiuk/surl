create table users(
    username varchar(30) primary key,
    password varchar(256) not null
);

create table urls(
    username varchar(30) references users(username) not null,
    hash varchar(8) not null,
    origin text not null,
    primary key(username, hash)
);
