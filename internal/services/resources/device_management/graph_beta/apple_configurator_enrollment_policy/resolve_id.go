package graphBetaAppleConfiguratorEnrollmentPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/device_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// resolveDepOnboardingSettingsId determines the depOnboardingSetting id to use.
// this id is the 'intuneAccountId' in the /deviceManagement endpoint.
func (r *AppleConfiguratorEnrollmentPolicyResource) resolveDepOnboardingSettingsId(ctx context.Context, provided types.String) (string, error) {
	return device_management.ResolveDepOnboardingSettingsId(ctx, r.client, provided)
}
