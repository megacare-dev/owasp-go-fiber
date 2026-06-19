# A09 — Security Logging & Monitoring Failures Demo

`POST /login` ที่ล้มเหลว **ไม่ถูกบันทึกเลย** → ถูก brute-force ก็มองไม่เห็น ตรวจจับ/แจ้งเตือนไม่ได้
(เทียบ `/audit` ของทีม security: โหมดมีช่องโหว่ = ว่างเปล่า)

## รัน

```bash
go run .            # ❌ VULNERABLE — ไม่ log เหตุการณ์ความปลอดภัย
./exploit.sh        # ยิง login ผิด 3 ครั้ง → /audit ว่าง (events:0)

SECURE=1 go run .   # ✅ SECURE — บันทึก audit log (ไม่เก็บรหัสผ่าน)
./exploit.sh        # /audit เห็น 3 เหตุการณ์ FAILED_LOGIN

go test ./...       # ✅ test: secure บันทึกครบทุกครั้ง · vulnerable ไม่บันทึกเลย · log ไม่มีรหัสผ่าน
```

หรือดูเทียบทั้งสองโหมดในคำสั่งเดียว: `./demo.sh`

## เห็นอะไร

| โหมด | `/audit` หลังโดนโจมตี |
|------|----------------------|
| VULNERABLE | `{"events":0,...}` ← มองไม่เห็นการโจมตี |
| SECURE | `{"events":3,...}` ← มีร่องรอย ตรวจจับได้ |

## บทเรียน

- ❌ ไม่ log เหตุการณ์ security (login ล้มเหลว, access denied, input ผิดปกติ) = ตาบอด
- ❌ **อย่า** log ข้อมูลอ่อนไหว (รหัสผ่าน, token, เลขบัตร) ลงไปด้วย
- ✅ log แบบมีโครงสร้าง + ส่งเข้า centralized log/SIEM + ตั้ง alert (เช่น failed login ถี่ผิดปกติ)

## 💬 vibe coding

> "เพิ่ม audit logging ให้ endpoint login: บันทึกความพยายาม login ล้มเหลว (user, ip, เวลา) แบบ structured log โดยห้ามบันทึกรหัสผ่าน แล้วเขียน test ว่ามี log เกิดขึ้นเมื่อ login ผิด"
