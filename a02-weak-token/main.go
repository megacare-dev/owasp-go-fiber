// A02 — Cryptographic Failures: password-reset token ที่ "เดาได้"
//
//	go run .            → VULNERABLE (token = md5(username) → attacker คำนวณเองได้)
//	SECURE=1 go run .   → SECURE     (token = crypto/rand 32 bytes → เดาไม่ได้)
//
// ยิง ./exploit.sh: attacker คำนวณ token ของ bob เองแล้วยึดบัญชีทันที (ไม่ต้อง brute-force)
package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

var resetTokens = map[string]string{} // token -> user

func makeToken(user string, secure bool) string {
	if secure {
		b := make([]byte, 32)
		rand.Read(b) // ✅ crypto/rand: ไม่ผูกกับข้อมูลใด เดาไม่ได้
		return hex.EncodeToString(b)
	}
	// ❌ VULNERABLE: token = md5(username) → ใครรู้ username ก็คำนวณ token ได้
	h := md5.Sum([]byte(user))
	return hex.EncodeToString(h[:])
}

func main() {
	secure := os.Getenv("SECURE") == "1"
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	// ขอ reset token (ปกติส่งเข้าอีเมล — attacker ไม่เห็นค่านี้)
	app.Post("/forgot", func(c *fiber.Ctx) error {
		user := c.Query("user", "bob")
		resetTokens[makeToken(user, secure)] = user
		return c.JSON(fiber.Map{"sent": true, "user": user}) // ไม่คืน token
	})

	// ใช้ token รีเซ็ตรหัส = ยึดบัญชี ถ้า token ถูก
	app.Post("/reset", func(c *fiber.Ctx) error {
		if user, ok := resetTokens[c.Query("token")]; ok {
			return c.JSON(fiber.Map{"ok": true, "took_over": user})
		}
		return c.Status(401).JSON(fiber.Map{"ok": false})
	})

	mode := "VULNERABLE ❌ (md5 username)"
	if secure {
		mode = "SECURE ✅ (crypto/rand)"
	}
	log.Printf("A02 weak-token demo [%s] → http://localhost:3000", mode)
	log.Fatal(app.Listen(":3000"))
}
