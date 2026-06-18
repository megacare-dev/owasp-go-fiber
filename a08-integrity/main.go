// A08 — Software & Data Integrity Failures: webhook ไม่ verify ลายเซ็น
//
//	go run .            → VULNERABLE (เชื่อ payload โดยไม่ตรวจลายเซ็น → ปลอม event ได้)
//	SECURE=1 go run .   → SECURE     (verify HMAC-SHA256(body, secret) ก่อนเชื่อ)
//
// ยิง ./exploit.sh: ส่ง webhook ปลอม {"status":"paid"} โดยไม่มีลายเซ็นที่ถูกต้อง
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// secret ที่แชร์กับผู้ส่ง webhook (เช่น payment gateway)
var webhookSecret = []byte("shared-webhook-secret")

func main() {
	secure := os.Getenv("SECURE") == "1"
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/webhook", func(c *fiber.Ctx) error {
		body := c.Body()

		if secure {
			// ✅ ตรวจ integrity: ผู้ส่งต้องแนบ X-Signature = HMAC-SHA256(body, secret)
			mac := hmac.New(sha256.New, webhookSecret)
			mac.Write(body)
			expect := hex.EncodeToString(mac.Sum(nil))
			if !hmac.Equal([]byte(expect), []byte(c.Get("X-Signature"))) {
				return c.Status(401).JSON(fiber.Map{"error": "invalid signature"})
			}
		}
		// ❌ VULNERABLE: ประมวลผล payload โดยไม่ตรวจว่ามาจากแหล่งที่เชื่อถือจริง
		return c.JSON(fiber.Map{"processed": true, "raw": string(body)})
	})

	mode := "VULNERABLE ❌ (no signature check)"
	if secure {
		mode = "SECURE ✅ (HMAC verified)"
	}
	log.Printf("A08 integrity demo [%s] → POST /webhook", mode)
	log.Fatal(app.Listen(":3000"))
}
