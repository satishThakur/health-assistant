# Go Live Plan: Health Assistant

**Goal:** App running on your iPhone/iPad, backend on AWS, Garmin syncing hourly.

**Last Updated:** 2026-02-22
**Status:** Planning

---

## Prerequisites Checklist

| Item | Status | Notes |
|------|--------|-------|
| Apple Developer account | ❌ Need to create | developer.apple.com, $99/yr, **approval takes 24-48h — do this first** |
| AWS account | ✅ Exists | |
| Domain name | ❌ Need to purchase | ~$12/yr for .com via Route53 or any registrar |
| Google OAuth (production) | ❌ Need iOS client | Currently dev mode (empty GOOGLE_CLIENT_ID) |
| Garmin scraper | ✅ Working locally | Just needs to move to AWS |

---

## Architecture (Target State)

```
iPhone/iPad
    │
    │  HTTPS
    ▼
Route53 → ALB (or Nginx) → EC2 t3.small
                               ├── Go backend (port 8083)
                               ├── TimescaleDB (Docker)
                               └── Garmin scraper (cron, hourly)

GitHub Actions (macOS runner)
    └── flutter build ipa → TestFlight → iPhone/iPad
```

**Why EC2 over ECS?** Single t3.small (~$15/mo) running docker-compose is simplest for a personal project. Can migrate to ECS later if needed.

---

## Phase 0: Do These First (Blocking)

### 0a. Apple Developer Account
1. Go to developer.apple.com/programs/enroll
2. Sign in with your Apple ID
3. Enroll as Individual ($99/yr)
4. Wait 24-48h for approval email
- **This blocks all iOS work — start immediately**

### 0b. Purchase a Domain
- Recommended: AWS Route53 → register domain (e.g. `healthassistant.app`, `yourinitials-health.com`)
- Or any registrar (Namecheap, Google Domains) and point DNS to Route53 hosted zone later
- ~$12-15/yr

---

## Phase 1: Google OAuth (Production)

Currently backend runs with `GOOGLE_CLIENT_ID=""` (dev mode, skips audience check). Prod needs real credentials.

### Steps
1. **Google Cloud Console** → select your project (or create one)
2. APIs & Services → Credentials → Create OAuth 2.0 Client ID
   - Type: **iOS**
   - Bundle ID: `com.yourdomain.healthassistant` (must match `mobile_app/ios/Runner.xcodeproj`)
3. Download `GoogleService-Info.plist` → place in `mobile_app/ios/Runner/` (gitignored)
4. Note the **client ID** from the plist → this becomes `GOOGLE_CLIENT_ID` env var on backend
5. APIs & Services → OAuth consent screen → add your domain to authorized domains (once you have it)

### Files affected
- `mobile_app/ios/Runner/GoogleService-Info.plist` (gitignored, download from console)
- Backend env var: `GOOGLE_CLIENT_ID=<ios-client-id>`

---

## Phase 2: AWS Infrastructure

### 2a. EC2 Instance
```bash
# Launch t3.small, Ubuntu 22.04 LTS
# Security group inbound rules:
#   22  (SSH)    — your IP only
#   80  (HTTP)   — 0.0.0.0/0
#   443 (HTTPS)  — 0.0.0.0/0
```
1. Launch t3.small in us-east-1 (or nearest region)
2. Allocate Elastic IP → attach to instance
3. Install Docker + docker-compose on instance

### 2b. Domain DNS
1. Create Route53 hosted zone for your domain
2. Add A record: `api.yourdomain.com` → Elastic IP
3. Update domain registrar NS records to Route53 (if domain not in Route53)

### 2c. SSL Certificate (Nginx + Let's Encrypt)
```bash
# On EC2
sudo apt install nginx certbot python3-certbot-nginx
sudo certbot --nginx -d api.yourdomain.com
# Auto-renews via systemd timer
```

### 2d. SSM Parameter Store (secrets)
Store secrets in AWS SSM (free), not in env files:
```bash
aws ssm put-parameter --name /health-assistant/DATABASE_URL --value "..." --type SecureString
aws ssm put-parameter --name /health-assistant/JWT_SECRET --value "..." --type SecureString
aws ssm put-parameter --name /health-assistant/GOOGLE_CLIENT_ID --value "..." --type SecureString
aws ssm put-parameter --name /health-assistant/GARMIN_INGEST_SECRET --value "..." --type SecureString
aws ssm put-parameter --name /health-assistant/GARMIN_USERNAME --value "..." --type SecureString
aws ssm put-parameter --name /health-assistant/GARMIN_PASSWORD --value "..." --type SecureString
```

---

## Phase 3: Backend Deployment

### 3a. Dockerfile (create if not exists)
`backend/Dockerfile` — multi-stage build, scratch or alpine final image.

### 3b. docker-compose.prod.yml on EC2
```yaml
services:
  db:
    image: timescale/timescaledb:latest-pg15
    volumes: [pgdata:/var/lib/postgresql/data]
    environment:
      POSTGRES_PASSWORD: ${DB_PASSWORD}
  backend:
    image: <ecr-url>/health-assistant-backend:latest
    ports: ["127.0.0.1:8083:8083"]  # only localhost; Nginx proxies
    environment:
      DATABASE_URL: ${DATABASE_URL}
      JWT_SECRET: ${JWT_SECRET}
      GOOGLE_CLIENT_ID: ${GOOGLE_CLIENT_ID}
      GARMIN_INGEST_SECRET: ${GARMIN_INGEST_SECRET}
    depends_on: [db]
```

### 3c. GitHub Actions: Backend CI/CD
`.github/workflows/deploy-backend.yml`
- Trigger: push to `main` with changes in `backend/`
- Steps: Go build → Docker build → push to ECR → SSH to EC2 → `docker compose pull && docker compose up -d`

### 3d. Run Migrations
```bash
# On EC2, one-time and after each migration
psql $DATABASE_URL -f scripts/db/migrations/001_*.sql
# ... etc
```

---

## Phase 4: Garmin Scraper on AWS

Simplest approach: **cron job on the same EC2 instance**.

### Steps
1. Create Python venv on EC2, install scraper dependencies
2. Create `/opt/garmin-sync/` directory with scraper + `.env` (pulls from SSM)
3. Add crontab entry:
```cron
0 * * * * cd /opt/garmin-sync && python garmin_sync.py >> /var/log/garmin-sync.log 2>&1
```
4. Monitor via `tail -f /var/log/garmin-sync.log`

### Future improvement
- Move to EventBridge + ECS task or Lambda (better observability, no EC2 dependency)

---

## Phase 5: iOS Build & TestFlight

Since you're on Linux, builds happen in **GitHub Actions macOS runners**.

### 5a. Apple Developer Setup (after account approved)
1. Xcode → Preferences → Accounts → add Apple ID
   - **OR** do this via App Store Connect web + GitHub Actions secrets (preferred for CI)
2. Create App ID: `com.yourdomain.healthassistant`
3. Create Distribution Certificate (`.p12`) + Provisioning Profile (`.mobileprovision`)
4. Export both → store as GitHub Secrets (base64 encoded)

### 5b. GitHub Actions: iOS CI/CD
`.github/workflows/deploy-ios.yml`
- Trigger: manual dispatch (or push to `main` with `mobile_app/` changes)
- Runner: `macos-latest`
- Steps:
  1. `flutter pub get`
  2. Decode + install signing cert + provisioning profile
  3. Set `GoogleService-Info.plist` from GitHub Secret
  4. Set `API_BASE_URL=https://api.yourdomain.com` via `--dart-define`
  5. `flutter build ipa --release`
  6. Upload to TestFlight via `xcrun altool` or `fastlane pilot`

### 5c. TestFlight Distribution
- App Store Connect → add yourself as internal tester
- Install TestFlight app on iPhone/iPad
- Accept invite → install app

### 5d. App Store (later, optional)
- Add screenshots, description, privacy policy URL
- Submit for review (~1-2 days)

---

## Phase 6: Flutter App Config Update

### `lib/core/config/app_config.dart`
```dart
static const String baseUrl = String.fromEnvironment(
  'API_BASE_URL',
  defaultValue: 'http://localhost:8083',  // dev
);
```
Pass `--dart-define=API_BASE_URL=https://api.yourdomain.com` at build time.
No code changes needed — already set up correctly.

---

## Effort & Cost Estimate

| Item | One-time Cost | Ongoing/mo |
|------|--------------|------------|
| Apple Developer account | $99/yr | — |
| Domain | ~$12/yr | — |
| EC2 t3.small | — | ~$15 |
| Elastic IP (attached) | — | Free |
| ACM / Let's Encrypt | Free | Free |
| ECR | — | ~$0.10 |
| SSM Parameter Store | Free | Free |
| GitHub Actions macOS | Free (2000 min/mo) | Free |
| **Total** | **~$111 first year** | **~$15/mo** |

---

## Recommended Order of Work

1. **Today:** Sign up for Apple Developer account (unblocks everything iOS)
2. **Today:** Purchase domain
3. **Week 1:** Google OAuth production setup + EC2 + backend deployment
4. **Week 1:** Garmin scraper cron on EC2
5. **Week 2:** GitHub Actions iOS workflow + TestFlight
6. **Week 2:** Install on iPhone/iPad via TestFlight ✅

---

## Open Questions

1. What bundle ID do you want? (e.g. `com.satishthakur.healthassistant`) — needed for Apple ID + Google OAuth
2. What domain? (needed for OAuth consent screen + SSL)
3. Do you want a separate `staging` environment or go straight to prod?

---

**Outcome:** App installed on your iPhone/iPad via TestFlight, backend on AWS at `api.yourdomain.com`, Garmin syncing hourly.
