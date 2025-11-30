-- Seed products for development/testing
-- Idempotent using INSERT IGNORE (works when primary key already exists)

INSERT IGNORE INTO products (id, name, price, category) VALUES
  ('3f6b5b2a-7f66-4b3f-9a1b-000000000000', 'Chicken Waffle', 190, 'Waffle'),
  ('3f6b5b2a-7f66-4b3f-9a1b-111111111111', 'Pizza Margherita', 150, 'pizza'),
  ('e2c1a9b0-1d2e-4c3a-9f44-222222222222', 'Spaghetti Carbonara', 180, 'pasta'),
  ('7a9d8c3f-3333-4444-5555-333333333333', 'Caesar Salad', 90.50, 'salad');
