-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabel users
CREATE TABLE IF NOT EXISTS users (
    user_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    nama_lengkap VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    nomor_telepon VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    status_akun VARCHAR(20) NOT NULL DEFAULT 'aktif',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Tabel terminals
CREATE TABLE IF NOT EXISTS terminals (
    terminal_id SERIAL PRIMARY KEY,
    nama_terminal VARCHAR(100) NOT NULL UNIQUE,
    lokasi TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Tabel gates
CREATE TABLE IF NOT EXISTS gates (
    gate_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    terminal_id INT NOT NULL REFERENCES terminals(terminal_id),
    nama_gate VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'online',
    last_seen TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Tabel cards
CREATE TABLE IF NOT EXISTS cards (
    card_id UUID PRIMARY KEY,
    saldo DECIMAL(10, 2) NOT NULL DEFAULT 0,
    status_kartu VARCHAR(20) NOT NULL DEFAULT 'aktif',
    registered_user_id UUID REFERENCES users(user_id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Tabel trips
CREATE TABLE IF NOT EXISTS trips (
    trip_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    card_id UUID NOT NULL REFERENCES cards(card_id),
    gate_checkin_id UUID NOT NULL REFERENCES gates(gate_id),
    waktu_checkin TIMESTAMPTZ NOT NULL,
    gate_checkout_id UUID REFERENCES gates(gate_id),
    waktu_checkout TIMESTAMPTZ,
    tarif DECIMAL(10, 2),
    status_perjalanan VARCHAR(20) NOT NULL DEFAULT 'berjalan',
    is_offline_tx BOOLEAN NOT NULL DEFAULT false
);

-- Tabel balance_transactions
CREATE TABLE IF NOT EXISTS balance_transactions (
    transaction_id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    card_id UUID NOT NULL REFERENCES cards(card_id),
    tipe_transaksi VARCHAR(20) NOT NULL,
    jumlah DECIMAL(10, 2) NOT NULL,
    keterangan TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Insert satu user dengan password: "password123"
-- Password hash dibuat dengan bcrypt dengan cost 10
INSERT INTO users (nama_lengkap, email, nomor_telepon, password_hash, status_akun)
VALUES (
    'Admin User',
    'admin@example.com',
    '081234567890',
    '$2a$10$wPKxd9Ow9djIxOsS4yYLJeDCegajG1fztmvmLghOx84XpLLFxlemK', -- password123
    'aktif'
); 