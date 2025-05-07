
# Otomatik İş Başvuru Sistemi

Bu proje, şirketlerini otomatik olarak arayıp, bu şirketlerin e-posta adreslerini tarayarak iş başvurusu göndermeyi amaçlamaktadır. Şirket detaylarını ve e-posta adreslerini almak için çeşitli public API'ler kullanılmaktadır.

## Özellikler
- Google Custom Search API ve diğer public API'ler kullanılarak şirketler taranır.
- Şirketlerin internet sitelerinden e-posta adresleri çıkarılır.
- E-posta adresleri veritabanına kaydedilir ve gönderilme durumu takip edilir.
- Geçerli e-posta adreslerine iş başvuru e-postaları gönderilir.

## Gereksinimler

- Go 1.16+
- Google Custom Search API Anahtarı
- Clearbit API Anahtarı (Ekstra şirket bilgileri için)
- Türkiye bazlı şirket verisi için geçerli TOBB API

## Kurulum

### 1. Depoyu Klonlayın

```bash
git clone https://github.com/Gambitdutku/auto_job_applicator/
cd auto_job_applicator
```

### 2. Bağımlılıkları Yükleyin

Gerekli bağımlılıkları yüklemek için aşağıdaki komutu çalıştırın:

```bash
go get github.com/mattn/go-sqlite3
go get github.com/gocolly/colly
```

### 3. Çevre Değişkenlerini Ayarlayın

Aşağıdaki çevre değişkenlerini ayarlamanız gerekecek:

- `GOOGLE_API_KEY`: Google Custom Search API anahtarınız.
### 4. Uygulamayı Çalıştırın

Kurulumdan sonra, uygulamayı şu komutla çalıştırabilirsiniz:

```bash
go run main.god
```

## Kullanılan Kütüphaneler

- [Go-SQLite3](https://github.com/mattn/go-sqlite3): Go için SQLite3 veritabanı kütüphanesi.
- [Colly](https://github.com/gocolly/colly): Go için güçlü bir web scraping framework'ü.
- [Clearbit API](https://clearbit.com/docs): Şirket ve iletişim bilgisi sağlayan bir API.
- [Google Custom Search API](https://developers.google.com/custom-search/v1/overview): Google'dan arama sonuçları almak için kullanılır.

## Katkıda Bulunma

Bu depoyu çatallayabilir ve geliştirmeler veya düzeltmeler yapmak için pull request gönderebilirsiniz.

## Lisans

Bu proje MIT Lisansı ile lisanslanmıştır - detaylar için [LICENSE](LICENSE) dosyasına bakabilirsiniz.
