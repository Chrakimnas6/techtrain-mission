CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   name varchar (255) NOT NULL,
   token varchar (255) NOT NULL
);