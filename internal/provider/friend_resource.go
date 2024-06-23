package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// "github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sunwupark/hashicups-client-go"
)

var (
	_ resource.Resource              = &friendResource{}
	_ resource.ResourceWithConfigure = &friendResource{}
)

func NewFriendResource() resource.Resource {
	return &friendResource{}
}

type friendResource struct {
	client *hashicups.Client
}

type friendResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Address     types.String `tfsdk:"address"`
	Description types.String `tfsdk:"description"`
	Image       types.String `tfsdk:"image"`
}

func (r *friendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_friend"
}

func (r *friendResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"address": schema.StringAttribute{
				Optional: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"image": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (r *friendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan friendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	friend := hashicups.Friend{
		Name:        plan.Name.ValueString(),
		Address:     plan.Address.ValueString(),
		Description: plan.Description.ValueString(),
		Image:       plan.Image.ValueString(),
	}

	createdFriend, err := r.client.CreateFriend([]hashicups.Friend{friend})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating friend",
			"Could not create friend, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(createdFriend.ID))
	plan.Name = types.StringValue(createdFriend.Name)
	plan.Address = types.StringValue(createdFriend.Address)
	plan.Description = types.StringValue(createdFriend.Description)
	plan.Image = types.StringValue(createdFriend.Image)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *friendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state friendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	friendID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting friend ID",
			"Could not convert friend ID to integer: "+err.Error(),
		)
		return
	}

	// Assume Getfriend now returns a list of friends
	friends, err := r.client.GetFriend(strconv.Itoa(friendID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups friend",
			err.Error(),
		)
		return
	}

	if len(friends) == 0 {
		resp.Diagnostics.AddError(
			"friend Not Found",
			"No friend found with the given ID",
		)
		return
	}

	friend := friends[0]

	// Map response body to model
	state.ID = types.StringValue(strconv.Itoa(friend.ID))
	state.Address = types.StringValue(friend.Address)
	state.Image = types.StringValue(friend.Image)
	state.Name = types.StringValue(friend.Name)
	state.Description = types.StringValue(friend.Description)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *friendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan friendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert the ID from string to int
	friendID, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Converting friend ID",
			"Could not convert friend ID to integer: "+err.Error(),
		)
		return
	}

	// Create a friend object
	friend := hashicups.Friend{
		ID:          friendID, // ID is an int
		Name:        plan.Name.ValueString(),
		Address:     plan.Address.ValueString(),
		Description: plan.Description.ValueString(),
		Image:       plan.Image.ValueString(),
	}

	// Update the existing friend
	updatedFriend, err := r.client.UpdateFriend(plan.ID.ValueString(), []hashicups.Friend{friend})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating HashiCups friend",
			"Could not update friend, unexpected error: "+err.Error(),
		)
		return
	}

	// Update resource state with updated items
	plan.ID = types.StringValue(strconv.Itoa(updatedFriend.ID))
	plan.Name = types.StringValue(updatedFriend.Name)
	plan.Address = types.StringValue(updatedFriend.Address)
	plan.Description = types.StringValue(updatedFriend.Description)
	plan.Image = types.StringValue(updatedFriend.Image)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *friendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state friendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	friendID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups friend",
			"Could not convert friend ID to integer: "+err.Error(),
		)
		return
	}

	err = r.client.DeleteFriend(strconv.Itoa(friendID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Friend",
			"Could not delete friend, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *friendResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*hashicups.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
