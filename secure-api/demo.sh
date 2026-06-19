#!/usr/bin/env bash
# Secure API — ยิงทดสอบครบทุกชั้นป้องกัน (Module 5 verify) · รัน: ./demo.sh
set -u
cd "$(dirname "$0")"
PORT=8080
BIN="/tmp/secure-api.$$"

if lsof -ti "tcp:$PORT" >/dev/null 2>&1; then
  echo "⚠️  พอร์ต $PORT ค้างอยู่ — เคลียร์ให้..."
  lsof -ti "tcp:$PORT" | xargs kill -9 2>/dev/null || true; sleep 1
fi

echo "🔨 build..."
go build -o "$BIN" . || { echo "❌ build ไม่ผ่าน"; exit 1; }
trap 'rm -f "$BIN"' EXIT
JWT_SECRET=demo-secret "$BIN" >/dev/null 2>&1 &
SRV=$!
for _ in $(seq 1 40); do curl -s -o /dev/null "localhost:$PORT/health" && break; sleep 0.25; done

B="localhost:$PORT"
echo
echo "1) GET /health (public)                 → "; curl -s "$B/health"; echo
echo "2) POST /users ไม่มี token              → "; curl -s -o /dev/null -w "HTTP %{http_code}\n" -X POST "$B/users"
UTOK=$(curl -s -X POST "$B/login?user=alice" | sed -E 's/.*"token":"([^"]+)".*/\1/')
echo "3) POST /users ด้วย token user ธรรมดา   → "; curl -s -o /dev/null -w "HTTP %{http_code}\n" -X POST "$B/users" -H "Authorization: Bearer $UTOK"
ATOK=$(curl -s -X POST "$B/login?user=admin" | sed -E 's/.*"token":"([^"]+)".*/\1/')
echo "4) POST /users (admin) email ผิด        → "; curl -s -w " [HTTP %{http_code}]\n" -X POST "$B/users" -H "Authorization: Bearer $ATOK" -H "Content-Type: application/json" -d '{"email":"not-an-email","role":"user"}'
echo "5) POST /users (admin) ถูกต้อง          → "; curl -s -w " [HTTP %{http_code}]\n" -X POST "$B/users" -H "Authorization: Bearer $ATOK" -H "Content-Type: application/json" -d '{"email":"new@corp.com","role":"user"}'
echo "6) security headers (helmet)            → "; curl -sI "$B/health" | grep -iE "x-frame-options|x-content-type-options|x-xss|content-security|strict" | sed 's/^/   /'

kill "$SRV" 2>/dev/null; wait "$SRV" 2>/dev/null || true
echo
echo "📌 สรุป: health=200 · no-token=401 · user=403 · bad-input=422 · admin-ok=201 · มี security headers"
