#!/usr/bin/env bash
# 🔍 SAST scan — gosec ตรวจ SQL injection
echo "== gosec ./... =="
gosec -quiet -fmt=text ./... 2>/dev/null | grep -E "G201|G202" | head
echo
if gosec -quiet ./... 2>/dev/null | grep -qE "G201|G202"; then
  echo "🔍 พบ! G201 SQL string formatting (CWE-89) → 🔧 ใช้ parameterized query (SECURE=1)"
else echo "✅ ไม่พบ SQLi pattern"; fi
