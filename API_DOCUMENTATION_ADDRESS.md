# Vietnam Administrative API Documentation

## Tổng quan

Vietnam Administrative API cung cấp dữ liệu đầy đủ về các đơn vị hành chính của Việt Nam, bao gồm tỉnh thành và phường/xã/thị trấn. API này được thiết kế để hỗ trợ các ứng dụng cần thông tin địa danh chính xác và cập nhật.

### Đặc điểm chính
- ✅ **Miễn phí**: Không cần API key hoặc authentication
- 🚀 **Nhanh chóng**: Response time dưới 100ms
- 📊 **Đầy đủ**: Hơn 60+ tỉnh thành và 10,000+ phường/xã
- 🔍 **Tìm kiếm mạnh mẽ**: Hỗ trợ tìm kiếm theo tên, slug, và mã
- 📱 **RESTful**: Tuân thủ chuẩn REST API
- 🌐 **CORS**: Hỗ trợ cross-origin requests

## Base URL

```
https://your-domain.com/api/v1
```

## Thông tin chung

### Content Type
Tất cả responses đều trả về `application/json`

### Pagination
Các endpoint trả về danh sách đều hỗ trợ pagination:
- `limit`: Số lượng items trả về (default: 50, max: 1000)
- `offset`: Vị trí bắt đầu (default: 0)

### Search
Hỗ trợ tìm kiếm qua parameter:
- `search`: Tìm kiếm theo tên
- `type`: Lọc theo loại (tỉnh, thành phố, phường, xã, etc.)

## Response Format

### Success Response
```json
{
  "success": true,
  "data": [...],
  "message": "Optional message"
}
```

### Paginated Response
```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "total": 63,
    "limit": 20,
    "offset": 0,
    "pages": 4
  }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description"
}
```

## Endpoints

### 1. Provinces (Tỉnh thành)

#### GET /provinces
Lấy danh sách tất cả tỉnh thành

**Parameters:**
- `search` (string, optional): Tìm kiếm theo tên
- `type` (string, optional): Lọc theo loại (tinh, thanh-pho)
- `limit` (int, optional): Số lượng trả về (default: 50)
- `offset` (int, optional): Vị trí bắt đầu (default: 0)

**Example Request:**
```bash
GET /api/v1/provinces?search=hà&limit=10
```

**Example Response:**
```json
{
  "success": true,
  "data": [
    {
      "code": "01",
      "name": "Hà Nội",
      "slug": "ha-noi",
      "type": "thanh-pho",
      "name_with_type": "Thành phố Hà Nội",
      "code_name": "ha_noi"
    },
    {
      "code": "17",
      "name": "Hòa Bình",
      "slug": "hoa-binh", 
      "type": "tinh",
      "name_with_type": "Tỉnh Hòa Bình",
      "code_name": "hoa_binh"
    }
  ],
  "pagination": {
    "total": 2,
    "limit": 10,
    "offset": 0,
    "pages": 1
  }
}
```

#### GET /provinces/:code
Lấy thông tin chi tiết của một tỉnh

**Parameters:**
- `code` (string, required): Mã tỉnh (01, 02, ...)

**Example Request:**
```bash
GET /api/v1/provinces/01
```

**Example Response:**
```json
{
  "success": true,
  "data": {
    "code": "01",
    "name": "Hà Nội",
    "slug": "ha-noi",
    "type": "thanh-pho",
    "name_with_type": "Thành phố Hà Nội",
    "code_name": "ha_noi"
  }
}
```

#### GET /provinces/:code/wards
Lấy danh sách phường/xã của một tỉnh

**Parameters:**
- `code` (string, required): Mã tỉnh
- `search` (string, optional): Tìm kiếm theo tên
- `type` (string, optional): Lọc theo loại
- `limit` (int, optional): Số lượng trả về
- `offset` (int, optional): Vị trí bắt đầu

**Example Request:**
```bash
GET /api/v1/provinces/01/wards?search=hoàng&limit=5
```

#### GET /provinces/types
Lấy danh sách các loại tỉnh thành

**Example Response:**
```json
{
  "success": true,
  "data": [
    "tinh",
    "thanh-pho"
  ]
}
```

### 2. Wards (Phường/Xã)

#### GET /wards
Lấy danh sách tất cả phường/xã

**Parameters:**
- `search` (string, optional): Tìm kiếm theo tên
- `type` (string, optional): Lọc theo loại (phuong, xa, thi-tran)
- `province_code` (string, optional): Lọc theo mã tỉnh
- `limit` (int, optional): Số lượng trả về (default: 50)
- `offset` (int, optional): Vị trí bắt đầu (default: 0)

**Example Request:**
```bash
GET /api/v1/wards?province_code=01&type=phuong&limit=10
```

**Example Response:**
```json
{
  "success": true,
  "data": [
    {
      "code": "00001",
      "name": "Phúc Xá",
      "slug": "phuc-xa",
      "type": "phuong",
      "name_with_type": "Phường Phúc Xá",
      "path": "Phúc Xá, Ba Đình, Hà Nội",
      "path_with_type": "Phường Phúc Xá, Quận Ba Đình, Thành phố Hà Nội",
      "parent_code": "001"
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

#### GET /wards/:code
Lấy thông tin chi tiết của một phường/xã

**Parameters:**
- `code` (string, required): Mã phường/xã

**Example Request:**
```bash
GET /api/v1/wards/00001
```

**Example Response:**
```json
{
  "success": true,
  "data": {
    "code": "00001",
    "name": "Phúc Xá",
    "slug": "phuc-xa",
    "type": "phuong",
    "name_with_type": "Phường Phúc Xá",
    "path": "Phúc Xá, Ba Đình, Hà Nội",
    "path_with_type": "Phường Phúc Xá, Quận Ba Đình, Thành phố Hà Nội",
    "parent_code": "001",
    "province": {
      "code": "01",
      "name": "Hà Nội",
      "name_with_type": "Thành phố Hà Nội",
      "type": "thanh-pho"
    }
  }
}
```

#### GET /wards/types
Lấy danh sách các loại phường/xã

**Example Response:**
```json
{
  "success": true,
  "data": [
    "phuong",
    "xa",
    "thi-tran",
    "dac-khu"
  ]
}
```

### 3. Search (Tìm kiếm)

#### GET /search
Tìm kiếm toàn cục qua tất cả tỉnh và phường/xã

**Parameters:**
- `q` (string, required): Từ khóa tìm kiếm (tối thiểu 2 ký tự)
- `entity` (string, optional): Loại entity cần tìm (province, ward, all) - default: all
- `limit` (int, optional): Số lượng trả về (default: 20, max: 100)

**Example Request:**
```bash
GET /api/v1/search?q=hà nội&limit=10
```

**Example Response:**
```json
{
  "success": true,
  "data": {
    "provinces": [
      {
        "code": "01",
        "name": "Hà Nội",
        "slug": "ha-noi",
        "type": "thanh-pho",
        "name_with_type": "Thành phố Hà Nội"
      }
    ],
    "wards": [
      {
        "code": "00001",
        "name": "Phúc Xá",
        "path": "Phúc Xá, Ba Đình, Hà Nội",
        "type": "phuong"
      }
    ]
  },
  "query": "hà nội"
}
```

### 4. Address Validation

#### POST /address/validate
Kiểm tra tính hợp lệ của một địa chỉ

**Request Body:**
```json
{
  "province_code": "01",
  "ward_code": "00001"
}
```

**Example Response:**
```json
{
  "success": true,
  "valid": true,
  "message": "Address is valid",
  "data": {
    "code": "00001",
    "name": "Phúc Xá",
    "path": "Phúc Xá, Ba Đình, Hà Nội",
    "type": "phuong"
  }
}
```

### 5. System Information

#### GET /health
Kiểm tra tình trạng hoạt động của API

**Example Response:**
```json
{
  "success": true,
  "status": "healthy",
  "timestamp": "2025-01-04T10:30:00Z",
  "services": {
    "data_loader": "healthy"
  },
  "version": "1.0.0",
  "uptime": "2h15m30s"
}
```

#### GET /stats
Lấy thống kê về dữ liệu

**Example Response:**
```json
{
  "success": true,
  "data": {
    "provinces": 63,
    "wards": 10960,
    "province_types": {
      "tinh": 58,
      "thanh-pho": 5
    },
    "ward_types": {
      "xa": 8045,
      "phuong": 1968,
      "thi-tran": 614,
      "dac-khu": 333
    },
    "uptime": "2h15m30s",
    "version": "1.0.0"
  }
}
```

### 6. Admin Endpoints

#### POST /admin/reload
Tải lại dữ liệu từ file JSON (không cần authentication hiện tại)

**Example Response:**
```json
{
  "success": true,
  "message": "Data reloaded successfully",
  "data": {
    "reload_time": "2025-01-04T10:30:00Z",
    "stats": {
      "provinces": 63,
      "wards": 10960
    }
  }
}
```

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request - Invalid parameters |
| 404 | Not Found - Resource doesn't exist |
| 503 | Service Unavailable - Data not loaded |
| 500 | Internal Server Error |

## Rate Limiting

Hiện tại API chưa có rate limiting. Nếu cần thiết, sẽ được thêm vào với các limit hợp lý.

## CORS Support

API hỗ trợ CORS với cấu hình:
- **Origin**: `*` (tất cả domains)
- **Methods**: `GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS`
- **Headers**: `Origin, Content-Length, Content-Type, Authorization, X-Requested-With`

## Examples

### Frontend JavaScript

```javascript
// Lấy danh sách tỉnh
const provinces = await fetch('/api/v1/provinces?limit=10')
  .then(res => res.json());

// Tìm kiếm
const searchResults = await fetch('/api/v1/search?q=hà nội')
  .then(res => res.json());

// Validate địa chỉ
const validation = await fetch('/api/v1/address/validate', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    province_code: '01',
    ward_code: '00001'
  })
}).then(res => res.json());
```

### cURL Examples

```bash
# Lấy tất cả tỉnh
curl "http://localhost:8080/api/v1/provinces"

# Tìm kiếm phường/xã
curl "http://localhost:8080/api/v1/search?q=phúc xá"

# Validate địa chỉ
curl -X POST "http://localhost:8080/api/v1/address/validate" \
  -H "Content-Type: application/json" \
  -d '{"province_code":"01","ward_code":"00001"}'

# Health check
curl "http://localhost:8080/api/v1/health"
```

## Support

Nếu có vấn đề hoặc câu hỏi, vui lòng:
1. Kiểm tra endpoint `/health` để đảm bảo service đang hoạt động
2. Xem `/stats` để biết thông tin về dữ liệu hiện tại
3. Đảm bảo request format đúng theo tài liệu

---

**Version**: 1.0.0  
**Last Updated**: 04/01/2025 