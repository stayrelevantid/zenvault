# 🔐 ZenVault

Ultra-lightweight Golang API running on k3d with automated Secret Management (ESO + AWS Secrets Manager) and TLS (cert-manager + ZeroSSL).

---

## ⚡ Highlights

- **Go Static Binary** — built from `scratch`, image size ~6MB
- **Secret Automation** — AWS Secrets Manager → External Secrets Operator → K8s Secret
- **TLS Automation** — ZeroSSL certificates via cert-manager + Traefik Ingress
- **Minimal Footprint** — zero OS dependencies, ultra-low RAM usage

## 🏗️ Architecture

```
┌─────────────┐       ┌──────────────┐       ┌─────────────────┐
│  ZeroSSL    │──TLS──│  Traefik     │──────▶│  ZenVault API   │
│  (CA)       │       │  Ingress     │       │  (scratch image) │
└─────────────┘       └──────────────┘       └────────┬────────┘
                                                      │ env
       ┌──────────────────┐       ┌───────────────────┘
       │   AWS Secrets    │──ESO──│  K8s Secret
       │   Manager        │       │  (auto-synced)
       └──────────────────┘       └───────────────────
```

## 📁 Project Structure

```
zenvault/
├── app/
│   ├── main.go            # Go API (endpoints: /healthz, /v1/debug)
│   ├── Dockerfile          # Multi-stage build (golang:alpine → scratch)
│   ├── go.mod              # Go module
│   └── .env.example        # Environment variable template
├── manifests/
│   ├── 01-eso-setup.yaml   # SecretStore & ExternalSecret
│   ├── 02-cert-manager.yaml# ClusterIssuer ZeroSSL
│   ├── 03-deployment.yaml  # App Deployment & Service
│   └── 04-ingress.yaml     # Ingress with TLS Annotation
├── scripts/
│   └── setup-local-dns.sh  # Script edit /etc/hosts
├── PRD.md                  # Product Requirements Document
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Go 1.25+
- Docker
- k3d
- kubectl
- Helm

### 1. Local Development

```bash
# Clone repo
git clone https://github.com/stayrelevantid/zenvault.git
cd zenvault/app

# Setup env
cp .env.example .env

# Run locally
go run main.go
```

```bash
# Test endpoints
curl http://localhost:8080/healthz
curl http://localhost:8080/v1/debug
```

### 2. Docker Build & Run

```bash
cd app

# Build image
docker build -t zenvault:latest .

# Run container
docker run -p 8080:8080 -e APP_DEBUG_KEY=my-secret zenvault:latest
```

### 3. Kubernetes (k3d)

```bash
# Create cluster
k3d cluster create zenvault -p "80:80@loadbalancer" -p "443:443@loadbalancer"

# Install operators
helm install eso external-secrets/external-secrets -n external-secrets --create-namespace
helm install cert-manager jetstack/cert-manager -n cert-manager --create-namespace --set installCRDs=true

# Deploy manifests
kubectl apply -f manifests/
```

## 🔗 API Endpoints

| Endpoint | Method | Description |
|---|---|---|
| `/healthz` | GET | Health check |
| `/v1/debug` | GET | Returns `APP_DEBUG_KEY` value from environment |

### Example Response

```json
// GET /v1/debug
{
  "app": "zenvault",
  "app_debug_key": "your-secret-value"
}
```

## 🏷️ AWS Resource Prefix

All AWS Secrets Manager resources use this prefix:

| Prefix | Value |
|---|---|
| `zenvault/` | `zenvault/app_debug_key`, `zenvault/zerossl_eab_hmac` |

## ⚔️ GCP vs AWS: Perbandingan Pengalaman Lab

| Aspek | 🔵 GCP (yang coba pertama) | 🟠 AWS (yang berhasil) |
|---|---|---|
| **Secret Storage** | Google Secret Manager | AWS Secrets Manager |
| **Auth Method** | Service Account JSON Key | IAM User Access Key + Secret Key |
| **Hambatan Utama** | Org Policy blokir pembuatan SA JSON Key | Tidak ada hambatan |
| **WIF Attempt** | ✅ Setup OIDC pool & provider → ❌ Gagal (ESO hardcode GCE Metadata Server) | N/A |
| **ESO Provider** | `gcpsm` (terikat lingkungan GKE) | `aws` (fleksibel, bisa lokal) |
| **Credential di k8s** | JSON Key file → K8s Secret | Access Key ID + Secret Key → K8s Secret |
| **Integrasi Lokal** | ❌ WIF tidak bisa di luar GCE/GKE | ✅ Static Credentials bekerja sempurna |
| **ACME / ZeroSSL** | ✅ Sama-sama berhasil registrasi account | ✅ Berhasil registrasi + fetch HMAC dari AWS |
| **TLS Sertifikat** | ❌ `.local` domain ditolak CA publik | ❌ Sama, karena bukan masalah cloud-nya |
| **Waktu Setup** | Lebih lama (troubleshooting WIF) | Lebih cepat dan mulus |

> **Kesimpulan Perbandingan:** Untuk *local development* dengan k3d, AWS lebih unggul karena IAM User static credentials tidak memerlukan infrastruktur khusus cloud. GCP WIF jauh lebih aman untuk *production* di GKE, namun tidak bisa dipakai di luar ekosistem GCP.

## 🏁 Kesimpulan

Lab ini berhasil membangun API Golang ultra-ringan (~6MB Docker image) yang berjalan di cluster Kubernetes lokal (`k3d`), dengan automasi penuh untuk manajemen rahasia via **AWS Secrets Manager + ESO** dan automasi TLS via **cert-manager + ZeroSSL**.

Pipeline sinkronisasi rahasia berjalan sepenuhnya:
```
AWS Secrets Manager → ESO → Kubernetes Secret → Pod Env → /v1/debug ✅
```

## 📚 Lessons Learned

### 1. 🔒 Org Policy Lebih Berkuasa dari Siapapun
Kebijakan keamanan level organisasi (contoh: `iam.disableServiceAccountKeyCreation` di GCP) bisa memblokir siapapun. Kita harus selalu siap **beradaptasi** — dalam kasus ini, dengan beralih ke AWS.

### 2. 🔬 WIF Lokal Itu Sangat Tricky
Provider `gcpsm` pada ESO secara *hardcoded* memanggil GCE Metadata Server (`169.254.169.254`) yang hanya ada di VM GCP/GKE. WIF tidak kompatibel di environment lokal seperti `k3d`.

### 3. ⚡ ESO = Zero-Code Secret Management
Aplikasi Golang (`main.go`) tidak punya satu baris pun kode AWS SDK. Namun bisa membaca rahasia dari cloud — semua dikerjakan ESO di balik layar secara *zero-ops*.

### 4. 🌐 HTTP-01 ACME Butuh IP Publik Nyata
ZeroSSL/Let's Encrypt butuh server yang bisa dijangkau dari internet saat validasi. Domain `.local` atau laptop di balik NAT/WiFi tidak bisa divalidasi. Sertifikat hijau baru muncul jika cluster berjalan di VPS dengan IP Publik + DNS A record.

### 5. 🏗️ Adaptasi adalah Skill Utama DevOps
Lab ini adalah simulasi *production-grade* yang realistis — lengkap dengan *roadblock* nyata. Pelajaran terpentingnya bukan hanya tentang tools, tapi tentang **ketangguhan adaptasi**.

## 📄 License

MIT

---

> Built with ☕ by [@stayrelevantid](https://github.com/stayrelevantid)
