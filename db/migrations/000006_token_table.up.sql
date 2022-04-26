CREATE TABLE IF NOT EXISTS tokens(
   id serial PRIMARY KEY,
   symbol varchar (255) NOT NULL,
   address varchar (255) NOT NULL
);