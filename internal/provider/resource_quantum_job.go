package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vr-ibm/terraform-quantum-aerospace/internal/client"
)

var _ resource.Resource = &QuantumJobResource{}
var _ resource.ResourceWithConfigure = &QuantumJobResource{}

type QuantumJobResource struct {
	client *client.IonQClient
}

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

func (r *QuantumJobResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*client.IonQClient)
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

	// Build a simple circuit payload for IonQ
	// TODO: parse OpenQASM into IonQ native format
	circuitJSON := json.RawMessage(`{"qubits": 4, "circuit": []}`)
	jobResp, err := r.client.CreateJob(ctx, client.JobInput{
		Target:  data.Backend.ValueString(),
		Shots:   data.Shots.ValueInt64(),
		Circuit: circuitJSON,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to create quantum job", err.Error())
		return
	}

	data.ID = types.StringValue(jobResp.ID)
	data.Status = types.StringValue(jobResp.Status)
	data.Results = types.StringValue("{}")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuantumJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data QuantumJobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jobResp, err := r.client.GetJob(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read quantum job", err.Error())
		return
	}

	data.Status = types.StringValue(jobResp.Status)
	if jobResp.Results != nil {
		data.Results = types.StringValue(string(jobResp.Results))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QuantumJobResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// Jobs are immutable — circuit and backend both RequiresReplace
}

func (r *QuantumJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data QuantumJobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CancelJob(ctx, data.ID.ValueString())
	if err != nil {
		// Job may already be completed — not a fatal error
		resp.Diagnostics.AddWarning("Could not cancel job", err.Error())
	}
}
