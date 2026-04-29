ALTER TABLE catatan_shift
ADD COLUMN IF NOT EXISTS suasana_hati VARCHAR(20)
    CHECK (suasana_hati IN ('Tenang', 'Gelisah', 'Lesu'));

ALTER TABLE catatan_shift
ADD COLUMN IF NOT EXISTS nafsu_makan VARCHAR(20)
    CHECK (nafsu_makan IN ('Baik', 'Sedang', 'Buruk'));

ALTER TABLE catatan_shift
ADD COLUMN IF NOT EXISTS aktivitas VARCHAR(20)
    CHECK (aktivitas IN ('Aktif', 'Istirahat'));

ALTER TABLE catatan_shift
ADD COLUMN IF NOT EXISTS status_jurnal VARCHAR(20)
    DEFAULT 'draf'
    CHECK (status_jurnal IN ('draf', 'terkirim'));

ALTER TABLE catatan_shift
DROP CONSTRAINT IF EXISTS catatan_shift_shift_check;

ALTER TABLE catatan_shift
ADD CONSTRAINT catatan_shift_shift_check
    CHECK (shift IN ('Pagi', 'Siang', 'Malam'));

CREATE INDEX IF NOT EXISTS idx_catatan_shift_lansia_created
    ON catatan_shift (lansia_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_catatan_shift_status
    ON catatan_shift (status_jurnal);
