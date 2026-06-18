#!/usr/bin/env bash
# 🔍 scan — แสดงว่า SAST มองไม่เห็น authz/IDOR (logic flaw)
echo "== gosec ./... (SAST) =="
out=$(gosec -quiet -fmt=text ./... 2>/dev/null | grep -E "G[0-9]+ \(CWE" | grep -v "G104")
[ -z "$out" ] && echo "(ไม่พบ security finding ที่เกี่ยวข้อง)" || echo "$out"
echo
echo "⚠️ SAST blind spot: IDOR/Broken Access เป็น logic flaw — gosec มองไม่เห็น"
echo "🔍 ต้องใช้ DAST → bash ./exploit.sh  (alice อ่าน order ของ bob ได้)"
echo "🔧 fix: ตรวจ owner ก่อนคืนข้อมูล → SECURE=1 go run ."
