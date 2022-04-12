create database if not exists devbook;

use devbook;

drop table if exists publications;
drop table if exists followers;
drop table if exists users;

create table users (
    id int auto_increment primary key,
    name varchar(50) not null,
    nick varchar(50) not null unique,
    email varchar(50) not null unique,
    password varchar(200) not null,
    createdAt timestamp default current_timestamp() 
) engine=innodb;

create table followers (
    user_id int not null,
    follower_id int not null,

    foreign key (user_id) references users(id) on delete cascade,
    foreign key (follower_id) references users(id) on delete cascade,
    
    primary key (user_id, follower_id)
) engine=innodb;

create table publications (
    id int auto_increment primary key,
    title varchar(50) not null,
    content varchar(50) not null,
    author_id int not null,
    likes int not null default 0,
    createdAt timestamp default current_timestamp() ,

    foreign key (author_id) references users(id) on delete cascade
) engine=innodb;
