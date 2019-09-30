create table if not exists users
(
	id serial not null
		constraint users_pkey
			primary key,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone,
	deleted_at timestamp with time zone,
	uuid uuid not null,
	access_level bigint not null,
	first_name text not null,
	last_name text not null,
	email text not null
		constraint users_email_key
			unique,
	password text not null,
	date_of_birth timestamp with time zone
);

alter table users owner to root;

create table if not exists votes
(
	id serial not null
		constraint votes_pkey
			primary key,
	uuid uuid not null,
	title text,
	description text,
	uuid_vote varchar(100) [],
	start_date timestamp with time zone,
	end_date timestamp with time zone,
	created_at timestamp with time zone not null,
	updated_at timestamp with time zone,
	deleted_at timestamp with time zone
);

alter table votes owner to root;


INSERT INTO users (created_at, updated_at, deleted_at, uuid, access_level, first_name, last_name, email, password, date_of_birth) VALUES ('2019-09-18 20:52:01.529000', null, null, '584cbd52-8f2a-46b6-8007-8988d3de1f1b', 1, 'admin', 'admin', 'admin@admin.com', 'MfemXjFVhqwZi9eYtmKc5JA9CJlHbVdBqfMuLlIbamY=', '1991-09-29 20:54:50.297000');