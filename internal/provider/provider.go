package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-confluence/internal/helpers"
)

// Ensure ConfluenceProvider satisfies various provider interfaces.
var _ provider.Provider = &ConfluenceProvider{}

// ConfluenceProvider defines the provider implementation.
type ConfluenceProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ConfluenceProviderModel describes the provider data model.
type ConfluenceProviderModel struct {
	Site          types.String `tfsdk:"site"`
	SiteTLS       types.Bool   `tfsdk:"site_tls"`
	PublicSite    types.String `tfsdk:"public_site"`
	PublicSiteTLS types.Bool   `tfsdk:"public_site_tls"`
	Context       types.String `tfsdk:"context"`

	Username types.String `tfsdk:"user"`
	Token    types.String `tfsdk:"token"`
}

func (p *ConfluenceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "confluence"
	resp.Version = p.version
}

func (p *ConfluenceProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"site": schema.StringAttribute{
				MarkdownDescription: "Confluence hostname (<name>.atlassian.net if using Cloud Confluence, otherwise hostname)",
				Required:            true,
			},
			"site_tls": schema.BoolAttribute{
				MarkdownDescription: "Use https for API calls",
				Optional:            true,
			},
			"public_site": schema.StringAttribute{
				MarkdownDescription: "Optional public Confluence Server hostname if different than API hostname",
				Optional:            true,
			},
			"public_site_tls": schema.BoolAttribute{
				MarkdownDescription: "Use https for public site URLs",
				Optional:            true,
			},
			"context": schema.StringAttribute{
				MarkdownDescription: "Confluence path context (Will default to /wiki if using an atlassian.net hostname)",
				Optional:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "User's email address for Cloud Confluence or username for Confluence Server",
				Required:            true,
				Sensitive:           true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Confluence API Token for Cloud Confluence or password for Confluence Server/Cloud",
				Required:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *ConfluenceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ConfluenceProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	site := ""
	siteTls := true
	publicSite := ""
	publicSiteTls := true
	context := ""
	username := "user"
	password := "password"

	if !data.Site.IsNull() {
		site = data.Site.ValueString()
		if strings.HasSuffix(site, ".atlassian.net") && data.Context.IsNull() {
			context = "/wiki"
		}
		publicSite = data.Site.ValueString()
	}

	if !data.SiteTLS.IsNull() {
		siteTls = data.SiteTLS.ValueBool()
		publicSiteTls = data.SiteTLS.ValueBool()
	}

	if !data.PublicSite.IsNull() {
		publicSite = data.PublicSite.ValueString()
	}

	if !data.PublicSiteTLS.IsNull() {
		publicSiteTls = data.PublicSiteTLS.ValueBool()
	}

	if !data.Context.IsNull() {
		context = data.Context.ValueString()
	}

	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if !data.Token.IsNull() {
		password = data.Token.ValueString()
	}

	// Example client configuration for data sources and resources
	client := helpers.NewClient(&helpers.NewClientInput{
		Site:             site,
		SiteUseTLS:       siteTls,
		PublicSite:       publicSite,
		PublicSiteUseTLS: publicSiteTls,
		Context:          context,
		Username:         username,
		Password:         password,
	})

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ConfluenceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGroupResource,
		NewSpaceResource,
		NewSpacePermissionResource,
	}
}

func (p *ConfluenceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPrivilegesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ConfluenceProvider{
			version: version,
		}
	}
}
