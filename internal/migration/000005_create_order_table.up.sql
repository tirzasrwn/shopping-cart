CREATE TABLE "order"(
  id serial PRIMARY KEY,
  cart_id INT,
  product_id INT,
  quantity FLOAT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_cart FOREIGN KEY(cart_id) REFERENCES "cart"(id),
  CONSTRAINT fk_product FOREIGN KEY(product_id) REFERENCES "product"(id) 
);

