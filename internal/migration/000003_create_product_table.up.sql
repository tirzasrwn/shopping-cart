create table "product" (
  id serial PRIMARY KEY,
  category_id int,
  name VARCHAR(255) NOT NULL,
  description VARCHAR(255),
  price FLOAT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES category(id) ON DELETE CASCADE
);
