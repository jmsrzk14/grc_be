INSERT INTO assessment_sessions (id, tenant_id, title, period_year, status, created_at)
VALUES 
    ('880e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 'Self Assessment TI 2026', 2026, 'In_Progress', NOW()),
    ('880e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'Audit Compliance Semester 1', 2026, 'Completed', NOW()),
    ('880e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'Mitigasi Risiko 2026', 2026, 'Draft', NOW()),
    ('880e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440004', 'Kesiapan Keamanan Data', 2026, 'In_Progress', NOW()),
    ('880e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440005', 'Sertifikasi ISO 27001', 2026, 'In_Progress', NOW()),
    ('880e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440006', 'Evaluasi Vendor TI', 2026, 'Draft', NOW()),
    ('880e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440007', 'Audit internal Tata Kelola', 2026, 'Completed', NOW()),
    ('880e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440008', 'Kepatuhan Syariah Q1', 2026, 'Completed', NOW()),
    ('880e8400-e29b-41d4-a716-446655440008', '550e8400-e29b-41d4-a716-446655440009', 'Assessment Infrastruktur', 2026, 'In_Progress', NOW()),
    ('880e8400-e29b-41d4-a716-446655440009', '550e8400-e29b-41d4-a716-446655440000', 'Review Kebijakan Privasi', 2026, 'Draft', NOW())
ON CONFLICT (id) DO NOTHING;
