package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vr-ibm/terraform-quantum-aerospace/internal/client"
)

var _ datasource.DataSource = &QuantumBackendDataSource{}
var _ datasource.DataSourceWithConfigure = &QuantumBackendDataSource{}

type QuantumBackendDataSource struct {
	client *client.IonQClient
}

type QuantumBackendDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Qubits        types.Int64  `tfsdk:"qubits"`
	Status        types.String `tfsdk:"status"`
}

func NewQuantumBackendDataSource() datasource.DataSource {
	return &QuantumBackendDataSource{}
}

func (d *QuantumBackendDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*client.IonQClient)
}

func (d *QuantumBackendDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backend"
}

func (d *QuantumBackendDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches metadata about a quantum cloud backend from IonQ.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Backend name (e.g. simulator, ionq.aria-1, ionq.forte-1).",
			},
			"cloud_provider": schema.StringAttribute{
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
		},
	}
}

func (d *QuantumBackendDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data QuantumBackendDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backend, err := d.client.GetBackend(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to read backend %s", data.Name.ValueString()),
			err.Error(),
		)
		return
	}

	data.ID = data.Name
	data.Status = types.StringValue(backend.Status)
	data.Qubits = types.Int64Value(backend.Qubits)
	data.CloudProvider = types.StringValue("ionq")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
