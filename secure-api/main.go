// Secure Fiber REST API — capstone (Module 5: สอนแบบทำให้ดู / live demo)
//
// รวมการป้องกัน OWASP ที่เรียนมาทั้งวันไว้ในไฟล์เดียว:
//   A01 RBAC · A02 secret crypto/rand · A04 validate input ·
//   A05 secure error + helmet headers · A07 JWT (HS256+exp, กัน alg:none) ·
//   A09 logging · limiter กัน brute-force
//
//   go run .          → ฟังที่ :8080
//   ./demo.sh         → ยิงทดสอบครบ (200/401/403/422/201 + headers)
//   ดู RUNSHEET.md     → ลำดับสอนทีละ step
package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = []byte(getSecret())
	validate  = validator.New()
)

// ✅ A02: secret สุ่มจาก crypto/rand ถ้าไม่ได้ตั้ง env (ห้าม hardcode)
func getSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ✅ A07: ออก JWT HS256 + claims sub/role/exp (หมดอายุ 1 ชม.)
func issueToken(sub, role string) (string, error) {
	claims := jwt.MapClaims{"sub": sub, "role": role, "exp": time.Now().Add(time.Hour).Unix()}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

// ✅ A07: middleware ตรวจ JWT แล้ว set role ลง Locals (บังคับ HMAC กัน alg:none)
func authJWT(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return c.Status(401).JSON(fiber.Map{"error": "missing token"})
	}
	tok, err := jwt.Parse(strings.TrimPrefix(auth, "Bearer "), func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil || !tok.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
	}
	claims := tok.Claims.(jwt.MapClaims)
	c.Locals("role", claims["role"])
	c.Locals("sub", claims["sub"])
	return c.Next()
}

// ✅ A01: RBAC — ตรวจ role ก่อนเข้า handler
func requireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("role") != role {
			return c.Status(403).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}

// ✅ A04: validate input ด้วย struct tag
type CreateUser struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=user admin"`
}

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		// ✅ A05: secure error handler — log ฝั่ง server, ไม่ leak ให้ client
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("internal error: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		},
	})

	// ✅ ลำดับ middleware สำคัญ: recover → helmet → logger → limiter
	app.Use(recover.New())                                                       // กัน panic ทำ server ล่ม
	app.Use(helmet.New())                                                        // ✅ A05: security headers
	app.Use(logger.New())                                                        // ✅ A09: logging
	app.Use(limiter.New(limiter.Config{Max: 20, Expiration: 30 * time.Second})) // กัน brute-force/DoS

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// login (demo): ออก token — เป็น admin เมื่อ user=admin
	app.Post("/login", func(c *fiber.Ctx) error {
		user := c.Query("user", "guest")
		role := "user"
		if user == "admin" {
			role = "admin"
		}
		t, _ := issueToken(user, role)
		return c.JSON(fiber.Map{"token": t, "role": role})
	})

	// ✅ ป้องกันครบชั้น: ต้องมี JWT + เป็น admin + ผ่าน validate
	app.Post("/users", authJWT, requireRole("admin"), func(c *fiber.Ctx) error {
		var in CreateUser
		if err := c.BodyParser(&in); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "bad request"})
		}
		if err := validate.Struct(in); err != nil {
			return c.Status(422).JSON(fiber.Map{"error": "invalid input"})
		}
		return c.Status(201).JSON(fiber.Map{"created": in.Email, "role": in.Role})
	})

	log.Println("Secure API on http://localhost:8080  (prod: ใส่ TLS → :8443)")
	log.Fatal(app.Listen(":8080"))
}
