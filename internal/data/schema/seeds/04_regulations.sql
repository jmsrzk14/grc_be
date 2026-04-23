INSERT INTO regulations (id, title, regulation_type, issued_date, status, category)
VALUES 
    ('770e8400-e29b-41d4-a716-446655440000', 'POJK No. 11/POJK.03/2022 tentang Penyelenggaraan Teknologi Informasi oleh Bank Umum', 'POJK', '2022-09-01', 'Active', 'External'),
    ('770e8400-e29b-41d4-a716-446655440001', 'SEOJK No. 14/SEOJK.03/2022 tentang Ketahanan dan Keamanan Siber bagi Bank Umum', 'SEOJK', '2022-12-15', 'Active', 'External'),
    ('770e8400-e29b-41d4-a716-446655440002', 'UU No. 27 Tahun 2022 tentang Pelindungan Data Pribadi', 'UU', '2022-10-17', 'Active', 'External'),
    ('770e8400-e29b-41d4-a716-446655440003', 'POJK No. 18/POJK.03/2016 tentang Penerapan Manajemen Risiko bagi Bank Umum', 'POJK', '2016-03-16', 'Active', 'External'),
    ('770e8400-e29b-41d4-a716-446655440004', 'PBI No. 23/6/PBI/2021 tentang Penyedia Jasa Pembayaran', 'PBI', '2021-07-01', 'Active', 'External'),
    ('770e8400-e29b-41d4-a716-446655440005', 'POJK No. 13/POJK.03/2021 tentang Penyelenggaraan Produk Bank Umum', 'POJK', '2021-08-30', 'Active', 'External'),
    ('770e8400-e29b-41d4-a716-446655440006', 'INPRES No. 2 Tahun 2022 tentang Percepatan Peningkatan Penggunaan Produk Dalam Negeri', 'INPRES', '2022-03-30', 'Active', 'External'),
    ('770e8400-e29b-41d4-a716-446655440007', 'PERMEN Kominfo No. 5 Tahun 2020 tentang Penyelenggara Sistem Elektronik Lingkup Privat', 'PERMEN', '2020-11-24', 'Active', 'External'),
    ('770e8400-e29b-41d4-a716-446655440008', 'SEOJK No. 1/SEOJK.07/2023 tentang Tata Cara Pelaksanaan Edukasi Keuangan', 'SEOJK', '2023-01-31', 'Active', 'External'),
    ('770e8400-e29b-41d4-a716-446655440009', 'Kebijakan Keamanan Informasi Internal Dimensi Kreasi Nusantara', 'INTERNAL', '2024-01-01', 'Active', 'Internal')
ON CONFLICT (id) DO NOTHING;
