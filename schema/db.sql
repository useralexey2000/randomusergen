CREATE TABLE IF NOT EXISTS users (
  id serial not null,
  title character varying (100),
  first_name character varying (100),
  last_name character varying (100),
  gender character varying (10),
  email character varying (100),
  phone character varying (30),
  cell character varying (30),
  nat character varying (5),
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS locations (
  user_id integer not null,
  id serial not null,
  city character varying (100),
  state character varying (100),
  country character varying (100),
  postcode character varying (100),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS streets (
  location_id integer not null,
  id serial not null,
  number integer,
  name character varying (100),
  PRIMARY KEY (id),
  FOREIGN KEY (location_id) REFERENCES locations(id)
);

CREATE TABLE IF NOT EXISTS coordinates (
  location_id integer not null,
  id serial not null,
  latitude character varying (20),
  longitude character varying (20),
  PRIMARY KEY (id),
  FOREIGN KEY (location_id) REFERENCES locations(id)
);

CREATE TABLE IF NOT EXISTS timezones (
  location_id integer not null,
  id serial not null,
  offset_time character varying (20),
  description character varying (100),
  PRIMARY KEY (id),
  FOREIGN KEY (location_id) REFERENCES locations(id)
);

CREATE TABLE IF NOT EXISTS logins (
  user_id integer not null,
  id serial not null,
  uuid character varying (100),
  username character varying (100),
  password character varying (100),
  salt character varying (100),
  md5 character varying (100),
  sha1 character varying (100),
  sha256 character varying (100),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS dob (
  user_id integer not null,
  id serial not null,
  date timestamp(0) with time zone,
  age integer,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS registered (
  user_id integer not null,
  id serial not null,
  date timestamp(0) with time zone,
  age integer,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS id (
  user_id integer not null,
  id serial not null,
  name character varying (100),
  value character varying (100),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS picture (
  user_id integer not null,
  id serial not null,
  large character varying (100),
  medium character varying (100),
  thumbnail character varying (100),
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- another alternative---
-------------------------

-- WITH user_key AS
--         (INSERT INTO users (title, first_name, last_name, gender, email, phone, cell, nat) VALUES (
-- 			'tit1','Laurenz', 'lname', 'male', 'lname@example.com', '22222', '33333', 'AR') RETURNING id),
--        locations_key AS 
-- 	(INSERT INTO locations(user_id, city, state, country, postcode) VALUES (
-- 			(SELECT id FROM user_key), 'somecity', 'somestate', 'somecountry', 123) RETURNING id)
-- INSERT INTO streets (location_id, number, name) VALUES (
-- 	(SELECT id from locations_key), 45, 'streetname');

