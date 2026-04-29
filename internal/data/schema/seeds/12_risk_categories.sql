INSERT INTO risk_categories (id, title, appetite, tolerance)
VALUES 
    (gen_random_uuid(), 'Likuiditas', '70', '10'),
    (gen_random_uuid(), 'Operasional', '65', '5'),
    (gen_random_uuid(), 'Compliance', '55', '10'),
    (gen_random_uuid(), 'Reputasi', '60', '7'),
    (gen_random_uuid(), 'Strategis', '65', '8')
ON CONFLICT (title) DO NOTHING;
