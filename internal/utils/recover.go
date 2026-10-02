package utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// RecoverToDiagnostic converts a panic into an error diagnostic. Use it with
// defer at the top of a CRUD method so that a malformed API response (the API
// client indexes response arrays without length checks) fails one operation
// instead of crashing the provider process and every other resource with it.
func RecoverToDiagnostic(diags *diag.Diagnostics, operation string) {
	if r := recover(); r != nil {
		diags.AddError(
			"Unexpected Provider Error",
			fmt.Sprintf("%s failed because of an unexpected response from the Elestio API (%v). "+
				"Check the resource on the Elestio dashboard, then retry. If it persists, report it to the provider developers.", operation, r),
		)
	}
}
