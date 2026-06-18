# A07 — Auth Failure (JWT alg:none Forge) Demo

server เชื่อ payload ของ JWT โดย **ไม่ verify signature** → attacker ปลอม `{"role":"admin"}` ด้วย `alg:none` เข้า /admin ได้ โดยไม่ต้องรู้ secret

## รัน

```bash
go run .            # ❌ VULNERABLE — ไม่ตรวจ signature
./scan.sh           # 🔍 Scan หาช่องโหว่ (gosec/govulncheck)
./exploit.sh        # ปลอม token alg:none role=admin → 200 welcome admin

SECURE=1 go run .   # ✅ SECURE — บังคับ HS256 + verify HMAC
./exploit.sh        # token ปลอมถูกปฏิเสธ → 401
```

## เห็นอะไร

| โหมด | ผลของ token ปลอม |
|------|------------------|
| VULNERABLE | `200 {"msg":"welcome admin","secret":"ALL_USER_DATA"}` ← เข้าได้ |
| SECURE | `401 {"error":"unauthorized"}` ← ปฏิเสธ |

## บทเรียน

- ❌ อย่าอ่าน claims จาก JWT โดยไม่ verify signature
- ❌ อย่ายอมรับ `alg:none` / อย่าให้ client เลือก alg
- ✅ บังคับ algorithm ที่กำหนด (HS256/RS256) + verify signature + ตรวจ `exp`
- ใช้ library ที่ดี (`golang-jwt/jwt/v5`) และระบุ allowed methods เสมอ

## 💬 vibe coding

> "เปลี่ยนการตรวจ JWT นี้ไปใช้ golang-jwt/jwt/v5 แบบ verify HS256 ด้วย secret จาก env, ปฏิเสธ alg:none, ตรวจ exp แล้วเขียน test ว่า token alg:none และ token เซ็นผิด secret เข้าไม่ได้"
