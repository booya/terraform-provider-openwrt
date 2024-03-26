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
	_ datasource.DataSource              = &BoardInfoDataSource{}
	_ datasource.DataSourceWithConfigure = &BoardInfoDataSource{}
)

func NewBoardInfoDataSource() datasource.DataSource {
	return &BoardInfoDataSource{}
}

// BoardInfoDataSource defines the data source implementation.
type BoardInfoDataSource struct {
	client *gowrt.Client
}

type boardInfoModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type boardInfoNetwork struct {
	Device   types.String   `tfsdk:"device"`
	Ports    []types.String `tfsdk:"ports"`
	Protocol types.String   `tfsdk:"protocol"`
}

type boardInfoLed struct {
	Name   types.String `tfsdk:"name"`
	Sysfs  types.String `tfsdk:"sysfs"`
	Device types.String `tfsdk:"device"`
	Type   types.String `tfsdk:"type"`
	Mode   types.String `tfsdk:"mode"`
}

type boardInfoSystem struct {
	CompatVersion string `tfsdk:"compat_version"`
}

// boardInfoDataSourceModel describes the data source data model.
type boardInfoDataSourceModel struct {
	Id        types.String `tfsdk:"id"`
	ModelId   types.String `tfsdk:"model_id"`
	ModelName types.String `tfsdk:"model_name"`
	// Model   boardInfoModel              `tfsdk:"model"`
	// Led     map[string]boardInfoLed     `tfsdk:"led"`
	// Network map[string]boardInfoNetwork `tfsdk:"network"`
	// System  map[string]boardInfoSystem  `tfsdk:"system"`
}

func (d *BoardInfoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_board_info"
}

func (d *BoardInfoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Board Information data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				// Optional: true,
			},
			"model_id": schema.StringAttribute{
				MarkdownDescription: "Model ID",
				// Computed:            true,
				Optional: true,
			},
			"model_name": schema.StringAttribute{
				MarkdownDescription: "Model Name",
				// Computed:            true,
				Optional: true,
			},
			// "network": schema.MapAttribute{
			// 	MarkdownDescription: "Network information",
			// 	ElementType:         types.StringType,
			// 	// Computed:            true,
			// 	Optional: true,
			// },
			// "led": schema.MapAttribute{
			// 	MarkdownDescription: "Model information",
			// 	ElementType:         types.StringType,
			// 	// Computed:            true,
			// 	Optional: true,
			// },
			// "system": schema.MapAttribute{
			// 	MarkdownDescription: "System information",
			// 	ElementType:         types.StringType,
			// 	// Computed:            true,
			// 	Optional: true,
			// },
		},
	}
}

func (d *BoardInfoDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BoardInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data boardInfoDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := d.client.GetBoardInfo()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read board info: %s", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("API RESPONSE: %#v", apiResp))

	// Hard Coding a Id value to save into the Terraform state.
	data.Id = types.StringValue("board-info")
	data.ModelId = types.StringValue(apiResp.Model.Id)
	data.ModelName = types.StringValue(apiResp.Model.Name)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
