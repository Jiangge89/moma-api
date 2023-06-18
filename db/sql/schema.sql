create database moma_api;
use moma_api;

create table IF NOT EXISTS rate (
    id int not null AUTO_INCREMENT,
    from_code varchar(8) not null,
    to_code varchar(8) not null,
    rate float not null,
    created_at int not null,
    updated_at int not null,
    primary key (`id`),
    unique key `unique_code` (`from_code`,`to_code`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='rate table';