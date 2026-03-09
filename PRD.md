# Final PRD: Project ZenVault

**Goal:** Membangun API Golang ultra-ringan di k3d dengan otomatisasi Secret (ESO + GSM) dan TLS (cert-manager + ZeroSSL).

---

## 1. Metadata & Labeling (Best Practice)
Semua resource di Google Secret Manager (GSM) wajib memiliki label berikut:
- `project: zenvault`
- `environment: production`
- `managed-by: eso`
- `owner: stayrelevantid`

## 2. Arsitektur Teknis
- **Runtime:** Go 1.22+ (Static Binary)
- **Container:** `scratch` image (Size: ~10MB)
- **Secret:** Google Secret Manager ditarik oleh External Secrets Operator (ESO)
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

### Fase 2: Cloud Setup (GCP & ZeroSSL)
Persiapan "Brankas" di Cloud.
- **GSM Setup:** Buat secret `APP_DEBUG_KEY` dan `ZEROSSL_EAB_HMAC`. Berikan label `owner: stayrelevantid`.
- **IAM:** Buat Service Account dengan role *Secret Manager Secret Accessor*.
- **ZeroSSL:** Ambil `EAB KID` dan `HMAC Key` dari dashboard ZeroSSL.

### Fase 3: Cluster Orchestration (k3d & Operator)
Menyiapkan "Kurir" di dalam Kubernetes.
- **k3d Up:** Jalankan cluster dengan port `80` & `443` terbuka.
- **Install ESO:** Pasang External Secrets Operator via Helm.
- **Install cert-manager:** Pasang cert-manager via Helm.
- **Connect GCP:** Masukkan JSON Key Service Account ke k3d sebagai Secret untuk digunakan oleh ESO.

### Fase 4: Integration & Automation
Menghubungkan semua titik.
- **Deploy SecretStore & ExternalSecret:** ESO akan menarik key dari GSM dan mengubahnya menjadi K8s Secret.
- **Deploy ClusterIssuer:** Konfigurasi cert-manager agar mengenali ZeroSSL menggunakan EAB dari GSM.
- **Deploy App & Ingress:**
  - Deployment memanggil env dari K8s Secret.
  - Ingress meminta sertifikat TLS ke cert-manager.

### Fase 5: Verification (The Testing)
- **Check HTTPS:** Akses `https://zenvault.local` (Pastikan sertifikat dari ZeroSSL).
- **Check Secret:** Akses `/v1/debug` (Pastikan data muncul dari GSM).
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