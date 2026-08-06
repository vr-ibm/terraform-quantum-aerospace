package provider

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &QuantumCircuitResource{}

type QuantumCircuitResource struct{}

type QuantumCircuitResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Body        types.String `tfsdk:"body"`
	Qubits      types.Int64  `tfsdk:"qubits"`
	Hash        types.String `tfsdk:"hash"`
}

func NewQuantumCircuitResource() resource.Resource {
	return &QuantumCircuitResource{}
}

func (r *QuantumCircuitResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_circuit"
}

func (r *QuantumCircuitResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Defines a reusable quantum circuit. This is a config-only resource — it has no backing API and exists only in Terraform state. It enables circuit composition and separation of concerns between circuit authoring and job submission.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Human-readable circuit name for reference.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of what this circuit computes.",
			},
			"body": schema.StringAttribute{
				Required:    true,
				Description: "Circuit definition (OpenQASM 3.0).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"qubits": schema.Int64Attribute{
				Required:    true,
				Description: "Number of qubits used by this circuit.",
			},
			"hash": schema.StringAttribute{
				Computed:    true,
				Description: "SHA-256 hash of the circuit body. Changes trigger replacement of dependent jobs.",
			},
		},
	}
}

func (r *QuantumCircuitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data QuantumCircuitResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	hash := sha256.Sum256([]byte(data.Body.ValueString()))
	data.ID = types.StringValue(fmt.Sprintf("%x", hash[:8]))
	data.Hash = types.StringValue(fmt.Sprintf("%x", hash))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuantumCircuitResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// Config-only resource — state is always authoritative, nothing to refresh
}

func (r *QuantumCircuitResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Only name/description/qubits can change without replace (body triggers RequiresReplace)
	var data QuantumCircuitResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	hash := sha256.Sum256([]byte(data.Body.ValueString()))
	data.ID = types.StringValue(fmt.Sprintf("%x", hash[:8]))
	data.Hash = types.StringValue(fmt.Sprintf("%x", hash))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuantumCircuitResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Config-only — nothing to clean up
}
