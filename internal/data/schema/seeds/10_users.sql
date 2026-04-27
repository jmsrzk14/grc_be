INSERT INTO users (id, username, password, email, full_name, tenant_id, role, created_at, updated_at)
VALUES 
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin@grc.com', 'System Administrator', '550e8400-e29b-41d4-a716-446655440000', 'SuperAdmin', NOW(), NOW()),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'user_sentosa', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'sentosa@grc.com', 'Sentosa User', '550e8400-e29b-41d4-a716-446655440000', 'User', NOW(), NOW()),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'user_amanah', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'amanah@grc.com', 'Amanah User', '550e8400-e29b-41d4-a716-446655440001', 'User', NOW(), NOW()),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'user_maju', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'maju@grc.com', 'Maju User', '550e8400-e29b-41d4-a716-446655440002', 'User', NOW(), NOW()),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'user_kencana', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'kencana@grc.com', 'Kencana User', '550e8400-e29b-41d4-a716-446655440005', 'User', NOW(), NOW())
ON CONFLICT (username) DO NOTHING;
