CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(150) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    role        VARCHAR(20)  NOT NULL CHECK (role IN ('pengelola', 'pengurus', 'keluarga')),
    is_verified BOOLEAN      NOT NULL DEFAULT FALSE,
    panti_id    UUID,          -- FK ke tabel panti (nullable saat pertama register)
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS panti (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nama         VARCHAR(150) NOT NULL,
    alamat       TEXT,
    telepon      VARCHAR(20),
    kode_undangan VARCHAR(10) NOT NULL UNIQUE, -- kode unik untuk join panti
    pengelola_id UUID        NOT NULL,          -- owner/pengelola panti
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD CONSTRAINT fk_users_panti
    FOREIGN KEY (panti_id) REFERENCES panti(id) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS kamar (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    panti_id   UUID        NOT NULL REFERENCES panti(id) ON DELETE CASCADE,
    nama_kamar VARCHAR(50) NOT NULL,
    kapasitas  INT         NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS lansia (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    panti_id        UUID        NOT NULL REFERENCES panti(id) ON DELETE CASCADE,
    kamar_id        UUID        REFERENCES kamar(id) ON DELETE SET NULL,
    nama            VARCHAR(100) NOT NULL,
    nik             VARCHAR(16),
    tanggal_lahir   DATE,
    jenis_kelamin   VARCHAR(10) CHECK (jenis_kelamin IN ('L', 'P')),
    alamat_asal     TEXT,
    golongan_darah  VARCHAR(5),
    riwayat_penyakit TEXT,
    alergi          TEXT,
    foto_url        TEXT,
    tanggal_masuk   DATE        NOT NULL DEFAULT CURRENT_DATE,
    status          VARCHAR(20) NOT NULL DEFAULT 'aktif' CHECK (status IN ('aktif', 'keluar', 'meninggal')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS pengurus_profil (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    kamar_id   UUID REFERENCES kamar(id) ON DELETE SET NULL,
    jabatan    VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS shift (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nama_shift  VARCHAR(50) NOT NULL,   -- contoh: Pagi, Siang, Malam
    jam_mulai   TIME        NOT NULL,
    jam_selesai TIME        NOT NULL,
    hari        TEXT        NOT NULL,   -- contoh: "Senin,Selasa,Rabu"
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS aktivitas_lansia (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lansia_id   UUID        NOT NULL REFERENCES lansia(id) ON DELETE CASCADE,
    nama        VARCHAR(100) NOT NULL,   -- contoh: Senam Pagi
    deskripsi   TEXT,
    jam         TIME,
    hari        TEXT,                    -- contoh: "Senin,Rabu,Jumat"
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS jadwal_checkup (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lansia_id   UUID        NOT NULL REFERENCES lansia(id) ON DELETE CASCADE,
    tanggal     DATE        NOT NULL,
    keterangan  TEXT,
    status      VARCHAR(20) NOT NULL DEFAULT 'terjadwal' CHECK (status IN ('terjadwal', 'selesai', 'dibatalkan')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS obat (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lansia_id    UUID        NOT NULL REFERENCES lansia(id) ON DELETE CASCADE,
    nama_obat    VARCHAR(100) NOT NULL,
    dosis        VARCHAR(50),
    keterangan   TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS jadwal_obat (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    obat_id    UUID        NOT NULL REFERENCES obat(id) ON DELETE CASCADE,
    jam        TIME        NOT NULL,
    hari       TEXT        NOT NULL DEFAULT 'setiap hari',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS log_pemberian_obat (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    jadwal_obat_id UUID        NOT NULL REFERENCES jadwal_obat(id) ON DELETE CASCADE,
    pengurus_id    UUID        NOT NULL REFERENCES users(id),
    diberikan_pada TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    catatan        TEXT
);

CREATE TABLE IF NOT EXISTS catatan_shift (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lansia_id   UUID        NOT NULL REFERENCES lansia(id) ON DELETE CASCADE,
    pengurus_id UUID        NOT NULL REFERENCES users(id),
    isi_catatan TEXT        NOT NULL,
    shift       VARCHAR(20),             -- Pagi / Siang / Malam
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pemeriksaan_kesehatan (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lansia_id       UUID        NOT NULL REFERENCES lansia(id) ON DELETE CASCADE,
    pengurus_id     UUID        NOT NULL REFERENCES users(id),
    tekanan_darah   VARCHAR(20),   -- contoh: "120/80"
    gula_darah      DECIMAL(5,2),
    suhu_tubuh      DECIMAL(4,1),
    berat_badan     DECIMAL(5,2),
    keluhan         TEXT,
    rekomendasi     TEXT,          -- output dari AI / logic
    status_darurat  VARCHAR(20) DEFAULT 'normal' CHECK (status_darurat IN ('normal', 'perhatian', 'darurat')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS kunjungan_keluarga (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    lansia_id        UUID        NOT NULL REFERENCES lansia(id) ON DELETE CASCADE,
    pengurus_id      UUID        NOT NULL REFERENCES users(id),
    tanggal_kunjungan TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    foto_url         TEXT,
    catatan          TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS keluarga_lansia (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lansia_id  UUID NOT NULL REFERENCES lansia(id) ON DELETE CASCADE,
    hubungan   VARCHAR(50),   -- contoh: anak, cucu, keponakan
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, lansia_id)
);

CREATE INDEX IF NOT EXISTS idx_users_email        ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_panti        ON users(panti_id);
CREATE INDEX IF NOT EXISTS idx_lansia_panti       ON lansia(panti_id);
CREATE INDEX IF NOT EXISTS idx_catatan_lansia     ON catatan_shift(lansia_id);
CREATE INDEX IF NOT EXISTS idx_obat_lansia        ON obat(lansia_id);
CREATE INDEX IF NOT EXISTS idx_pemeriksaan_lansia ON pemeriksaan_kesehatan(lansia_id);
CREATE INDEX IF NOT EXISTS idx_kunjungan_lansia   ON kunjungan_keluarga(lansia_id);

