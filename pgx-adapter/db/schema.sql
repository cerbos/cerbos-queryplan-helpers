CREATE SCHEMA IF NOT EXISTS cerbforce;

SET search_path TO cerbforce;

CREATE TABLE IF NOT EXISTS companies (
    id SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL PRIMARY KEY,
    username VARCHAR(128) NOT NULL,
    email VARCHAR(128) NOT NULL,
    name VARCHAR(128) NOT NULL,
    role VARCHAR(128) NOT NULL,
    department VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    first_name VARCHAR(128) NOT NULL,
    last_name VARCHAR(128) NOT NULL,
    owner_id INT NOT NULL REFERENCES cerbforce.users(id),
    company_id INT NULL REFERENCES cerbforce.companies(id),
    active BOOLEAN default false,
    marketing_opt_in BOOLEAN default false
);
DO $$
BEGIN
    CREATE USER cerbforce_user WITH PASSWORD 'cerb';
EXCEPTION WHEN duplicate_object THEN
END$$;
GRANT CONNECT ON DATABASE postgres TO cerbforce_user;
GRANT USAGE ON SCHEMA cerbforce TO cerbforce_user;
GRANT SELECT,INSERT,UPDATE,DELETE ON cerbforce.users, cerbforce.companies, cerbforce.contacts TO cerbforce_user;
