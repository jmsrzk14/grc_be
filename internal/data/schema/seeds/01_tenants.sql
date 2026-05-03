INSERT INTO tenants (id, name, type, status, created_at) VALUES
  ('550e8400-e29b-41d4-a716-446655440000','BPR Sentosa Jaya',      'BPR', 'Active',   '2024-01-15 08:00:00+07'),
  ('550e8400-e29b-41d4-a716-446655440001','BPRS Amanah Ummah',     'BPRS','Active',   '2024-02-01 08:00:00+07'),
  ('550e8400-e29b-41d4-a716-446655440002','BPR Maju Mapan',        'BPR', 'Active',   '2024-02-20 08:00:00+07'),
  ('550e8400-e29b-41d4-a716-446655440003','BPR Sejahtera Mandiri', 'BPR', 'Inactive', '2024-03-05 08:00:00+07'),
  ('550e8400-e29b-41d4-a716-446655440004','BPRS Barokah Utama',    'BPRS','Active',   '2024-03-15 08:00:00+07'),
  ('550e8400-e29b-41d4-a716-446655440005','BPR Kencana Artha',     'BPR', 'Active',   '2024-04-01 08:00:00+07'),
  ('550e8400-e29b-41d4-a716-446655440006','BPR Cahaya Abadi',      'BPR', 'Active',   '2024-04-15 08:00:00+07'),
  ('550e8400-e29b-41d4-a716-446655440007','BPR Mulia Pratama',     'BPR', 'Active',   '2024-05-01 08:00:00+07'),
  ('550e8400-e29b-41d4-a716-446655440008','BPRS Syariah Madani',   'BPRS','Active',   '2024-05-10 08:00:00+07'),
  ('550e8400-e29b-41d4-a716-446655440009','BPR Dana Swadaya',      'BPR', 'Active',   '2024-06-01 08:00:00+07')
ON CONFLICT (id) DO NOTHING;
