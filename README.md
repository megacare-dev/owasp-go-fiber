# OWASP Top 10 — Go Fiber Runnable Demos 🧵

Runnable, **break → see → fix → verify** demos of the OWASP Top 10 for the
[Go Fiber](https://github.com/gofiber/fiber) web framework.

แต่ละบทมี **โค้ดที่มีช่องโหว่จริง** (`go run .`) คู่กับ **เวอร์ชันที่แก้แล้ว**
(`SECURE=1 go run .`) อยู่ในไฟล์เดียว สลับด้วย env — ยิง exploit เดิมซ้ำเห็นว่าโดนบล็อก
**ไม่ต้องเขียน Go เป็น**: แค่ `./demo.sh` แล้วดูผลเทียบสองโหมดในคำสั่งเดียว

## Quick start

```bash
git clone https://github.com/megacare-dev/owasp-go-fiber.git
cd owasp-go-fiber/a01-broken-access
./demo.sh
```

`./demo.sh` จะ: build → รันโหมดมีช่องโหว่ + ยิง exploit → รันโหมดแก้แล้ว + ยิงซ้ำ → สรุปผลเทียบ

ต้องมี: **Go 1.22+**. ส่วน scanner (ไม่บังคับ) สำหรับ `./scan.sh`:

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

## บทเรียน (OWASP Top 10)

| # | OWASP | โฟลเดอร์ | "เห็นมันพัง" ยังไง |
|---|-------|---------|---------------------|
| A01 | Broken Access Control | [`a01-broken-access/`](a01-broken-access/) | IDOR: alice อ่าน order ของ bob + non-admin ลบ user |
| A02 | Cryptographic Failures | [`a02-weak-token/`](a02-weak-token/) | คำนวณ `md5(username)` ยึดบัญชีทันที |
| A03 | Injection | [`a03-injection/`](a03-injection/) | `' OR '1'='1` ดึงผู้ใช้ทั้งตาราง (SQLi) |
| A04 | Insecure Design | [`a04-insecure-design/`](a04-insecure-design/) | ส่ง `price=1` ซื้อ iPhone 1 บาท |
| A05 | Security Misconfiguration | [`a05-misconfig/`](a05-misconfig/) | error ดิบ leak DB connection string |
| A06 | Vulnerable & Outdated Components | [`a06-vulnerable-components/`](a06-vulnerable-components/) | fiber v2.52.0 มี CVE → govulncheck |
| A07 | Identification & Auth Failures | [`a07-jwt/`](a07-jwt/) | ปลอม JWT `alg:none` เข้า /admin |
| A08 | Software & Data Integrity Failures | [`a08-integrity/`](a08-integrity/) | webhook ปลอม (ไม่ verify ลายเซ็น) |
| A09 | Logging & Monitoring Failures | [`a09-logging/`](a09-logging/) | brute-force แล้วไม่มี log ตรวจจับ |
| A10 | Server-Side Request Forgery | [`a10-ssrf/`](a10-ssrf/) | ดึง metadata ภายใน (169.254.169.254) |
| — | Bonus: CVE-2024-25124 | [`cors-cve/`](cors-cve/) | evil.com อ่าน response + credentials |

## โครงสร้างแต่ละ demo

```
<demo>/
├── main.go      โค้ด: โหมด vulnerable/secure ในไฟล์เดียว (สลับด้วย env SECURE=1)
├── exploit.sh   ยิงโจมตีจริง (DAST) — ใช้ยืนยันช่องโหว่
├── scan.sh      สแกน SAST/SCA (gosec / govulncheck)
├── demo.sh      รันคำสั่งเดียวจบ: เทียบ vulnerable ↔ secure
└── README.md    อธิบาย + บทเรียน + prompt สำหรับ vibe coding
```

## ⚠️ ความปลอดภัย

- demo **จงใจมีช่องโหว่** เพื่อการเรียนรู้ — รันบน **localhost เท่านั้น ห้าม deploy ออก public**
- ข้อมูลทั้งหมดเป็น **ของปลอม** (alice/bob, secret demo) ไม่ใช่ของจริง
- SSRF demo ดึง metadata จาก **mock server ในเครื่อง** ไม่ยิงออกภายนอก

## License

[MIT](LICENSE) — ใช้เป็นสื่อการสอน/อ้างอิงได้อิสระ
