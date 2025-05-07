# OpenAPI YAML Operation ID Fixer

เครื่องมือสำหรับแก้ไข `operationId` ที่ซ้ำกันในไฟล์ OpenAPI YAML

## คุณสมบัติ

- ตรวจสอบและแก้ไข `operationId` ที่ซ้ำกันในไฟล์ OpenAPI YAML
- แยกตรวจสอบตาม Tag ของ API
- สร้างไฟล์ใหม่โดยไม่แก้ไขไฟล์ต้นฉบับ
- แสดงรายงานการแก้ไขที่เกิดขึ้น

## การติดตั้ง

1. ตรวจสอบให้แน่ใจว่าคุณมี Go ติดตั้งแล้วบนเครื่อง
2. Clone repository นี้:
   ```
   git clone https://github.com/yourusername/openapi-operation-id-fixer.git
   cd openapi-operation-id-fixer
   ```

## การใช้งาน

```
go run fix_operation_id.go <openapi.yaml>
```

ตัวอย่าง:
```
go run fix_operation_id.go api-spec.yaml
```

หลังจากทำงานเสร็จ โปรแกรมจะสร้างไฟล์ใหม่ชื่อ `<openapi>_fixed.yaml`

