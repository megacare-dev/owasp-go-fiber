# A06 — Vulnerable & Outdated Components Demo

โปรเจกต์ pin `gofiber/fiber v2.52.0` ที่มีช่องโหว่ที่รู้แล้ว (เช่น **CVE-2024-25124** ใน CORS middleware)
→ แอปสืบทอดช่องโหว่ของ dependency โดยไม่รู้ตัว ตราบใดที่ไม่สแกน + ไม่อัปเดต

> ต่างจาก demo อื่น: A06 ไม่ใช่ bug ในโค้ดเรา แต่อยู่ใน **dependency** → ตรวจด้วย **SCA** (govulncheck) ไม่ใช่ exploit runtime

## รัน

```bash
./demo.sh           # สรุป: เวอร์ชันที่ใช้ + ผล govulncheck
# หรือ
./scan.sh           # 🔍 govulncheck ./... → รายงาน CVE ใน fiber v2.52.0
```

ติดตั้ง scanner ก่อน (ถ้ายังไม่มี):

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## วิธีแก้ (Fix)

```bash
go get github.com/gofiber/fiber/v2@v2.52.5
go mod tidy
./scan.sh           # สแกนซ้ำ → สะอาด
```

## บทเรียน

- ❌ ใช้ dependency เวอร์ชันเก่าที่มี CVE = เปิดช่องโหว่ทั้งที่โค้ดเราไม่ผิด
- ✅ สแกน dependency สม่ำเสมอ (`govulncheck`, Dependabot, `go list -m -u all`)
- ✅ อัปเดต/patch เป็นเวอร์ชันที่แก้แล้ว และตรึงเวอร์ชันด้วย `go.sum`
- 👉 ช่องโหว่ตัวจริงที่รันได้ของ CVE นี้ดูที่ [`../cors-cve/`](../cors-cve/)

## 💬 vibe coding

> "สแกนโปรเจกต์นี้ด้วย govulncheck บอกว่ามี CVE อะไรใน dependency บ้าง แล้วอัปเดต gofiber/fiber เป็นเวอร์ชันที่แก้แล้วให้ พร้อมอธิบายว่า CVE-2024-25124 คืออะไร"
