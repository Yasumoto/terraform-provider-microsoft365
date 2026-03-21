package graphBetaDepMacOSEnrollmentProfile

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructEnrollmentProfile(ctx context.Context, data *DepMacOSEnrollmentProfileResourceModel) (graphmodels.EnrollmentProfileable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDepMacOSEnrollmentProfile()

	// --- EnrollmentProfile base fields ---
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(data.RequiresUserAuthentication, requestBody.SetRequiresUserAuthentication)
	convert.FrameworkToGraphBool(data.EnableAuthenticationViaCompanyPortal, requestBody.SetEnableAuthenticationViaCompanyPortal)
	convert.FrameworkToGraphBool(data.RequireCompanyPortalOnSetupAssistantEnrolledDevices, requestBody.SetRequireCompanyPortalOnSetupAssistantEnrolledDevices)

	// --- DepEnrollmentBaseProfile fields ---
	convert.FrameworkToGraphBool(data.ConfigurationWebUrl, requestBody.SetConfigurationWebUrl)
	convert.FrameworkToGraphString(data.DeviceNameTemplate, requestBody.SetDeviceNameTemplate)

	// --- Setup Assistant screen control booleans (base profile) ---
	convert.FrameworkToGraphBool(data.AppleIdDisabled, requestBody.SetAppleIdDisabled)
	convert.FrameworkToGraphBool(data.ApplePayDisabled, requestBody.SetApplePayDisabled)
	convert.FrameworkToGraphBool(data.DiagnosticsDisabled, requestBody.SetDiagnosticsDisabled)
	convert.FrameworkToGraphBool(data.DisplayToneSetupDisabled, requestBody.SetDisplayToneSetupDisabled)
	convert.FrameworkToGraphBool(data.LocationDisabled, requestBody.SetLocationDisabled)
	convert.FrameworkToGraphBool(data.PrivacyPaneDisabled, requestBody.SetPrivacyPaneDisabled)
	convert.FrameworkToGraphBool(data.RestoreBlocked, requestBody.SetRestoreBlocked)
	convert.FrameworkToGraphBool(data.ScreenTimeScreenDisabled, requestBody.SetScreenTimeScreenDisabled)
	convert.FrameworkToGraphBool(data.SiriDisabled, requestBody.SetSiriDisabled)
	convert.FrameworkToGraphBool(data.TermsAndConditionsDisabled, requestBody.SetTermsAndConditionsDisabled)
	convert.FrameworkToGraphBool(data.TouchIdDisabled, requestBody.SetTouchIdDisabled)

	// Auto-generate enabledSkipKeys array from Setup Assistant boolean fields.
	// Note: Privacy and Registration screens are excluded from the skip keys array
	// because the Microsoft Graph API rejects these strings, though their boolean
	// properties (privacyPaneDisabled, registrationDisabled) work correctly.
	requestBody.SetEnabledSkipKeys(buildEnabledSkipKeys(data))

	// Note: enrollmentTimeAzureAdGroupIds is read-only for legacy Apple DEP profiles.
	// The Graph API accepts it in PATCH/POST but silently drops the value.
	// Enrollment time grouping support for Apple ADE is expected in a future release.

	convert.FrameworkToGraphBool(data.IsMandatory, requestBody.SetIsMandatory)
	convert.FrameworkToGraphBool(data.ProfileRemovalDisabled, requestBody.SetProfileRemovalDisabled)
	convert.FrameworkToGraphBool(data.SupervisedModeEnabled, requestBody.SetSupervisedModeEnabled)
	convert.FrameworkToGraphString(data.SupportDepartment, requestBody.SetSupportDepartment)
	convert.FrameworkToGraphString(data.SupportPhoneNumber, requestBody.SetSupportPhoneNumber)
	convert.FrameworkToGraphBool(data.WaitForDeviceConfiguredConfirmation, requestBody.SetWaitForDeviceConfiguredConfirmation)

	// --- DepMacOSEnrollmentProfile: Setup Assistant screen control booleans (macOS-specific) ---
	convert.FrameworkToGraphBool(data.AccessibilityScreenDisabled, requestBody.SetAccessibilityScreenDisabled)
	convert.FrameworkToGraphBool(data.AutoUnlockWithWatchDisabled, requestBody.SetAutoUnlockWithWatchDisabled)
	convert.FrameworkToGraphBool(data.ChooseYourLockScreenDisabled, requestBody.SetChooseYourLockScreenDisabled)
	convert.FrameworkToGraphBool(data.FileVaultDisabled, requestBody.SetFileVaultDisabled)
	convert.FrameworkToGraphBool(data.ICloudDiagnosticsDisabled, requestBody.SetICloudDiagnosticsDisabled)
	convert.FrameworkToGraphBool(data.ICloudStorageDisabled, requestBody.SetICloudStorageDisabled)
	convert.FrameworkToGraphBool(data.PassCodeDisabled, requestBody.SetPassCodeDisabled)
	convert.FrameworkToGraphBool(data.RegistrationDisabled, requestBody.SetRegistrationDisabled)
	convert.FrameworkToGraphBool(data.ZoomDisabled, requestBody.SetZoomDisabled)

	// --- DepMacOSEnrollmentProfile: Account and enrollment behavior fields ---
	// Admin account fields: explicitly send nil when null OR empty string to clear them in the API
	convert.FrameworkToGraphStringOrNil(data.AdminAccountFullName, requestBody.SetAdminAccountFullName)
	convert.FrameworkToGraphStringOrNil(data.AdminAccountPassword, requestBody.SetAdminAccountPassword)
	convert.FrameworkToGraphStringOrNil(data.AdminAccountUserName, requestBody.SetAdminAccountUserName)

	convert.FrameworkToGraphBool(data.AutoAdvanceSetupEnabled, requestBody.SetAutoAdvanceSetupEnabled)
	convert.FrameworkToGraphBool(data.DontAutoPopulatePrimaryAccountInfo, requestBody.SetDontAutoPopulatePrimaryAccountInfo)
	convert.FrameworkToGraphBool(data.EnableRestrictEditing, requestBody.SetEnableRestrictEditing)
	convert.FrameworkToGraphBool(data.HideAdminAccount, requestBody.SetHideAdminAccount)
	convert.FrameworkToGraphString(data.PrimaryAccountFullName, requestBody.SetPrimaryAccountFullName)
	convert.FrameworkToGraphString(data.PrimaryAccountUserName, requestBody.SetPrimaryAccountUserName)
	convert.FrameworkToGraphBool(data.RequestRequiresNetworkTether, requestBody.SetRequestRequiresNetworkTether)
	convert.FrameworkToGraphBool(data.SetPrimarySetupAccountAsRegularUser, requestBody.SetSetPrimarySetupAccountAsRegularUser)
	convert.FrameworkToGraphBool(data.SkipPrimarySetupAccountCreation, requestBody.SetSkipPrimarySetupAccountCreation)

	// --- LAPS rotation setting nested object ---
	// Explicitly set to nil when null to clear LAPS configuration in the API
	if data.DepProfileAdminAccountPasswordRotationSetting.IsNull() {
		requestBody.SetDepProfileAdminAccountPasswordRotationSetting(nil)
	} else if !data.DepProfileAdminAccountPasswordRotationSetting.IsUnknown() {
		var rotationModel DepProfileAdminAccountPasswordRotationSettingModel
		diags := data.DepProfileAdminAccountPasswordRotationSetting.As(ctx, &rotationModel, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("extracting dep_profile_admin_account_password_rotation_setting: %s", diagsToString(diags))
		}

		rotationSetting := graphmodels.NewDepProfileAdminAccountPasswordRotationSetting()
		convert.FrameworkToGraphInt32(rotationModel.AutoRotationPeriodInDays, rotationSetting.SetAutoRotationPeriodInDays)

		// Map delay auto-rotation sub-object
		if !rotationModel.DepProfileDelayAutoRotationSetting.IsNull() && !rotationModel.DepProfileDelayAutoRotationSetting.IsUnknown() {
			var delayModel DepProfileDelayAutoRotationSettingModel
			diags := rotationModel.DepProfileDelayAutoRotationSetting.As(ctx, &delayModel, basetypes.ObjectAsOptions{})
			if diags.HasError() {
				return nil, fmt.Errorf("extracting dep_profile_delay_auto_rotation_setting: %s", diagsToString(diags))
			}
			delaySetting := graphmodels.NewDepProfileDelayAutoRotationSetting()
			convert.FrameworkToGraphBool(delayModel.OnRetrievalAutoRotatePasswordEnabled, delaySetting.SetOnRetrievalAutoRotatePasswordEnabled)
			convert.FrameworkToGraphInt32(delayModel.OnRetrievalDelayAutoRotatePasswordInHours, delaySetting.SetOnRetrievalDelayAutoRotatePasswordInHours)
			rotationSetting.SetDepProfileDelayAutoRotationSetting(delaySetting)
		}

		requestBody.SetDepProfileAdminAccountPasswordRotationSetting(rotationSetting)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// skipKeyMapping maps a boolean model field to the corresponding Microsoft Graph API skip key string.
type skipKeyMapping struct {
	disabled func() types.Bool
	key      string
}

// buildEnabledSkipKeys constructs the enabledSkipKeys array from boolean fields.
// This is called internally when constructing API requests - users don't set skip keys directly.
//
// Note: Privacy and Registration are intentionally excluded because Microsoft Graph API
// rejects these as skip key values, even though the corresponding boolean properties
// (privacyPaneDisabled, registrationDisabled) work correctly. This is a known API limitation.
func buildEnabledSkipKeys(data *DepMacOSEnrollmentProfileResourceModel) []string {
	mappings := []skipKeyMapping{
		// Base profile screens (excluding Privacy - API rejects "Privacy" skip key)
		{func() types.Bool { return data.AppleIdDisabled }, "AppleID"},
		{func() types.Bool { return data.ApplePayDisabled }, "Payment"},
		{func() types.Bool { return data.DiagnosticsDisabled }, "Diagnostics"},
		{func() types.Bool { return data.DisplayToneSetupDisabled }, "DisplayTone"},
		{func() types.Bool { return data.LocationDisabled }, "Location"},
		{func() types.Bool { return data.RestoreBlocked }, "Restore"},
		{func() types.Bool { return data.ScreenTimeScreenDisabled }, "ScreenTime"},
		{func() types.Bool { return data.SiriDisabled }, "Siri"},
		{func() types.Bool { return data.TermsAndConditionsDisabled }, "TOS"},
		{func() types.Bool { return data.TouchIdDisabled }, "Biometric"},
		// macOS-specific screens (excluding Registration - API rejects "Registration" skip key)
		{func() types.Bool { return data.WelcomeScreenDisabled }, "Welcome"},
		{func() types.Bool { return data.AccessibilityScreenDisabled }, "Accessibility"},
		{func() types.Bool { return data.AutoUnlockWithWatchDisabled }, "UnlockWithWatch"},
		{func() types.Bool { return data.ChooseYourLockScreenDisabled }, "Wallpaper"},
		{func() types.Bool { return data.FileVaultDisabled }, "FileVault"},
		{func() types.Bool { return data.ICloudDiagnosticsDisabled }, "iCloudDiagnostics"},
		{func() types.Bool { return data.ICloudStorageDisabled }, "iCloudStorage"},
		{func() types.Bool { return data.PassCodeDisabled }, "Passcode"},
		{func() types.Bool { return data.ZoomDisabled }, "Zoom"},
	}

	skipKeys := make([]string, 0, len(mappings))
	for _, m := range mappings {
		if m.disabled().ValueBool() {
			skipKeys = append(skipKeys, m.key)
		}
	}
	return skipKeys
}

// diagsToString converts diagnostics to a single error string for wrapping.
func diagsToString(diags diag.Diagnostics) string {
	var sb strings.Builder
	for _, d := range diags {
		sb.WriteString(d.Summary())
		sb.WriteString(": ")
		sb.WriteString(d.Detail())
		sb.WriteString("; ")
	}
	return sb.String()
}
