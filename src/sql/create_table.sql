CREATE TABLE articles (
	id varchar(16) primary key,
	title varchar(255) not null,
	slug varchar(255) not null,
	description varchar(255),
	date varchar(10),
	md text not null
);
