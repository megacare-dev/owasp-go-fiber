# A01 — Broken Access Control (IDOR + RBAC) Demo

มีสองช่องโหว่ในตระกูล A01:

1. **IDOR** — ผู้ใช้ `alice` อ่าน order ของ `bob` ได้ เพราะ endpoint คืนข้อมูลตาม `:id` โดย **ไม่ตรวจเจ้าของ**
2. **Missing RBAC** — `DELETE /admin/users/:id` **ลืมใส่ role check** → ใครก็ลบ user ได้ (privilege escalation) จนกว่าจะใส่ middleware `RequireRole("admin")`

> **identity มาจากไหน?** โหมด VULNERABLE เชื่อ header `X-User-Id` / `X-Role` ที่ client ส่งมาตรงๆ → ปลอมได้ง่าย
> โหมด SECURE ดึง identity/role จาก **session ที่ verify ฝั่ง server** (จำลองด้วย token `Authorization: tok-admin`) → ปลอม header ไม่มีผล (ของจริงใช้ JWT/session store)

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

| โหมด | IDOR (ปลอม `X-User-Id: alice`) | RBAC (ปลอม `X-Role: admin`) |
|------|---------------------------|--------------------------------|
| VULNERABLE | `HTTP 200` + ข้อมูล bob รั่ว | `HTTP 200` + `{"deleted":"42"}` ← ปลอม header ก็ลบได้ |
| SECURE | `HTTP 403` + `{"error":"forbidden"}` | `HTTP 403` + `{"error":"forbidden"}` ← header ปลอมไม่มีผล |

## จุดที่ต่างกัน (ดูใน `main.go`)

**identity source — ไม่เชื่อ header ดิบ**

```go
if secure {
    // ✅ identity/role มาจาก session ที่ verify ฝั่ง server เท่านั้น
    if s, ok := sessions[c.Get("Authorization")]; ok {
        c.Locals("user", s.user)
        c.Locals("role", s.role)
    }
} else {
    // ❌ vulnerable: เชื่อ header ที่ client ส่งมา → ปลอมได้
    c.Locals("user", c.Get("X-User-Id"))
    c.Locals("role", c.Get("X-Role"))
}
```

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
