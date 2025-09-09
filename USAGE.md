# Penggunaan Bot Discord

## Interaksi dengan Bot

Bot ini memiliki beberapa cara interaksi yang berbeda:

### 1. Private Messages (DMs)
- Kirim pesan apa pun ke bot secara langsung dalam bentuk chat pribadi
- Pesan akan langsung diteruskan ke AI tanpa perlu prefix perintah
- AI akan merespons secara langsung

### 2. Perintah di Server Discord
Gunakan prefix `/` diikuti dengan perintah di channel server:

### Perintah Bantuan
- `/help` - Menampilkan pesan bantuan dengan daftar perintah yang tersedia

### Perintah AI
- `/ai <pertanyaan>` atau `/ask <pertanyaan>` - Bertanya kepada AI dengan teks
- Contoh: `/ai Apa itu machine learning?`

### Perintah Download
- `/download <url> [-a]` atau `/dl <url> [-a]` - Mendownload video/audio dari URL
  - Gunakan `-a` atau `--audio` untuk mendownload audio saja
  - Contoh: `/download https://youtube.com/watch?v=example`
  - Contoh: `/download -a https://youtube.com/watch?v=example`

### Perintah Music Player
- `/play <url>` - Memutar audio dari URL
  - Contoh: `/play https://youtube.com/watch?v=example`
- `/pause` - Menjeda pemutaran
- `/resume` - Melanjutkan pemutaran
- `/skip` atau `/next` - Melewati ke track berikutnya
- `/stop` - Menghentikan pemutaran dan mengosongkan antrian
- `/queue` - Menampilkan antrian pemutaran saat ini
- `/volume [level]` - Menampilkan atau mengatur volume (0-100)
  - Contoh: `/volume` (menampilkan volume saat ini)
  - Contoh: `/volume 50` (mengatur volume ke 50%)

### Interaksi Proaktif
- Bot akan secara otomatis memberikan respons ke dalam percakapan setiap 10 pesan di server
- Respons ini akan berupa komentar atau pertanyaan yang relevan berdasarkan riwayat percakapan

## Sistem Keamanan

Bot ini dilengkapi dengan sistem keamanan yang mencakup:

### Rate Limiting
- Setiap pengguna dibatasi hingga 5 permintaan per menit
- Jika melebihi batas, pengguna akan menerima pesan "You are being rate limited. Please wait before sending more commands."

### Sanitasi Input
- Semua input dari pengguna di sanitasi untuk mencegah karakter berbahaya
- Panjang input dibatasi hingga 2000 karakter

### Validasi URL
- Semua URL yang dimasukkan divalidasi untuk mencegah akses ke IP internal atau localhost

## Tools yang Dapat Digunakan oleh AI

AI menggunakan model `openrouter/sonoma-dusk-alpha` yang cepat, akurat, dan kuat untuk merespon pengguna di Discord. AI memiliki kemampuan untuk memanggil tools/functions secara otomatis berdasarkan permintaan pengguna:

1. **download_video** - Mendownload video atau audio dari URL
2. **play_music** - Memutar musik dari URL
3. **get_video_info** - Mendapatkan informasi tentang video
4. **search_web** - Mencari informasi di web dengan flow sebagai berikut:
   - AI memanggil tool karena kekurangan informasi real time atau permintaan pengguna
   - AI memberikan kata kunci untuk hal yang dicari
   - Sistem mengambil cuplikan website teratas terkait pencarian
   - Melakukan web scraping di 4 website teratas
   - Hasil scraping diteruskan ke AI dengan prompt 'tolong rangkum hasil web search ini dengan rapi'
   - Hasil rangkuman AI diteruskan ke pengguna

## Voice Channel Management

### Join to Create
- Fitur "Join to Create" memungkinkan pengguna untuk membuat voice channel dinamis
- Ketika pengguna bergabung ke channel khusus, bot akan membuat channel baru untuk mereka
- Channel akan dihapus secara otomatis ketika kosong

## Konfigurasi

### Environment Variables
Bot ini menggunakan file `.env` untuk konfigurasi. Salin `.env.example` ke `.env` dan isi nilai yang sesuai:

```env
# Token bot Discord (diperlukan)
DISCORD_TOKEN=your_discord_bot_token_here

# API key OpenRouter (diperlukan untuk fitur AI)
OPENROUTER_API_KEY=your_openrouter_api_key_here

# Google Custom Search API Key (diperlukan untuk fitur web search)
GOOGLE_SEARCH_API_KEY=your_google_search_api_key_here

# Google Custom Search Engine ID (diperlukan untuk fitur web search)
GOOGLE_SEARCH_ENGINE_ID=your_search_engine_id_here

# Prefix perintah bot (opsional, default: /)
BOT_PREFIX=/

# Maksimal download bersamaan (opsional, default: 3)
MAX_CONCURRENT_DOWNLOADS=3

# Ukuran maksimal file dalam MB (opsional, default: 100)
MAX_FILE_SIZE=100
```

## Pengembangan

### Menjalankan Test
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
# Membangun image Docker
make docker-build

# Menjalankan container
make docker-run
```

### Menggunakan Docker Compose
```bash
# Menjalankan dengan docker-compose
make docker-compose-up

# Menghentikan docker-compose
make docker-compose-down
```

## Fitur yang Tersedia

Bot ini telah memiliki fitur-fitur berikut:

1. **Integrasi OpenRouter dengan Tools/Functions Calling**
   - Dukungan untuk model AI yang dapat memanggil tools
   - Kemampuan memproses berbagai jenis input
   - Respons yang lebih interaktif dan fungsional

2. **Downloader YT-DLP yang Lengkap**
   - Mendownload video dan audio dari berbagai platform
   - Konfigurasi fleksibel untuk berbagai kebutuhan

3. **Music Player dengan Kontrol Lengkap**
   - Sistem antrian musik
   - Berbagai kontrol pemutaran
   - Pengaturan volume

4. **Interaksi yang Cerdas**
   - Respons proaktif dalam percakapan
   - Kemampuan memahami konteks percakapan
   - Pengalaman pengguna yang lebih alami

5. **Voice Channel Management**
   - Fitur "Join to Create" untuk voice channel dinamis
   - Manajemen channel otomatis