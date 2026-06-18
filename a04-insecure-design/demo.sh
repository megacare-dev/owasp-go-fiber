#!/usr/bin/env bash
# demo คำสั่งเดียวจบ (สำหรับผู้เรียนที่ไม่เคยเขียน Go) — รัน: ./demo.sh
# มันจะ: รันโหมดมีช่องโหว่ → ยิง exploit → รันโหมดแก้แล้ว → ยิง exploit เดิมซ้ำ → เทียบผล
set -u
cd "$(dirname "$0")"

PORT=3000
BIN="/tmp/owasp-demo.$$"

if lsof -ti "tcp:$PORT" >/dev/null 2>&1; then
  echo "⚠️  พอร์ต $PORT มีโปรเซสค้างอยู่ — กำลังเคลียร์ให้..."
  lsof -ti "tcp:$PORT" | xargs kill -9 2>/dev/null || true
  sleep 1
fi

echo "🔨 กำลัง build (ครั้งแรกอาจโหลด dependency สักครู่)..."
go build -o "$BIN" . || { echo "❌ build ไม่ผ่าน"; exit 1; }
trap 'rm -f "$BIN"' EXIT

wait_up() {
  for _ in $(seq 1 40); do
    curl -s "localhost:$PORT/" >/dev/null 2>&1 && return 0
    sleep 0.25
  done
  return 0
}

run() {
  SECURE="$1" "$BIN" >/dev/null 2>&1 &
  local srv=$!
  wait_up
  bash ./exploit.sh
  kill "$srv" 2>/dev/null; wait "$srv" 2>/dev/null || true
}

echo
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " ❌ โหมดมีช่องโหว่ (ยังไม่แก้)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
run ""

echo
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo " ✅ โหมดแก้แล้ว (SECURE=1)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
run "1"

echo
echo "📌 เทียบด้านบน: โหมด ❌ เห็นช่องโหว่จริง · โหมด ✅ โดนบล็อก = ปลอดภัยแล้ว"
