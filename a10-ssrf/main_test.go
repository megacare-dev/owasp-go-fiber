package main

import (
	"net"
	"testing"
)

// ปลายทางภายในทุกชนิดต้องถูกบล็อก (รวม cloud metadata 169.254.169.254)
func TestBlocksInternalIPs(t *testing.T) {
	for _, s := range []string{
		"127.0.0.1",       // loopback
		"169.254.169.254", // cloud metadata (link-local)
		"10.0.0.5",        // private
		"192.168.1.1",     // private
		"0.0.0.0",         // unspecified
	} {
		if !isBlockedIP(net.ParseIP(s)) {
			t.Errorf("ต้องบล็อก %s", s)
		}
	}
}

// ปลายทางสาธารณะต้องไม่ถูกบล็อก
func TestAllowsPublicIP(t *testing.T) {
	if isBlockedIP(net.ParseIP("8.8.8.8")) {
		t.Fatal("8.8.8.8 (public) ไม่ควรถูกบล็อก")
	}
}

// host ที่ resolve เป็น loopback ต้องถูกบล็อก (ไม่ต้องใช้เน็ต)
func TestBlocksLoopbackHost(t *testing.T) {
	if !isBlockedHost("127.0.0.1") {
		t.Fatal("127.0.0.1 ต้องถูกบล็อก")
	}
	if !isBlockedHost("localhost") {
		t.Fatal("localhost ต้องถูกบล็อก")
	}
}

// host ที่ resolve ไม่ได้ต้อง fail-closed (บล็อก)
func TestBlocksUnresolvableHost(t *testing.T) {
	if !isBlockedHost("no-such-host.invalid") {
		t.Fatal("host ที่ resolve ไม่ได้ต้องถูกบล็อก (fail-closed)")
	}
}
