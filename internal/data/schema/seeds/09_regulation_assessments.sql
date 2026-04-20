INSERT INTO regulation_assesments (id, regulation_id, session_id, amount_pass, amount_fail, amount_na)
VALUES 
    (gen_random_uuid(), '770e8400-e29b-41d4-a716-446655440000'::uuid, '880e8400-e29b-41d4-a716-446655440000'::uuid, 5, 0, 0),
    (gen_random_uuid(), '770e8400-e29b-41d4-a716-446655440001'::uuid, '880e8400-e29b-41d4-a716-446655440001'::uuid, 3, 2, 0),
    (gen_random_uuid(), '770e8400-e29b-41d4-a716-446655440002'::uuid, '880e8400-e29b-41d4-a716-446655440003'::uuid, 10, 0, 1),
    (gen_random_uuid(), '770e8400-e29b-41d4-a716-446655440003'::uuid, '880e8400-e29b-41d4-a716-446655440004'::uuid, 4, 1, 0),
    (gen_random_uuid(), '770e8400-e29b-41d4-a716-446655440004'::uuid, '880e8400-e29b-41d4-a716-446655440006'::uuid, 7, 0, 0),
    (gen_random_uuid(), '770e8400-e29b-41d4-a716-446655440005'::uuid, '880e8400-e29b-41d4-a716-446655440007'::uuid, 2, 5, 0),
    (gen_random_uuid(), '770e8400-e29b-41d4-a716-446655440006'::uuid, '880e8400-e29b-41d4-a716-446655440008'::uuid, 9, 0, 0),
    (gen_random_uuid(), '770e8400-e29b-41d4-a716-446655440007'::uuid, '880e8400-e29b-41d4-a716-446655440000'::uuid, 6, 2, 1),
    (gen_random_uuid(), '770e8400-e29b-41d4-a716-446655440008'::uuid, '880e8400-e29b-41d4-a716-446655440001'::uuid, 1, 0, 8),
    (gen_random_uuid(), '770e8400-e29b-41d4-a716-446655440009'::uuid, '880e8400-e29b-41d4-a716-446655440003'::uuid, 12, 1, 0)
ON CONFLICT (regulation_id, session_id) DO NOTHING;
