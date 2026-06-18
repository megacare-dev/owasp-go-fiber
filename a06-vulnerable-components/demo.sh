#!/usr/bin/env bash
# A06 demo — ตรวจ dependency ที่มีช่องโหว่ (เป็น scan ไม่ใช่ runtime toggle)
set -u
cd "$(dirname "$0")"

echo "📦 dependency ปัจจุบัน (go.mod):"
grep "gofiber/fiber" go.mod
echo
echo "🔍 สแกนหา CVE:"
bash ./scan.sh
echo
echo "📌 A06 = อย่าใช้ component ที่มีช่องโหว่ — สแกน (govulncheck) + อัปเดตให้เป็นเวอร์ชันที่ patch แล้ว"
echo "   ดูช่องโหว่ตัวจริงของ fiber CVE-2024-25124 ที่รันได้ ใน ../cors-cve/"
