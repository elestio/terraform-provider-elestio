package validators

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"golang.org/x/crypto/ssh"
)

type isSSHPublicKeyValidator struct{}

func (v isSSHPublicKeyValidator) Description(ctx context.Context) string {
	return "string should be a valid SSH key"
}

func (v isSSHPublicKeyValidator) MarkdownDescription(ctx context.Context) string {
	return "string should be a valid SSH key"
}

func (v isSSHPublicKeyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
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

	if strings.Contains(attributeValue, "\n") || strings.Contains(attributeValue, "\r") {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"Your SSH public key must be on a single line. You can use the `chomp()` terraform function to remove newlines from your key. Example: key_data = chomp(file(\"~/.ssh/id_rsa.pub\"))",
		)
		return
	}

	fields := strings.Fields(attributeValue)
	if len(fields) < 2 {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"Expected a string with at least two fields separated by a space.",
		)
		return
	}

	bytes, err := base64.StdEncoding.DecodeString(fields[1])
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"Error decoding the key data. Your SSH public key does not seem to be base64 encoded.",
		)
		return
	}

	pubKey, err := ssh.ParsePublicKey(bytes)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"Error parsing the key data. Your SSH public key does not seem to be valid. It may be corrupted or altered.",
		)
		return
	}

	if pubKey.Type() != ssh.KeyAlgoRSA {
		resp.Diagnostics.AddAttributeError(
			attributePath,
			"Invalid Attribute Configuration",
			"Only RSA keys are supported by Elestio.",
		)
		return
	}

	// TODO: Enable this check to support only strong keys https://www.ibm.com/docs/en/zos/2.3.0?topic=certificates-size-considerations-public-private-keys
	//
	//	if len(pubKey.Marshal()) < 2048 {
	//		resp.Diagnostics.AddAttributeError(
	//			attributePath,
	//			"Invalid Attribute Configuration",
	//			"Only RSA keys with a length of 2048 bits or more are supported by Elestio.",
	//		)
	//		return
	//	}
}

func IsSSHPublicKey() isSSHPublicKeyValidator {
	return isSSHPublicKeyValidator{}
}
