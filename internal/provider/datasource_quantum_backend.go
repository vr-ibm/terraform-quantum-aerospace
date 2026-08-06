package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &QuantumBackendDataSource{}

type QuantumBackendDataSource struct{}

type QuantumBackendDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Provider    types.String `tfsdk:"provider"`
	Qubits      types.Int64  `tfsdk:"qubits"`
	Status      types.String `tfsdk:"status"`
	MaxShots    types.Int64  `tfsdk:"max_shots"`
	Description types.String `tfsdk:"description"`
}

func NewQuantumBackendDataSource() datasource.DataSource {
	return &QuantumBackendDataSource{}
}

func (d *QuantumBackendDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backend"
}

func (d *QuantumBackendDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches metadata about a quantum cloud backend (IonQ, Braket) including qubit count, status, and shot limits.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Backend name (e.g. ionq.harmony, ionq.aria, braket.sv1).",
			},
			"provider": schema.StringAttribute{
				Computed:    true,
				Description: "Cloud provider hosting the backend (ionq, braket).",
			},
			"qubits": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of qubits available on this backend.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current availability status of the backend.",
			},
			"max_shots": schema.Int64Attribute{
				Computed:    true,
				Description: "Maximum number of shots supported per job.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable description of the backend.",
			},
		},
	}
}

func (d *QuantumBackendDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data QuantumBackendDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// TODO: fetch backend metadata from API
	// Stubbed with realistic IonQ Aria values for now
	data.ID = data.Name
	data.Provider = types.StringValue("ionq")
	data.Qubits = types.Int64Value(25)
	data.Status = types.StringValue("available")
	data.MaxShots = types.Int64Value(10000)
	data.Description = types.StringValue("IonQ Aria — 25-qubit trapped-ion quantum computer.")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
