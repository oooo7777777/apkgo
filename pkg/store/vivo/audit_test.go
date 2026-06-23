package vivo

import (
	"testing"

	"github.com/KevinGong2013/apkgo/v3/pkg/store"
)

func TestMapVivoAuditState_UsesUnPassReasonForRejected(t *testing.T) {
	got, detail := mapVivoAuditState(4, "应用截图不符合规范")
	if got != store.AuditRejected {
		t.Fatalf("state = %q, want %q", got, store.AuditRejected)
	}
	if detail != "应用截图不符合规范" {
		t.Fatalf("detail = %q, want unPassReason", detail)
	}
}
