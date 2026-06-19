// A04 — Insecure Design: เชื่อราคาที่ client ส่งมา (business-logic flaw)
//
//	go run .            → VULNERABLE (คิดเงินตาม price ที่ client ส่ง → จ่าย 1 บาทก็ได้)
//	SECURE=1 go run .   → SECURE     (ดึงราคาจาก catalog ฝั่ง server เสมอ)
//
// ยิง ./exploit.sh: สั่งซื้อ iphone โดยแนบ price=1 → vulnerable คิดเงินแค่ 1 บาท
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// ราคาจริงฝั่ง server (บาท) — แหล่งความจริงเพียงแหล่งเดียว
var catalog = map[string]int{"iphone": 39900, "macbook": 89900}

type Order struct {
	Item  string `json:"item"`
	Qty   int    `json:"qty"`
	Price int    `json:"price"` // ราคาที่ client ส่งมา — ไม่ควรเชื่อ
}

func main() {
	secure := os.Getenv("SECURE") == "1"
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/checkout", func(c *fiber.Ctx) error {
		var o Order
		if err := c.BodyParser(&o); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "bad request"})
		}

		// ✅ validate qty เสมอ: กัน qty<=0 (ของฟรี/ยอดติดลบ) — business-logic flaw
		if o.Qty <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "invalid qty"})
		}

		unit := o.Price // ❌ VULNERABLE: ใช้ราคาที่ client ส่งมา → ตั้ง 1 บาทก็ได้
		if secure {
			// ✅ ออกแบบให้ปลอดภัย: ราคามาจาก server เท่านั้น ไม่สนค่าที่ client ส่ง
			p, ok := catalog[o.Item]
			if !ok {
				return c.Status(400).JSON(fiber.Map{"error": "unknown item"})
			}
			unit = p
		}
		return c.JSON(fiber.Map{"item": o.Item, "qty": o.Qty, "charged": unit * o.Qty})
	})

	mode := "VULNERABLE ❌ (trust client price)"
	if secure {
		mode = "SECURE ✅ (server-side price)"
	}
	log.Printf("A04 insecure-design demo [%s] → POST /checkout", mode)
	log.Fatal(app.Listen(":3000"))
}
