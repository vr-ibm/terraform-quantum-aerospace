package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vr-ibm/terraform-quantum-aerospace/internal/client"
)

var _ provider.Provider = &QuantumAerospaceProvider{}

type QuantumAerospaceProvider struct {
	version string
}

type QuantumAerospaceProviderModel struct {
	Backend types.String `tfsdk:"backend"`
	APIKey  types.String `tfsdk:"api_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &QuantumAerospaceProvider{
			version: version,
		}
	}
}

func (p *QuantumAerospaceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "quantum"
	resp.Version = p.version
}

func (p *QuantumAerospaceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for quantum cloud backends applied to aerospace optimization problems.",
		Attributes: map[string]schema.Attribute{
			"backend": schema.StringAttribute{
				Description: "Quantum backend to use (ionq, braket).",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "API key for the quantum backend. Can also be set via IONQ_API_KEY env var.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *QuantumAerospaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config QuantumAerospaceProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve API key: config > env var
	apiKey := config.APIKey.ValueString()
	if apiKey == "" {
		apiKey = os.Getenv("IONQ_API_KEY")
	}
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"Set api_key in the provider block or the IONQ_API_KEY environment variable.",
		)
		return
	}

	ionqClient := client.NewIonQClient(apiKey)
	resp.DataSourceData = ionqClient
	resp.ResourceData = ionqClient
}

func (p *QuantumAerospaceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewQuantumJobResource,
	}
}

func (p *QuantumAerospaceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewQuantumBackendDataSource,
	}
}
