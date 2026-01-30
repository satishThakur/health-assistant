# GitHub Actions & CI/CD

This directory contains automated workflows for testing and deployment.

## 📁 Directory Structure

```
.github/
├── workflows/
│   ├── ci.yml           # Runs on every PR (no secrets needed)
│   ├── integration.yml  # Runs on main branch (requires secrets)
│   └── deploy.yml       # Deploys to AWS (requires AWS credentials)
├── CICD_SETUP.md        # Detailed setup instructions
└── README.md            # This file
```

## 🚀 Quick Start

### 1. For Contributors (PRs)

Just create a PR! The CI workflow will automatically:
- Run all tests
- Check code formatting
- Validate builds
- Test SQL migrations

**No configuration needed** - works out of the box.

### 2. For Maintainers (Integration Tests)

Add these secrets to run integration tests on main branch:

```
Settings → Secrets → Actions:
- GARMIN_TEST_EMAIL
- GARMIN_TEST_PASSWORD
```

See [CICD_SETUP.md](CICD_SETUP.md) for detailed instructions.

### 3. For Deployment (AWS)

Additional secrets required:
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_ACCOUNT_ID`
- `INGESTION_SERVICE_URL`
- `SCHEDULER_SERVICE_URL`

## 📊 Workflow Status

Check workflow runs: https://github.com/satishThakur/health-assistant/actions

### CI Workflow
- ✅ Go tests, linting, builds
- ✅ Python tests, linting
- ✅ Docker builds
- ✅ SQL schema validation

### Integration Workflow
- 🔒 Requires Garmin test credentials
- 🔍 Tests end-to-end data flow
- 📊 Validates sync audit logs

### Deploy Workflow
- 🏗️ Builds & pushes Docker images
- 🚀 Deploys to AWS ECS
- 🏥 Runs smoke tests

## 🛠️ Local Testing

Before pushing, test locally:

```bash
# Run Go tests
cd backend && go test ./...

# Run Python tests
cd services/garmin-scheduler && pytest

# Run integration test
./scripts/test-integration.sh
```

## 📚 Documentation

- **[CICD_SETUP.md](CICD_SETUP.md)** - Complete CI/CD setup guide
- **[../scripts/README.md](../scripts/README.md)** - Integration test script docs
- **[../QUICK_START.md](../QUICK_START.md)** - Manual setup reference

## 🔐 Security Notes

1. **Secrets are encrypted** and only accessible to workflows in this repo
2. **Forks don't have secret access** - integration tests are skipped on forks
3. **Use test accounts** for integration tests, not production credentials
4. **Review workflow changes** in PRs carefully

## 🐛 Common Issues

### "Integration Tests" workflow doesn't run

- **On fork?** Integration tests only run on main repository
- **Secrets set?** Check Settings → Secrets → Actions
- **On main branch?** Integration tests only run on main

### Tests pass locally but fail in CI

- **Different environment:** CI uses fresh Ubuntu container
- **Missing dependencies:** Check workflow file for installed packages
- **Cache issues:** Try re-running workflow

### Deployment fails

- **AWS credentials:** Verify IAM permissions
- **ECS services:** Check services exist in cluster
- **ECR repositories:** Ensure repositories are created

## 💡 Tips

- **Use draft PRs** to test CI without notifying reviewers
- **Re-run failed jobs** instead of pushing empty commits
- **Check workflow logs** for detailed error messages
- **Test integration locally** with `./scripts/test-integration.sh`

## 🔄 Continuous Improvement

These workflows are designed to:
- ✅ Catch bugs before they reach main
- ✅ Ensure code quality and consistency
- ✅ Automate repetitive tasks
- ✅ Enable confident deployments
- ✅ Provide fast feedback to developers

Feel free to propose improvements via PR!
