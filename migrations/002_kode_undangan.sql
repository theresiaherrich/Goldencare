CREATE TABLE IF NOT EXISTS kode_undangan (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kode         VARCHAR(10)  NOT NULL UNIQUE,
    untuk_role   VARCHAR(20)  NOT NULL CHECK (untuk_role IN ('pengelola', 'pengurus', 'keluarga')),
    dibuat_oleh  UUID         REFERENCES users(id) ON DELETE SET NULL,
    panti_id     UUID         REFERENCES panti(id) ON DELETE CASCADE,
    tipe         VARCHAR(20)  NOT NULL DEFAULT 'single_use'
                              CHECK (tipe IN ('single_use', 'multi_use')),
    dipakai_count INT         NOT NULL DEFAULT 0,
    maks_pakai   INT,
    expired_at   TIMESTAMPTZ,
    is_aktif     BOOLEAN      NOT NULL DEFAULT TRUE,
    catatan      TEXT,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS log_pemakaian_kode (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kode_undangan_id UUID        NOT NULL REFERENCES kode_undangan(id),
    user_id          UUID        NOT NULL REFERENCES users(id),
    dipakai_pada     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS superadmin (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email      VARCHAR(150) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    name       VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD COLUMN IF NOT EXISTS
    kode_dipakai VARCHAR(10); 

CREATE INDEX IF NOT EXISTS idx_kode_undangan_kode     ON kode_undangan(kode);
CREATE INDEX IF NOT EXISTS idx_kode_undangan_panti    ON kode_undangan(panti_id);
CREATE INDEX IF NOT EXISTS idx_kode_undangan_role     ON kode_undangan(untuk_role);
CREATE INDEX IF NOT EXISTS idx_log_kode_user          ON log_pemakaian_kode(user_id);

INSERT INTO superadmin (id, email, password, name)
VALUES (
    uuid_generate_v4(),
    'superadmin@rawatkasih.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: "password"
    'Super Admin Rawat Kasih'
) ON CONFLICT (email) DO NOTHING;
