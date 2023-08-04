create table "cart" (
  id serial PRIMARY KEY,
  user_id int,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id) 
);

