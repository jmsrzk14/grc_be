INSERT INTO tenants (id, name, type, status, created_at)
VALUES 
    ('550e8400-e29b-41d4-a716-446655440000', 'BPR Sentosa Jaya', 'BPR', 'Active', NOW()),
    ('550e8400-e29b-41d4-a716-446655440001', 'BPRS Amanah Ummah', 'BPRS', 'Active', NOW()),
    ('550e8400-e29b-41d4-a716-446655440002', 'BPR Maju Mapan', 'BPR', 'Active', NOW()),
    ('550e8400-e29b-41d4-a716-446655440003', 'BPR Sejahtera Mandiri', 'BPR', 'Inactive', NOW()),
    ('550e8400-e29b-41d4-a716-446655440004', 'BPRS Barokah Utama', 'BPRS', 'Active', NOW()),
    ('550e8400-e29b-41d4-a716-446655440005', 'BPR Kencana Artha', 'BPR', 'Active', NOW()),
    ('550e8400-e29b-41d4-a716-446655440006', 'BPR Cahaya Abadi', 'BPR', 'Active', NOW()),
    ('550e8400-e29b-41d4-a716-446655440007', 'BPR Mulia Pratama', 'BPR', 'Active', NOW()),
    ('550e8400-e29b-41d4-a716-446655440008', 'BPRS Syariah Madani', 'BPRS', 'Active', NOW()),
    ('550e8400-e29b-41d4-a716-446655440009', 'BPR Dana Swadaya', 'BPR', 'Active', NOW())
ON CONFLICT (id) DO NOTHING;
