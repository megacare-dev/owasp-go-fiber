#!/usr/bin/env bash
# 🔍 SAST scan — gosec ตรวจ weak crypto
echo "== gosec ./... =="
gosec -quiet -fmt=text ./... 2>/dev/null | grep -E "G401|G501" | head
echo
if gosec -quiet ./... 2>/dev/null | grep -qE "G401|G501"; then
  echo "🔍 พบ! G401/G501 weak crypto (crypto/md5) → 🔧 เปลี่ยนไปใช้ crypto/rand (SECURE=1)"
else echo "✅ ไม่พบ weak crypto"; fi
