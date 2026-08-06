package provider

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vr-ibm/terraform-quantum-aerospace/internal/client"
)

var _ resource.Resource = &QuantumJobResource{}
var _ resource.ResourceWithConfigure = &QuantumJobResource{}
var _ resource.ResourceWithImportState = &QuantumJobResource{}

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
	Timeout types.String `tfsdk:"timeout"`
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
		Description: "Submits a quantum circuit to a cloud backend for execution. Designed for aviation optimization workloads including route selection, fuel burn minimization, and emissions reduction using QAOA.",
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
			"timeout": schema.StringAttribute{
				Optional:    true,
				Description: "Maximum time to wait for job completion (e.g. 5m, 10m). Defaults to 10m.",
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

	// Build circuit payload
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

	// Parse timeout
	pollConfig := client.DefaultPollConfig()
	if !data.Timeout.IsNull() && data.Timeout.ValueString() != "" {
		duration, err := time.ParseDuration(data.Timeout.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid timeout", err.Error())
			return
		}
		pollConfig.Timeout = duration
	}

	// Poll until complete
	finalResp, err := r.client.PollJobUntilComplete(ctx, jobResp.ID, pollConfig)
	if err != nil {
		// Job was submitted but didn't complete in time — save state anyway
		data.Status = types.StringValue(jobResp.Status)
		data.Results = types.StringValue("{}")
		resp.Diagnostics.AddWarning("Job submitted but not yet complete", err.Error())
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	data.Status = types.StringValue(finalResp.Status)
	if finalResp.Results != nil {
		data.Results = types.StringValue(string(finalResp.Results))
	} else {
		data.Results = types.StringValue("{}")
	}
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

func (r *QuantumJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	jobResp, err := r.client.GetJob(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to import quantum job", err.Error())
		return
	}

	var data QuantumJobResourceModel
	data.ID = types.StringValue(jobResp.ID)
	data.Backend = types.StringValue(jobResp.Target)
	data.Status = types.StringValue(jobResp.Status)

	// Circuit is not returned by the IonQ API — it's a write-only attribute.
	// This is a fundamental tension: immutable resources lose their input after submission.
	data.Circuit = types.StringValue("")

	// Shots not returned in GET — default to unknown.
	data.Shots = types.Int64Value(0)
	data.Timeout = types.StringNull()

	if jobResp.Results != nil {
		data.Results = types.StringValue(string(jobResp.Results))
	} else {
		data.Results = types.StringValue("{}")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
