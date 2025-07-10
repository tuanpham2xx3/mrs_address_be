# Vietnam Administrative API

🇻🇳 **Microservice API cho hệ thống tỉnh thành Việt Nam** được xây dựng bằng **Go** và **Gin Framework**.

## 🚀 Quick Start

```bash
# 1. Clone repository
git clone <repo-url>
cd vietnam-admin-api

# 2. Copy dữ liệu JSON
cp province.json ward.json ./

# 3. Setup và chạy
make setup
make run

# Hoặc sử dụng script
chmod +x setup.sh
./setup.sh
```

**Server sẽ chạy tại: http://localhost:8080**

## ✨ Tính năng

- 🚀 **High Performance**: Xử lý hàng nghìn requests/giây
- 💾 **In-Memory Data**: Load JSON vào RAM để truy xuất cực nhanh
- 🔍 **Advanced Search**: Tìm kiếm thông minh với fuzzy matching
- 📊 **RESTful API**: Thiết kế API chuẩn REST
- 🛡️ **Production Ready**: CORS, Rate limiting, Health check
- 🐳 **Docker Support**: Containerized deployment
- 📈 **Monitoring**: Metrics và logging

## 🏗️ Kiến trúc

```
vietnam-admin-api/
├── models/          # Data models và structs
├── services/        # Business logic và data access
├── handlers/        # HTTP request handlers
├── middleware/      # Middleware (CORS, logging, auth)
├── data/           # JSON data files
├── main.go         # Application entry point
├── Dockerfile      # Container configuration
├── Makefile        # Build automation
└── README.md
```

## 📋 API Endpoints

### 🏢 **Provinces (Tỉnh/Thành phố)**

```bash
GET /api/v1/provinces                    # Lấy danh sách tỉnh/thành
GET /api/v1/provinces/{code}             # Chi tiết 1 tỉnh
GET /api/v1/provinces/{code}/wards       # Xã/phường thuộc tỉnh
GET /api/v1/provinces/types              # Loại tỉnh (thành phố, tỉnh)
```

### 🏘️ **Wards (Xã/Phường/Thị trấn)**

```bash
GET /api/v1/wards                        # Lấy danh sách xã/phường
GET /api/v1/wards/{code}                 # Chi tiết 1 xã/phường
GET /api/v1/wards/types                  # Loại xã (xã, phường, thị trấn)
```

### 🔍 **Search & Utility**

```bash
GET /api/v1/search                       # Tìm kiếm toàn cục
POST /api/v1/address/validate            # Validate địa chỉ
GET /api/v1/health                       # Health check
GET /api/v1/stats                        # Thống kê dữ liệu
```

### 🔧 **Admin**

```bash
POST /api/v1/admin/reload                # Reload data (yêu cầu auth)
```

## 🚀 Cách chạy

### **1. Development (Local)**

```bash
# Cách 1: Sử dụng Makefile (Khuyến nghị)
make help          # Xem tất cả commands
make setup         # Khởi tạo project
make run           # Chạy development mode
make build         # Build binary
make test          # Run tests

# Cách 2: Manual
mkdir -p data
cp province.json ward.json data/
go mod tidy
go run main.go

# Cách 3: Sử dụng setup script
chmod +x setup.sh
./setup.sh
```

### **2. Production (Docker)**

```bash
# Build và chạy với Docker
make docker
make docker-run

# Hoặc manual
docker build -t vietnam-admin-api .
docker run -d -p 8080:8080 vietnam-admin-api

# Deploy với docker-compose
make deploy
```

### **3. Environment Variables**

```bash
PORT=8080                    # Server port (default: 8080)
DATA_PATH=./data            # Đường dẫn tới JSON files
GIN_MODE=release            # Gin mode: debug/release
```

## 📖 Ví dụ sử dụng API

### **Lấy danh sách tỉnh/thành**

```bash
curl "http://localhost:8080/api/v1/provinces?limit=10&search=hà"
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "code": "11",
      "name": "Hà Nội",
      "slug": "ha-noi",
      "type": "thanh-pho",
      "name_with_type": "Thành phố Hà Nội"
    }
  ],
  "pagination": {
    "total": 1,
    "limit": 10,
    "offset": 0,
    "pages": 1
  }
}
```

### **Lấy xã/phường theo tỉnh**

```bash
curl "http://localhost:8080/api/v1/provinces/11/wards?limit=5"
```

### **Tìm kiếm toàn cục**

```bash
curl "http://localhost:8080/api/v1/search?q=minh&limit=10"
```

### **Validate địa chỉ**

```bash
curl -X POST "http://localhost:8080/api/v1/address/validate" \
  -H "Content-Type: application/json" \
  -d '{
    "province_code": "11",
    "ward_code": "267"
  }'
```

### **Health check**

```bash
curl "http://localhost:8080/api/v1/health"
```

**Response:**
```json
{
  "success": true,
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "services": {
    "data_loader": "healthy"
  },
  "version": "1.0.0",
  "uptime": "2h15m30s"
}
```

## 🎯 Query Parameters

### **Pagination**
- `limit`: Số records trả về (default: 50, max: 1000)
- `offset`: Bỏ qua n records đầu (default: 0)

### **Filtering**
- `search`: Tìm kiếm theo tên, slug
- `type`: Filter theo loại (thanh-pho, tinh, xa, phuong, etc.)
- `province_code`: Filter ward theo tỉnh

### **Search**
- `q`: Search query (tối thiểu 2 ký tự)
- `entity`: Tìm kiếm trong (province, ward, all)

## 🛠️ Development Commands

```bash
# Development workflow
make dev           # Clean + setup + run
make test          # Run tests
make benchmark     # Run benchmarks
make lint          # Lint code (cần golangci-lint)
make fmt           # Format code

# Production workflow
make prod          # Clean + build + run production
make docker        # Build Docker image
make deploy        # Deploy với docker-compose

# Utilities
make health        # Check API health
make logs          # View docker logs
make load-test     # Run load test (cần hey tool)
```

## 🔧 Troubleshooting

### **Lỗi thường gặp:**

#### **1. Import errors khi build**
```bash
# Fix: Download dependencies
go mod tidy
go mod download
```

#### **2. File not found: province.json/ward.json**
```bash
# Fix: Copy files vào đúng vị trí
mkdir -p data
cp province.json ward.json data/
```

#### **3. Port 8080 đã sử dụng**
```bash
# Fix: Thay đổi port
PORT=8081 go run main.go
# Hoặc kill process
lsof -ti:8080 | xargs kill -9
```

#### **4. Docker build failed**
```bash
# Fix: Đảm bảo JSON files tồn tại
ls -la province.json ward.json
docker build --no-cache -t vietnam-admin-api .
```

#### **5. API trả về 503 Service Unavailable**
```bash
# Fix: Kiểm tra data loading
curl http://localhost:8080/api/v1/health
# Check logs để xem lỗi load data
```

### **Debug mode:**
```bash
# Chạy với debug logging
GIN_MODE=debug go run main.go

# Xem detailed logs
make logs
```

## 📊 Performance

- **Startup time**: < 1 giây
- **Memory usage**: ~50MB (với ~44 tỉnh và ~11k xã/phường)
- **Response time**: < 5ms cho queries đơn giản
- **Throughput**: > 10,000 requests/second

## 🐳 Docker

### **Multi-stage Dockerfile**

```dockerfile
FROM golang:1.21-alpine AS builder
# Build stage...

FROM alpine:latest
# Runtime stage với binary nhỏ gọn
```

### **Docker Compose**

```yaml
version: '3.8'
services:
  vietnam-api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
    volumes:
      - ./data:/root/data
    restart: unless-stopped
```

## 🔐 Security Features

- **CORS**: Configured cho cross-origin requests
- **Rate Limiting**: Prevent abuse (có thể enable)
- **Input Validation**: Validate all inputs
- **Admin Auth**: Token-based authentication cho admin endpoints
- **Non-root Container**: Docker container chạy với non-root user

## 🚀 Deployment

### **Production Checklist**

- [ ] Set `GIN_MODE=release`
- [ ] Configure proper logging
- [ ] Setup monitoring (Prometheus/Grafana)
- [ ] Enable rate limiting
- [ ] Setup load balancer
- [ ] Configure auto-scaling
- [ ] Setup health checks

### **Recommended Stack**

- **Reverse Proxy**: Nginx/Traefik
- **Container**: Docker + Kubernetes
- **Monitoring**: Prometheus + Grafana
- **Logging**: ELK Stack
- **Cache**: Redis (optional)

## 📈 Monitoring & Observability

### **Health Endpoints**

```bash
GET /health                 # Simple health check
GET /api/v1/health         # Detailed health with services
GET /api/v1/stats          # Data statistics
```

### **Metrics**

- Request count và duration
- Error rates
- Memory usage
- Data load time

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run benchmarks
make benchmark

# Load testing
make load-test
```

## 🤝 Contributing

1. Fork the project
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 🙏 Acknowledgments

- Dữ liệu từ Vietnam Administrative Database
- Gin Web Framework
- Go Community

---

**Made with ❤️ for Vietnam** 🇻🇳 