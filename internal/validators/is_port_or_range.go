package validators

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type isPortOrRangeValidator struct{}

func (v isPortOrRangeValidator) Description(ctx context.Context) string {
	return "string should be a valid port number (1-65535) or port range (e.g., 8000-9000)"
}

func (v isPortOrRangeValidator) MarkdownDescription(ctx context.Context) string {
	return "string should be a valid port number (1-65535) or port range (e.g., 8000-9000)"
}

func (v isPortOrRangeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	if strings.Contains(value, "-") {
		validatePortRange(value, req.Path, resp)
	} else {
		validateSinglePort(value, req.Path, resp)
	}
}

func validatePortRange(value string, path path.Path, resp *validator.StringResponse) {
	rangeRegex := regexp.MustCompile(`^(\d+)-(\d+)$`)
	matches := rangeRegex.FindStringSubmatch(value)

	if matches == nil {
		resp.Diagnostics.AddAttributeError(
			path,
			"Invalid Port Range",
			fmt.Sprintf("Port range must be in format 'start-end' (e.g., '8000-9000'), got: %s", value),
		)
		return
	}

	startPort, err1 := strconv.Atoi(matches[1])
	endPort, err2 := strconv.Atoi(matches[2])

	if err1 != nil || err2 != nil {
		resp.Diagnostics.AddAttributeError(
			path,
			"Invalid Port Range",
			fmt.Sprintf("Port range must contain valid numbers, got: %s", value),
		)
		return
	}

	if !isValidPort(startPort) {
		resp.Diagnostics.AddAttributeError(
			path,
			"Invalid Port Range",
			fmt.Sprintf("Start port must be between 1 and 65535, got: %d", startPort),
		)
		return
	}

	if !isValidPort(endPort) {
		resp.Diagnostics.AddAttributeError(
			path,
			"Invalid Port Range",
			fmt.Sprintf("End port must be between 1 and 65535, got: %d", endPort),
		)
		return
	}

	if startPort >= endPort {
		resp.Diagnostics.AddAttributeError(
			path,
			"Invalid Port Range",
			fmt.Sprintf("Start port must be less than end port, got: %s", value),
		)
	}
}

func validateSinglePort(value string, path path.Path, resp *validator.StringResponse) {
	port, err := strconv.Atoi(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path,
			"Invalid Port Number",
			fmt.Sprintf("Port must be a valid number, got: %s", value),
		)
		return
	}

	if !isValidPort(port) {
		resp.Diagnostics.AddAttributeError(
			path,
			"Invalid Port Number",
			fmt.Sprintf("Port must be between 1 and 65535, got: %d", port),
		)
	}
}

func isValidPort(port int) bool {
	return port >= 1 && port <= 65535
}

func IsPortOrRange() isPortOrRangeValidator {
	return isPortOrRangeValidator{}
}
