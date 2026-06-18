// A06 — Vulnerable & Outdated Components: ใช้ dependency ที่มี CVE
//
// demo นี้ pin gofiber/fiber v2.52.0 ซึ่งมีช่องโหว่ที่รู้แล้ว (เช่น CVE-2024-25124
// ที่ CORS middleware) — แอปสืบทอดช่องโหว่ของ lib ทันทีโดยไม่รู้ตัว
//
//	❌ ปัญหา : ไม่ได้อัปเดต/สแกน dependency
//	✅ วิธีแก้: สแกนด้วย govulncheck แล้วอัปเดตเป็นเวอร์ชันที่แก้แล้ว (≥ v2.52.5)
//
// ตรวจด้วย: ./scan.sh   (govulncheck)   หรือดูสรุป: ./demo.sh
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	// ใช้ CORS middleware ของ fiber v2.52.0 — โค้ดนี้แตะฟังก์ชันที่มี CVE-2024-25124
	// (govulncheck จะตามรอย call path เจอช่องโหว่ผ่านการเรียกนี้)
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"fiber": "v2.52.0 (outdated) — รัน ./scan.sh เพื่อดู CVE"})
	})

	log.Printf("A06 outdated-components demo → ใช้ govulncheck ตรวจ (ดู ./scan.sh)")
	log.Fatal(app.Listen(":3000"))
}
