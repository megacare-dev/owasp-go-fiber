# A08 — Software & Data Integrity Failures (Unverified Webhook) Demo

endpoint `POST /webhook` ประมวลผล payload โดย **ไม่ตรวจลายเซ็น** → ใครก็ปลอม event
(เช่น `{"status":"paid"}`) ส่งมาหลอกให้ระบบคิดว่าจ่ายเงินแล้วได้

## รัน

```bash
go run .            # ❌ VULNERABLE — เชื่อ payload โดยไม่ตรวจลายเซ็น
./exploit.sh        # ส่ง webhook ปลอม → {"processed":true}

SECURE=1 go run .   # ✅ SECURE — verify HMAC-SHA256(body, secret)
./exploit.sh        # ลายเซ็นไม่ถูก → 401 invalid signature
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
- ✅ ใช้ secret จาก env/secret manager ไม่ hardcode (demo นี้ hardcode เพื่อความง่าย)

## 💬 vibe coding

> "endpoint webhook นี้ไม่ verify ลายเซ็น ช่วยเพิ่มการตรวจ HMAC-SHA256 จาก header X-Signature เทียบด้วย hmac.Equal และอ่าน secret จาก env แล้วเขียน test ว่า payload ที่ไม่มีลายเซ็นถูกปฏิเสธ"
