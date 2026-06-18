#!/usr/bin/env bash
# 🔍 SAST scan — gosec ตรวจ SSRF (variable URL ใน http request)
echo "== gosec ./... =="
gosec -quiet -fmt=text ./... 2>/dev/null | grep -E "G107" | head
echo
if gosec -quiet ./... 2>/dev/null | grep -qE "G107"; then
  echo "🔍 พบ! G107 HTTP request ด้วย variable url (CWE-88) → 🔧 บล็อก internal IP (SECURE=1)"
else echo "✅ ไม่พบ SSRF pattern"; fi
