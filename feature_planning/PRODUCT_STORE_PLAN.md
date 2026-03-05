# Rencana Pengembangan Fitur Produk & Toko (Multi-vendor Catalog)

Dokumen ini merinci perencanaan teknis untuk mengimplementasikan fitur manajemen toko (vendor) dan katalog produk pada aplikasi Swiftly.

## 1. Arsitektur Database (Schema Design)

Pengembangan akan mencakup penambahan tabel-tabel berikut dalam migrasi database PostgreSQL:

### A. Tabel `stores` (Informasi Toko)
*   `id`: UUID (Primary Key)
*   `user_id`: UUID (Foreign Key ke `users`, Unique - 1 user 1 toko)
*   `name`: VARCHAR(255)
*   `slug`: VARCHAR(255) (Unique, untuk URL `/store/:slug`)
*   `description`: TEXT
*   `logo_url`: VARCHAR(255) (Path di MinIO)
*   `banner_url`: VARCHAR(255) (Path di MinIO)
*   `theme_color`: VARCHAR(7) (Hex code untuk kustomisasi landing page)
*   `is_active`: BOOLEAN (Default: true)
*   `created_at`, `updated_at`: TIMESTAMP

### B. Tabel `store_addresses` (Alamat Toko)
*   `id`: UUID (Primary Key)
*   `store_id`: UUID (Foreign Key ke `stores`)
*   `province`: VARCHAR(100)
*   `city`: VARCHAR(100)
*   `district`: VARCHAR(100)
*   `postal_code`: VARCHAR(10)
*   `full_address`: TEXT
*   `latitude`, `longitude`: DECIMAL(10,8) (Untuk integrasi peta/ongkir)

### C. Tabel `store_notes` (Catatan Dinamis: FAQ, Kebijakan, dll)
*   `id`: UUID (Primary Key)
*   `store_id`: UUID (Foreign Key ke `stores`)
*   `title`: VARCHAR(100) (Contoh: "FAQ", "Ketentuan Pengiriman")
*   `slug`: VARCHAR(100)
*   `content`: TEXT (Mendukung Markdown)
*   `is_visible`: BOOLEAN (Default: true)
*   `order_priority`: INTEGER (Urutan tampilan)

### D. Tabel `global_categories` (Kategori Platform)
*   `id`: SERIAL (Primary Key)
*   `name`: VARCHAR(100)
*   `slug`: VARCHAR(100) (Unique)
*   `icon_url`: VARCHAR(255)

### E. Tabel `etalases` (Kategori Internal Toko)
*   `id`: UUID (Primary Key)
*   `store_id`: UUID (Foreign Key ke `stores`)
*   `name`: VARCHAR(100)
*   `slug`: VARCHAR(100)

### F. Tabel `products` (Data Utama Produk)
*   `id`: UUID (Primary Key)
*   `store_id`: UUID (Foreign Key ke `stores`)
*   `category_id`: INTEGER (Foreign Key ke `global_categories`)
*   `etalase_id`: UUID (Foreign Key ke `etalases`, Nullable)
*   `name`: VARCHAR(255)
*   `slug`: VARCHAR(255) (Unique per toko)
*   `description`: TEXT
*   `base_price`: DECIMAL(12,2)
*   `stock`: INTEGER (Stok global jika tidak ada varian)
*   `weight`: INTEGER (Dalam gram, untuk ongkir)
*   `is_active`: BOOLEAN (Default: true)
*   `created_at`, `updated_at`: TIMESTAMP

### G. Tabel `product_images` (Galeri Produk)
*   `id`: UUID (Primary Key)
*   `product_id`: UUID (Foreign Key ke `products`)
*   `image_url`: VARCHAR(255)
*   `is_primary`: BOOLEAN (Gambar utama untuk thumbnail)

### H. Tabel `product_variants` (Opsi Produk - Size/Color)
*   `id`: UUID (Primary Key)
*   `product_id`: UUID (Foreign Key ke `products`)
*   `name`: VARCHAR(100)
*   `sku`: VARCHAR(100) (Unique)
*   `price_adjustment`: DECIMAL(12,2)
*   `stock`: INTEGER
*   `image_url`: VARCHAR(255)

---

## 2. Fitur Kustomisasi Landing Page Toko (Vendor Control)

Vendor dapat memodifikasi tampilan halaman toko mereka melalui dashboard:
1.  **Layout Elements**: Menentukan urutan section (misal: Banner -> Produk Terbaru -> Catatan Toko).
2.  **Visual Branding**: Mengatur `theme_color` yang akan diaplikasikan pada button dan accent di landing page toko.
3.  **Dynamic Content**: Mengelola `store_notes` untuk memberikan informasi tambahan kepada pembeli.
4.  **Header Styling**: Memilih antara tampilan banner minimalis atau full-width.

---

## 3. Definisi API (Backend Routes)

### Public API (Pembeli)
*   `GET /api/v1/stores/:store_slug/notes`: List catatan publik toko.
*   `GET /api/v1/stores/:store_slug/notes/:note_slug`: Detail isi catatan.
*   `GET /api/v1/stores/:store_slug/address`: Informasi lokasi toko.

### Vendor API (Penjual)
*   `POST/PUT /api/v1/vendor/store/address`: Manajemen alamat toko.
*   `GET /api/v1/vendor/store/notes`: List semua catatan toko (termasuk yang hidden).
*   `POST /api/v1/vendor/store/notes`: Tambah catatan baru (FAQ, SOP, dll).
*   `PUT /api/v1/vendor/store/notes/:id`: Update judul atau isi catatan.
*   `DELETE /api/v1/vendor/store/notes/:id`: Hapus catatan.

---

## 4. Implementasi Frontend (Vue 3 + Tailwind)

### A. Vendor Dashboard (`/src/views/vendor/settings`)
1.  **Address Form**: Integrasi input alamat lengkap.
2.  **Notes Manager**: UI untuk mengelola list catatan dengan editor teks (Markdown).
3.  **Appearance Settings**: Color picker untuk `theme_color` dan toggle visibility untuk banner/description.

### B. Public Store Page (`/store/:slug`)
1.  **Identity Section**: Menampilkan Banner, Logo, dan Nama Toko dengan aksen warna sesuai `theme_color`.
2.  **Information Sidebar**: Navigasi ke `store_notes` (FAQ, SOP) dan Alamat Toko.
3.  **Custom Grid**: Produk disusun berdasarkan preferensi kustomisasi vendor.

---

## 5. Rencana Pengujian & Edge Cases

*   **Note Title Collision**: Dua catatan dengan judul sama di satu toko (Sistem harus unikkan slug).
*   **Missing Address**: Validasi saat checkout jika alamat toko belum diset (mempengaruhi ongkir).
*   **Excessive Customization**: Input warna hex yang tidak valid atau kontras yang buruk (Sistem harus punya fallback color).
*   **Note Order Priority**: Memastikan catatan tampil berurutan sesuai prioritas yang diset vendor.

---

## 6. Roadmap Eksekusi

1.  **Fase 1**: Database Migration (Penambahan tabel `addresses` dan `notes`).
2.  **Fase 2**: Backend API untuk Store Management lengkap (Address, Notes, Custom Settings).
3.  **Fase 3**: Implementasi Dashboard Vendor (UI Alamat & Notes).
4.  **Fase 4**: Pembuatan Dynamic Landing Page Toko untuk pembeli.
5.  **Fase 5**: Final Polish & Testing Edge Cases.
