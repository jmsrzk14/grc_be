INSERT INTO regulation_assesments (id, regulation_id, session_id, amount_pass, amount_fail, amount_na) VALUES
  (gen_random_uuid(),'770e8400-e29b-41d4-a716-446655440000','880e8400-e29b-41d4-a716-446655440000', 2, 1, 0),
  (gen_random_uuid(),'770e8400-e29b-41d4-a716-446655440001','880e8400-e29b-41d4-a716-446655440001', 2, 1, 0),
  (gen_random_uuid(),'770e8400-e29b-41d4-a716-446655440002','880e8400-e29b-41d4-a716-446655440003', 1, 1, 1),
  (gen_random_uuid(),'770e8400-e29b-41d4-a716-446655440003','880e8400-e29b-41d4-a716-446655440004', 2, 0, 1),
  (gen_random_uuid(),'770e8400-e29b-41d4-a716-446655440004','880e8400-e29b-41d4-a716-446655440006', 3, 0, 0),
  (gen_random_uuid(),'770e8400-e29b-41d4-a716-446655440005','880e8400-e29b-41d4-a716-446655440007', 2, 1, 0),
  (gen_random_uuid(),'770e8400-e29b-41d4-a716-446655440009','880e8400-e29b-41d4-a716-446655440008', 1, 0, 1)
ON CONFLICT (regulation_id, session_id) DO NOTHING;
