INSERT INTO regulation_items (id, regulation_id, reference_number, content)
VALUES 
    ('330e8400-e29b-41d4-a716-446655440000', '770e8400-e29b-41d4-a716-446655440000', 'Pasal 1 ayat 1', 'BPR wajib memiliki rencana strategis teknologi informasi.'),
    ('330e8400-e29b-41d4-a716-446655440001', '770e8400-e29b-41d4-a716-446655440000', 'Pasal 5 ayat 2', 'Infrastruktur TI harus memiliki backup data secara periodik.'),
    ('330e8400-e29b-41d4-a716-446655440002', '770e8400-e29b-41d4-a716-446655440001', 'Pasal 2', 'Seluruh data nasabah wajib dienkripsi saat dikirimkan.'),
    ('330e8400-e29b-41d4-a716-446655440003', '770e8400-e29b-41d4-a716-446655440002', 'Pasal 12', 'Pengendali Data Pribadi wajib menjamin keamanan data pribadi.'),
    ('330e8400-e29b-41d4-a716-446655440004', '770e8400-e29b-41d4-a716-446655440003', 'Pasal 3', 'Bank wajib menerapkan manajemen risiko secara efektif.'),
    ('330e8400-e29b-41d4-a716-446655440005', '770e8400-e29b-41d4-a716-446655440004', 'Pasal 8', 'Penyelenggara jasa pembayaran wajib menjaga ketersediaan sistem.'),
    ('330e8400-e29b-41d4-a716-446655440006', '770e8400-e29b-41d4-a716-446655440005', 'Pasal 1', 'Lembaga keuangan wajib melaporkan transaksi keuangan mencurigakan.'),
    ('330e8400-e29b-41d4-a716-446655440007', '770e8400-e29b-41d4-a716-446655440006', 'Pasal 10', 'Pengadaan barang/jasa pemerintah wajib menggunakan produk dalam negeri.'),
    ('330e8400-e29b-41d4-a716-446655440008', '770e8400-e29b-41d4-a716-446655440007', 'Pasal 4', 'PSE wajib melakukan pendaftaran pada sistem kementerian.'),
    ('330e8400-e29b-41d4-a716-446655440009', '770e8400-e29b-41d4-a716-446655440008', 'Pasal 6', 'Edukasi keuangan harus dilakukan secara inklusif.')
ON CONFLICT (id) DO NOTHING;
