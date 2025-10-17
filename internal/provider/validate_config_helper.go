package provider

import (
	"context"
	"fmt"

	"github.com/elestio/elestio-go-api-client/v2"
)

func validateProviderConfig(
	ctx context.Context,
	client *elestio.Client,
	templateId int64,
	providerName, datacenter, serverType string,
) error {
	isConfigValid, err := client.Service.ValidateConfig(elestio.ValidateConfigRequest{
		TemplateId:   templateId,
		ProviderName: providerName,
		Datacenter:   datacenter,
		ServerType:   serverType,
	})

	if err != nil || !isConfigValid {
		errorMsg := ""
		if err != nil {
			errorMsg = fmt.Sprintf("%s\n\n", err)
		}

		errorMsg += fmt.Sprintf("Configuration provided:\n"+
			"  provider_name: %s\n"+
			"  datacenter: %s\n"+
			"  server_type: %s\n\n"+
			"For a complete list of valid combinations, see: https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/providers_datacenters_server_types",
			providerName, datacenter, serverType)

		return fmt.Errorf("%s", errorMsg)
	}

	return nil
}
