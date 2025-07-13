# 🐧 Tác động lên VPS Debian

## 🔍 **Tổng quan tác động**

### ✅ **Hoàn toàn tương thích:**
- **Docker containers**: Debian hỗ trợ đầy đủ Docker
- **GitHub Actions**: Chạy trên GitHub servers, không ảnh hưởng VPS
- **Webhook HTTP requests**: Hoạt động bình thường
- **Go applications**: Tương thích 100%
- **JSON API endpoints**: Không có tác động

### ⚠️ **Cần điều chỉnh:**
- **Đường dẫn deployment**: `/opt/` → `/home/deploy/`
- **Package manager**: `yum`/`dnf` → `apt`
- **Service management**: Sử dụng `systemd` (đã tương thích)
- **User permissions**: Cần tạo user `deploy` riêng

## 📦 **Dependencies cần cài đặt:**

```bash
# Cơ bản
sudo apt update && sudo apt upgrade -y
sudo apt install -y git curl wget jq openssl build-essential

# Docker
sudo apt install -y docker-ce docker-ce-cli containerd.io
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Security
sudo apt install -y ufw fail2ban unattended-upgrades
```

## 🔧 **Thay đổi cấu hình:**

### 1. **Webhook Environment Variables**
```bash
# Thay đổi từ:
WORK_DIR_OWNER_VIETNAM_ADMIN_API=/opt/vietnam-admin-api

# Thành:
WORK_DIR_OWNER_VIETNAM_ADMIN_API=/home/deploy/vietnam-admin-api
```

### 2. **Deployment Commands**
```bash
# Debian-specific path
DEPLOY_COMMANDS_OWNER_VIETNAM_ADMIN_API=cd /home/deploy/vietnam-admin-api;git pull origin main;docker-compose down;docker-compose pull;docker-compose up -d --build;docker system prune -f
```

### 3. **Service Configuration**
```bash
# Tạo systemd service
sudo nano /etc/systemd/system/webhook-server.service

[Unit]
Description=Webhook Deployment Server
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
User=deploy
WorkingDirectory=/home/deploy
ExecStart=/usr/local/bin/webhook-server
Restart=always
RestartSec=3
Environment=PORT=8300
EnvironmentFile=/etc/environment

[Install]
WantedBy=multi-user.target
```

## 🛡️ **Security Requirements:**

### 1. **Firewall (UFW)**
```bash
sudo ufw allow 8300/tcp  # Webhook server
sudo ufw allow 8100/tcp  # Production API
sudo ufw allow 8101/tcp  # Staging API
sudo ufw enable
```

### 2. **User Setup**
```bash
# Tạo user deployment riêng
sudo useradd -m -s /bin/bash deploy
sudo usermod -aG docker deploy
```

### 3. **SSH Security**
```bash
# Disable root login
sudo nano /etc/ssh/sshd_config
# PermitRootLogin no
# PasswordAuthentication no
```

## 🚨 **Lưu ý quan trọng:**

### 1. **Resource Requirements**
- **RAM**: Tối thiểu 2GB (khuyến nghị 4GB)
- **Storage**: Tối thiểu 20GB (Docker images + logs)
- **Network**: Ports 8300, 8100, 8101 cần mở

### 2. **Performance Optimization**
```bash
# Cấu hình swap
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Docker log rotation
sudo nano /etc/docker/daemon.json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```

### 3. **Monitoring Setup**
```bash
# Health check script
#!/bin/bash
if ! curl -f http://localhost:8300/health > /dev/null 2>&1; then
    sudo systemctl restart webhook-server
fi

if ! docker ps | grep -q vietnam-admin-api; then
    cd /home/deploy/vietnam-admin-api
    docker-compose up -d
fi
```

## 🔄 **Backup Strategy:**

```bash
# Daily backup script
#!/bin/bash
BACKUP_DIR="/home/deploy/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Backup application
tar -czf $BACKUP_DIR/vietnam-admin-api_$DATE.tar.gz -C /home/deploy vietnam-admin-api

# Backup Docker volumes
docker run --rm -v vietnam-admin-api_data:/data -v $BACKUP_DIR:/backup ubuntu tar czf /backup/volumes_$DATE.tar.gz -C /data .

# Auto-cleanup (7 days)
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

## 📊 **Testing Commands:**

```bash
# Test webhook connectivity
curl -X POST http://localhost:8300/health

# Test API endpoints
curl -f http://localhost:8100/api/v1/provinces
curl -f http://localhost:8101/api/v1/provinces  # Staging

# Check container status
docker ps -a
docker logs vietnam-admin-api
docker logs vietnam-admin-api-staging

# Check system resources
df -h
free -h
docker system df
```

## 🚀 **Deployment Process:**

1. **Setup server** (theo `deployment/debian-setup.md`)
2. **Configure webhook** (theo `deployment/webhook-config.env`)
3. **Test connectivity** (dùng `scripts/test-webhook.sh`)
4. **Deploy application** (GitHub Actions trigger)
5. **Verify deployment** (health checks)

## 📋 **Checklist sau khi setup:**

- [ ] Docker đã cài đặt và chạy
- [ ] User `deploy` đã tạo với quyền docker
- [ ] Webhook server service đã chạy
- [ ] Firewall đã cấu hình ports
- [ ] SSH security đã setup
- [ ] Backup script đã cấu hình
- [ ] Health monitoring đã setup
- [ ] Test webhook thành công
- [ ] Application containers chạy được
- [ ] API endpoints trả về kết quả

## 🔗 **Liên kết tài liệu:**

- [Debian Setup Guide](debian-setup.md)
- [Webhook Configuration](webhook-config.env)
- [Main Deployment Guide](../DEPLOYMENT.md)
- [Test Scripts](../scripts/)

---

**Kết luận**: Hệ thống CI/CD hoàn toàn tương thích với Debian VPS, chỉ cần điều chỉnh một số đường dẫn và cài đặt dependencies cơ bản. Không có tác động tiêu cực nào đến hiệu suất hay bảo mật của VPS. 