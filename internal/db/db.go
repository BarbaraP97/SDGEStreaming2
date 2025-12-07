// internal/db/db.go
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error al abrir la base de datos: %w", err)
	}
	
	DB.SetMaxOpenConns(1)

	// Ejecutar el schema SQL
	schema := `
CREATE TABLE IF NOT EXISTS plans (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    price REAL NOT NULL,
    max_quality TEXT DEFAULT 'HD',          
    max_devices INTEGER DEFAULT 1
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    age INTEGER NOT NULL,
    plan_id INTEGER NOT NULL DEFAULT 1,
    age_rating TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_login DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_id) REFERENCES plans(id)
);

CREATE TABLE IF NOT EXISTS payment_methods (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    card_holder_name TEXT NOT NULL,
    card_number TEXT NOT NULL,
    expiration_date TEXT NOT NULL,
    cvv TEXT NOT NULL,
    card_number_last4 TEXT NOT NULL,
    expiry_month INTEGER NOT NULL,
    expiry_year INTEGER NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS audiovisual_content (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    type TEXT NOT NULL,
    genre TEXT NOT NULL,
    duration INTEGER NOT NULL,
    age_rating TEXT NOT NULL,
    synopsis TEXT,
    release_year INTEGER,
    director TEXT,
    average_rating REAL DEFAULT 0.0,
    is_available BOOLEAN NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS audio_content (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    type TEXT NOT NULL,
    genre TEXT NOT NULL,
    duration INTEGER NOT NULL,
    age_rating TEXT NOT NULL,
    artist TEXT,
    album TEXT,
    track_number INTEGER,
    average_rating REAL DEFAULT 0.0,
    is_available BOOLEAN NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS user_ratings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content_id INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    rating REAL NOT NULL,
    rated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, content_id, content_type)
);

CREATE TABLE IF NOT EXISTS playback_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content_id INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    progress_seconds INTEGER DEFAULT 0,
    watched_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS favorites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    content_id INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, content_id, content_type)
);
`
	_, err = DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("error en la migración: %w", err)
	}

	// Insertar planes
	_, err = DB.Exec(`
INSERT OR IGNORE INTO plans (id, name, price, max_quality, max_devices)
VALUES
    (1, 'Free', 0.0, 'SD', 1),
    (2, 'Estándar', 9.99, 'HD', 2),
    (3, 'Premium 4K', 15.99, '4K', 4);
`)
	if err != nil {
		return fmt.Errorf("error al insertar planes: %w", err)
	}

	return nil
}

func GetDB() *sql.DB {
	return DB
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
