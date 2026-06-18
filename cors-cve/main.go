// M4 / CVE-2024-25124 — CORS misconfiguration: reflect any origin + credentials
//
//	go run .            → VULNERABLE (สะท้อน Origin ใดก็ได้ + Allow-Credentials:true)
//	SECURE=1 go run .   → SECURE     (allow-list เฉพาะ origin ที่เชื่อถือ)
//
// ยิง ./exploit.sh ด้วย Origin: https://evil.com → ดู response CORS headers
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	secure := os.Getenv("SECURE") == "1"
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	if secure {
		// ✅ allow-list เฉพาะ origin ที่เชื่อถือ
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "https://app.example.com",
			AllowCredentials: true,
		}))
	} else {
		// ❌ VULNERABLE (CVE-2024-25124 pattern): สะท้อน origin ใดก็ได้ + ส่ง credentials
		app.Use(cors.New(cors.Config{
			AllowOriginsFunc: func(origin string) bool { return true }, // ยอมรับทุก origin
			AllowCredentials: true,
		}))
	}

	app.Get("/api/me", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"user": "bob", "balance": 100000})
	})

	mode := "VULNERABLE ❌ (reflect any origin + credentials)"
	if secure {
		mode = "SECURE ✅ (allow-list origin)"
	}
	log.Printf("CORS CVE demo [%s] → /api/me", mode)
	log.Fatal(app.Listen(":3000"))
}
