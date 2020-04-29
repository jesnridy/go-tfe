package tfe

import (
	"context"
	"fmt"
)

// Compile-time proof of interface implementation.
var _ Registry = (*registry)(nil)

// Registry describes all the registry module related methods that the Terraform
// Enterprise API supports.
//
// TFE API docs: https://www.terraform.io/docs/cloud/api/modules.html
type Registry interface {
	// Publish a module to the TFE private registry
	Publish(ctx context.Context, options ModulePublishOptions) (*Module, error)

	// Create a module on the TFE Private registry
	CreateModule(ctx context.Context, organization string, options ModuleCreateOptions) (*Module, error)

	// Create a module version on TFE Private registry
	CreateModuleVersion(ctx context.Context, organizationName, moduleName, providerName string, options ModuleCreateVersionOptions) (*ModuleVersion, error)

	// Delete a module on the TFE private registry
	DeleteModule(ctx context.Context, organizationName, moduleName string) error

	// Delete a specific module version on the TFE private registry
	DeleteModuleVersion(ctx context.Context, organizationName, moduleName, provider, version string) error

	// Delete a specific module provider on the TFE private registry
	DeleteModuleProvider(ctx context.Context, organizationName, moduleName, provider string) error
}

// registry implements Registry.
type registry struct {
	client *Client
}

// ModuleCreateOptions is used when creating a module
type ModuleCreateOptions struct {
	Name     string `jsonapi:"attr,name"`
	Provider string `jsonapi:"attr,provider"`
}

// ModuleCreateVersionOptions is used when creating a module version
type ModuleCreateVersionOptions struct {
	Version string `jsonapi:"attr,version"`
}

// ModulePublishOptions options for publishing a registry module
type ModulePublishOptions struct {
	// VCS connection information to import a module to the registry
	ModuleVCSOptions *ModuleVCSOptions `jsonapi:"attr,vcs-repo"`
}

// ModuleVCSOptions contains the configuration of a VCS integration.
type ModuleVCSOptions struct {
	Identifier        string `json:"identifier"`
	OAuthTokenID      string `json:"oauth-token-id"`
	DisplayIdentifier string `json:"display_identifier"`
}

// Module represents a registry module
type Module struct {
	ID   string `jsonapi:"primary,registry-modules"`
	Type string `json:"type"`

	Name      string `jsonapi:"attr,name"`
	Provider  string `jsonapi:"attr,provider"`
	Status    string `jsonapi:"attr,status"`
	CreatedAt string `jsonapi:"attr,created-at"`
	UpdatedAt string `jsonapi:"attr,updated-at"`

	// TODO:
	// version-statuses
	// permissions

	// Relations
	Organization *Organization `jsonapi:"relation,organization"`

	// Links
	// TODO
}

// ModuleVersion represents a registry module version
type ModuleVersion struct {
	ID   string `jsonapi:"primary,registry-module-versions"`
	Type string `json:"type"`

	Source    string `jsonapi:"attr,Source"`
	Status    string `jsonapi:"attr,status"`
	Version   string `jsonapi:"attr,version"`
	CreatedAt string `jsonapi:"attr,created-at"`
	UpdatedAt string `jsonapi:"attr,updated-at"`

	// TODO:
	// version-statuses
	// permissions

	// Relations

	// Links
	// TODO
}

// Publish is used to publish a new module to the TFE private registry
func (r *registry) Publish(ctx context.Context, options ModulePublishOptions) (*Module, error) {
	req, err := r.client.newRequest("POST", "registry-modules", &options)
	if err != nil {
		return nil, err
	}

	m := &Module{}
	err = r.client.do(ctx, req, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *registry) CreateModule(ctx context.Context, organizationName string, options ModuleCreateOptions) (*Module, error) {
	path := fmt.Sprintf("organizations/%s/registry-modules", organizationName)
	req, err := r.client.newRequest("POST", path, &options)
	if err != nil {
		return nil, err
	}

	m := &Module{}
	err = r.client.do(ctx, req, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *registry) CreateModuleVersion(ctx context.Context, organizationName, moduleName, providerName string, options ModuleCreateVersionOptions) (*ModuleVersion, error) {
	path := fmt.Sprintf("registry-modules/%s/%s/%s/versions", organizationName, moduleName, providerName)
	req, err := r.client.newRequest("POST", path, &options)
	if err != nil {
		return nil, err
	}

	m := &ModuleVersion{}
	err = r.client.do(ctx, req, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// DeleteModule is used to delete the entire module on the TFE private registry
func (r *registry) DeleteModule(ctx context.Context, organizationName, moduleName string) error {
	path := fmt.Sprintf("registry-modules/actions/delete/%s/%s", organizationName, moduleName)
	req, err := r.client.newRequest("POST", path, nil)
	if err != nil {
		return err
	}

	err = r.client.do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// DeleteModuleVersion is used to delete the specific module version on the TFE private registry
func (r *registry) DeleteModuleVersion(ctx context.Context, organizationName, moduleName, provider, version string) error {
	path := fmt.Sprintf("registry-modules/actions/delete/%s/%s/%s/%s", organizationName, moduleName, provider, version)
	req, err := r.client.newRequest("POST", path, nil)
	if err != nil {
		return err
	}

	err = r.client.do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

// DeleteModuleProvider is used to delete the specific module provider on the TFE private registry
func (r *registry) DeleteModuleProvider(ctx context.Context, organizationName, moduleName, provider string) error {
	path := fmt.Sprintf("registry-modules/actions/delete/%s/%s/%s", organizationName, moduleName, provider)
	req, err := r.client.newRequest("POST", path, nil)
	if err != nil {
		return err
	}

	err = r.client.do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}