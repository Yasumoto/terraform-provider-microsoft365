package device_management

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// ResolveDepOnboardingSettingsId determines the depOnboardingSetting ID to use.
// This ID is the 'intuneAccountId' in the /deviceManagement endpoint.
// If a value is provided, it is returned as-is. Otherwise, the function fetches
// the intuneAccountId from the /deviceManagement endpoint.
func ResolveDepOnboardingSettingsId(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, provided types.String) (string, error) {
	if !provided.IsNull() && !provided.IsUnknown() && provided.ValueString() != "" {
		return provided.ValueString(), nil
	}

	reqConfig := &devicemanagement.DeviceManagementRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.DeviceManagementRequestBuilderGetQueryParameters{
			Select: []string{"intuneAccountId"},
		},
	}

	dm, err := client.
		DeviceManagement().
		Get(ctx, reqConfig)

	if err != nil {
		return "", fmt.Errorf("failed to GET /deviceManagement to resolve intuneAccountId: %w", err)
	}

	if dm == nil {
		return "", fmt.Errorf("deviceManagement response is nil")
	}

	intuneAccountId := dm.GetIntuneAccountId()
	if intuneAccountId == nil {
		return "", fmt.Errorf("intuneAccountId is nil in deviceManagement response")
	}

	return intuneAccountId.String(), nil
}
