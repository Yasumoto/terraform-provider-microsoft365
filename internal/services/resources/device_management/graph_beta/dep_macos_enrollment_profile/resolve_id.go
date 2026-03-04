package graphBetaDepMacOSEnrollmentProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/device_management"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// getOrResolveDepId resolves the DEP onboarding settings ID, either from the
// provided value or by fetching the intuneAccountId from /deviceManagement.
// Handles diagnostics automatically. Returns the ID and a boolean indicating success.
func (r *DepMacOSEnrollmentProfileResource) getOrResolveDepId(ctx context.Context, provided types.String, diags *diag.Diagnostics) (string, bool) {
	depId, err := device_management.ResolveDepOnboardingSettingsId(ctx, r.client, provided)
	if err != nil {
		diags.AddError("Failed to resolve dep_onboarding_settings_id", err.Error())
		return "", false
	}
	return depId, true
}
