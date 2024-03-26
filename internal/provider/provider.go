// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/booya/gowrt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure OpenWrtProvider satisfies various provider interfaces.
var _ provider.Provider = &OpenWrtProvider{}
var _ provider.ProviderWithFunctions = &OpenWrtProvider{}

// OpenWrtProvider defines the provider implementation.
type OpenWrtProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// openWrtProviderModel maps provider schema data to a Go type.
type openWrtProviderModel struct {
	Host        types.String `tfsdk:"host"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	InsecureTls types.Bool   `tfsdk:"insecure_tls"`
}

func (p *OpenWrtProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openwrt"
	resp.Version = p.version
}

func (p *OpenWrtProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Host/ip of your OpenWrt router",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for login",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for login",
				Required:            true,
				Sensitive:           true,
			},
			"insecure_tls": schema.BoolAttribute{
				MarkdownDescription: "Disable TLS certificate verification for the connection",
				Optional:            true,
			},
		},
	}
}

func (p *OpenWrtProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config openWrtProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown OpenWrt host",
			"The provider cannot create the OpenWrt API client as there is an unknown configuration value for the OpenWrt API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OpenWrt_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unnknown OpenWrt username",
			"The provider cannot create the OpenWrt API client as there is an unknown configuration value for the OpenWrt API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OpenWrt_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unnknown OpenWrt password",
			"The provider cannot create the OpenWrt API client as there is an unknown configuration value for the OpenWrt API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OpenWrt_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var clientOpts []gowrt.Option
	if config.InsecureTls.ValueBool() {
		clientOpts = append(clientOpts, gowrt.WithInsecureTls())
	}
	client := gowrt.New(config.Host.ValueString(), clientOpts...)
	err := client.Login(config.Username.ValueString(), config.Password.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to login to OpenWrt API: "+err.Error(),
			"Failed to login to OpenWrt API with the given credentials. Check the host, username and password configuration values.",
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OpenWrtProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *OpenWrtProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
		NewBoardInfoDataSource,
		NewNetworkInterfaceDataSource,
	}
}

func (p *OpenWrtProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpenWrtProvider{
			version: version,
		}
	}
}
