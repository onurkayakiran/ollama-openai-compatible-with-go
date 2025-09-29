# OpenAI Uyumlu API

Bu proje, Ollama destekli ve OpenAI API formatına tam uyumlu bir REST API sağlar. Go Fiber framework'ü kullanılarak geliştirilmiştir.

## Özellikler

- ✅ OpenAI API formatına tam uyumluluk
- ✅ Chat Completions endpoint (`/v1/chat/completions`)
- ✅ Text Completions endpoint (`/v1/completions`)
- ✅ Models endpoint (`/v1/models`)
- ✅ Streaming response desteği (Server-Sent Events)
- ✅ API Key authentication
- ✅ CORS desteği
- ✅ Hata yönetimi ve logging
- ✅ .env dosyası ile environment variable konfigürasyonu

## Gereksinimler

- Go 1.21 veya üzeri
- Ollama kurulu ve çalışır durumda

## Kurulum

1. Repoyu klonlayın:
```bash
git clone <repo-url>
cd openai-compatible
```

2. Bağımlılıkları yükleyin:
```bash
go mod tidy
```

3. Environment konfigürasyonunuzu oluşturun:
```bash
cp .env.example .env
```

4. `.env` dosyasını kendi konfigürasyonunuza göre düzenleyin:
```bash
# Port ayarları
PORT=8080

# Authentication için API Key
API_KEY=sizin-gizli-api-anahtariniz

# Ollama sunucu ayarları
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=llama3.2:latest
```

5. Ollama'nın çalıştığından emin olun:
```bash
ollama serve
```

6. Uygulamayı çalıştırın:
```bash
go run main.go
```

## Konfigürasyon

Uygulama, konfigürasyon için environment variable'ları kullanır. Bu değişkenleri `.env` dosyasında ayarlayabilirsiniz:

| Değişken | Açıklama | Varsayılan |
|----------|----------|------------|
| `PORT` | Sunucu portu | 8080 |
| `API_KEY` | Kimlik doğrulama için API anahtarı | sk-your-secret-api-key-here |
| `OLLAMA_URL` | Ollama sunucu URL'i | http://localhost:11434 |
| `OLLAMA_MODEL` | Kullanılacak model | llama3.2:latest |

**Not:** Güvenlik nedeniyle `.env` dosyası `.gitignore` ile versiyon kontrolünden hariç tutulmuştur. Her zaman `.env.example` dosyasını şablon olarak kullanın.

## API Kullanımı

### Chat Completions

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sizin-gizli-api-anahtariniz" \
  -d '{
    "model": "llama3.2:latest",
    "messages": [
      {"role": "user", "content": "Merhaba!"}
    ]
  }'
```

### Streaming Chat Completions

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sizin-gizli-api-anahtariniz" \
  -d '{
    "model": "llama3.2:latest",
    "messages": [
      {"role": "user", "content": "Merhaba!"}
    ],
    "stream": true
  }'
```

### Text Completions

```bash
curl -X POST http://localhost:8080/v1/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sizin-gizli-api-anahtariniz" \
  -d '{
    "model": "llama3.2:latest",
    "prompt": "Merhaba dünya",
    "max_tokens": 100
  }'
```

### Model Listesi

```bash
curl -X GET http://localhost:8080/v1/models \
  -H "Authorization: Bearer sizin-gizli-api-anahtariniz"
```

### Sağlık Kontrolü

```bash
curl -X GET http://localhost:8080/health
```

## Proje Yapısı

```
openai-compatible/
├── main.go                 # Ana uygulama
├── go.mod                  # Go modül dosyası
├── go.sum                  # Go bağımlılık checksum'ları
├── .env                    # Environment variables (git'te yok)
├── .env.example            # Örnek environment variables
├── .gitignore             # Git ignore dosyası
├── config/
│   └── config.go          # Konfigürasyon yönetimi
├── handlers/
│   ├── chat.go            # Chat completions handler
│   ├── completions.go     # Text completions handler
│   └── models.go          # Models handler
├── middleware/
│   └── auth.go            # Authentication middleware
├── models/
│   ├── openai.go          # OpenAI API yapıları
│   └── ollama.go          # Ollama API yapıları
└── services/
    └── ollama.go          # Ollama servis entegrasyonu
```

## Desteklenen Parametreler

### Chat Completions
- `model`: Model adı
- `messages`: Mesaj dizisi
- `max_tokens`: Maksimum token sayısı
- `temperature`: Yaratıcılık seviyesi (0.0-2.0)
- `top_p`: Nucleus sampling
- `stream`: Streaming aktif/pasif
- `stop`: Durma dizileri
- `presence_penalty`: Presence penalty
- `frequency_penalty`: Frequency penalty

### Text Completions
- `model`: Model adı
- `prompt`: Metin prompt'u
- `max_tokens`: Maksimum token sayısı
- `temperature`: Yaratıcılık seviyesi
- `top_p`: Nucleus sampling
- `stop`: Durma dizileri

## Güvenlik

- API anahtarları authentication middleware üzerinden doğrulanır
- Hassas bilgiler içeren `.env` dosyası versiyon kontrolünden hariç tutulmuştur
- Üretim ortamlarında her zaman güçlü ve benzersiz API anahtarları kullanın

## Production için Build

```bash
# Binary oluştur
go build -o openai-compatible

# Binary'yi çalıştır
./openai-compatible
```

## Lisans

MIT License

---

[🇬🇧 Click for English README](README.md)
