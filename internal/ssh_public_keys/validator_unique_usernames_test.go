package ssh_public_keys

import (
	"context"
	"testing"

	"github.com/elestio/terraform-provider-elestio/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUniqueUsernames(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.Set
		expectError bool
	}

	sshPublicKeyAttrTypes := models.SSHPublicKeyAttrTypes

	tests := map[string]testCase{
		"unknown": {
			val: types.SetUnknown(types.ObjectType{
				AttrTypes: sshPublicKeyAttrTypes,
			}),
			expectError: false,
		},
		"null": {
			val: types.SetNull(types.ObjectType{
				AttrTypes: sshPublicKeyAttrTypes,
			}),
			expectError: false,
		},
		"valid-empty": {
			val: types.SetValueMust(types.ObjectType{
				AttrTypes: sshPublicKeyAttrTypes,
			}, []attr.Value{}),
			expectError: false,
		},
		"valid-single": {
			val: types.SetValueMust(types.ObjectType{
				AttrTypes: sshPublicKeyAttrTypes,
			}, []attr.Value{
				types.ObjectValueMust(sshPublicKeyAttrTypes, map[string]attr.Value{
					"username": types.StringValue("user1"),
					"key_data": types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."),
				}),
			}),
			expectError: false,
		},
		"valid-multiple-unique": {
			val: types.SetValueMust(types.ObjectType{
				AttrTypes: sshPublicKeyAttrTypes,
			}, []attr.Value{
				types.ObjectValueMust(sshPublicKeyAttrTypes, map[string]attr.Value{
					"username": types.StringValue("user1"),
					"key_data": types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."),
				}),
				types.ObjectValueMust(sshPublicKeyAttrTypes, map[string]attr.Value{
					"username": types.StringValue("user2"),
					"key_data": types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."),
				}),
			}),
			expectError: false,
		},
		"invalid-duplicate-username": {
			val: types.SetValueMust(types.ObjectType{
				AttrTypes: sshPublicKeyAttrTypes,
			}, []attr.Value{
				types.ObjectValueMust(sshPublicKeyAttrTypes, map[string]attr.Value{
					"username": types.StringValue("user1"),
					"key_data": types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."),
				}),
				types.ObjectValueMust(sshPublicKeyAttrTypes, map[string]attr.Value{
					"username": types.StringValue("user1"),
					"key_data": types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."),
				}),
			}),
			expectError: true,
		},
		"invalid-multiple-duplicates": {
			val: types.SetValueMust(types.ObjectType{
				AttrTypes: sshPublicKeyAttrTypes,
			}, []attr.Value{
				types.ObjectValueMust(sshPublicKeyAttrTypes, map[string]attr.Value{
					"username": types.StringValue("user1"),
					"key_data": types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."),
				}),
				types.ObjectValueMust(sshPublicKeyAttrTypes, map[string]attr.Value{
					"username": types.StringValue("user1"),
					"key_data": types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."),
				}),
				types.ObjectValueMust(sshPublicKeyAttrTypes, map[string]attr.Value{
					"username": types.StringValue("user2"),
					"key_data": types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."),
				}),
				types.ObjectValueMust(sshPublicKeyAttrTypes, map[string]attr.Value{
					"username": types.StringValue("user2"),
					"key_data": types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."),
				}),
			}),
			expectError: true,
		},
	}

	for name, test := range tests {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.SetRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}
			response := validator.SetResponse{}
			UniqueUsernames().ValidateSet(context.Background(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}
}

