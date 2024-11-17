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

create table IF NOT EXISTS account (
    id int not null AUTO_INCREMENT,
    user_id varchar(256) NOT NULL,
    user_name varchar(256) NOT NULL,
    user_email varchar(512) NOT NULL,
    created_at int not null,
    updated_at int not null,
    primary key (`id`),
    unique key `unique_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='account info';


