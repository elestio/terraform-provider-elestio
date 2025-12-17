package ssh_public_keys

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"golang.org/x/crypto/ssh"
)

type isValidKeyValidator struct{}

func (v isValidKeyValidator) Description(ctx context.Context) string {
	return "string should be a valid SSH key"
}

func (v isValidKeyValidator) MarkdownDescription(ctx context.Context) string {
	return "string should be a valid SSH key"
}

func (v isValidKeyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	attributePath := req.Path
	attributeValue := req.ConfigValue.ValueString()

	if strings.TrimSpace(attributeValue) == "" {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"Expected a non-empty or whitespace string.",
		)
		return
	}

	docURL := "https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/ssh_keys"

	if strings.Contains(attributeValue, "\n") || strings.Contains(attributeValue, "\r") {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"Your SSH public key must be on a single line. "+
				"Use `provider::elestio::parse_ssh_key_data(file(\"~/.ssh/id_rsa.pub\"))` to read and parse your key file. "+
				"Read the guide: "+docURL,
		)
		return
	}

	fields := strings.Fields(attributeValue)
	if len(fields) < 2 {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"Invalid SSH key format: expected at least key type and key data separated by a space. "+
				"Use `provider::elestio::parse_ssh_key_data(file(\"~/.ssh/id_rsa.pub\"))` to read and parse your key file. "+
				"Read the guide: "+docURL,
		)
		return
	}

	bytes, err := base64.StdEncoding.DecodeString(fields[1])
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"Error decoding the key data. Your SSH public key does not seem to be base64 encoded. "+
				"Read the guide: "+docURL,
		)
		return
	}

	pubKey, err := ssh.ParsePublicKey(bytes)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"Error parsing the key data. Your SSH public key does not seem to be valid. It may be corrupted or altered. "+
				"Read the guide: "+docURL,
		)
		return
	}

	// Supported key types: ssh-rsa, ssh-ed25519, ecdsa-sha2-nistp256, ecdsa-sha2-nistp384, ecdsa-sha2-nistp521, ssh-dss
	supportedTypes := map[string]bool{
		ssh.KeyAlgoRSA:      true,
		ssh.KeyAlgoED25519:  true,
		ssh.KeyAlgoECDSA256: true,
		ssh.KeyAlgoECDSA384: true,
		ssh.KeyAlgoECDSA521: true,
		ssh.KeyAlgoDSA:      true,
	}

	if !supportedTypes[pubKey.Type()] {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"Unsupported SSH key type. Supported types: ssh-rsa, ssh-ed25519, ecdsa-sha2-nistp256, ecdsa-sha2-nistp384, ecdsa-sha2-nistp521, ssh-dss. "+
				"Read the guide: "+docURL,
		)
		return
	}

	if len(fields) > 2 {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"SSH public key comments are not allowed. "+
				"Use `provider::elestio::parse_ssh_key_data()` to remove the comment, "+
				"or `provider::elestio::parse_ssh_key()` to extract both the key and username from the comment. "+
				"Read the guide: "+docURL,
		)
		return
	}
}

// IsValidKey returns a validator that ensures the SSH public key is valid
func IsValidKey() isValidKeyValidator {
	return isValidKeyValidator{}
}
