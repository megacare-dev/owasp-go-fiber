// A03 — Injection (SQL Injection) demo · SQLite pure-Go (ไม่ต้อง CGO/Docker)
//
//	go run .            → VULNERABLE (ต่อ string เข้า query → ' OR '1'='1 ดึงทั้งตาราง)
//	SECURE=1 go run .   → SECURE     (parameterized query → input เป็นข้อมูล ไม่ใช่ SQL)
//
// ยิง ./exploit.sh: ส่ง name=' OR '1'='1 → vulnerable คืน user ทุกคน + secret
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "modernc.org/sqlite"
)

func main() {
	secure := os.Getenv("SECURE") == "1"

	db, _ := sql.Open("sqlite", ":memory:")
	db.Exec(`CREATE TABLE users(id INTEGER, name TEXT, secret TEXT)`)
	db.Exec(`INSERT INTO users VALUES
		(1,'alice','alice_api_key_AAA'),
		(2,'bob','bob_api_key_BBB'),
		(3,'admin','ADMIN_ROOT_TOKEN_9f3x')`)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Get("/search", func(c *fiber.Ctx) error {
		name := c.Query("name")
		var rows *sql.Rows
		var err error

		if secure {
			// ✅ parameterized: input ถูกส่งแยกจาก SQL — ไม่ถูกตีความเป็นคำสั่ง
			rows, err = db.Query("SELECT id,name,secret FROM users WHERE name = ?", name)
		} else {
			// ❌ VULNERABLE: ต่อ string ตรงๆ → attacker แทรก SQL ได้
			q := fmt.Sprintf("SELECT id,name,secret FROM users WHERE name = '%s'", name)
			rows, err = db.Query(q)
		}
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		out := []fiber.Map{}
		for rows.Next() {
			var id int
			var n, s string
			rows.Scan(&id, &n, &s)
			out = append(out, fiber.Map{"id": id, "name": n, "secret": s})
		}
		return c.JSON(out)
	})

	mode := "VULNERABLE ❌ (string concat)"
	if secure {
		mode = "SECURE ✅ (parameterized)"
	}
	log.Printf("A03 SQLi demo [%s] → http://localhost:3000/search?name=alice", mode)
	log.Fatal(app.Listen(":3000"))
}
