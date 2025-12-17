package ssh_public_keys

import (
	"context"
	"fmt"

	"github.com/elestio/terraform-provider-elestio/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type uniqueUsernamesValidator struct{}

func (v uniqueUsernamesValidator) Description(ctx context.Context) string {
	return "SSH public key usernames must be unique within the set"
}

func (v uniqueUsernamesValidator) MarkdownDescription(ctx context.Context) string {
	return "SSH public key usernames must be unique within the set"
}

func (v uniqueUsernamesValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	// Convert the set elements to the SSH public key model
	var keys []models.SSHPublicKeyModel
	diags := req.ConfigValue.ElementsAs(ctx, &keys, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Track seen usernames
	usernameMap := make(map[string]bool)
	var duplicates []string

	for _, key := range keys {
		username := key.Username.ValueString()
		if usernameMap[username] {
			// Found a duplicate, add to duplicates slice if not already added
			found := false
			for _, duplicate := range duplicates {
				if duplicate == username {
					found = true
					break
				}
			}
			if !found {
				duplicates = append(duplicates, username)
			}
		} else {
			usernameMap[username] = true
		}
	}

	// Report all duplicates
	if len(duplicates) > 0 {
		var message string
		if len(duplicates) == 1 {
			message = fmt.Sprintf("SSH Public Key username must be unique per resource. The following username is duplicated: %s", duplicates[0])
		} else {
			message = fmt.Sprintf("SSH Public Key usernames must be unique per resource. The following usernames are duplicated: %v", duplicates)
		}

		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Configuration",
			message,
		)
	}
}

// UniqueUsernames returns a validator that ensures SSH public key usernames are unique
func UniqueUsernames() uniqueUsernamesValidator {
	return uniqueUsernamesValidator{}
}

