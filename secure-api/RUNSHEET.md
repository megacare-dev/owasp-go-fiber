# 🎤 RUNSHEET — Module 5: สร้าง Secure REST API (สอนแบบทำให้ดู)

คู่มือผู้สอน **live demo** — พิมพ์/รันทีละ step ให้ผู้เรียนดู (ไม่ปล่อยทำเอง)
แต่ละ step = **พูดอะไร → prompt ให้ AI → โค้ดที่ได้ → รันยืนยัน**

> เตรียมก่อนเริ่ม: เปิด terminal 2 บาน · `cd owasp-go-fiber/secure-api` · เปิดไฟล์ `main.go` ค้างไว้
> ถ้าเวลาน้อย: เปิดโค้ดเสร็จแล้ว walk-through ทีละส่วน + รัน `./demo.sh` ปิดท้าย

---

## Step 0 · เกริ่น (1 นาที)
🗣️ "ช่วงนี้เราจะรวมทุกอย่างที่เรียนมา (A01–A10) มาสร้าง API ปลอดภัยจริง 1 ตัว — สั่ง AI ทีละชั้น แล้ว **review ทุกครั้ง**"
- เป้าหมาย: `POST /users` ที่ create user ได้เฉพาะ admin + ตรวจ input + ไม่ leak

## Step 1 · Scaffold + /health (5 นาที)
🗣️ "เริ่มจากโครงเปล่า + endpoint สุขภาพก่อน"

💬 prompt: *"สร้าง Fiber app ไฟล์เดียว ฟัง :8080 มี GET /health คืน {status:ok} พร้อม go.mod"*

```bash
go mod tidy
go run .
# อีก terminal:
curl -s localhost:8080/health        # → {"status":"ok"}
```
✅ ชี้: ได้ baseline ที่รันได้ก่อน ค่อยเติมความปลอดภัย

## Step 2 · Security middleware + headers (8 นาที)  ← A05
🗣️ "ก่อนทำ logic ใส่เกราะชั้นนอกก่อน — ลำดับสำคัญ"

💬 prompt: *"เพิ่ม Fiber middleware ตามลำดับ recover → helmet → logger → limiter (20 req/30s)"*

```go
app.Use(recover.New())   // กัน panic ล่ม
app.Use(helmet.New())    // security headers
app.Use(logger.New())    // A09 logging
app.Use(limiter.New(limiter.Config{Max: 20, Expiration: 30 * time.Second}))
```
```bash
curl -sI localhost:8080/health | grep -i "x-frame\|x-content-type"
# → เห็น X-Frame-Options, X-Content-Type-Options (helmet ใส่ให้)
```
✅ ชี้: headers โผล่มาเองจาก helmet · ❌ จับผิด AI: ถ้าวาง recover ไว้ท้าย จะกัน panic ไม่ครบ

## Step 3 · JWT login + verify (8 นาที)  ← A07
🗣️ "ต่อไปคือ 'รู้ว่าใคร' — ออก JWT ตอน login แล้ว verify ทุก request"

💬 prompt: *"เพิ่ม POST /login ออก JWT HS256 (claims sub, role, exp) secret จาก env, และ middleware authJWT ที่ verify แล้ว set c.Locals(\"role\") — บังคับ HMAC กัน alg:none"*

```go
tok, _ := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
    if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {  // ✅ กัน alg:none
        return nil, jwt.ErrSignatureInvalid
    }
    return jwtSecret, nil
})
```
```bash
curl -s -X POST "localhost:8080/login?user=admin"   # → {"token":"...","role":"admin"}
```
✅ ชี้: ย้อนกลับไป A07 — ถ้าไม่บังคับ HMAC จะโดน alg:none · secret อ่านจาก env (A02)

## Step 4 · RBAC (5 นาที)  ← A01
🗣️ "รู้ว่าใครแล้ว ต้องเช็คว่า 'มีสิทธิ์ไหม'"

💬 prompt: *"เพิ่ม middleware requireRole(role) อ่าน c.Locals(\"role\") ไม่ตรงตอบ 403 แล้วใช้กับ POST /users ที่ต้องเป็น admin"*

```go
app.Post("/users", authJWT, requireRole("admin"), createUser)
```
```bash
UTOK=$(curl -s -X POST "localhost:8080/login?user=alice" | sed -E 's/.*"token":"([^"]+)".*/\1/')
curl -s -o /dev/null -w "%{http_code}\n" -X POST localhost:8080/users -H "Authorization: Bearer $UTOK"   # → 403
```
✅ ชี้: user ธรรมดาโดน 403 (A01 privilege escalation ถูกกัน)

## Step 5 · Validate input (5 นาที)  ← A04
🗣️ "อย่าเชื่อ input — ตรวจก่อนใช้"

💬 prompt: *"ใน /users ใช้ go-playground/validator ตรวจ email (required,email) และ role (oneof=user admin) ผิดตอบ 422"*

```go
type CreateUser struct {
    Email string `json:"email" validate:"required,email"`
    Role  string `json:"role"  validate:"required,oneof=user admin"`
}
if err := validate.Struct(in); err != nil {
    return c.Status(422).JSON(fiber.Map{"error": "invalid input"})
}
```
```bash
ATOK=$(curl -s -X POST "localhost:8080/login?user=admin" | sed -E 's/.*"token":"([^"]+)".*/\1/')
curl -s -X POST localhost:8080/users -H "Authorization: Bearer $ATOK" \
  -H "Content-Type: application/json" -d '{"email":"bad","role":"user"}'   # → 422
```

## Step 6 · Secure error handler (4 นาที)  ← A05
🗣️ "ปิดท้ายความปลอดภัย — error ห้าม leak ภายใน"

💬 prompt: *"ตั้ง fiber.Config.ErrorHandler ที่ log ฝั่ง server แต่คืน client แค่ {error:\"internal server error\"} (500)"*

✅ ชี้: ย้อน A05 — default handler ของ Fiber ส่ง error ดิบ ต้อง override

## Step 7 · ยิงทดสอบครบ (5 นาที)
🗣️ "รวมทุกชั้น ทดสอบทีเดียว"
```bash
./demo.sh
# → health 200 · no-token 401 · user 403 · bad-input 422 · admin 201 · + headers
```

## Step 8 · Verify (Module 5 หน้า 2 · 20 นาที)
🗣️ "ขั้น Verify ของ vibe coding — อย่าจบแค่ 'AI บอกว่าเสร็จ'"
```bash
curl -sI localhost:8080/health | grep -iE 'x-frame|x-content-type|x-xss'   # headers
curl -i -X POST localhost:8080/users -H "Authorization: Bearer $UTOK"      # authz ต้อง 403
zap-baseline.py -t http://localhost:8080                                    # (ถ้ามี ZAP) baseline scan
```
💬 prompt (อ่านผล): *"อ่านผล OWASP ZAP นี้ จัดลำดับความเสี่ยงตาม OWASP แล้วเสนอวิธีแก้เป็นโค้ด Fiber"*

---

### ⏱ รวม ~55 นาที (พอดี Module 5 + เริ่ม Verify) · 🔑 ทุก step ย้ำ: สั่ง AI → **review/verify เอง**
