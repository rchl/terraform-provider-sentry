package sentry

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jianyuan/go-sentry/v2/sentry"
)

func dataSourceSentryProject() *schema.Resource {
	return &schema.Resource{
		Description: "Sentry Project data source.",

		ReadContext: dataSourceSentryProjectRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Description: "The slug of the organization for this project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "The unique URL slug for this project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"internal_id": {
				Description: "The internal ID for this project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"platform": {
				Description: "The platform for this project.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceSentryProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*sentry.Client)

	org := d.Get("organization").(string)
	slug := d.Get("slug").(string)

	tflog.Debug(ctx, "Reading project", map[string]interface{}{
		"projectSlug": slug,
		"org":         org,
	})
	proj, _, err := client.Projects.Get(ctx, org, slug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(proj.Slug)
	retErr := multierror.Append(
		d.Set("organization", org),
		d.Set("slug", proj.Slug),
		d.Set("internal_id", proj.ID),
		d.Set("platform", proj.Platform),
	)
	return diag.FromErr(retErr.ErrorOrNil())
}
