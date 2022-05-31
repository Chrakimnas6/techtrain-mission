CREATE TABLE IF NOT EXISTS transactions(
   id serial PRIMARY KEY,
   tx varchar (255) NOT NULL,
   status varchar (255) NOT NULL
);