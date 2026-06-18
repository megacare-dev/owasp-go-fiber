# A01 — Broken Access Control (IDOR + RBAC) Demo

มีสองช่องโหว่ในตระกูล A01:

1. **IDOR** — ผู้ใช้ `alice` อ่าน order ของ `bob` ได้ เพราะ endpoint คืนข้อมูลตาม `:id` โดย **ไม่ตรวจเจ้าของ**
2. **Missing RBAC** — `DELETE /admin/users/:id` **ลืมใส่ role check** → ใครก็ลบ user ได้ (privilege escalation) จนกว่าจะใส่ middleware `RequireRole("admin")`

## รัน

```bash
go run .            # ❌ VULNERABLE
./scan.sh           # 🔍 Scan หาช่องโหว่ (gosec/govulncheck)
./exploit.sh        # alice ขอ order 1002 (ของ bob) → 200 + ข้อมูล bob รั่ว
```

แล้วสลับเป็นโหมดแก้:

```bash
SECURE=1 go run .   # ✅ SECURE — ตรวจ owner ก่อนคืนข้อมูล
./exploit.sh        # ยิงซ้ำ → 403 forbidden
```

## เห็นอะไร

| โหมด | IDOR (`GET /orders/1002`) | RBAC (`DELETE /admin/users/42`) |
|------|---------------------------|--------------------------------|
| VULNERABLE | `HTTP 200` + ข้อมูล bob รั่ว | `HTTP 200` + `{"deleted":"42"}` ← non-admin ลบได้ |
| SECURE | `HTTP 403` + `{"error":"forbidden"}` | `HTTP 403` + `{"error":"forbidden"}` |

## จุดที่ต่างกัน (ดูใน `main.go`)

**IDOR — owner check**

```go
if secure && o.Owner != caller {
    return c.Status(403).JSON(fiber.Map{"error": "forbidden"})  // ✅ owner check
}
return c.JSON(o)   // ❌ vulnerable: ไม่เช็คว่าใครเรียก
```

**RBAC — RequireRole middleware**

```go
func RequireRole(role string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        if c.Locals("role") != role {
            return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
        }
        return c.Next()
    }
}

// ✅ SECURE: ใส่ RequireRole ก่อน handler
app.Delete("/admin/users/:id", RequireRole("admin"), deleteUser)
// ❌ VULNERABLE: ลืม RequireRole → ใครก็ลบได้
app.Delete("/admin/users/:id", deleteUser)
```

## 💬 ต่อยอดด้วย vibe coding

> "แก้ endpoint นี้ให้กัน IDOR โดยดึง user จริงจาก JWT (ไม่ใช่ header `X-User-Id`) แล้วเขียน test ที่ยืนยันว่า alice เข้าถึง order ของ bob ไม่ได้"
