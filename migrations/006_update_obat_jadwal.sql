ALTER TABLE obat
ADD COLUMN IF NOT EXISTS cara_pemberian VARCHAR(50) DEFAULT 'Oral'
    CHECK (cara_pemberian IN ('Oral', 'Inhaler', 'Injeksi', 'Tetes', 'Topikal'));

ALTER TABLE obat
ADD COLUMN IF NOT EXISTS is_aktif BOOLEAN DEFAULT true;

ALTER TABLE jadwal_obat
ADD COLUMN IF NOT EXISTS shift VARCHAR(20) DEFAULT 'Pagi'
    CHECK (shift IN ('Pagi', 'Siang', 'Sore'));

ALTER TABLE jadwal_obat
ADD COLUMN IF NOT EXISTS jam VARCHAR(5); 

CREATE INDEX IF NOT EXISTS idx_jadwal_obat_obat_id
    ON jadwal_obat (obat_id);

CREATE INDEX IF NOT EXISTS idx_log_pemberian_jadwal_tanggal
    ON log_pemberian_obat (jadwal_obat_id, diberikan_pada DESC);

CREATE INDEX IF NOT EXISTS idx_obat_lansia_aktif
    ON obat (lansia_id, is_aktif);
