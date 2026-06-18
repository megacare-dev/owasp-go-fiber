// A10 — Server-Side Request Forgery (SSRF) demo
//
//	go run .            → VULNERABLE (fetch URL ที่ client ส่งมาตรงๆ → เข้าถึง internal ได้)
//	SECURE=1 go run .   → SECURE     (บล็อก loopback/private/metadata ก่อน fetch)
//
// มี mock "internal metadata" server บน 127.0.0.1:9999 (จำลอง 169.254.169.254)
// ยิง ./exploit.sh: หลอกให้ server ไปดึง credential ภายในมาให้
package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// บล็อก loopback / private / link-local(metadata 169.254.169.254)
func isBlocked(host string) bool {
	ips, err := net.LookupIP(host)
	if err != nil {
		return true // resolve ไม่ได้ = บล็อกไว้ก่อน
	}
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() {
			return true
		}
	}
	return false
}

func main() {
	secure := os.Getenv("SECURE") == "1"

	// mock internal metadata service (แทน cloud metadata 169.254.169.254)
	go func() {
		m := http.NewServeMux()
		m.HandleFunc("/latest/meta-data/iam/credentials", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/EXAMPLE/INTERNAL")
		})
		http.ListenAndServe("127.0.0.1:9999", m)
	}()
	time.Sleep(150 * time.Millisecond)

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Get("/fetch", func(c *fiber.Ctx) error {
		target := c.Query("url")
		host := target
		if i := strings.Index(host, "://"); i >= 0 {
			host = host[i+3:]
		}
		if j := strings.IndexAny(host, ":/"); j >= 0 {
			host = host[:j]
		}

		if secure && isBlocked(host) {
			// ✅ บล็อกปลายทางภายในก่อนยิง request
			return c.Status(400).JSON(fiber.Map{"error": "blocked: internal address"})
		}

		// ❌ VULNERABLE: ยิงไปยัง URL ที่ client กำหนด โดยไม่ตรวจปลายทาง
		resp, err := http.Get(target)
		if err != nil {
			return c.Status(502).JSON(fiber.Map{"error": err.Error()})
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return c.SendString(string(body))
	})

	mode := "VULNERABLE ❌ (no destination check)"
	if secure {
		mode = "SECURE ✅ (block internal)"
	}
	log.Printf("A10 SSRF demo [%s] → /fetch?url=...", mode)
	log.Fatal(app.Listen(":3000"))
}
