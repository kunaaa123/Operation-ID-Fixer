# OpenAPI YAML Operation ID Fixer

เครื่องมือสำหรับแก้ไข `operationId` ที่ซ้ำกันในไฟล์ OpenAPI YAML

## คุณสมบัติ

- ตรวจสอบและแก้ไข `operationId` ที่ซ้ำกันในไฟล์ OpenAPI YAML
- แยกตรวจสอบตาม Tag ของ API
- สร้างไฟล์ใหม่โดยไม่แก้ไขไฟล์ต้นฉบับ
- แสดงรายงานการแก้ไขที่เกิดขึ้น

## การติดตั้ง

1. Clone repository นี้:
   ```
   git clone https://github.com/kunaaa123/OpenAPI-YAML-Operation-ID-Fixer.git
   cd openapi-operation-id-fixer
   ```

## การใช้งาน

```
go run fix_operation_id.go <openapi.yaml>
```



หลังจากทำงานเสร็จ โปรแกรมจะสร้างไฟล์ใหม่ชื่อ `<openapi>_fixed.yaml`

