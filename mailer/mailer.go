package mailer

import (
	"auto_job_applicator/db"
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
)

// SMTP tanımlama
type MailConfig struct {
	SMTPHost    string
	SMTPPort    string
	SenderEmail string
	SenderPass  string
}

func SendEmail(config MailConfig, record db.Record, cvPath string) error {
	// Alıcı ve Gönderici adresleri
	to := []string{record.Email}
	from := config.SenderEmail

	// SMTP sunucu ayarları
	auth := smtp.PlainAuth("", config.SenderEmail, config.SenderPass, config.SMTPHost)

	// E-posta başlıkları
	subject := "İş Başvurusu: Yazılım Geliştirici"
	header := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nMIME-Version: 1.0\nContent-Type: multipart/mixed; boundary=boundary42\n\n", from, record.Email, subject)

	// Mesaj gövdesi
	body := "--boundary42\n" +
		"Content-Type: text/plain; charset=utf-8\n\n" +
		"Merhaba,\n\n" +
		"Ek'te yer alan CV'mi incelemenizi rica ederim.\n\n" +
		"Saygılarımla,\n" +
		"Adınız Soyadınız\n\n"

	// CV ekleme işlemi
	attachment, err := os.ReadFile(cvPath)
	if err != nil {
		log.Printf("CV dosyası okunamadı: %v", err)
		return err
	}
	baseName := filepath.Base(cvPath)
	attachmentPart := "--boundary42\n" +
		"Content-Type: application/octet-stream\n" +
		"Content-Disposition: attachment; filename=\"" + baseName + "\"\n" +
		"Content-Transfer-Encoding: base64\n\n" +
		encodeBase64(attachment) + "\n--boundary42--"

	// Son mesaj formatı
	message := header + body + attachmentPart

	// SMTP üzerinden gönderim
	err = smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, from, to, []byte(message))
	if err != nil {
		log.Printf("E-posta gönderilemedi: %s - Hata: %v", record.Email, err)
		return err
	}

	log.Printf("E-posta başarıyla gönderildi: %s", record.Email)
	return nil
}

// encodeBase64 encodes the file to Base64 format for attachment
func encodeBase64(data []byte) string {
	var buf bytes.Buffer
	for i := 0; i < len(data); i += 3 {
		buf.WriteString(fmt.Sprintf("%02X", data[i]))
		if i+1 < len(data) {
			buf.WriteString(fmt.Sprintf("%02X", data[i+1]))
		}
		if i+2 < len(data) {
			buf.WriteString(fmt.Sprintf("%02X", data[i+2]))
		}
	}
	return buf.String()
}
