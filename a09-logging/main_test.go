package main

import "testing"

// โหมด SECURE ต้องบันทึกเหตุการณ์ login ล้มเหลว (ทีม security จึงตรวจจับได้)
func TestSecureRecordsFailedLogin(t *testing.T) {
	auditLog = nil
	if !recordFailedLogin("admin", "1.2.3.4", true) {
		t.Fatal("โหมด secure ต้องบันทึก login ล้มเหลว")
	}
	if len(auditLog) != 1 {
		t.Fatalf("ต้องมี 1 เหตุการณ์ใน audit log, ได้ %d", len(auditLog))
	}
}

// โหมด VULNERABLE ต้องไม่บันทึกอะไรเลย (คือตัวช่องโหว่ — มองไม่เห็นการโจมตี)
func TestVulnerableRecordsNothing(t *testing.T) {
	auditLog = nil
	recordFailedLogin("admin", "1.2.3.4", false)
	if len(auditLog) != 0 {
		t.Fatalf("โหมด vulnerable ต้องไม่บันทึกอะไร, ได้ %d เหตุการณ์", len(auditLog))
	}
}

// brute-force หลายครั้ง → ต้องเห็นร่องรอยครบทุกครั้ง
func TestSecureRecordsEveryAttempt(t *testing.T) {
	auditLog = nil
	for i := 0; i < 5; i++ {
		recordFailedLogin("admin", "1.2.3.4", true)
	}
	if len(auditLog) != 5 {
		t.Fatalf("ต้องบันทึกครบทุกครั้ง (5), ได้ %d", len(auditLog))
	}
}

// audit log ต้องเป็น structured format ที่มีแค่ user+ip — รหัสผ่านไม่ถูกส่งเข้าฟังก์ชันด้วยซ้ำ
func TestAuditLogFormatHasNoSecret(t *testing.T) {
	auditLog = nil
	recordFailedLogin("admin", "1.2.3.4", true)
	want := "FAILED_LOGIN user=admin ip=1.2.3.4"
	if auditLog[0] != want {
		t.Fatalf("รูปแบบ log ผิด: got %q want %q", auditLog[0], want)
	}
}
