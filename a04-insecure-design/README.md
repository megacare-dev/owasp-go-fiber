# A04 — Insecure Design (Trust Client Price) Demo

endpoint `POST /checkout` คิดเงินตาม `price` ที่ **client ส่งมา** → ผู้ใช้ตั้งราคาเองเป็น 1 บาทได้
นี่ไม่ใช่ bug การเขียน แต่เป็น **การออกแบบที่ผิด** (เชื่อ input ที่ไม่ควรเชื่อ)

## รัน

```bash
go run .            # ❌ VULNERABLE — เชื่อ price จาก client
./exploit.sh        # สั่งซื้อ iphone ด้วย price=1 → charged=1

SECURE=1 go run .   # ✅ SECURE — ดึงราคาจาก catalog ฝั่ง server
./exploit.sh        # charged=39900 (ราคาจริง)
```

หรือดูเทียบทั้งสองโหมดในคำสั่งเดียว: `./demo.sh`

## เห็นอะไร

| โหมด | ผล |
|------|-----|
| VULNERABLE | `{"charged":1}` ← ซื้อ iphone 1 บาท |
| SECURE | `{"charged":39900}` ← ราคาจริงจาก server |

## บทเรียน

- ❌ อย่าเชื่อข้อมูลที่กระทบ business logic จาก client (ราคา, สิทธิ์, ส่วนลด, ยอดคงเหลือ)
- ✅ คำนวณ/ตรวจสอบฝั่ง server เสมอ จาก source of truth (DB/catalog)
- ✅ Insecure Design แก้ที่ "การออกแบบ flow" ไม่ใช่แค่ patch ทีละจุด — คิด threat model ตั้งแต่ออกแบบ

## 💬 vibe coding

> "endpoint checkout นี้คิดเงินจาก price ที่ client ส่งมา ช่วยแก้ให้ดึงราคาจาก catalog ฝั่ง server เท่านั้น แล้วเขียน test ว่าส่ง price ปลอมมาก็ไม่มีผล"
