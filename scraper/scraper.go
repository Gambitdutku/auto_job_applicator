package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type PublicAPIResponse struct {
	Name     string `json:"name"`
	Domain   string `json:"domain"`
	Location string `json:"location"`
}

type GoogleSearchResponse struct {
	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
}

// SearchCompanies performs Google search and fetches companies
func SearchCompanies(query string) []string {
	log.Printf("Google araması başlatılıyor: %s\n", query)

	// Google araması yap
	googleResults := googleSearch(query)

	// Public API'lerden çekilen sonuçlar
	publicAPIResults := fetchFromPublicAPIs(query)

	// Sonuçları birleştir
	allResults := append(googleResults, publicAPIResults...)

	return allResults
}

// Google araması için API isteği yapar
func googleSearch(query string) []string {
	apiKey := "" // Google API Key
	cx := ""     // Google Custom Search CX ID

	// URL encode edilmesi gereken sorgu parametresi
	encodedQuery := url.QueryEscape(query)

	// API isteği URL'si
	url := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&key=%s&cx=%s", encodedQuery, apiKey, cx)

	// API isteğini yap
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Google araması sırasında hata oluştu: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	// Yanıtın içerik tipini kontrol et
	contentType := resp.Header.Get("Content-Type")
	log.Printf("Yanıt İçerik Türü: %s\n", contentType)

	// API yanıtını oku
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Yanıt okunamadı: %v\n", err)
		return nil
	}

	// Yanıtı logla
	log.Printf("Google API Yanıtı: %s\n", string(body))

	// Yanıtın JSON formatında olup olmadığını kontrol et
	if resp.StatusCode != 200 {
		log.Printf("Google API Hatası: %v\n", string(body))
		return nil
	}

	var googleResponse GoogleSearchResponse
	if err := json.Unmarshal(body, &googleResponse); err != nil {
		log.Printf("JSON çözümleme hatası: %v\n", err)
		return nil
	}

	// Arama sonuçlarını topla
	var results []string
	for _, item := range googleResponse.Items {
		results = append(results, item.Link)
	}

	return results
}

// fetchFromPublicAPIs fetches company data from public APIs
func fetchFromPublicAPIs(query string) []string {
	var results []string

	// Clearbit API örneği
	url := fmt.Sprintf("https://company.clearbit.com/v2/companies/find?domain=%s", query)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer YOUR_API_KEY") // Buraya API Key gelecek

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Clearbit API hatası: %v\n", err)
		return results
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var response PublicAPIResponse
		body, _ := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &response); err == nil {
			if response.Domain != "" {
				results = append(results, "https://"+response.Domain)
			}
		} else {
			log.Printf("Clearbit API yanıtı çözümlenemedi: %v", err)
		}
	}

	// Türkiye bazlı API - TOBB Sanayi Veritabanı (Örnek URL)
	tobbURL := "https://tobb.org.tr/sanayi/veritabanı/api?city=istanbul"
	tobbResp, err := http.Get(tobbURL)
	if err != nil {
		log.Printf("TOBB API hatası: %v\n", err)
		return results
	}
	defer tobbResp.Body.Close()

	if tobbResp.StatusCode == 200 {
		var tobbData []PublicAPIResponse
		body, _ := ioutil.ReadAll(tobbResp.Body)
		if err := json.Unmarshal(body, &tobbData); err == nil {
			for _, company := range tobbData {
				if company.Domain != "" {
					results = append(results, "https://"+company.Domain)
				}
			}
		} else {
			log.Printf("TOBB API yanıtı çözümlenemedi: %v", err)
		}
	}

	return results
}

// ExtractEmails finds email addresses on the provided URL
func ExtractEmails(url string) []string {
	log.Printf("E-posta aranıyor: %s\n", url)

	// HTTP GET isteği
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("URL alınamadı: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	// Sayfa içeriğini oku
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Sayfa okunamadı: %v\n", err)
		return nil
	}

	// E-posta adreslerini regex ile bul
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	emails := emailRegex.FindAllString(string(body), -1)

	// Aynı e-posta adreslerini filtrele
	uniqueEmails := make(map[string]struct{})
	for _, email := range emails {
		uniqueEmails[email] = struct{}{}
	}

	// Tekrarları kaldırılmış listeyi döndür
	var finalEmails []string
	for email := range uniqueEmails {
		if strings.Contains(email, "example") {
			continue // Örnek domainler alınmasın
		}
		finalEmails = append(finalEmails, email)
	}

	log.Printf("%d adet e-posta bulundu: %v\n", len(finalEmails), finalEmails)
	return finalEmails
}
