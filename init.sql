CREATE TABLE IF NOT EXISTS user (
	id character varying(32) PRIMARY KEY,
	name character varying(64) NOT NULL,
	email character varying(320) UNIQUE,
	avatar text,
);
