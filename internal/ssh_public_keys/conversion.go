package ssh_public_keys

import (
	"context"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/elestio/terraform-provider-elestio/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ConvertElestioToTerraform converts SSH keys from Elestio API format to Terraform format
func ConvertElestioToTerraform(ctx context.Context, elestioSSHKeys []elestio.ServiceSSHPublicKey, diags *diag.Diagnostics) types.Set {
	sshPublicKeys := make([]models.SSHPublicKeyModel, len(elestioSSHKeys))
	for i, s := range elestioSSHKeys {
		sshPublicKeys[i] = models.SSHPublicKeyModel{
			Username: types.StringValue(s.Name),
			KeyData:  types.StringValue(s.Key),
		}
	}

	setSSHPublicKeys, d := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: models.SSHPublicKeyAttrTypes}, sshPublicKeys)
	diags.Append(d...)

	return setSSHPublicKeys
}
