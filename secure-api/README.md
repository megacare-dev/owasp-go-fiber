# Secure Fiber REST API — Capstone (Module 5)

โค้ดอ้างอิงสำหรับช่วง **"สร้าง Secure REST API ด้วย Vibe Coding"** — รวมการป้องกัน OWASP
ที่เรียนมาทั้งวันไว้ในไฟล์เดียว เพื่อให้ผู้สอน **ทำให้ดู (live demo)** ทีละ step

| ชั้นป้องกัน | OWASP | ในโค้ด |
|------------|-------|--------|
| RBAC (ตรวจ role) | A01 | `requireRole("admin")` |
| secret สุ่ม | A02 | `crypto/rand` ใน `getSecret()` |
| validate input | A04 | `go-playground/validator` |
| secure error + headers | A05 | custom `ErrorHandler` + `helmet` |
| JWT (HS256+exp, กัน alg:none) | A07 | `authJWT` |
| logging | A09 | `logger` middleware |
| กัน brute-force | — | `limiter` middleware |

## รัน

```bash
cd owasp-go-fiber/secure-api
go run .          # http://localhost:8080
./demo.sh         # ยิงทดสอบครบ: 200 / 401 / 403 / 422 / 201 + security headers
```

## พฤติกรรมที่ถูกต้อง

| คำขอ | ผล |
|------|-----|
| `GET /health` | `200 {"status":"ok"}` |
| `POST /users` ไม่มี token | `401` |
| `POST /users` token user ธรรมดา | `403 forbidden` |
| `POST /users` admin + email ผิด | `422 invalid input` |
| `POST /users` admin + ถูกต้อง | `201 created` |

> สอน step-by-step ดู [`RUNSHEET.md`](RUNSHEET.md) · prod จริงให้เพิ่ม TLS (`:8443`)
