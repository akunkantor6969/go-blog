#go-blog API — Backend Santuy Rasa Proyek Tongkrongan

Yo! Selamat datang di repo **go-blog** 
Ini bukan proyek korporat kaku berdasi, tapi murni eksperimen ngoding bareng sirkel buat bikin REST API pakai **Go (Golang)** + **MariaDB/MySQL** tanpa *framework* macem-macem (pure `net/http` bos, biar kerasa ototnya).

---

## Fitur Utama (Udah Kayak Startup Tier-1)
- **CRUD Artikel**: Bisa *create*, baca semua (udah ada *pagination*-nya dong!), baca detail, edit, sampai hapus artikel.
- **Satpam JWT (`BearerAuth`)**: Yang mau bikin/ngubah artikel wajib login dulu biar dapet tiket masuk berupa token JWT. Nggak bisa asal nyelonong!
- **Swagger UI Interaktif**: Males buka Bruno/Postman? Tenang, sekarang dokumentasinya udah *live* di browser! Tinggal klik-klik *Try it out*, kelar urusan.

---

## Tech Stack
- **Bahasa**: Go (Golang)
- **Database**: MariaDB / MySQL
- **Auth**: JWT (JSON Web Token)
- **Dokumentasi**: Swaggo (`swag` + `http-swagger`)

---

## Cara Running di Laptop Masing-Masing

1. **Clone dulu repo ini, jangan cuma diliatin**:
   ```bash
   git clone <url-repo-kamu>
   cd go-blog