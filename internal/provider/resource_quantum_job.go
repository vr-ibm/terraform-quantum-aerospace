package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &QuantumJobResource{}

type QuantumJobResource struct{}

type QuantumJobResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Backend types.String `tfsdk:"backend"`
	Circuit types.String `tfsdk:"circuit"`
	Shots   types.Int64  `tfsdk:"shots"`
	Status  types.String `tfsdk:"status"`
	Results types.String `tfsdk:"results"`
}

func NewQuantumJobResource() resource.Resource {
	return &QuantumJobResource{}
}

func (r *QuantumJobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job"
}

func (r *QuantumJobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a quantum job submitted to a cloud quantum backend.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"backend": schema.StringAttribute{
				Required:    true,
				Description: "Quantum backend to submit the job to (ionq, braket).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"circuit": schema.StringAttribute{
				Required:    true,
				Description: "Quantum circuit definition (OpenQASM 3.0 string).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"shots": schema.Int64Attribute{
				Required:    true,
				Description: "Number of times to execute the circuit.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Job status returned by the backend (queued, running, completed, failed).",
			},
			"results": schema.StringAttribute{
				Computed:    true,
				Description: "JSON-encoded measurement results from the backend.",
			},
		},
	}
}

func (r *QuantumJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data QuantumJobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// TODO: submit job to backend API
	data.ID = types.StringValue("placeholder-job-id")
	data.Status = types.StringValue("queued")
	data.Results = types.StringValue("{}")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuantumJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data QuantumJobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// TODO: poll job status from backend API
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuantumJobResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// Jobs are immutable — circuit and backend both RequiresReplace
}

func (r *QuantumJobResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// TODO: cancel job if still running
}
