#!/usr/bin/env bash
# 🔍 SCA scan — govulncheck ตรวจ dependency ที่มี CVE
if ! command -v govulncheck >/dev/null 2>&1; then
  echo "ℹ️  ยังไม่มี govulncheck — ติดตั้งก่อน:"
  echo "    go install golang.org/x/vuln/cmd/govulncheck@latest"
  echo
  echo "(ทางลัด: ดูว่ามีเวอร์ชันใหม่กว่าไหมด้วย)"
  go list -m -u github.com/gofiber/fiber/v2 2>/dev/null
  exit 0
fi
echo "== govulncheck ./... =="
govulncheck ./... 2>&1 | grep -iE "vulnerability|GO-[0-9]|fiber|fixed in|more info" | head -20
echo
echo "🔧 วิธีแก้: go get github.com/gofiber/fiber/v2@v2.52.5 && go mod tidy → สแกนซ้ำ = สะอาด"
