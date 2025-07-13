# 🚀 Deployment Guide - Vietnam Admin API

This guide covers the complete CI/CD setup for the Vietnam Admin API using GitHub Actions and webhook deployment.

## 📋 Table of Contents
- [Overview](#overview)
- [GitHub Actions CI/CD Setup](#github-actions-cicd-setup)
- [Webhook Deployment Setup](#webhook-deployment-setup)
- [Environment Configuration](#environment-configuration)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)

## 🔍 Overview

The deployment system consists of:
- **GitHub Actions**: CI/CD pipeline for automated testing, building, and deployment
- **Webhook Service**: Automated deployment service that responds to GitHub webhooks
- **Docker**: Containerization for consistent deployment across environments
- **Multi-environment support**: Staging and production environments

## 🛠️ GitHub Actions CI/CD Setup

### 1. Workflow Files

Two main workflow files are included:

#### `/.github/workflows/ci-cd.yml`
- **Purpose**: Complete CI/CD pipeline with testing, building, and deployment
- **Triggers**: Push to `main` and `develop` branches, pull requests to `main`
- **Features**:
  - Go testing and linting
  - Binary building with artifact upload
  - Docker image building and pushing to GitHub Container Registry
  - Automated deployment to staging/production
  - Webhook connectivity testing

#### `/.github/workflows/webhook-deploy.yml`
- **Purpose**: Dedicated webhook deployment workflow
- **Triggers**: Push to `main` and `develop` branches, manual dispatch
- **Features**:
  - Environment-specific deployment
  - Proper webhook payload generation
  - Health check verification
  - Deployment summary reporting

### 2. Repository Secrets Configuration

Configure these secrets in your GitHub repository (Settings > Secrets and variables > Actions):

```bash
# Required secrets
WEBHOOK_URL=https://webhook1.iceteadev.site/deploy
WEBHOOK_SECRET=du_an_cua_tuan

# Optional: For custom registry
REGISTRY_USERNAME=your_registry_username
REGISTRY_PASSWORD=your_registry_password
```

### 3. Environment Protection Rules

Set up environment protection rules in GitHub:

1. Go to **Settings** > **Environments**
2. Create environments: `staging` and `production`
3. Configure protection rules:
   - **Staging**: No restrictions (auto-deploy on develop branch)
   - **Production**: Require reviews, restrict to main branch

## 🔗 Webhook Deployment Setup

### 1. Webhook Server Configuration

Configure the webhook server with environment variables following the pattern from `/deployment/webhook-config.env`:

```bash
# Production environment
DEPLOY_COMMANDS_OWNER_VIETNAM_ADMIN_API=cd /opt/vietnam-admin-api;git pull origin main;docker-compose down;docker-compose pull;docker-compose up -d --build;docker system prune -f
WORK_DIR_OWNER_VIETNAM_ADMIN_API=/opt/vietnam-admin-api

# Staging environment
DEPLOY_COMMANDS_OWNER_VIETNAM_ADMIN_API_STAGING=cd /opt/vietnam-admin-api-staging;git pull origin develop;docker-compose -f docker-compose.staging.yml down;docker-compose -f docker-compose.staging.yml pull;docker-compose -f docker-compose.staging.yml up -d --build;docker system prune -f
WORK_DIR_OWNER_VIETNAM_ADMIN_API_STAGING=/opt/vietnam-admin-api-staging

# Webhook server configuration
WEBHOOK_SECRET=du_an_cua_tuan
PORT=8300
```

### 2. Server Directory Setup

Create deployment directories on your server:

```bash
# Create directories
sudo mkdir -p /opt/vietnam-admin-api
sudo mkdir -p /opt/vietnam-admin-api-staging

# Set ownership
sudo chown -R $USER:$USER /opt/vietnam-admin-api
sudo chown -R $USER:$USER /opt/vietnam-admin-api-staging

# Clone repositories
git clone https://github.com/your-org/vietnam-admin-api.git /opt/vietnam-admin-api
git clone https://github.com/your-org/vietnam-admin-api.git /opt/vietnam-admin-api-staging

# Setup staging branch
cd /opt/vietnam-admin-api-staging
git checkout -b develop origin/develop
```

### 3. GitHub Webhook Configuration

1. Go to your repository on GitHub
2. Navigate to **Settings** > **Webhooks**
3. Click **Add webhook**
4. Configure:
   - **Payload URL**: `https://webhook1.iceteadev.site/`
   - **Content type**: `application/json`
   - **Secret**: Same as `WEBHOOK_SECRET` in server config
   - **Events**: Select "Just the push event"
   - **Active**: ✅ Checked

## 🌍 Environment Configuration

### Production Environment
- **Branch**: `main`
- **Port**: `8100`
- **Docker Compose**: `docker-compose.yml`
- **Domain**: `api.vietnam-admin.your-domain.com`
- **Container**: `vietnam-admin-api`

### Staging Environment
- **Branch**: `develop`
- **Port**: `8101`
- **Docker Compose**: `docker-compose.staging.yml`
- **Domain**: `staging-api.vietnam-admin.your-domain.com`
- **Container**: `vietnam-admin-api-staging`

## 🧪 Testing

### 1. Test Webhook Connectivity

Use the provided test scripts to verify webhook functionality:

#### Using Go Script:
```bash
# Test staging deployment
go run scripts/test-webhook.go https://webhook1.iceteadev.site/ your_webhook_secret staging

# Test production deployment  
go run scripts/test-webhook.go https://webhook1.iceteadev.site/ your_webhook_secret production
```

#### Using Bash Script:
```bash
# Make script executable (Linux/macOS)
chmod +x scripts/test-webhook.sh

# Test staging deployment
./scripts/test-webhook.sh https://webhook1.iceteadev.site/ your_webhook_secret staging

# Test production deployment
./scripts/test-webhook.sh https://webhook1.iceteadev.site/ your_webhook_secret production
```

### 2. Manual Deployment Trigger

Trigger deployment manually from GitHub:

1. Go to **Actions** tab
2. Select **Webhook Deployment** workflow
3. Click **Run workflow**
4. Select environment and branch
5. Click **Run workflow**

### 3. API Health Check

After deployment, verify the API is running:

```bash
# Production
curl -f https://api.vietnam-admin.your-domain.com/health

# Staging
curl -f https://staging-api.vietnam-admin.your-domain.com/health

# Local/Direct
curl -f http://your-server-ip:8100/health
curl -f http://your-server-ip:8101/health  # Staging
```

## 🔍 Troubleshooting

### Common Issues

#### 1. Webhook Not Triggering
- **Check webhook URL**: Ensure `https://webhook1.iceteadev.site/` is accessible
- **Verify secret**: Ensure `WEBHOOK_SECRET` matches between GitHub and server
- **Check IP restrictions**: Verify GitHub IPs are allowed
- **Review payload**: Check GitHub webhook delivery logs

#### 2. Deployment Failures
- **Check server logs**: Review webhook server logs for errors
- **Verify permissions**: Ensure deployment user has proper permissions
- **Check Docker**: Verify Docker service is running
- **Review environment variables**: Ensure all required vars are set

#### 3. Application Not Starting
- **Check Docker logs**: `docker logs vietnam-admin-api`
- **Verify ports**: Ensure ports 8100/8101 are available
- **Check data files**: Verify `data/province.json` and `data/ward.json` exist
- **Review environment**: Check environment variables in compose file

### Debugging Commands

```bash
# Check webhook server status
systemctl status webhook-server

# View webhook server logs
journalctl -u webhook-server -f

# Check application logs
docker logs vietnam-admin-api -f
docker logs vietnam-admin-api-staging -f

# Check Docker Compose status
docker-compose ps
docker-compose -f docker-compose.staging.yml ps

# Test API endpoints
curl -v http://localhost:8100/api/v1/provinces
curl -v http://localhost:8101/api/v1/provinces  # Staging
```

### Log Locations

- **Webhook Server**: `/var/log/webhook-server/`
- **Application**: Docker container logs
- **Nginx**: `/var/log/nginx/`
- **System**: `/var/log/syslog`

## 📊 Monitoring

### Key Metrics to Monitor

1. **Deployment Success Rate**
   - Monitor GitHub Actions workflow success/failure
   - Track webhook response times

2. **Application Health**
   - API response times
   - Error rates
   - Container resource usage

3. **Infrastructure**
   - Server resource usage
   - Docker container status
   - Network connectivity

### Alerting Setup

Consider setting up alerts for:
- Failed deployments
- API downtime
- High error rates
- Resource exhaustion

## 🔄 Rollback Procedures

### Quick Rollback Options

1. **Revert Git Commit**:
   ```bash
   git revert HEAD
   git push origin main  # Triggers automatic redeployment
   ```

2. **Manual Docker Rollback**:
   ```bash
   docker pull ghcr.io/your-org/vietnam-admin-api:previous-tag
   docker-compose down
   docker-compose up -d
   ```

3. **Emergency Stop**:
   ```bash
   docker-compose down
   # Fix issues, then restart
   docker-compose up -d
   ```

## 🚀 Advanced Configuration

### Custom Deployment Strategies

Modify `/deployment/webhook-config.env` for different deployment approaches:

- **Binary Deployment**: Deploy Go binary directly
- **Kubernetes**: Use kubectl for K8s deployments
- **Blue-Green**: Implement blue-green deployment strategy
- **Canary**: Progressive rollout configuration

### Security Enhancements

1. **Network Security**:
   - Use VPN for webhook communication
   - Implement IP whitelisting
   - Use SSL/TLS certificates

2. **Access Control**:
   - Use dedicated deployment user
   - Implement role-based access
   - Regular credential rotation

3. **Audit Logging**:
   - Log all deployment activities
   - Monitor access patterns
   - Set up security alerts

## 📚 Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Webhook API Documentation](webhook.md)
- [API Documentation](API_DOCUMENTATION_ADDRESS.md)

## 🤝 Support

For deployment issues:
1. Check this troubleshooting guide
2. Review server logs
3. Test webhook connectivity
4. Verify environment configuration
5. Contact system administrator

---

**Last Updated**: $(date)
**Version**: 1.0.0 