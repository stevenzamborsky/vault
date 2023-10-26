// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	paths "path"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathPreauthTest() *framework.Path {
	return &framework.Path{
		Pattern: "preauth-test",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
		},
		Fields: map[string]*framework.FieldSchema{
			"accessor": {
				Type: framework.TypeString,
			},
			"path": {
				Type: framework.TypeString,
			},
			"username": {
				Type: framework.TypeString,
			},
			"password": {
				Type: framework.TypeString,
			},
			"loop": {
				Type:    framework.TypeBool,
				Default: false,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.handlePreauthTest,
				Summary:  "Returns the size of the active cache",
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "cache-configuration",
				},
			},
		},
	}
}

func (b *backend) handlePreauthTest(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if d.Get("loop").(bool) {
		da := logical.NewDelegatedAuthenticationRequest(d.Get("accessor").(string), paths.Join(d.Get("path").(string), d.Get("username").(string)),
			map[string]interface{}{
				"password": d.Get("password").(string),
			},
			nil)
		return nil, da
	}
	if req.ClientTokenSource != logical.ClientTokenFromInternalAuth {
		da := logical.NewDelegatedAuthenticationRequest(d.Get("accessor").(string), paths.Join(d.Get("path").(string), d.Get("username").(string)),
			map[string]interface{}{
				"password": d.Get("password").(string),
			},
			nil)
		return nil, da
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"success": true,
		},
	}, nil
}