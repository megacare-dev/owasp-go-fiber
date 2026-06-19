// M4 / CVE-2024-25124 — CORS misconfiguration + Fiber security middleware stack
//
// CORS demo (เดิม):
//	go run .            → VULNERABLE (สะท้อน Origin ใดก็ได้ + Allow-Credentials:true)
//	SECURE=1 go run .   → SECURE     (allow-list เฉพาะ origin ที่เชื่อถือ)
//
// Middleware stack (เพิ่มใหม่): recover · helmet · session(30m) · keyauth(validAPIKey)
//	GET /api/me          → public (ใช้สาธิต CORS เดิม)
//	GET /api/secure/me   → ต้องมี header X-API-Key ที่ valid (keyauth) + นับ visit ผ่าน session
package main

import (
	"crypto/subtle"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
)

// ✅ keyauth validator: เทียบ API key แบบ constant-time (กัน timing attack)
// คีย์อ้างอิงมาจาก env API_KEY (ถ้าไม่ตั้ง ใช้ค่า demo — prod ต้องตั้ง env เสมอ)
func validAPIKey(c *fiber.Ctx, key string) (bool, error) {
	expected := os.Getenv("API_KEY")
	if expected == "" {
		expected = "demo-secret-key" // demo เท่านั้น
	}
	if subtle.ConstantTimeCompare([]byte(key), []byte(expected)) == 1 {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}

func main() {
	secure := os.Getenv("SECURE") == "1"
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	// ✅ middleware ระดับ global (ลำดับสำคัญ: recover ก่อน เพื่อครอบ panic ของ handler/มิดเดิลแวร์ถัดไป)
	app.Use(recover.New()) // กัน panic ทำ server ล่ม
	app.Use(helmet.New())  // security headers (X-Frame-Options, nosniff, ฯลฯ)

	if secure {
		// ✅ allow-list เฉพาะ origin ที่เชื่อถือ (อ่านจาก env, ไม่ hardcode/ไม่ wildcard)
		// ตั้งค่าได้: ALLOWED_ORIGINS="https://app.example.com,https://admin.example.com"
		allowed := os.Getenv("ALLOWED_ORIGINS")
		if allowed == "" {
			allowed = "https://app.example.com"
		}
		app.Use(cors.New(cors.Config{
			// ปลอดภัย: ระบุ origin ชัดเจน + credentials (ห้ามใช้ "*" คู่ credentials → CVE-2024-25124)
			AllowOrigins:     allowed,
			AllowCredentials: true,
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		}))
	} else {
		// ❌ VULNERABLE (CVE-2024-25124 pattern): สะท้อน origin ใดก็ได้ + ส่ง credentials
		app.Use(cors.New(cors.Config{
			AllowOriginsFunc: func(origin string) bool { return true }, // ยอมรับทุก origin
			AllowCredentials: true,
		}))
	}

	// ✅ session store — หมดอายุ 30 นาที (idle timeout)
	store := session.New(session.Config{
		Expiration: 30 * time.Minute,
	})

	// public — ใช้สาธิต CORS เดิม (exploit.sh/demo.sh ยิงเส้นนี้)
	app.Get("/api/me", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"user": "bob", "balance": 100000})
	})

	// ✅ protected group — ต้องมี X-API-Key ที่ valid (keyauth) แล้วนับ visit ลง session
	secured := app.Group("/api/secure", keyauth.New(keyauth.Config{
		KeyLookup: "header:X-API-Key",
		Validator: validAPIKey,
	}))
	secured.Get("/me", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return err
		}
		visits, _ := sess.Get("visits").(int)
		visits++
		sess.Set("visits", visits)
		if err := sess.Save(); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"user": "bob", "balance": 100000, "visits": visits})
	})

	mode := "VULNERABLE ❌ (reflect any origin + credentials)"
	if secure {
		mode = "SECURE ✅ (allow-list origin)"
	}
	log.Printf("CORS CVE demo [%s] → /api/me", mode)
	log.Printf("middleware: recover · helmet · session(30m) · keyauth → /api/secure/me (header X-API-Key)")
	log.Fatal(app.Listen(":3000"))
}
