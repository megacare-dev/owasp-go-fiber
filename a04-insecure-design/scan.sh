#!/usr/bin/env bash
echo "== gosec ./... =="
gosec -quiet ./... 2>/dev/null | tail -4
echo
echo "ℹ️  gosec (SAST) มองไม่เห็นช่องโหว่นี้ — เป็น business-logic / design flaw"
echo "   ยืนยันด้วย ./exploit.sh (DAST) แทน"
