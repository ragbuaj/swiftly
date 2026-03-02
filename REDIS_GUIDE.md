# Redis Learning Guide for Swiftly 🚀

Redis (**Remote Dictionary Server**) adalah penyimpanan struktur data di dalam memori (RAM) yang sangat cepat. Di proyek ini, kita menggunakannya untuk meningkatkan keamanan dan performa.

---

## 1. Kenapa Redis? (Teori Dasar)
Berbeda dengan PostgreSQL yang menyimpan data di Disk (Harddrive), Redis menyimpan data di RAM.
*   **Kecepatan:** Redis dapat menangani >100.000 request per detik.
*   **Data Model:** Bukan tabel, tapi **Key-Value Pair** (seperti kamus atau objek JavaScript).
*   **Temporary:** Sangat cocok untuk data yang punya masa berlaku (TTL - Time To Live).

---

## 2. Cara Instalasi di Proyek Ini
Karena Swiftly menggunakan Docker, Anda tidak perlu menginstal Redis di Windows Anda.

### A. Konfigurasi Docker
Layanan Redis sudah ditambahkan di `docker-compose.yaml`:
```yaml
  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
    networks:
      - swiftly-network
```

### B. Menjalankan Redis
Jalankan perintah ini di terminal root:
```bash
docker compose up -d redis
```

---

## 3. Cara Akses Redis (Manual)
Anda bisa "masuk" ke dalam server Redis untuk melihat data yang tersimpan secara langsung.

1.  **Masuk ke CLI Redis:**
    ```bash
    docker compose exec redis redis-cli
    ```
2.  **Perintah Dasar di dalam CLI:**
    *   `PING` -> Harus dijawab `PONG` (Cek koneksi).
    *   `SET nama "Budi"` -> Menyimpan data.
    *   `GET nama` -> Mengambil data.
    *   `KEYS *` -> Melihat semua kunci yang ada.
    *   `DEL nama` -> Menghapus data.
    *   `EXPIRE nama 60` -> Memberi waktu hidup 60 detik pada kunci 'nama'.
    *   `TTL nama` -> Melihat sisa waktu hidup kunci tersebut.
    *   `QUIT` -> Keluar dari CLI.

---

## 4. Implementasi di Backend Go
Di proyek Swiftly, kita menggunakan library `go-redis/v9`.

### A. Inisialisasi (di `internal/database/redis.go`)
Kita membuat koneksi tunggal (*Singleton*) yang bisa dipakai di mana saja.
```go
rdb := redis.NewClient(&redis.Options{
    Addr: "redis:6379", // Nama service di docker-compose
})
```

### B. Contoh Penggunaan (Set & Get)
```go
ctx := context.Background()
rdb := database.GetRedis()

// Menyimpan data dengan waktu kadaluwarsa 1 jam
err := rdb.Set(ctx, "session:123", "active", time.Hour).Err()

// Mengambil data
val, err := rdb.Get(ctx, "session:123").Result()
```

---

## 5. Studi Kasus: Token Blacklisting (Sudah Terpasang)
Fitur ini digunakan agar user yang sudah Logout tidak bisa masuk lagi dengan token lama.

**Alurnya:**
1.  User klik **Logout**.
2.  Backend mengambil Token JWT tersebut.
3.  Backend simpan ke Redis: `SET blacklist:TOKEN_LAMA "true" EX 86400` (86400 detik = 24 jam).
4.  Setiap User akses API, `AuthMiddleware` mengecek: `EXISTS blacklist:TOKEN_DITERIMA`.
5.  Jika ada di Redis, akses ditolak meskipun Token JWT-nya asli.

---

## 6. Apa yang Bisa Dipelajari Selanjutnya?
*   **Caching Database:** Simpan hasil query SQL yang berat ke Redis agar tidak membebani PostgreSQL.
*   **Rate Limiting:** Gunakan Redis untuk menghitung berapa kali seorang User mencoba login (mencegah brute-force).
*   **Pub/Sub:** Mengirim notifikasi antar layanan secara real-time.

---

*Selamat belajar! Jika ada error saat menjalankan Redis, pastikan port 6379 di komputer Anda tidak sedang dipakai oleh aplikasi lain.*
