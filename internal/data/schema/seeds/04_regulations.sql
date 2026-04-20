INSERT INTO regulations (id, title, regulation_type, issued_date, status)
VALUES 
    ('770e8400-e29b-41d4-a716-446655440000', 'POJK No 11/2022 tentang Penyelenggaraan Teknologi Informasi', 'POJK', '2022-07-01', 'Active'),
    ('770e8400-e29b-41d4-a716-446655440001', 'SEOJK tentang Standar Pengamanan Data', 'SEOJK', '2023-01-15', 'Active'),
    ('770e8400-e29b-41d4-a716-446655440002', 'UU No 27/2022 tentang Perlindungan Data Pribadi', 'UU', '2022-10-17', 'Active'),
    ('770e8400-e29b-41d4-a716-446655440003', 'POJK No 4/2021 tentang Manajemen Risiko', 'POJK', '2021-03-12', 'Active'),
    ('770e8400-e29b-41d4-a716-446655440004', 'PBI No 23/2021 tentang SPW', 'PBI', '2021-06-30', 'Active'),
    ('770e8400-e29b-41d4-a716-446655440005', 'Peraturan OJK No 13/2020', 'POJK', '2020-05-20', 'Revoked'),
    ('770e8400-e29b-41d4-a716-446655440006', 'Instruksi Presiden No 2/2022', 'INPRES', '2022-03-30', 'Active'),
    ('770e8400-e29b-41d4-a716-446655440007', 'Peraturan Menteri Kominfo No 5/2020', 'PERMEN', '2020-11-24', 'Active'),
    ('770e8400-e29b-41d4-a716-446655440008', 'SEOJK tentang Literasi Keuangan', 'SEOJK', '2023-04-10', 'Active'),
    ('770e8400-e29b-41d4-a716-446655440009', 'UU Perbankan No 10/1998', 'UU', '1998-11-10', 'Active')
ON CONFLICT (id) DO NOTHING;
