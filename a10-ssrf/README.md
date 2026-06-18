# A10 — SSRF Demo

`/fetch?url=` ดึง URL ที่ client ส่งมาฝั่ง server โดยไม่ตรวจปลายทาง → หลอกให้ดึง **cloud metadata / internal service** มาให้ได้

> demo มี mock metadata server บน `127.0.0.1:9999` (จำลอง `169.254.169.254`) — ไม่ยิงออกเน็ตจริง

## รัน

```bash
go run .            # ❌ VULNERABLE — fetch ตรงๆ
./scan.sh           # 🔍 Scan หาช่องโหว่ (gosec/govulncheck)
./exploit.sh        # /fetch?url=http://127.0.0.1:9999/... → ได้ AWS_SECRET ภายใน

SECURE=1 go run .   # ✅ SECURE — บล็อก loopback/private/link-local
./exploit.sh        # 400 blocked: internal address
```

## เห็นอะไร

| โหมด | ผล |
|------|-----|
| VULNERABLE | `200 AWS_SECRET_ACCESS_KEY=...` ← ดึง credential ภายใน |
| SECURE | `400 {"error":"blocked: internal address"}` |

## บทเรียน

- ตรวจปลายทางก่อนยิง outbound: บล็อก **loopback / private / link-local (169.254.169.254)**
- ใช้ allow-list domain ที่อนุญาต ดีกว่า block-list
- ระวัง DNS rebinding (resolve แล้วค่อยเช็ค IP จริง) + redirect ตาม

```go
if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() { /* block */ }
```

## 💬 vibe coding

> "endpoint /fetch นี้มี SSRF ช่วยเพิ่มการตรวจปลายทาง: resolve host แล้วบล็อก private/loopback/metadata IP, ใช้ allow-list domain, ปิด redirect แล้วเขียน test ยิง 127.0.0.1 และ 169.254.169.254 ว่าโดนบล็อก"
