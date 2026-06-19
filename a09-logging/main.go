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

// auditLog จำลอง centralized log / SIEM — ในงานจริงส่งเข้า SIEM แล้วตั้ง alert
var auditLog []string

// recordFailedLogin บันทึกเหตุการณ์ login ล้มเหลว "เฉพาะโหมด secure"
// เก็บ user + ip แต่ ❗ไม่รับ/ไม่เก็บรหัสผ่าน — คืน true ถ้าบันทึกจริง
// โหมด vulnerable ❌ จะล้มเหลวเงียบๆ ไม่มีร่องรอย → ตรวจจับ brute-force ไม่ได้
func recordFailedLogin(user, ip string, secure bool) bool {
	if !secure {
		return false
	}
	auditLog = append(auditLog, fmt.Sprintf("FAILED_LOGIN user=%s ip=%s", user, ip))
	return true
}

func main() {
	secure := os.Getenv("SECURE") == "1"
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/login", func(c *fiber.Ctx) error {
		user := c.Query("user", "admin")
		// demo: สมมติรหัสผ่านผิดทุกครั้ง (จำลองการ brute-force)
		recordFailedLogin(user, c.IP(), secure)
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
