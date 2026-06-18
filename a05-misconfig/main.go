// A05 — Security Misconfiguration: verbose error leak
//
//	go run .            → VULNERABLE (ส่ง error ดิบกลับ → leak DSN/รหัสผ่าน DB ภายใน)
//	SECURE=1 go run .   → SECURE     (custom ErrorHandler → ข้อความกลางๆ + log ฝั่ง server)
//
// ยิง ./exploit.sh: เรียก endpoint ที่ error → vulnerable เผย connection string ภายใน
package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	secure := os.Getenv("SECURE") == "1"

	cfg := fiber.Config{DisableStartupMessage: true}
	if secure {
		// ✅ custom ErrorHandler: ไม่ส่งรายละเอียดภายในกลับให้ client
		cfg.ErrorHandler = func(c *fiber.Ctx, err error) error {
			log.Printf("internal error: %v", err) // log ฝั่ง server เท่านั้น
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		}
	}
	app := fiber.New(cfg)

	app.Get("/report/:id", func(c *fiber.Ctx) error {
		// จำลอง error จากชั้น DB ที่มี connection string ติดมาด้วย
		return errors.New(fmt.Sprintf(
			"db connect failed dsn=postgres://app:S3cretDbP@ss@10.0.0.5:5432/prod query report_id=%s",
			c.Params("id")))
		// ❌ VULNERABLE: default ErrorHandler ของ Fiber ส่ง err.Error() กลับไปตรงๆ
	})

	mode := "VULNERABLE ❌ (default error handler)"
	if secure {
		mode = "SECURE ✅ (safe error handler)"
	}
	log.Printf("A05 misconfig demo [%s] → http://localhost:3000/report/42", mode)
	log.Fatal(app.Listen(":3000"))
}
