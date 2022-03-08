CREATE TABLE IF NOT EXISTS user_characters(
   id serial PRIMARY KEY,
   user_id int,
   character_id int,
   character_rank VARCHAR (50),
   name VARCHAR (50),
   level int DEFAULT 1
);