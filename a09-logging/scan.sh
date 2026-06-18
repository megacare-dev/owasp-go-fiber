#!/usr/bin/env bash
echo "== gosec ./... =="
gosec -quiet ./... 2>/dev/null | tail -4
echo
echo "ℹ️  gosec (SAST) มองไม่เห็น — 'การไม่มี log' ตรวจด้วย static scan ไม่ได้"
echo "   ยืนยันด้วย ./exploit.sh (ยิงโจมตีแล้วดูว่ามีร่องรอยใน /audit ไหม)"
