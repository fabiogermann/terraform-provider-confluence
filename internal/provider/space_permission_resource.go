package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"sort"
	"strconv"
	"strings"
	"terraform-provider-confluence/internal/helpers"
	"terraform-provider-confluence/internal/provider/transferobjects"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &SpacePermissionResource{}
var _ resource.ResourceWithImportState = &SpacePermissionResource{}
var validPermissions = []string{
	"create:page", "create:blogpost", "create:comment", "create:attachment",
	"read:space",
	"delete:space", "delete:page", "delete:blogpost", "delete:comment", "delete:attachment",
	"export:space",
	"administer:space",
	"archive:page",
	"restrict_content:space",
}

func NewSpacePermissionResource() resource.Resource {
	return &SpacePermissionResource{}
}

// SpacePermissionResource defines the resource implementation.
type SpacePermissionResource struct {
	client *helpers.Client
}

// SpacePermissionResourceModel describes the resource data model.
type SpacePermissionResourceModel struct {
	Key          types.String `tfsdk:"key"`
	Operations   types.List   `tfsdk:"operations"`
	OperationIds types.Map    `tfsdk:"operation_ids"`
	Group        types.String `tfsdk:"group"`
	Id           types.String `tfsdk:"id"`
}

func (r *SpacePermissionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space_permission"
}

func (r *SpacePermissionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Space permission resource",

		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				MarkdownDescription: "The space key of the confluence space (all caps)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"operations": schema.ListAttribute{
				MarkdownDescription: "The operations allowed for the group",
				ElementType:         types.StringType,
				Required:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"operation_ids": schema.MapAttribute{
				MarkdownDescription: "The operation's ids for the group",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "The group that is allowed",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Resource identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SpacePermissionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*helpers.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *SpacePermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SpacePermissionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	permissionRequests := spacePermissionMappingFromResourceModel(ctx, data)

	// Create the rule through API
	var createdIds []string
	var elements = make(map[string]attr.Value)
	for _, body := range permissionRequests {
		var response transferobjects.SpacePermission
		path := fmt.Sprintf("/rest/api/space/%s/permission", data.Key.ValueString())
		if err := r.client.Post(path, body, &response, []string{}); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: \n%s", err))
			return
		}
		createdIds = append(createdIds, response.Id.String())
		key := fmt.Sprintf("%s:%s", body.Operation.Key, body.Operation.Target)
		elements[key] = types.StringValue(response.Id.String())
	}

	data.OperationIds, _ = types.MapValue(types.StringType, elements)

	// Save id into the Terraform state.
	sort.Strings(createdIds)
	data.Id = types.StringValue(strings.Join(createdIds[:], ":"))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpacePermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SpacePermissionResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rule through the API
	var response transferobjects.SummarySpacePermissions
	path := fmt.Sprintf("/rest/api/space/%s?expand=permissions", data.Key.ValueString())
	if err := r.client.Get(path, &response); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		return
	}

	resourceId, operationIds := generateIdFromSummaryResponse(ctx, data, &response)
	data.OperationIds, _ = types.MapValue(types.StringType, operationIds)
	data.Id = types.StringValue(resourceId)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpacePermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SpacePermissionResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Update not supported
	resp.Diagnostics.AddError("Schema Error", fmt.Sprintf("UPDATE operation not supported"))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpacePermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SpacePermissionResourceModel
	var permissions map[string]string

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.OperationIds.ElementsAs(ctx, &permissions, false)

	// Get the rule through the API
	for permission, permissionId := range permissions {
		path := fmt.Sprintf("/rest/api/space/%s/permission/%s", data.Key.ValueString(), permissionId)
		if err := r.client.Delete(path); err != nil {
			errorMsg := fmt.Sprintf("Error while deleting permission [%s][%s]: %s", permission, permissionId, err.Error())
			tflog.Warn(ctx, errorMsg)
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Error during request, got error: %s", err))
		}
	}

}

func (r *SpacePermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func spacePermissionMappingFromResourceModel(ctx context.Context, data *SpacePermissionResourceModel) []*transferobjects.SpacePermission {
	var collection []*transferobjects.SpacePermission
	var permissions []string

	data.Operations.ElementsAs(ctx, &permissions, false)

	if helpers.Contains(permissions, "read:space") && len(permissions) > 1 && permissions[0] != "read:space" {
		permissions = helpers.MoveToFirstPositionOfSlice(permissions, "read:space")
	}
	for _, permission := range permissions {
		permissionParts := strings.Split(permission, ":")
		subject := &transferobjects.Subject{
			Type:       "group",
			Identifier: data.Group.ValueString(),
		}
		operation := &transferobjects.Operation{
			Key:    permissionParts[0],
			Target: permissionParts[1],
		}
		spacePermission := &transferobjects.SpacePermission{
			Id:        0,
			Subject:   subject,
			Operation: operation,
		}
		collection = append(collection, spacePermission)
	}
	return collection
}
func generateIdFromSummaryResponse(ctx context.Context, data *SpacePermissionResourceModel, spacePermissions *transferobjects.SummarySpacePermissions) (string, map[string]attr.Value) {
	var permissionIds []string
	var elements = make(map[string]attr.Value)
	for _, permission := range spacePermissions.Permissions {
		if permission.Subjects.Group != nil && len(permission.Subjects.Group.Results) > 0 {
			for _, group := range permission.Subjects.Group.Results {
				if data.Group.ValueString() == group.Name {
					permissionIds = append(permissionIds, strconv.Itoa(permission.ID))
					key := fmt.Sprintf("%s:%s", permission.Operation.Operation, permission.Operation.TargetType)
					elements[key] = types.StringValue(strconv.Itoa(permission.ID))
				}
			}
		}
	}
	sort.Strings(permissionIds)
	return strings.Join(permissionIds[:], ":"), elements
}
