// A07 — Identification & Authentication Failures: forge JWT ด้วย alg:none
//
//	go run .            → VULNERABLE (เชื่อ payload โดยไม่ verify signature → alg:none ผ่าน)
//	SECURE=1 go run .   → SECURE     (บังคับ HS256 + verify HMAC ด้วย secret)
//
// ยิง ./exploit.sh: ปลอม token {"role":"admin"} แบบไม่เซ็น → เข้า /admin ได้
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var secret = []byte("super-secret-key")

func b64dec(s string) []byte {
	b, _ := base64.RawURLEncoding.DecodeString(s)
	return b
}

// verifyToken คืน (role, ok)
func verifyToken(tok string, secure bool) (string, bool) {
	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		return "", false
	}
	var hdr struct{ Alg string }
	json.Unmarshal(b64dec(parts[0]), &hdr)

	if secure {
		// ✅ บังคับ alg HS256 และ verify signature
		if hdr.Alg != "HS256" {
			return "", false
		}
		mac := hmac.New(sha256.New, secret)
		mac.Write([]byte(parts[0] + "." + parts[1]))
		expect := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(expect), []byte(parts[2])) {
			return "", false
		}
	}
	// ❌ VULNERABLE (secure=false): เชื่อ payload โดยไม่ verify signature เลย
	var claims struct{ Role string }
	json.Unmarshal(b64dec(parts[1]), &claims)
	return claims.Role, true
}

func main() {
	secure := os.Getenv("SECURE") == "1"
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Get("/admin", func(c *fiber.Ctx) error {
		tok := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		role, ok := verifyToken(tok, secure)
		if !ok || role != "admin" {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}
		return c.JSON(fiber.Map{"ok": true, "msg": "welcome admin", "secret": "ALL_USER_DATA"})
	})

	mode := "VULNERABLE ❌ (no signature check)"
	if secure {
		mode = "SECURE ✅ (HS256 verified)"
	}
	log.Printf("A07 JWT demo [%s] → GET /admin (Bearer token)", mode)
	log.Fatal(app.Listen(":3000"))
}
