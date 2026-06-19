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

// secret อ่านจาก env (production) — มี default ไว้ให้ demo รันได้ทันที
// ห้าม hardcode secret จริงในโค้ด: ใช้ env / secret manager เสมอ
func secret() []byte {
	if s := os.Getenv("WEBHOOK_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte("shared-webhook-secret")
}

// sign คืน HMAC-SHA256(body, secret) เป็น hex — ผู้ส่ง webhook ที่ถูกต้อง
// (เช่น payment gateway) ต้องแนบค่านี้มาใน header X-Signature
func sign(body []byte) string {
	mac := hmac.New(sha256.New, secret())
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// verifySignature ✅ — ตรวจ integrity ของ payload ด้วย hmac.Equal (constant-time)
// ลายเซ็นที่หายไป/ผิด/body ถูกแก้ จะไม่ผ่าน → กันทั้ง event ปลอมและ timing attack
func verifySignature(body []byte, sig string) bool {
	return hmac.Equal([]byte(sign(body)), []byte(sig))
}

// verifyWebhook เลือกพฤติกรรมตามโหมด:
//   - secure       → เชื่อเฉพาะ payload ที่มีลายเซ็น HMAC ถูกต้องเท่านั้น
//   - vulnerable ❌ → เชื่อ payload ทันที โดยไม่ตรวจว่ามาจากแหล่งที่เชื่อถือจริง
func verifyWebhook(body []byte, sig string, secure bool) bool {
	if secure {
		return verifySignature(body, sig)
	}
	return true
}

func main() {
	secure := os.Getenv("SECURE") == "1"
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/webhook", func(c *fiber.Ctx) error {
		body := c.Body()
		if !verifyWebhook(body, c.Get("X-Signature"), secure) {
			return c.Status(401).JSON(fiber.Map{"error": "invalid signature"})
		}
		return c.JSON(fiber.Map{"processed": true, "raw": string(body)})
	})

	mode := "VULNERABLE ❌ (no signature check)"
	if secure {
		mode = "SECURE ✅ (HMAC verified)"
	}
	log.Printf("A08 integrity demo [%s] → POST /webhook", mode)
	log.Fatal(app.Listen(":3000"))
}
