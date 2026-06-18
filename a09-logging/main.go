// A09 — Security Logging & Monitoring Failures: ไม่บันทึกเหตุการณ์ด้านความปลอดภัย
//
//	go run .            → VULNERABLE (login ล้มเหลวเงียบๆ → ตรวจจับการโจมตีไม่ได้)
//	SECURE=1 go run .   → SECURE     (บันทึก audit log การ login ล้มเหลว ไม่เก็บรหัสผ่าน)
//
// ยิง ./exploit.sh: ลอง login ผิดหลายครั้ง แล้วเช็ค /audit ว่ามีร่องรอยไหม
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// ในงานจริงคือ centralized log / SIEM — ที่นี่จำลองในหน่วยความจำ
var auditLog []string

func main() {
	secure := os.Getenv("SECURE") == "1"
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/login", func(c *fiber.Ctx) error {
		user := c.Query("user", "admin")
		// demo: สมมติรหัสผ่านผิดทุกครั้ง (จำลองการ brute-force)
		if secure {
			// ✅ บันทึกเหตุการณ์ login ล้มเหลว (ไม่เก็บรหัสผ่าน!) เพื่อให้ตรวจจับการโจมตีได้
			auditLog = append(auditLog, fmt.Sprintf("FAILED_LOGIN user=%s ip=%s", user, c.IP()))
		}
		// ❌ VULNERABLE: ล้มเหลวเงียบๆ ไม่บันทึกอะไรเลย → มองไม่เห็นการโจมตี
		return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
	})

	// จำลองหน้า monitoring/SIEM ที่ทีม security ใช้ดูเหตุการณ์
	app.Get("/audit", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"events": len(auditLog), "log": auditLog})
	})

	mode := "VULNERABLE ❌ (no security logging)"
	if secure {
		mode = "SECURE ✅ (audit log)"
	}
	log.Printf("A09 logging demo [%s] → POST /login, GET /audit", mode)
	log.Fatal(app.Listen(":3000"))
}
