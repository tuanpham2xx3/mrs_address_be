# 🐧 Debian VPS Setup Guide

Hướng dẫn cấu hình VPS Debian cho webhook deployment system.

## 📦 **Cài đặt Dependencies**

### 1. Cập nhật hệ thống
```bash
sudo apt update && sudo apt upgrade -y
```

### 2. Cài đặt Docker
```bash
# Cài đặt dependencies
sudo apt install -y apt-transport-https ca-certificates curl gnupg lsb-release

# Thêm Docker GPG key
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Thêm Docker repository
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Cài đặt Docker
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io

# Cài đặt Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Thêm user vào docker group
sudo usermod -aG docker $USER
```

### 3. Cài đặt Git và các tools cần thiết
```bash
sudo apt install -y git curl wget jq openssl build-essential
```

### 4. Cài đặt Go (nếu cần)
```bash
# Download Go
wget https://golang.org/dl/go1.21.5.linux-amd64.tar.gz

# Cài đặt Go
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Thêm vào PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

## 🔧 **Cấu hình Webhook Server**

### 1. Tạo user deployment
```bash
# Tạo user riêng cho deployment
sudo useradd -m -s /bin/bash deploy
sudo usermod -aG docker deploy

# Chuyển sang user deploy
sudo su - deploy
```

### 2. Tạo thư mục deployment
```bash
# Tạo thư mục
mkdir -p /home/deploy/vietnam-admin-api
mkdir -p /home/deploy/vietnam-admin-api-staging

# Clone repository
git clone https://github.com/your-org/vietnam-admin-api.git /home/deploy/vietnam-admin-api
git clone https://github.com/your-org/vietnam-admin-api.git /home/deploy/vietnam-admin-api-staging

# Setup staging branch
cd /home/deploy/vietnam-admin-api-staging
git checkout -b develop origin/develop
```

### 3. Cấu hình environment variables
```bash
# Tạo file environment
sudo nano /etc/environment

# Thêm vào file:
DEPLOY_COMMANDS_OWNER_VIETNAM_ADMIN_API="cd /home/deploy/vietnam-admin-api;git pull origin main;docker-compose down;docker-compose pull;docker-compose up -d --build;docker system prune -f"
WORK_DIR_OWNER_VIETNAM_ADMIN_API="/home/deploy/vietnam-admin-api"
DEPLOY_COMMANDS_OWNER_VIETNAM_ADMIN_API_STAGING="cd /home/deploy/vietnam-admin-api-staging;git pull origin develop;docker-compose -f docker-compose.staging.yml down;docker-compose -f docker-compose.staging.yml pull;docker-compose -f docker-compose.staging.yml up -d --build;docker system prune -f"
WORK_DIR_OWNER_VIETNAM_ADMIN_API_STAGING="/home/deploy/vietnam-admin-api-staging"
WEBHOOK_SECRET="your_webhook_secret_here"
PORT="8300"
```

## 🔒 **Cấu hình Security**

### 1. Cấu hình Firewall (UFW)
```bash
# Cài đặt UFW
sudo apt install -y ufw

# Cấu hình rules
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 8300/tcp  # Webhook server
sudo ufw allow 8100/tcp  # Production API
sudo ufw allow 8101/tcp  # Staging API

# Kích hoạt firewall
sudo ufw enable
```

### 2. Cấu hình SSH Security
```bash
# Backup SSH config
sudo cp /etc/ssh/sshd_config /etc/ssh/sshd_config.bak

# Chỉnh sửa SSH config
sudo nano /etc/ssh/sshd_config

# Thêm/sửa các dòng:
# PermitRootLogin no
# PasswordAuthentication no
# PubkeyAuthentication yes
# Port 2222  # Đổi port SSH

# Restart SSH
sudo systemctl restart sshd
```

## 🚀 **Cấu hình Webhook Server**

### 1. Tạo systemd service
```bash
sudo nano /etc/systemd/system/webhook-server.service
```

```ini
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

### 2. Kích hoạt service
```bash
sudo systemctl daemon-reload
sudo systemctl enable webhook-server
sudo systemctl start webhook-server
```

## 📊 **Monitoring và Logs**

### 1. Cấu hình log rotation
```bash
sudo nano /etc/logrotate.d/webhook-server
```

```
/var/log/webhook-server/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    copytruncate
    su deploy deploy
}
```

### 2. Monitoring scripts
```bash
# Tạo script check health
nano /home/deploy/check-health.sh
```

```bash
#!/bin/bash
# Health check script

# Check webhook server
if ! curl -f http://localhost:8300/health > /dev/null 2>&1; then
    echo "Webhook server is down!"
    sudo systemctl restart webhook-server
fi

# Check application containers
if ! docker ps | grep -q vietnam-admin-api; then
    echo "Vietnam Admin API container is down!"
    cd /home/deploy/vietnam-admin-api
    docker-compose up -d
fi
```

## 🔄 **Backup và Recovery**

### 1. Backup script
```bash
nano /home/deploy/backup.sh
```

```bash
#!/bin/bash
BACKUP_DIR="/home/deploy/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Tạo backup directory
mkdir -p $BACKUP_DIR

# Backup application data
tar -czf $BACKUP_DIR/vietnam-admin-api_$DATE.tar.gz -C /home/deploy vietnam-admin-api

# Backup Docker volumes
docker run --rm -v vietnam-admin-api_data:/data -v $BACKUP_DIR:/backup ubuntu tar czf /backup/volumes_$DATE.tar.gz -C /data .

# Cleanup old backups (giữ lại 7 ngày)
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
```

### 2. Setup cron job
```bash
crontab -e

# Thêm dòng:
0 2 * * * /home/deploy/backup.sh
```

## 🛡️ **Security Best Practices**

### 1. Cấu hình fail2ban
```bash
sudo apt install -y fail2ban

# Cấu hình fail2ban
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
sudo nano /etc/fail2ban/jail.local

# Kích hoạt fail2ban
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

### 2. Cấu hình unattended-upgrades
```bash
sudo apt install -y unattended-upgrades

# Cấu hình auto-updates
sudo dpkg-reconfigure -plow unattended-upgrades
```

## 🚨 **Troubleshooting Commands cho Debian**

```bash
# Check system status
systemctl status webhook-server
systemctl status docker

# Check logs
journalctl -u webhook-server -f
journalctl -u docker -f

# Check Docker
docker ps -a
docker logs vietnam-admin-api
docker system df

# Check network
netstat -tlnp | grep :8300
ss -tlnp | grep :8100

# Check disk space
df -h
du -sh /var/lib/docker/

# Check memory
free -h
top -p $(pgrep -d',' docker)
```

## 🎯 **Optimization cho Debian**

### 1. Docker optimization
```bash
# Cấu hình Docker daemon
sudo nano /etc/docker/daemon.json
```

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "storage-driver": "overlay2"
}
```

### 2. System optimization
```bash
# Cấu hình swap
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Thêm vào /etc/fstab
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

## 📋 **Checklist sau khi setup**

- [ ] Docker và Docker Compose đã cài đặt
- [ ] User deploy đã tạo và có quyền docker
- [ ] Webhook server service đã chạy
- [ ] Firewall đã cấu hình
- [ ] SSH security đã setup
- [ ] Backup script đã cấu hình
- [ ] Monitoring đã setup
- [ ] Test webhook connectivity thành công
- [ ] Application containers chạy được
- [ ] Health checks hoạt động

---

**Lưu ý**: Thay đổi các giá trị như `your-org`, `your_webhook_secret_here` theo thông tin thực tế của bạn. 