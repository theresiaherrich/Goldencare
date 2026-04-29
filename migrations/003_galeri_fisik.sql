CREATE TABLE IF NOT EXISTS galeri_fisik (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lansia_id        UUID        NOT NULL REFERENCES lansia(id) ON DELETE CASCADE,
    pengurus_id      UUID        NOT NULL REFERENCES users(id),
    foto_url         TEXT        NOT NULL,
    lokasi_luka      VARCHAR(100),
    deskripsi        TEXT,
    analisis_ai      TEXT,
    tingkat_darurat  VARCHAR(20) DEFAULT 'ringan' CHECK (tingkat_darurat IN ('ringan', 'sedang', 'berat', 'kritis')),
    prediksi_penyakit TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_galeri_lansia ON galeri_fisik(lansia_id);
