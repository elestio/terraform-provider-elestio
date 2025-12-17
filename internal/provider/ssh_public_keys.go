package provider

import (
	"context"
	"fmt"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/elestio/terraform-provider-elestio/internal/models"
	"github.com/elestio/terraform-provider-elestio/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func SSHPublicKeysDefaultValue() types.Set {
	return types.SetValueMust(types.ObjectType{AttrTypes: models.SSHPublicKeyAttrTypes}, []attr.Value{})
}

var sshPublicKeysSchema = schema.SetNestedAttribute{
	MarkdownDescription: "You can add Public Keys to your resource to access it via the SSH protocol.",
	Optional:            true,
	Computed:            true,
	Default:             setdefault.StaticValue(SSHPublicKeysDefaultValue()),
	Validators: []validator.Set{
		validators.SSHPublicKeysUniqueUsernames(),
	},
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "The username is used to identify the Public Key among others. Must be unique (per resource).",
				Required:            true,
			},
			"key_data": schema.StringAttribute{
				MarkdownDescription: "The Public Key value without comment." +
					" Use `provider::elestio::parse_ssh_key_data(file(\"~/.ssh/id_rsa.pub\"))` to remove the comment from your key." +
					" Read the guide [\"How to use SSH keys with Elestio Terraform Provider\"](https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/ssh_keys).",
				Required: true,
				Validators: []validator.String{
					validators.IsSSHPublicKey(),
				},
			},
		},
	},
}

func compareSSHPublicKeys(ctx *context.Context, state *basetypes.SetValue, plan *basetypes.SetValue, diags *diag.Diagnostics) (toAdd []models.SSHPublicKeyModel, toUpdate []models.SSHPublicKeyModel, toRemove []models.SSHPublicKeyModel) {
	var stateKeys []models.SSHPublicKeyModel
	diags.Append(state.ElementsAs(*ctx, &stateKeys, true)...)
	if diags.HasError() {
		return nil, nil, nil
	}
	var planKeys []models.SSHPublicKeyModel
	diags.Append(plan.ElementsAs(*ctx, &planKeys, true)...)
	if diags.HasError() {
		return nil, nil, nil
	}

	toAdd, toUpdate, toRemove = []models.SSHPublicKeyModel{}, []models.SSHPublicKeyModel{}, []models.SSHPublicKeyModel{}

	if len(stateKeys) == 0 && len(planKeys) == 0 {
		return toAdd, toUpdate, toRemove
	}

	// Usernames are unique and can be used as map index
	stateKeysMap := make(map[string]models.SSHPublicKeyModel)
	for _, obj := range stateKeys {
		stateKeysMap[obj.Username.ValueString()] = obj
	}
	planKeysMap := make(map[string]models.SSHPublicKeyModel)
	for _, obj := range planKeys {
		planKeysMap[obj.Username.ValueString()] = obj
	}

	// Iterate over state and delete any objects that are not in the plan
	for _, stateKey := range stateKeys {
		if _, exists := planKeysMap[stateKey.Username.ValueString()]; !exists {
			toRemove = append(toRemove, stateKey)
		}
	}

	// Iterate over the plan and compare each key to the corresponding key in state
	for _, planKey := range planKeys {
		if stateKey, exists := stateKeysMap[planKey.Username.ValueString()]; exists {
			if !planKey.KeyData.Equal(stateKey.KeyData) {
				toUpdate = append(toUpdate, planKey)
			}
		} else {
			toAdd = append(toAdd, planKey)
		}
	}

	return toAdd, toUpdate, toRemove
}

func applySSHPublicKeyChanges(ctx context.Context, serviceID string, keysToAdd, keysToUpdate, keysToRemove []models.SSHPublicKeyModel, providerName string, client *elestio.Client, serviceResource *ServiceResource, service *elestio.Service) error {
	// Remove keys first
	for _, key := range keysToRemove {
		if err := client.Service.RemoveSSHPublicKey(serviceID, key.Username.ValueString()); err != nil {
			return fmt.Errorf("failed to remove ssh public key: %s", err)
		}
	}

	// Update keys (remove old, add new)
	for _, key := range keysToUpdate {
		if err := client.Service.RemoveSSHPublicKey(serviceID, key.Username.ValueString()); err != nil {
			return fmt.Errorf("failed to update (remove the old one) ssh public key: %s", err)
		}
		if err := client.Service.AddSSHPublicKey(serviceID, key.Username.ValueString(), key.KeyData.ValueString()); err != nil {
			return fmt.Errorf("failed to update (add the new one) ssh public key: %s", err)
		}
	}

	// Add new keys
	for _, key := range keysToAdd {
		if err := client.Service.AddSSHPublicKey(serviceID, key.Username.ValueString(), key.KeyData.ValueString()); err != nil {
			return fmt.Errorf("failed to add ssh public key: %s", err)
		}
	}

	keyWasUpdated := len(keysToAdd) > 0 || len(keysToUpdate) > 0 || len(keysToRemove) > 0

	// Scaleway does not support updating ssh keys without rebooting the server.
	if keyWasUpdated && providerName == "scaleway" {
		if err := client.Service.RebootServer(serviceID); err != nil {
			return fmt.Errorf("failed to reboot server to update scaleway ssh keys: %s", err)
		}
		if _, err := serviceResource.waitServerReboot(ctx, service); err != nil {
			return fmt.Errorf("failed to wait server reboot to update scaleway ssh keys: %s", err)
		}
	}

	return nil
}

// convertElestioSSHKeysToTerraform converts SSH keys from Elestio API format to Terraform format
func convertElestioSSHKeysToTerraform(ctx context.Context, elestioSSHKeys []elestio.ServiceSSHPublicKey, diags *diag.Diagnostics) types.Set {
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
