# Discord Bot Go

Bot Discord multifungsi yang dikembangkan dengan bahasa Go dengan fitur:
- Integrasi OpenRouter (Multimodal AI dengan Tools/Functions Calling)
- Downloader YT-DLP
- Music & Radio Player
- Sistem keamanan yang kuat
- Fitur "Join to Create" Voice Channel
- Interaksi AI yang proaktif

## Fitur Utama

### 1. Integrasi OpenRouter (Multimodal AI)
- Menggunakan package HTTP Go (`net/http`) untuk API calls
- Menggunakan model `openrouter/sonoma-dusk-alpha` yang cepat, akurat, dan kuat untuk merespon pengguna di Discord
- Mendukung tools/functions calling
- Menangani berbagai jenis input: teks, gambar, file
- Dapat memanggil tools/functions secara otomatis berdasarkan permintaan pengguna

### 2. Downloader YT-DLP
- Eksekusi yt-dlp sebagai subprocess menggunakan `os/exec`
- Konfigurasi tanpa cookie: `--no-cookies`, `--no-check-certificates`
- Mendukung berbagai format: video, audio, thumbnail
- Parallel processing untuk multiple requests

### 3. Music & Radio Player
- Sistem antrian dengan goroutines dan mutex untuk thread safety
- Mendukung streaming langsung dari URL
- Kontrol pemutaran lengkap (play, pause, resume, skip, stop)
- Pengaturan volume

### 4. Sistem Keamanan
- Environment variables untuk API keys
- Sistem permission per command
- Rate limiting per user dengan sliding window
- Sanitasi input untuk mencegah injection
- Validasi URL untuk mencegah akses internal

### 5. Fitur Interaksi yang Canggih
- Semua pesan private chat langsung diteruskan ke AI
- Setiap 10 pesan di server, bot akan memberikan respons proaktif
- Prefix perintah diubah dari `!` menjadi `/`
- Sistem tools/functions yang dapat dipanggil oleh AI

### 6. Voice Channel Management
- Fitur "Join to Create" untuk membuat voice channel dinamis

## Persyaratan Sistem

- Go 1.21+
- FFmpeg
- yt-dlp

## Instalasi

```bash
# Menginstal dependensi
make deps

# Atau menggunakan go secara langsung
go mod tidy
```

## Konfigurasi

1. Salin `.env.example` ke `.env`:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` dan isi dengan nilai yang sesuai:
   ```env
   DISCORD_TOKEN=value
   OPENROUTER_API_KEY=value
   GOOGLE_SEARCH_API_KEY=value
   GOOGLE_SEARCH_ENGINE_ID=value
   BOT_PREFIX=/
   MAX_CONCURRENT_DOWNLOADS=3
   MAX_FILE_SIZE=100
   ```

3. **PENTING**: Ganti `your_discord_bot_token_here` dengan token bot Discord Anda yang sebenarnya.

## Penggunaan

### Menjalankan Bot

```bash
# Menjalankan langsung
make run

# Atau menggunakan go secara langsung
go run cmd/bot/main.go
```

### Interaksi dengan Bot

1. **Private Messages (DMs)**: Kirim pesan apa pun ke bot secara langsung, akan diteruskan ke AI
2. **Perintah di Server**: Gunakan prefix `/` diikuti dengan perintah:
   - `/help` - Menampilkan bantuan
   - `/ai <pertanyaan>` - Bertanya kepada AI
   - `/download <url>` - Mendownload video/audio
   - `/play <url>` - Memutar audio
   - `/pause` - Menjeda pemutaran
   - `/resume` - Melanjutkan pemutaran
   - `/skip` - Melewati ke track berikutnya
   - `/stop` - Menghentikan pemutaran
   - `/queue` - Menampilkan antrian
   - `/volume [level]` - Mengatur volume

### Testing

```bash
# Menjalankan semua test
make test

# Menjalankan test dengan coverage
make test-cover

# Menjalankan test dengan race detector
make test-race
```

### Quality Checks

```bash
# Memformat kode
make fmt

# Memeriksa kode dengan vet
make vet

# Menjalankan semua quality checks
make check
```

## Deployment

### Menggunakan Docker

```bash
# Membangun image
make docker-build

# Menjalankan container
make docker-run
```

### Menggunakan Docker Compose

```bash
# Menjalankan dengan docker-compose
make docker-compose-up

# Menghentikan
make docker-compose-down
```

## Tools yang Tersedia untuk AI

AI dapat memanggil tools berikut secara otomatis:
- `download_video` - Mendownload video/audio
- `play_music` - Memutar musik
- `get_video_info` - Mendapatkan informasi video
- `search_web` - Mencari informasi di web dengan flow sebagai berikut:
  1. AI memanggil tool karena kekurangan informasi real time atau permintaan pengguna
  2. AI memberikan kata kunci untuk hal yang dicari
  3. Sistem mengambil cuplikan website teratas terkait pencarian
  4. Melakukan web scraping di 4 website teratas
  5. Hasil scraping diteruskan ke AI dengan prompt 'tolong rangkum hasil web search ini dengan rapi'
  6. Hasil rangkuman AI diteruskan ke pengguna

## Struktur Proyek

```
discord-bot/
├── cmd/
│   └── bot/           # Entry point aplikasi
├── internal/          # Package internal
│   ├── config/        # Konfigurasi aplikasi
│   ├── openrouter/    # Client OpenRouter
│   ├── ytdlp/         # Downloader YT-DLP
│   ├── music/         # Music player
│   └── security/      # Sistem keamanan
├── pkg/               # Package eksternal
├── deployments/       # File deployment
├── go.mod             # Dependensi Go
├── go.sum             # Checksum dependensi
├── Dockerfile         # Konfigurasi Docker
├── docker-compose.yml # Konfigurasi Docker Compose
├── Makefile           # Script build dan test
├── README.md          # Dokumentasi
├── USAGE.md           # Panduan penggunaan
├── .env.example       # Contoh file environment
└── .gitignore         # File yang diabaikan Git
```

## Pengembangan Lebih Lanjut

Untuk informasi lebih detail tentang penggunaan dan pengembangan, lihat [USAGE.md](USAGE.md).
