package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"terraform-provider-confluence/internal/helpers"
	"terraform-provider-confluence/internal/provider/transferobjects"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &GroupMembershipDataSource{}

func NewPrivilegesDataSource() datasource.DataSource {
	return &GroupMembershipDataSource{}
}

// GroupMembershipDataSource defines the data source implementation.
type GroupMembershipDataSource struct {
	client *helpers.Client
}

// GroupMembershipDataSourceModel describes the data source data model.
type GroupMembershipDataSourceModel struct {
	GroupId      types.String `tfsdk:"group_id"`
	GroupName    types.String `tfsdk:"group_name"`
	GroupMembers types.Map    `tfsdk:"group_members"`
	Id           types.String `tfsdk:"id"`
}

func (d *GroupMembershipDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_membership"
}

func (d *GroupMembershipDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Group membership data source",

		Attributes: map[string]schema.Attribute{
			"group_id": schema.StringAttribute{
				MarkdownDescription: "The Id of the group",
				Optional:            true,
			},
			"group_name": schema.StringAttribute{
				MarkdownDescription: "The name of the group",
				Optional:            true,
			},
			"group_members": schema.MapAttribute{
				MarkdownDescription: "The members of the group",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Privileges identifier",
				Computed:            true,
			},
		},
	}
}

func (d *GroupMembershipDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*helpers.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *GroupMembershipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GroupMembershipDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.GroupName.IsNull() && data.GroupId.IsNull() {
		resp.Diagnostics.AddError("Validation error", "Please provide either group_id or group_name")
		return
	}
	if !data.GroupName.IsNull() && !data.GroupId.IsNull() {
		resp.Diagnostics.AddError("Validation error", "Please provide only one - either group_id or group_name")
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.GroupName.IsNull() && data.GroupId.IsNull() {
		var response transferobjects.Group
		path := fmt.Sprintf("/rest/api/group/by-name?name=%s", data.GroupName.ValueString())
		if err := d.client.Get(path, &response); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
			return
		}
		data.GroupId = types.StringValue(response.Id)
	}

	// Get the privileges through the API
	var response transferobjects.GroupMembersResponse
	path := fmt.Sprintf("/rest/api/group/%s/membersByGroupId", data.GroupId.ValueString())
	if err := d.client.Get(path, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}

	// Save id into the Terraform state.
	data.Id = types.StringValue(helpers.Sha256String(data.GroupId.ValueString()))

	var elements = make(map[string]attr.Value)

	for _, member := range response.Members {
		elements[member.Email] = types.StringValue(member.AccountID)
	}

	data.GroupMembers, _ = types.MapValue(types.StringType, elements)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
