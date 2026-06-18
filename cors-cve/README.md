# CORS Misconfiguration Demo (CVE-2024-25124)

CORS ที่สะท้อน **origin ใดก็ได้** ร่วมกับ `Allow-Credentials: true` → เว็บของ attacker (evil.com) อ่าน response ที่มี cookie/credential ของเหยื่อได้

## รัน

```bash
go run .            # ❌ VULNERABLE — AllowOriginsFunc=true + AllowCredentials
./scan.sh           # 🔍 Scan หาช่องโหว่ (gosec/govulncheck)
./exploit.sh        # Origin: https://evil.com → ถูกสะท้อนกลับ + credentials:true

SECURE=1 go run .   # ✅ SECURE — allow-list https://app.example.com
./exploit.sh        # evil.com ไม่ถูกสะท้อน
```

## เห็นอะไร

| โหมด | response header สำหรับ Origin: evil.com |
|------|------------------------------------------|
| VULNERABLE | `Access-Control-Allow-Origin: https://evil.com` + `Allow-Credentials: true` ← อันตราย |
| SECURE | ไม่มี header สำหรับ evil.com ← บล็อก |

## บทเรียน (CVE-2024-25124)

- ❌ ห้ามใช้ `AllowOrigins:"*"` (หรือสะท้อนทุก origin) คู่กับ `AllowCredentials:true`
- ✅ ใช้ **allow-list** origin ที่เชื่อถือเท่านั้น
- อัปเดต gofiber/fiber เป็นเวอร์ชันที่แก้ CVE แล้ว + รัน `govulncheck`

## 💬 vibe coding

> "config CORS ของ Fiber app นี้สะท้อนทุก origin พร้อม credentials ช่วยแก้เป็น allow-list เฉพาะ domain ที่กำหนดใน env, ห้าม wildcard+credentials แล้วอธิบายว่า CVE-2024-25124 อันตรายยังไง"
