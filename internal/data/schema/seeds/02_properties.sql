INSERT INTO properties (id, name, description) VALUES
  ('660e8400-e29b-41d4-a716-446655440000','Aset Teknologi',  'Infrastruktur IT, Server, dan Aplikasi'),
  ('660e8400-e29b-41d4-a716-446655440001','SDM',             'Sumber Daya Manusia dan Kompetensi'),
  ('660e8400-e29b-41d4-a716-446655440002','Fisik',           'Gedung, Kantor, dan Dokumen Fisik'),
  ('660e8400-e29b-41d4-a716-446655440003','Keuangan',        'Aset Keuangan, Kas, dan Investasi'),
  ('660e8400-e29b-41d4-a716-446655440004','Reputasi',        'Brand, Kepercayaan Publik, dan Goodwill'),
  ('660e8400-e29b-41d4-a716-446655440005','Data Nasabah',    'Database Nasabah, PII, dan Rekam Medis'),
  ('660e8400-e29b-41d4-a716-446655440006','Operasional',     'Proses Bisnis, Flowchart, dan SOP'),
  ('660e8400-e29b-41d4-a716-446655440007','Compliance',      'Dokumen Legal, Izin Usaha, dan Sertifikasi'),
  ('660e8400-e29b-41d4-a716-446655440008','Vendor',          'Kontrak Pihak Ketiga dan Outsourcing'),
  ('660e8400-e29b-41d4-a716-446655440009','Inovasi',         'Hak Kekayaan Intelektual dan R&D')
ON CONFLICT (id) DO NOTHING;
