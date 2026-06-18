#!/usr/bin/env bash
# 🔍 scan — SAST เจอแค่ผิวเผิน ไม่เจอ auth-bypass จริง
echo "== gosec ./... (SAST) =="
gosec -quiet -fmt=text ./... 2>/dev/null | grep -E "G101" | head
echo
echo "⚠️ gosec เจอแค่ G101 (hardcoded secret) — แต่ bug จริงคือ 'ไม่ verify signature'"
echo "   logic นี้ SAST มองไม่เห็น"
echo "🔍 ต้องใช้ DAST → bash ./exploit.sh  (ปลอม token alg:none เข้า /admin)"
echo "🔧 fix: บังคับ HS256 + verify HMAC → SECURE=1 go run ."
