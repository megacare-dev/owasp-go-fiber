# A02 — Cryptographic Failures (Predictable Token) Demo

reset token สร้างจาก `md5(username)` → ใครรู้ username ก็ **คำนวณ token เองได้** ยึดบัญชีทันที

## รัน

```bash
go run .            # ❌ VULNERABLE — token = md5(username)
./scan.sh           # 🔍 Scan หาช่องโหว่ (gosec/govulncheck)
./exploit.sh        # attacker คำนวณ md5("bob") → ยึดบัญชี bob (ยิงครั้งเดียว)

SECURE=1 go run .   # ✅ SECURE — token = crypto/rand 32 bytes
./exploit.sh        # token เดาเดิมใช้ไม่ได้ → 401
```

## เห็นอะไร

| โหมด | ผล |
|------|-----|
| VULNERABLE | `{"ok":true,"took_over":"bob"}` ← ยึดบัญชีด้วย token ที่คำนวณเอง |
| SECURE | `{"ok":false}` (401) ← token เดาไม่ได้ |

## บทเรียน

- ❌ อย่าสร้าง token จากข้อมูลที่เดาได้ (username, email, timestamp, sequential id)
- ❌ อย่าใช้ `math/rand` / `md5` / `sha1` กับ secret
- ✅ ใช้ `crypto/rand` ความยาวพอ (≥128-bit) สำหรับ token/secret ทุกชนิด

## 💬 vibe coding

> "endpoint reset password นี้สร้าง token จาก md5(username) ช่วยแก้ให้ใช้ crypto/rand 32 bytes, ผูกกับ user + เวลาหมดอายุ 15 นาที แล้วเขียน test ว่า token เดาจาก username ไม่ได้"
