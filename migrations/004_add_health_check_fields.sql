ALTER TABLE pemeriksaan_kesehatan
ADD COLUMN IF NOT EXISTS detak_jantung INT;

ALTER TABLE pemeriksaan_kesehatan
ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'STABIL' CHECK (status IN ('STABIL', 'OBSERVASI', 'WASPADA', 'DARURAT'));

ALTER TABLE pemeriksaan_kesehatan
DROP CONSTRAINT IF EXISTS pemeriksaan_kesehatan_status_darurat_check;

ALTER TABLE pemeriksaan_kesehatan
ADD CONSTRAINT pemeriksaan_kesehatan_status_darurat_check 
CHECK (status_darurat IN ('hijau', 'kuning', 'oranye', 'merah', 'normal', 'perhatian', 'darurat'));


