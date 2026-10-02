package utils

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestRecoverToDiagnostic(t *testing.T) {
	var diags diag.Diagnostics
	func() {
		defer RecoverToDiagnostic(&diags, "Create")
		var s []int
		_ = s[0]
	}()
	if !diags.HasError() {
		t.Fatal("panic must become an error diagnostic")
	}

	var ok diag.Diagnostics
	func() { defer RecoverToDiagnostic(&ok, "Create") }()
	if ok.HasError() {
		t.Fatal("no panic, no diagnostic")
	}
}
