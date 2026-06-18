// A01 — Broken Access Control demo (IDOR + missing RBAC)
//
//	go run .            → VULNERABLE (อ่าน order ของคนอื่นได้ + non-admin ลบ user ได้)
//	SECURE=1 go run .   → SECURE     (ตรวจ owner + บังคับ role ด้วย RequireRole)
//
// แล้วยิง ./exploit.sh เทียบผลทั้งสองโหมด
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

type Order struct {
	ID     string `json:"id"`
	Owner  string `json:"owner"`
	Item   string `json:"item"`
	Amount int    `json:"amount_thb"`
}

// ฐานข้อมูลจำลอง — order 1001 เป็นของ alice, 1002 เป็นของ bob
var orders = map[string]Order{
	"1001": {"1001", "alice", "iPhone 15", 39900},
	"1002": {"1002", "bob", "MacBook Pro", 89900},
}

// ผู้ใช้จำลองที่ admin มีสิทธิ์ลบ
var users = map[string]string{
	"42": "charlie",
	"43": "dave",
}

// RequireRole — middleware กันตามสิทธิ์: อ่าน role จาก c.Locals("role")
// (ปกติ auth layer ดึงมาจาก JWT/session แล้ว set ไว้) ถ้าไม่ตรงตอบ 403 forbidden
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("role") != role {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}

func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, ok := users[id]; !ok {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	delete(users, id)
	return c.JSON(fiber.Map{"deleted": id})
}

func main() {
	secure := os.Getenv("SECURE") == "1"
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	// จำลอง auth layer: ปกติ role มาจาก JWT/session — ที่นี่อ่านจาก header X-Role
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", c.Get("X-Role"))
		return c.Next()
	})

	app.Get("/orders/:id", func(c *fiber.Ctx) error {
		caller := c.Get("X-User-Id") // จำลองผู้ใช้ที่ล็อกอินอยู่
		o, ok := orders[c.Params("id")]
		if !ok {
			return c.Status(404).JSON(fiber.Map{"error": "not found"})
		}

		if secure && o.Owner != caller {
			// ✅ A01 fix: ตรวจว่าเป็นเจ้าของ order ก่อนคืนข้อมูล
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}

		// ❌ VULNERABLE: คืนข้อมูลตาม :id โดยไม่สนว่าใครเรียก (IDOR)
		return c.JSON(o)
	})

	// DELETE /admin/users/:id — ควรให้เฉพาะ admin
	if secure {
		// ✅ SECURE: ใส่ RequireRole("admin") ก่อนถึง handler
		app.Delete("/admin/users/:id", RequireRole("admin"), deleteUser)
	} else {
		// ❌ VULNERABLE: ลืมใส่ RequireRole → ใครก็ลบ user ได้ (privilege escalation)
		app.Delete("/admin/users/:id", deleteUser)
	}

	mode := "VULNERABLE ❌"
	if secure {
		mode = "SECURE ✅"
	}
	log.Printf("A01 demo [%s] → GET /orders/:id (IDOR) · DELETE /admin/users/:id (RBAC)", mode)
	log.Fatal(app.Listen(":3000"))
}
