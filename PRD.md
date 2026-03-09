# Final PRD: Project ZenVault

**Goal:** Membangun API Golang ultra-ringan di k3d dengan otomatisasi Secret (ESO + AWS Secrets Manager) dan TLS (cert-manager + ZeroSSL).

---

## 1. Metadata & Naming Convention
Semua resource di AWS Secrets Manager akan menggunakan prefix berikut (jika diperlukan):
- `zenvault/`
- Region: `ap-southeast-1`

## 2. Arsitektur Teknis
- **Runtime:** Go 1.22+ (Static Binary)
- **Container:** `scratch` image (Size: ~10MB)
- **Secret:** AWS Secrets Manager ditarik oleh External Secrets Operator (ESO)
- **TLS:** ZeroSSL diotomatisasi oleh cert-manager via Ingress Traefik

## 3. Fase Pembuatan (Roadmap Terstruktur)

### Fase 1: Development API & Dockerization
Fokus pada pembuatan aplikasi yang "Environment Aware".
- **Inisialisasi Go:** Buat `main.go` yang membaca variabel `APP_DEBUG_KEY` dari environment.
- **Local Testing:** Gunakan file `.env` untuk simulasi lokal.
- **Multi-stage Build:**
  - **Stage 1:** Build binary menggunakan `golang:alpine`.
  - **Stage 2:** Pindahkan binary ke `scratch` untuk hasil akhir yang super kecil.
- **Verifikasi:** Pastikan image Docker < 15MB.

### Fase 2: Cloud Setup (AWS & ZeroSSL)
Persiapan "Brankas" di Cloud.
- **AWS Secrets Setup:** Buat secret `zenvault/app_debug_key` dan `zenvault/zerossl_eab_hmac`.
- **IAM:** Buat IAM User `zenvault-eso` dengan policy khusus untuk membaca tipe Secret ini. Tulis Access Key dan Secret Access Key ke environment lokal.
- **ZeroSSL:** Ambil `EAB KID` dan `HMAC Key` dari dashboard ZeroSSL.

### Fase 3: Cluster Orchestration (k3d & Operator)
Menyiapkan "Kurir" di dalam Kubernetes.
- **k3d Up:** Jalankan cluster dengan port `80` & `443` terbuka.
- **Install cert-manager:** Pasang cert-manager via Helm.
- **Connect AWS:** Masukkan Access Key dan Secret Key AWS ke k3d sebagai Kubernetes Secret (`aws-secret`) untuk digunakan oleh ESO.

### Fase 4: Integration & Automation
Menghubungkan semua titik.
- **Deploy SecretStore & ExternalSecret:** ESO akan menarik key dari AWS Secrets Manager dan mengubahnya menjadi K8s Secret.
- **Deploy ClusterIssuer:** Konfigurasi cert-manager agar mengenali ZeroSSL menggunakan EAB dari AWS Secrets Manager.
- **Deploy App & Ingress:**
  - Deployment memanggil env dari K8s Secret.
  - Ingress meminta sertifikat TLS ke cert-manager.

### Fase 5: Verification (The Testing)
- **Check HTTPS:** Akses `https://zenvault.local` (Pastikan sertifikat dari ZeroSSL).
- **Check Secret:** Akses `/v1/debug` (Pastikan data muncul dari AWS).
- **Check Resource:** Jalankan `kubectl top pod` (Pastikan RAM usage minimal).

## 4. Struktur Repositori yang Disarankan

```plaintext
zenvault/
├── app/
│   ├── main.go
│   ├── Dockerfile
│   └── .env.example
├── manifests/
│   ├── 01-eso-setup.yaml        # SecretStore & ExternalSecret
│   ├── 02-cert-manager.yaml     # ClusterIssuer ZeroSSL
│   ├── 03-deployment.yaml       # App Deployment & Service
│   └── 04-ingress.yaml          # Ingress with TLS Annotation
└── scripts/
    └── setup-local-dns.sh       # Script edit /etc/hosts
```