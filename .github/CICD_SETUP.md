# CI/CD Setup Guide

This document explains how to set up GitHub Actions for automated testing and deployment.

## Overview

We have three main workflows:

1. **CI (Continuous Integration)** - Runs on every PR
2. **Integration Tests** - Runs on main branch only (requires secrets)
3. **Deploy** - Deploys to AWS after successful tests

## 🔐 Required GitHub Secrets

### For Integration Tests

Go to: **Settings → Secrets and variables → Actions → New repository secret**

| Secret Name | Description | Example |
|------------|-------------|---------|
| `GARMIN_TEST_EMAIL` | Test Garmin account email | `test@example.com` |
| `GARMIN_TEST_PASSWORD` | Test Garmin account password | `securepassword123` |

**⚠️ Important:**
- Use a **dedicated test account**, not your personal Garmin account
- This account should have sample data for testing
- Credentials are encrypted and only accessible to workflows in your repo (not forks)

### For AWS Deployment

| Secret Name | Description | How to Get |
|------------|-------------|-----------|
| `AWS_ACCESS_KEY_ID` | AWS access key | IAM User credentials |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | IAM User credentials |
| `AWS_ACCOUNT_ID` | AWS account ID | Top right in AWS Console |
| `INGESTION_SERVICE_URL` | Production ingestion URL | After ECS deployment |
| `SCHEDULER_SERVICE_URL` | Production scheduler URL | After ECS deployment |

## 📋 Workflows Explained

### 1. CI Workflow (`.github/workflows/ci.yml`)

**Triggers:** Every PR and push to main/develop

**What it does:**
- ✅ Runs Go unit tests with race detection
- ✅ Runs Go vet (linting)
- ✅ Builds all 4 Go services
- ✅ Runs Python tests with pytest
- ✅ Checks Python code formatting (black)
- ✅ Runs Python linting (flake8)
- ✅ Validates Docker builds
- ✅ Tests SQL schema and migrations
- ✅ Uploads code coverage to Codecov

**No secrets required** - Safe for public repos and forks

### 2. Integration Test Workflow (`.github/workflows/integration.yml`)

**Triggers:**
- Push to main branch
- Manual trigger via Actions tab

**What it does:**
- ✅ Creates .env with Garmin credentials
- ✅ Starts all Docker containers
- ✅ Waits for services to be healthy
- ✅ Runs full integration test script
- ✅ Validates data was synced
- ✅ Creates GitHub issue if tests fail

**Security:**
- Only runs on your repository (not forks)
- Credentials are encrypted secrets
- Automatically skipped if secrets are missing

### 3. Deploy Workflow (`.github/workflows/deploy.yml`)

**Triggers:**
- Push to main (after integration tests pass)
- Manual trigger with environment selection

**What it does:**
- 🏗️ Builds Docker images
- 📦 Pushes to AWS ECR
- 🚀 Updates ECS task definitions
- ⏳ Waits for deployment to stabilize
- 🏥 Runs smoke tests
- 💬 Comments on related PRs

**Requires:**
- AWS credentials
- ECS cluster and services set up
- ECR repositories created

## 🚀 Setup Steps

### Step 1: Add Integration Test Secrets

```bash
# Go to your GitHub repo:
https://github.com/satishThakur/health-assistant/settings/secrets/actions

# Click "New repository secret" and add:
Name: GARMIN_TEST_EMAIL
Value: your_test_email@example.com

Name: GARMIN_TEST_PASSWORD
Value: your_test_password
```

### Step 2: Enable GitHub Actions

1. Go to **Settings → Actions → General**
2. Under **Actions permissions**, select:
   - "Allow all actions and reusable workflows"
3. Under **Workflow permissions**, select:
   - "Read and write permissions"
   - Check "Allow GitHub Actions to create and approve pull requests"

### Step 3: Verify CI Works

```bash
# Create a test branch and PR
git checkout -b test/ci-setup
git commit --allow-empty -m "test: Trigger CI"
git push origin test/ci-setup

# Create PR and watch Actions tab
```

You should see:
- ✅ Go Tests & Linting
- ✅ Go Build Validation
- ✅ Python Tests & Linting
- ✅ Docker Build Validation
- ✅ SQL Schema Validation

### Step 4: Test Integration Workflow

```bash
# Merge PR to main (or push directly)
git checkout main
git merge test/ci-setup
git push origin main

# Watch Actions tab for "Integration Tests" workflow
```

If secrets are configured correctly, integration tests will run.

## 🔧 AWS Deployment Setup (Optional)

### Prerequisites

1. **AWS Account** with appropriate permissions
2. **ECR Repositories** created:
   ```bash
   aws ecr create-repository --repository-name health-assistant-ingestion
   aws ecr create-repository --repository-name health-assistant-scheduler
   ```

3. **ECS Cluster** and services:
   ```bash
   # Create cluster
   aws ecs create-cluster --cluster-name health-assistant-cluster

   # Create task definitions and services (see infra/aws/ for templates)
   ```

4. **IAM User** for GitHub Actions with policies:
   - `AmazonEC2ContainerRegistryPowerUser`
   - `AmazonECS_FullAccess`
   - Custom policy for task definitions

### Add AWS Secrets

```bash
# In GitHub repo secrets, add:
AWS_ACCESS_KEY_ID=AKIA...
AWS_SECRET_ACCESS_KEY=...
AWS_ACCOUNT_ID=123456789012
INGESTION_SERVICE_URL=https://ingestion.your-domain.com
SCHEDULER_SERVICE_URL=https://scheduler.your-domain.com
```

### Trigger Deployment

```bash
# Automatic: Push to main after tests pass
git push origin main

# Manual: Go to Actions → Deploy to AWS → Run workflow
# Select environment: staging or production
```

## 📊 Monitoring Workflows

### View Workflow Runs

```
https://github.com/satishThakur/health-assistant/actions
```

### Check Individual Jobs

Click on any workflow run → Click on job name → View logs

### Re-run Failed Workflows

Click "Re-run failed jobs" or "Re-run all jobs"

## 🐛 Troubleshooting

### CI Tests Fail on Fork

**Expected behavior.** Forks don't have access to repository secrets, so integration tests are automatically skipped.

### Integration Tests Fail

1. **Check secrets are set:**
   ```
   Settings → Secrets → Actions → GARMIN_TEST_EMAIL exists?
   ```

2. **Check test account:**
   - Can you log in to garminconnect.com?
   - Does the account have recent activity data?

3. **View logs:**
   - Go to failed workflow run
   - Click "integration-test" job
   - Expand "Collect logs on failure"

### Deployment Fails

1. **Verify AWS credentials:**
   ```bash
   # Test locally
   aws sts get-caller-identity
   ```

2. **Check ECS services exist:**
   ```bash
   aws ecs list-services --cluster health-assistant-cluster
   ```

3. **Check ECR repositories:**
   ```bash
   aws ecr describe-repositories
   ```

## 📝 Best Practices

### For Contributors

1. **Run tests locally first:**
   ```bash
   # Go tests
   cd backend && go test ./...

   # Python tests
   cd services/garmin-scheduler && pytest

   # Integration test
   ./scripts/test-integration.sh
   ```

2. **Keep CI green:** Fix failing tests before merging

3. **Don't commit secrets:** Use environment variables and .env files (gitignored)

### For Maintainers

1. **Review integration test failures** before merging to main

2. **Rotate test credentials** periodically

3. **Monitor AWS costs** - Integration tests run Docker containers

4. **Use staging environment** for testing deployments before production

## 🔄 Workflow Diagram

```
┌─────────────────────────────────────────────────────────┐
│                     PR Created                          │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
         ┌───────────────────────┐
         │   CI Workflow         │
         │   - Unit Tests        │
         │   - Linting           │
         │   - Build Validation  │
         │   - SQL Tests         │
         └───────────┬───────────┘
                     │
                     ▼
              ✅ Tests Pass
                     │
                     ▼
         ┌───────────────────────┐
         │   PR Merged to Main   │
         └───────────┬───────────┘
                     │
         ┌───────────┴───────────┐
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌─────────────────┐
│ Integration     │    │ Deploy to AWS   │
│ Tests           │    │ (after tests)   │
│ (with secrets)  │    │                 │
└────────┬────────┘    └────────┬────────┘
         │                      │
         ▼                      ▼
   ✅ Data Synced         🚀 Deployed
```

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [AWS ECS Deployment Guide](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/deploy-cloudformation.html)
- [Docker Multi-Stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [Codecov Setup](https://docs.codecov.com/docs/quick-start)

## 🆘 Need Help?

- Check [scripts/README.md](../scripts/README.md) for local testing
- Review [QUICK_START.md](../QUICK_START.md) for manual setup
- Open an issue with workflow logs attached
