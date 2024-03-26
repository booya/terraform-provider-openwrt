// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/booya/gowrt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &NetworkInterfaceDataSource{}
	_ datasource.DataSourceWithConfigure = &NetworkInterfaceDataSource{}
)

func NewNetworkInterfaceDataSource() datasource.DataSource {
	return &NetworkInterfaceDataSource{}
}

// NetworkInterfaceDataSource defines the data source implementation.
type NetworkInterfaceDataSource struct {
	client *gowrt.Client
}

// NetworkInterfaceDataSourceModel describes the data source data model.
type networkInterfaceModel struct {
	Id       types.Int64  `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Device   types.String `tfsdk:"device"`
	Proto    types.String `tfsdk:"proto"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	// Ipv6       types.String `tfsdk:"ipv6"`
	// Ip6Assign  types.String `tfsdk:"ip6assign"`
	// ReqAddress types.String `tfsdk:"reqaddress"`
	// ReqPrefix  types.String `tfsdk:"reqprefix"`
}

func (d *NetworkInterfaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_interface"
}

func (d *NetworkInterfaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Network Configuration data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				// Optional: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Interface Name",
				// Computed:            true,
				// Optional: true,
				Required: true,
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "Device Name",
				// Computed:            true,
				Optional: true,
			},
			"proto": schema.StringAttribute{
				MarkdownDescription: "Network proto",
				// Computed:            true,
				Optional: true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Interface Login Username",
				// Computed:            true,
				Optional: true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Interface Login Password",
				// Computed:            true,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (d *NetworkInterfaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*gowrt.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *gowrt.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *NetworkInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data networkInterfaceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read interface %s from api", data.Name.ValueString()))
	apiResp, err := d.client.GetInterfaceConfiguration(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read network interface: %s", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("Api Response: %#v", apiResp))

	// Hard Coding a Id value to save into the Terraform state.
	data.Id = types.Int64Value(int64(apiResp.Index))
	data.Name = types.StringValue(apiResp.Name)
	data.Device = types.StringValue(apiResp.Device)
	data.Proto = types.StringValue(apiResp.Proto)
	data.Username = types.StringValue(apiResp.Username)
	data.Password = types.StringValue(apiResp.Password)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "finished reading network interface data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
