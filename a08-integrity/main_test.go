package main

import "testing"

// payload คือ webhook ที่ exploit.sh ยิงจริง (event จ่ายเงิน)
const payload = `{"event":"payment","status":"paid","amount":0}`

// โหมด SECURE ต้องปฏิเสธ webhook ที่ไม่มีลายเซ็น (เคสที่ exploit.sh ยิง)
func TestSecureRejectsMissingSignature(t *testing.T) {
	if verifyWebhook([]byte(payload), "", true) {
		t.Fatal("webhook ที่ไม่มีลายเซ็นต้องถูกปฏิเสธในโหมด secure")
	}
}

// โหมด SECURE ต้องปฏิเสธลายเซ็นที่ผิด
func TestSecureRejectsWrongSignature(t *testing.T) {
	if verifyWebhook([]byte(payload), "deadbeef", true) {
		t.Fatal("ลายเซ็นที่ผิดต้องถูกปฏิเสธ")
	}
}

// โหมด SECURE ต้องปฏิเสธ body ที่ถูกแก้ไข แม้แนบลายเซ็นเดิมของ payload จริงมา
// (กันการดักแก้ยอดเงินแล้ว replay ลายเซ็นเก่า)
func TestSecureRejectsTamperedBody(t *testing.T) {
	sig := sign([]byte(payload))
	tampered := `{"event":"payment","status":"paid","amount":999999}`
	if verifyWebhook([]byte(tampered), sig, true) {
		t.Fatal("body ที่ถูกแก้ไขต้องไม่ผ่านลายเซ็นเดิม")
	}
}

// webhook ที่ลายเซ็นถูกต้อง (มาจากแหล่งที่รู้ secret) ต้องผ่าน
func TestSecureAcceptsValidSignature(t *testing.T) {
	sig := sign([]byte(payload))
	if !verifyWebhook([]byte(payload), sig, true) {
		t.Fatal("webhook ที่ลายเซ็นถูกต้องต้องผ่าน")
	}
}

// ยืนยันว่าโหมด VULNERABLE ยังเปิดช่องโหว่จริง (payload ปลอมไม่มีลายเซ็นก็เชื่อ)
func TestVulnerableAcceptsForgedWebhook(t *testing.T) {
	if !verifyWebhook([]byte(payload), "", false) {
		t.Fatal("โหมด vulnerable ควรเชื่อ webhook ปลอม (คือตัวช่องโหว่)")
	}
}
