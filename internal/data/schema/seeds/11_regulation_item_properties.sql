INSERT INTO regulation_item_properties (regulation_item_id, property_id)
VALUES 
    -- Item 0: BPR wajib memiliki rencana strategis teknologi informasi. -> Aset Teknologi
    ('330e8400-e29b-41d4-a716-446655440000', '660e8400-e29b-41d4-a716-446655440000'),
    
    -- Item 1: Infrastruktur TI harus memiliki backup data secara periodik. -> Aset Teknologi
    ('330e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440000'),
    
    -- Item 2: Seluruh data nasabah wajib dienkripsi saat dikirimkan. -> Data Nasabah
    ('330e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440005'),
    
    -- Item 3: Pengendali Data Pribadi wajib menjamin keamanan data pribadi. -> Data Nasabah
    ('330e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440005'),
    
    -- Item 4: Bank wajib menerapkan manajemen risiko secara efektif. -> Operasional
    ('330e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440006'),
    
    -- Item 5: Penyelenggara jasa pembayaran wajib menjaga ketersediaan sistem. -> Aset Teknologi
    ('330e8400-e29b-41d4-a716-446655440005', '660e8400-e29b-41d4-a716-446655440000')
ON CONFLICT (regulation_item_id, property_id) DO NOTHING;
