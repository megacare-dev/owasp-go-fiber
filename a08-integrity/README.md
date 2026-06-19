# A08 — Software & Data Integrity Failures (Unverified Webhook) Demo

endpoint `POST /webhook` ประมวลผล payload โดย **ไม่ตรวจลายเซ็น** → ใครก็ปลอม event
(เช่น `{"status":"paid"}`) ส่งมาหลอกให้ระบบคิดว่าจ่ายเงินแล้วได้

## รัน

```bash
go run .            # ❌ VULNERABLE — เชื่อ payload โดยไม่ตรวจลายเซ็น
./exploit.sh        # ส่ง webhook ปลอม → {"processed":true}

SECURE=1 go run .   # ✅ SECURE — verify HMAC-SHA256(body, secret)
./exploit.sh        # ลายเซ็นไม่ถูก → 401 invalid signature

WEBHOOK_SECRET=xxx SECURE=1 go run .   # secret มาจาก env (ไม่ hardcode)
go test ./...       # ✅ test: ไม่มีลายเซ็น / ลายเซ็นผิด / body ถูกแก้ → ปฏิเสธ
```

หรือดูเทียบทั้งสองโหมดในคำสั่งเดียว: `./demo.sh`

## เห็นอะไร

| โหมด | ผล |
|------|-----|
| VULNERABLE | `{"processed":true,...}` ← เชื่อ payload ปลอม |
| SECURE | `{"error":"invalid signature"}` (401) ← ปฏิเสธ |

## บทเรียน

- ❌ อย่าเชื่อข้อมูล/โค้ดจากภายนอกโดยไม่ตรวจ integrity (webhook, อัปเดต, deserialization, artifact)
- ✅ verify ลายเซ็น (HMAC / digital signature) ก่อนประมวลผลเสมอ — ใช้ `hmac.Equal` (constant-time)
- ✅ ลายเซ็นผูกกับ body ทั้งก้อน → body ที่ถูกแก้ (เช่น เปลี่ยนยอดเงิน) จะ verify ไม่ผ่าน
- ✅ secret อ่านจาก env (`WEBHOOK_SECRET`) ไม่ hardcode — demo มี default ให้รันได้ทันที

## 💬 vibe coding

> "endpoint webhook นี้ไม่ verify ลายเซ็น ช่วยเพิ่มการตรวจ HMAC-SHA256 จาก header X-Signature เทียบด้วย hmac.Equal และอ่าน secret จาก env แล้วเขียน test ว่า payload ที่ไม่มีลายเซ็นถูกปฏิเสธ"
