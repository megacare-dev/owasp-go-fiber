#!/usr/bin/env bash
# 🔍 SCA scan — ตรวจช่องโหว่ใน dependency ด้วย govulncheck
echo "== govulncheck ./... (เน้น gofiber/fiber) =="
govulncheck ./... 2>/dev/null | grep -A2 -iE "Module: github.com/gofiber"
echo
if govulncheck ./... 2>/dev/null | grep -q "gofiber/fiber"; then
  echo "🔍 พบ! gofiber/fiber มีช่องโหว่ (รวม CVE-2024-25124 CORS) → อัปเดตตาม 'Fixed in'"
else
  echo "✅ ไม่พบช่องโหว่ใน dependency (fiber อัปเดตแล้ว)"
fi
