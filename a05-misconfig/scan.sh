#!/usr/bin/env bash
# 🔍 scan — info-leak via error: SAST มองไม่เห็น
echo "== gosec ./... (SAST) =="
out=$(gosec -quiet -fmt=text ./... 2>/dev/null | grep -E "G[0-9]+ \(CWE" | grep -v "G104")
[ -z "$out" ] && echo "(ไม่พบ security finding ที่เกี่ยวข้อง)" || echo "$out"
echo
echo "⚠️ SAST blind spot: verbose-error leak เป็น runtime behavior — gosec มองไม่เห็น"
echo "🔍 ต้องใช้ DAST → bash ./exploit.sh  (error เผย DB connection string)"
echo "🔧 fix: custom ErrorHandler → SECURE=1 go run ."
