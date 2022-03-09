CREATE TABLE IF NOT EXISTS user_characters(
   id serial PRIMARY KEY,
   user_id int NOT NULL,
   character_id int NOT NULL,
   character_rank varchar (255) NOT NULL,
   name varchar (255) NOT NULL,
   level int DEFAULT 1
);