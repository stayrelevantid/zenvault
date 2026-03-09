# 🔐 ZenVault

Ultra-lightweight Golang API running on k3d with automated Secret Management (ESO + Google Secret Manager) and TLS (cert-manager + ZeroSSL).

---

## ⚡ Highlights

- **Go Static Binary** — built from `scratch`, image size ~6MB
- **Secret Automation** — Google Secret Manager → External Secrets Operator → K8s Secret
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
       │  Google Secret   │──ESO──│  K8s Secret
       │  Manager (GSM)   │       │  (auto-synced)
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

## 🏷️ GCP Resource Labels

All Google Secret Manager resources use these labels:

| Label | Value |
|---|---|
| `project` | `zenvault` |
| `environment` | `production` |
| `managed-by` | `eso` |
| `owner` | `stayrelevantid` |

## 📄 License

MIT

---

> Built with ☕ by [@stayrelevantid](https://github.com/stayrelevantid)
