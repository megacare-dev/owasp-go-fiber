// A10 — Server-Side Request Forgery (SSRF) demo
//
//	go run .            → VULNERABLE (fetch URL ที่ client ส่งมาตรงๆ → เข้าถึง internal ได้)
//	SECURE=1 go run .   → SECURE     (บล็อก loopback/private/metadata + ปิด redirect + กัน DNS rebinding)
//
// มี mock "internal metadata" server บน 127.0.0.1:9999 (จำลอง 169.254.169.254)
// ยิง ./exploit.sh: หลอกให้ server ไปดึง credential ภายในมาให้
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// isBlockedIP — true ถ้าเป็นปลายทางภายในที่ไม่ควรให้ server ยิงไปถึง
// (loopback / private / link-local เช่น cloud metadata 169.254.169.254 / unspecified)
func isBlockedIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified()
}

// isBlockedHost — resolve host แล้วเช็คทุก IP (fail-closed ถ้า resolve ไม่ได้)
func isBlockedHost(host string) bool {
	ips, err := net.LookupIP(host)
	if err != nil {
		return true // resolve ไม่ได้ = บล็อกไว้ก่อน
	}
	for _, ip := range ips {
		if isBlockedIP(ip) {
			return true
		}
	}
	return false
}

// secureClient ป้องกัน SSRF เชิงลึก 2 ชั้น:
//  1. ปิด redirect — กัน URL ที่อนุญาต 302 เด้งไปปลายทางภายใน
//  2. ตรวจ IP "จริง" ตอน dial — กัน DNS rebinding/TOCTOU (resolve ตอนเช็คกับตอนยิงเป็นคนละ IP)
func secureClient() *http.Client {
	dialer := &net.Dialer{Timeout: 3 * time.Second}
	return &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse // ไม่ตาม redirect
		},
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}
				ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
				if err != nil {
					return nil, err
				}
				for _, ip := range ips {
					if isBlockedIP(ip.IP) {
						return nil, fmt.Errorf("blocked: internal address %s", ip.IP)
					}
				}
				// เชื่อมต่อไปยัง IP ที่ตรวจแล้วเท่านั้น (ไม่ resolve ใหม่)
				return dialer.DialContext(ctx, network, net.JoinHostPort(ips[0].IP.String(), port))
			},
		},
	}
}

// sendResponse ส่ง body ของ response กลับให้ client
func sendResponse(c *fiber.Ctx, resp *http.Response) error {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return c.SendString(string(body))
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

		if secure {
			// ✅ parse URL ให้ถูกต้อง (กัน trick แบบ user@host) + อนุญาตเฉพาะ http/https
			u, err := url.Parse(target)
			if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
				return c.Status(400).JSON(fiber.Map{"error": "invalid url"})
			}
			// เช็คปลายทางก่อนยิง — defense ชั้นแรก (ชั้นสองอยู่ใน secureClient ตอน dial)
			if isBlockedHost(u.Hostname()) {
				return c.Status(400).JSON(fiber.Map{"error": "blocked: internal address"})
			}
			resp, err := secureClient().Get(target)
			if err != nil {
				return c.Status(502).JSON(fiber.Map{"error": err.Error()})
			}
			return sendResponse(c, resp)
		}

		// ❌ VULNERABLE: ยิงไปยัง URL ที่ client กำหนด โดยไม่ตรวจปลายทาง + ตาม redirect
		resp, err := http.Get(target)
		if err != nil {
			return c.Status(502).JSON(fiber.Map{"error": err.Error()})
		}
		return sendResponse(c, resp)
	})

	mode := "VULNERABLE ❌ (no destination check)"
	if secure {
		mode = "SECURE ✅ (block internal + no redirect + anti-rebinding)"
	}
	log.Printf("A10 SSRF demo [%s] → /fetch?url=...", mode)
	log.Fatal(app.Listen(":3000"))
}
