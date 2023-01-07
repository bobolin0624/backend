-- TODO index all the tables

CREATE TABLE IF NOT EXISTS users (
	id varchar(32) PRIMARY KEY,
	name varchar(64) NOT NULL,
	email varchar(320) UNIQUE,
	avatar_url text,

	google_id varchar(32) UNIQUE,

	active boolean NOT NULL DEFAULT true,

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS politicians (
	id SERIAL PRIMARY KEY,
	name varchar(64) NOT NULL,
	en_name varchar(64),
	avatar_url text,

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS legislators (
	politicians_id int NOT NULL REFERENCES politicians(id),
	term smallint NOT NULL,
	session smallint NOT NULL,
	-- create a committee table if needed
	committee varchar(64) NOT NULL,
	onboard_date date,
	resign_date date,
	resign_reason text,

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
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

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS staging_data (
	group_id uuid NOT NULL,
	table_name varchar(64) NOT NULL,
	data jsonb NOT NULL,
	status smallint NOT NULL DEFAULT 0,

	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
