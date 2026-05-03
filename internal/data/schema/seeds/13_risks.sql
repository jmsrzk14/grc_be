INSERT INTO risks (id, risk_title, risk_description, category_id,
  likelihood_inherent, impact_inherent,  -- score_inherent = L*I
  likelihood_residual, impact_residual,  -- score_residual = L*I
  mitigation_plan, mitigation_status) VALUES

  ('cc000001-0000-0000-0000-000000000001',
   'Risiko Kesenjangan Pendanaan Jangka Pendek',
   'Ketidaksesuaian antara jatuh tempo aset dan kewajiban jangka pendek yang dapat menyebabkan kesulitan likuiditas.',
   'bb000001-0000-0000-0000-000000000001', 4, 5, 2, 3,
   'Menetapkan Liquidity Coverage Ratio (LCR) minimum 100%, mempertahankan aset likuid berkualitas tinggi, dan membuat fasilitas repo dengan bank koresponden.',
   'sedang direncanakan'),

  ('cc000001-0000-0000-0000-000000000002',
   'Risiko Kegagalan Sistem Core Banking',
   'Potensi gangguan atau kegagalan sistem core banking yang berdampak pada operasional layanan perbankan kepada nasabah.',
   'bb000001-0000-0000-0000-000000000002', 3, 4, 1, 2,
   'Implementasi sistem disaster recovery dengan RPO < 1 jam dan RTO < 4 jam. Uji coba failover setiap 6 bulan. Kontrak SLA dengan vendor teknologi.',
   'sudah diimplementasikan'),

  ('cc000001-0000-0000-0000-000000000003',
   'Risiko Pelanggaran Regulasi OJK',
   'Risiko sanksi akibat ketidakpatuhan terhadap POJK, SEOJK, dan peraturan OJK lainnya yang berlaku bagi BPR.',
   'bb000001-0000-0000-0000-000000000003', 2, 5, 1, 3,
   'Membentuk unit kepatuhan internal, melakukan compliance review bulanan, mengikuti sosialisasi regulasi OJK, dan melaksanakan audit kepatuhan tahunan.',
   'sudah diimplementasikan'),

  ('cc000001-0000-0000-0000-000000000004',
   'Risiko Kebocoran Data Nasabah',
   'Risiko tereksposnya data pribadi nasabah akibat serangan siber atau kelalaian internal yang dapat merusak kepercayaan publik.',
   'bb000001-0000-0000-0000-000000000004', 3, 3, 2, 2,
   'Implementasi enkripsi end-to-end, pelatihan kesadaran keamanan data untuk seluruh pegawai, penerapan kebijakan akses data berbasis kebutuhan (need-to-know).',
   'sedang direncanakan'),

  ('cc000001-0000-0000-0000-000000000005',
   'Risiko Persaingan Fintech dan Digital Banking',
   'Tekanan kompetitif dari layanan keuangan digital yang dapat mengikis pangsa pasar dan menurunkan relevansi BPR.',
   'bb000001-0000-0000-0000-000000000005', 4, 4, 3, 3,
   'Mengembangkan layanan digital BPR (mobile banking), bermitra dengan fintech terpilih, dan berfokus pada segmen pasar lokal yang underserved.',
   'sedang direncanakan'),

  ('cc000001-0000-0000-0000-000000000006',
   'Risiko Fraud Internal',
   'Risiko tindak kecurangan yang dilakukan oleh pegawai internal seperti penggelapan dana, manipulasi data, atau penyalahgunaan wewenang.',
   'bb000001-0000-0000-0000-000000000002', 5, 3, 3, 2,
   'Penerapan four-eyes principle untuk semua transaksi di atas batas nilai tertentu, rotasi jabatan periodik, dan pelaporan anonim melalui whistleblowing system.',
   'sudah diimplementasikan'),

  ('cc000001-0000-0000-0000-000000000007',
   'Risiko Penarikan Dana Massal (Bank Run)',
   'Risiko penarikan dana secara besar-besaran oleh nasabah secara bersamaan yang berpotensi mengancam solvabilitas BPR.',
   'bb000001-0000-0000-0000-000000000001', 2, 4, 1, 2,
   'Mempertahankan rasio kecukupan modal (CAR) di atas batas minimum OJK, menyusun contingency funding plan, dan meningkatkan komunikasi proaktif kepada nasabah.',
   'sudah diimplementasikan'),

  ('cc000001-0000-0000-0000-000000000008',
   'Risiko Ketidakpatuhan APU-PPT',
   'Risiko kegagalan dalam penerapan program Anti Pencucian Uang dan Pencegahan Pendanaan Terorisme sesuai POJK.',
   'bb000001-0000-0000-0000-000000000003', 3, 4, 2, 3,
   'Mengimplementasikan sistem transaction monitoring otomatis, memperbarui database PEP/blacklist bulanan, dan melatih seluruh front-liner untuk Customer Due Diligence.',
   'sedang direncanakan')

ON CONFLICT (id) DO NOTHING;
