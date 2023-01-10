CREATE TABLE articles (
	id varchar(16) primary key,
	title varchar(255) not null,
	slug varchar(255) not null,
	description text,
	date varchar(10) not null,
	md text not null
);
