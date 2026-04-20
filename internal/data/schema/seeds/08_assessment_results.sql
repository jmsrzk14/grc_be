INSERT INTO assessment_results (id, session_id, regulation_item_id, compliance_status, evidence_link, remarks, updated_at)
VALUES 
    ('110e8400-e29b-41d4-a716-446655440000', '880e8400-e29b-41d4-a716-446655440000', '330e8400-e29b-41d4-a716-446655440000', 'YES', 'https://docs.bpr.com/ev1', 'Bukti lengkap', NOW()),
    ('110e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440001', '330e8400-e29b-41d4-a716-446655440001', 'NO', 'https://docs.bpr.com/ev2', 'Masih dalam proses perbaikan', NOW()),
    ('110e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440003', '330e8400-e29b-41d4-a716-446655440002', 'YES', 'https://docs.bpr.com/ev3', 'Sudah terenkripsi', NOW()),
    ('110e8400-e29b-41d4-a716-446655440003', '880e8400-e29b-41d4-a716-446655440004', '330e8400-e29b-41d4-a716-446655440003', 'N/A', '', 'Tidak relevan untuk lingkup ini', NOW()),
    ('110e8400-e29b-41d4-a716-446655440004', '880e8400-e29b-41d4-a716-446655440006', '330e8400-e29b-41d4-a716-446655440004', 'YES', 'https://docs.bpr.com/ev4', 'Dokumen manajemen risiko ada', NOW()),
    ('110e8400-e29b-41d4-a716-446655440005', '880e8400-e29b-41d4-a716-446655440007', '330e8400-e29b-41d4-a716-446655440005', 'YES', 'https://docs.bpr.com/ev5', 'Sistem tersedia 99.9%', NOW()),
    ('110e8400-e29b-41d4-a716-446655440006', '880e8400-e29b-41d4-a716-446655440008', '330e8400-e29b-41d4-a716-446655440006', 'NO', '', 'Belum melaporkan transaksi', NOW()),
    ('110e8400-e29b-41d4-a716-446655440007', '880e8400-e29b-41d4-a716-446655440000', '330e8400-e29b-41d4-a716-446655440007', 'YES', 'https://docs.bpr.com/ev6', 'Produk DN 100%', NOW()),
    ('110e8400-e29b-41d4-a716-446655440008', '880e8400-e29b-41d4-a716-446655440001', '330e8400-e29b-41d4-a716-446655440008', 'YES', 'https://docs.bpr.com/ev7', 'Sudah daftar PSE', NOW()),
    ('110e8400-e29b-41d4-a716-446655440009', '880e8400-e29b-41d4-a716-446655440003', '330e8400-e29b-41d4-a716-446655440009', 'YES', 'https://docs.bpr.com/ev8', 'Edukasi rutin bulanan', NOW())
ON CONFLICT (id) DO NOTHING;
