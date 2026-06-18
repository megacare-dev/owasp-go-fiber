# A03 — SQL Injection Demo (SQLite pure-Go)

`/search?name=` ต่อ input เข้า SQL ตรงๆ → payload `' OR '1'='1` ทำให้ดึง user ทุกคน + secret ออกมา

## รัน

```bash
go run .            # ❌ VULNERABLE — string concat
./scan.sh           # 🔍 Scan หาช่องโหว่ (gosec/govulncheck)
./exploit.sh        # name=' OR '1'='1 → ดึงทั้งตาราง + ADMIN_ROOT_TOKEN

SECURE=1 go run .   # ✅ SECURE — parameterized query (WHERE name = ?)
./exploit.sh        # payload ถูกมองเป็นชื่อธรรมดา → []
```

> ใช้ `modernc.org/sqlite` (pure-Go, ไม่ต้อง CGO/Docker) DB อยู่ใน `:memory:`

## เห็นอะไร

| โหมด | ผลของ `name=' OR '1'='1` |
|------|--------------------------|
| VULNERABLE | คืน **ทุกแถว** รวม `admin / ADMIN_ROOT_TOKEN_9f3x` ← รั่ว |
| SECURE | `[]` ← payload เป็นแค่ string ค้นหา |

## จุดที่ต่างกัน (`main.go`)

```go
// ❌ db.Query(fmt.Sprintf("... WHERE name = '%s'", name))
db.Query("SELECT ... WHERE name = ?", name)   // ✅ input แยกจาก SQL
```

## 💬 vibe coding

> "endpoint /search นี้ต่อ string เข้า SQL ช่วยแก้เป็น parameterized query ทั้งหมด, เพิ่ม input validation, แล้วเขียน test ที่ยิง `' OR '1'='1` และ `'; DROP TABLE users;--` ยืนยันว่ากันได้"
