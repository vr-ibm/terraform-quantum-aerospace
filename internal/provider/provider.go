package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &QuantumAerospaceProvider{}

type QuantumAerospaceProvider struct {
	version string
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
		},
	}
}

func (p *QuantumAerospaceProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
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
