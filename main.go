package main

import (
	"auto_job_applicator/db"
	"auto_job_applicator/mailer"
	"auto_job_applicator/scraper"
	"log"
	"time"
)

const cvPath = "cv_dosyasi.pdf" //gerçek bir dosya ile değiştirmeniz gerekmektedird

func main() {
	// Veritabanı başlatılıyor
	database, err := db.InitDB("auto_job_applicator.db")
	if err != nil {
		log.Fatal("Veritabanı başlatılamadı:", err)
	}
	defer database.Close()

	// SMTP Konfigürasyonu
	smtpConfig := mailer.MailConfig{
		SMTPHost:    "smtp.gmail.com",
		SMTPPort:    "587",
		SenderEmail: "ornek@gmail.com",
		SenderPass:  "ornek_sifre",
	}

	// Arama querysi
	query := "istanbul Yazılım Şirketleri"
	log.Printf("Google araması başlatılıyor: %s\n", query)

	// Şirket URL'lerini Google ve Public API'lerden çekiyoruz
	companyURLs := scraper.SearchCompanies(query)

	for _, companyURL := range companyURLs {
		log.Printf("E-posta aranıyor: %s\n", companyURL)
		emails := scraper.ExtractEmails(companyURL)

		for _, email := range emails {
			// Veritabanına kaydetme işlemi
			err := db.InsertCompany(database, companyURL, email)
			if err != nil {
				log.Printf("Veritabanına kaydedilemedi: %s - Hata: %v", email, err)
			}
		}
	}

	log.Println("Tüm şirketler eklendi, gönderilmeyi bekleyen mailler kontrol ediliyor...")

	// Veritabanından gönderilmemiş e-postaları alıyoruz
	pendingRecords, err := db.GetPendingEmails(database)
	if err != nil {
		log.Fatal("Bekleyen mailler alınamadı:", err)
	}

	// Her birini SMTP ile gönder
	for _, record := range pendingRecords {
		log.Printf("E-posta gönderiliyor: %s (%s)", record.Email, record.Company)
		err := mailer.SendEmail(smtpConfig, record, cvPath)
		if err != nil {
			log.Printf("Gönderim başarısız: %s", record.Email)
			db.MarkAsFailed(database, record.ID)
		} else {
			log.Printf("Gönderim başarılı: %s", record.Email)
			db.MarkAsSent(database, record.ID)
		}
		time.Sleep(2 * time.Second) // Anti-spam koruma
	}

	log.Println("İşlem tamamlandı.")
}
