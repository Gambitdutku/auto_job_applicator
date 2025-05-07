package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Record struct {
	ID      int
	Company string
	Email   string
	Status  string
}

func InitDB(path string) (*sql.DB, error) {
	// SQLite veritabanına bağlanma
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS companies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		company TEXT,
		email TEXT UNIQUE,
		status TEXT DEFAULT 'pending'
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}

	log.Println("Veritabanı başlatıldı ve tablo kontrol edildi.")
	return db, nil
}

func InsertCompany(db *sql.DB, company, email string) error {
	_, err := db.Exec("INSERT OR IGNORE INTO companies (company, email) VALUES (?, ?)", company, email)
	if err != nil {
		log.Printf("Kayıt eklenemedi: %s - %s\nHata: %v", company, email, err)
	} else {
		log.Printf("Yeni kayıt eklendi: %s - %s", company, email)
	}
	return err
}

func GetPendingEmails(db *sql.DB) ([]Record, error) {
	rows, err := db.Query("SELECT id, company, email FROM companies WHERE status = 'pending'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.Company, &r.Email); err != nil {
			log.Println("Satır okunamadı:", err)
			continue
		}
		records = append(records, r)
	}
	return records, nil
}

func MarkAsSent(db *sql.DB, id int) {
	_, err := db.Exec("UPDATE companies SET status = 'sent' WHERE id = ?", id)
	if err != nil {
		log.Printf("Kayıt durumu güncellenemedi (sent): %d - Hata: %v", id, err)
	} else {
		log.Printf("Kayıt başarıyla gönderildi: ID=%d", id)
	}
}

func MarkAsFailed(db *sql.DB, id int) {
	_, err := db.Exec("UPDATE companies SET status = 'failed' WHERE id = ?", id)
	if err != nil {
		log.Printf("Kayıt durumu güncellenemedi (failed): %d - Hata: %v", id, err)
	} else {
		log.Printf("Kayıt gönderilemedi olarak işaretlendi: ID=%d", id)
	}
}
