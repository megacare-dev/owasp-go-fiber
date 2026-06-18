# A05 — Security Misconfiguration (Verbose Error) Demo

default error handler ส่ง error ดิบกลับให้ client → leak connection string + รหัสผ่าน DB ภายใน

## รัน

```bash
go run .            # ❌ VULNERABLE — default ErrorHandler
./scan.sh           # 🔍 Scan หาช่องโหว่ (gosec/govulncheck)
./exploit.sh        # /report/42 → error เผย postgres://app:S3cretDbP@ss@10.0.0.5...

SECURE=1 go run .   # ✅ SECURE — custom ErrorHandler
./exploit.sh        # คืน {"error":"internal server error"} เท่านั้น
```

## เห็นอะไร

| โหมด | body ตอน error |
|------|----------------|
| VULNERABLE | `db connect failed dsn=postgres://app:S3cretDbP@ss@10.0.0.5:5432/prod ...` ← รั่ว |
| SECURE | `{"error":"internal server error"}` (log จริงอยู่ฝั่ง server) |

## บทเรียน

- ตั้ง custom `ErrorHandler` ที่ **ไม่ส่งรายละเอียดภายใน** กลับ client
- log error เต็มไว้ฝั่ง server (เข้า SIEM) — แต่ตอบ client แบบกลางๆ
- ปิด stack trace / debug mode ใน production

## 💬 vibe coding

> "ใส่ custom ErrorHandler ให้ Fiber app นี้ ที่ตอบ client เป็น `{error:'internal server error'}` + request id แต่ log error เต็มฝั่ง server, และเพิ่ม recover middleware กัน panic leak"
