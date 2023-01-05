-- TODO index all the tables

CREATE TABLE IF NOT EXISTS users (
	id varchar(32) PRIMARY KEY,
	name varchar(64) NOT NULL,
	email varchar(320) UNIQUE,
	avatar_url text,

	google_id varchar(32) UNIQUE,

	active boolean NOT NULL DEFAULT true,
);

CREATE TABLE IF NOT EXISTS politicians (
	id SERIAL PRIMARY KEY,
	name varchar(64) NOT NULL,
);

CREATE TABLE IF NOT EXISTS parties (
	id int PRIMARY KEY,
	name varchar(64) NOT NULL,
	
	chairman varchar(64),
	established_date date,
	filing_date date,

	main_office_address text,
	mailing_address text,

	phone_number varchar(32),
	status smallint NOT NULL DEFAULT 0,
);

CREATE TABLE IF NOT EXISTS staging_data (
	id SERIAL PRIMARY KEY,
	table_name varchar(64) NOT NULL,
	data json NOT NULL,
	timestamp timestamp NOT NULL DEFAULT NOW(),
	status smallint NOT NULL DEFAULT 0,
);
