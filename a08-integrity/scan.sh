#!/usr/bin/env bash
echo "== gosec ./... =="
gosec -quiet ./... 2>/dev/null | tail -4
echo
echo "ℹ️  gosec (SAST) มองไม่เห็น — เป็น logic flaw (ไม่ verify integrity)"
echo "   ยืนยันด้วย ./exploit.sh (DAST) แทน"
