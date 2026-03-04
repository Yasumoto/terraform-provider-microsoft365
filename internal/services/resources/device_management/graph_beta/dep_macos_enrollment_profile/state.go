package graphBetaDepMacOSEnrollmentProfile

import (
	"context"
	"fmt"
	"slices"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// delayRotationSettingAttrTypes defines the attribute types for the
// dep_profile_delay_auto_rotation_setting nested sub-object.
var delayRotationSettingAttrTypes = map[string]attr.Type{
	"on_retrieval_auto_rotate_password_enabled":        types.BoolType,
	"on_retrieval_delay_auto_rotate_password_in_hours": types.Int32Type,
}

// rotationSettingAttrTypes defines the attribute types for the
// dep_profile_admin_account_password_rotation_setting nested object.
// Used in both MapRemoteStateToTerraform and setExtendedFieldsToNull.
var rotationSettingAttrTypes = map[string]attr.Type{
	"auto_rotation_period_in_days":            types.Int32Type,
	"dep_profile_delay_auto_rotation_setting": types.ObjectType{AttrTypes: delayRotationSettingAttrTypes},
}

// MapRemoteStateToTerraform maps the EnrollmentProfile to Terraform state.
// depId, when non-empty, will be stated to DepOnboardingSettingsId.
// priorPassword preserves the admin_account_password from state since the API does not return it.
// Note: This function preserves the existing timeouts from the current state.
func MapRemoteStateToTerraform(ctx context.Context, data *DepMacOSEnrollmentProfileResourceModel, profile graphmodels.EnrollmentProfileable, depId string, priorPassword types.String) {
	if profile == nil {
		tflog.Debug(ctx, "Remote enrollmentProfile is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote enrollmentProfile to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(profile.GetId()),
	})

	// --- EnrollmentProfile base fields ---
	data.ID = convert.GraphToFrameworkString(profile.GetId())
	data.DisplayName = convert.GraphToFrameworkString(profile.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(profile.GetDescription())
	data.RequiresUserAuthentication = convert.GraphToFrameworkBool(profile.GetRequiresUserAuthentication())
	data.EnableAuthenticationViaCompanyPortal = convert.GraphToFrameworkBool(profile.GetEnableAuthenticationViaCompanyPortal())
	data.RequireCompanyPortalOnSetupAssistantEnrolledDevices = convert.GraphToFrameworkBool(profile.GetRequireCompanyPortalOnSetupAssistantEnrolledDevices())
	data.ConfigurationEndpointUrl = convert.GraphToFrameworkString(profile.GetConfigurationEndpointUrl())
	if depId != "" {
		data.DepOnboardingSettingsId = types.StringValue(depId)
	}

	// Type-assert to DepMacOSEnrollmentProfileable for extended fields.
	// The SDK's discriminator deserialization creates *DepMacOSEnrollmentProfile
	// when @odata.type is "#microsoft.graph.depMacOSEnrollmentProfile".
	macProfile, ok := profile.(graphmodels.DepMacOSEnrollmentProfileable)
	if !ok {
		tflog.Warn(ctx, "Profile is not a DepMacOSEnrollmentProfile; extended fields will be null")
		setExtendedFieldsToNull(ctx, data, priorPassword)
		return
	}

	// --- DepEnrollmentBaseProfile fields ---
	data.ConfigurationWebUrl = convert.GraphToFrameworkBool(macProfile.GetConfigurationWebUrl())
	data.DeviceNameTemplate = convert.GraphToFrameworkStringNullIfEmpty(macProfile.GetDeviceNameTemplate())

	// --- Setup Assistant screen control booleans (base profile) ---
	data.AppleIdDisabled = convert.GraphToFrameworkBool(macProfile.GetAppleIdDisabled())
	data.ApplePayDisabled = convert.GraphToFrameworkBool(macProfile.GetApplePayDisabled())
	data.DiagnosticsDisabled = convert.GraphToFrameworkBool(macProfile.GetDiagnosticsDisabled())
	data.DisplayToneSetupDisabled = convert.GraphToFrameworkBool(macProfile.GetDisplayToneSetupDisabled())
	data.LocationDisabled = convert.GraphToFrameworkBool(macProfile.GetLocationDisabled())
	data.PrivacyPaneDisabled = convert.GraphToFrameworkBool(macProfile.GetPrivacyPaneDisabled())
	data.RestoreBlocked = convert.GraphToFrameworkBool(macProfile.GetRestoreBlocked())
	data.ScreenTimeScreenDisabled = convert.GraphToFrameworkBool(macProfile.GetScreenTimeScreenDisabled())
	data.SiriDisabled = convert.GraphToFrameworkBool(macProfile.GetSiriDisabled())
	data.TermsAndConditionsDisabled = convert.GraphToFrameworkBool(macProfile.GetTermsAndConditionsDisabled())
	data.TouchIdDisabled = convert.GraphToFrameworkBool(macProfile.GetTouchIdDisabled())

	// Map enrollment time Azure AD group IDs (UUID collection → string set)
	groupIds := macProfile.GetEnrollmentTimeAzureAdGroupIds()
	if groupIds != nil {
		strs := make([]string, len(groupIds))
		for i, g := range groupIds {
			strs[i] = g.String()
		}
		data.EnrollmentTimeAzureAdGroupIds = convert.GraphToFrameworkStringSet(ctx, strs)
	} else {
		data.EnrollmentTimeAzureAdGroupIds = types.SetNull(types.StringType)
	}
	data.IsDefault = convert.GraphToFrameworkBool(macProfile.GetIsDefault())
	data.IsMandatory = convert.GraphToFrameworkBool(macProfile.GetIsMandatory())
	data.ProfileRemovalDisabled = convert.GraphToFrameworkBool(macProfile.GetProfileRemovalDisabled())
	data.SupervisedModeEnabled = convert.GraphToFrameworkBool(macProfile.GetSupervisedModeEnabled())
	data.SupportDepartment = convert.GraphToFrameworkStringNullIfEmpty(macProfile.GetSupportDepartment())
	data.SupportPhoneNumber = convert.GraphToFrameworkStringNullIfEmpty(macProfile.GetSupportPhoneNumber())
	data.WaitForDeviceConfiguredConfirmation = convert.GraphToFrameworkBool(macProfile.GetWaitForDeviceConfiguredConfirmation())

	// --- Setup Assistant screen control booleans (macOS-specific) ---
	// WelcomeScreenDisabled is skip-key-only (no API boolean property).
	// The Microsoft Graph API does not provide a boolean property for this screen.
	// We derive this from the "Welcome" key in the enabledSkipKeys array returned by the API.
	data.WelcomeScreenDisabled = types.BoolValue(false)
	if baseProfile, ok := macProfile.(interface{ GetEnabledSkipKeys() []string }); ok {
		data.WelcomeScreenDisabled = types.BoolValue(slices.Contains(baseProfile.GetEnabledSkipKeys(), "Welcome"))
	}
	data.AccessibilityScreenDisabled = convert.GraphToFrameworkBool(macProfile.GetAccessibilityScreenDisabled())
	data.AutoUnlockWithWatchDisabled = convert.GraphToFrameworkBool(macProfile.GetAutoUnlockWithWatchDisabled())
	data.ChooseYourLockScreenDisabled = convert.GraphToFrameworkBool(macProfile.GetChooseYourLockScreenDisabled())
	data.FileVaultDisabled = convert.GraphToFrameworkBool(macProfile.GetFileVaultDisabled())
	data.ICloudDiagnosticsDisabled = convert.GraphToFrameworkBool(macProfile.GetICloudDiagnosticsDisabled())
	data.ICloudStorageDisabled = convert.GraphToFrameworkBool(macProfile.GetICloudStorageDisabled())
	data.PassCodeDisabled = convert.GraphToFrameworkBool(macProfile.GetPassCodeDisabled())
	data.RegistrationDisabled = convert.GraphToFrameworkBool(macProfile.GetRegistrationDisabled())
	data.ZoomDisabled = convert.GraphToFrameworkBool(macProfile.GetZoomDisabled())

	// --- DepMacOSEnrollmentProfile: Account and enrollment behavior fields ---
	// Admin account fields: map API nil to empty string ("") so that the
	// explicit-clear round-trip works: user sets "" → construct sends nil →
	// API returns nil → state is "". UseStateForUnknownString preserves ""
	// on subsequent plans when the field is omitted, avoiding perpetual diffs.
	data.AdminAccountFullName = convert.GraphToFrameworkStringWithDefault(macProfile.GetAdminAccountFullName(), "")
	data.AdminAccountUserName = convert.GraphToFrameworkStringWithDefault(macProfile.GetAdminAccountUserName(), "")

	// Preserve admin_account_password from prior state since the API does not return it (write-only field).
	if v := macProfile.GetAdminAccountPassword(); v != nil {
		data.AdminAccountPassword = convert.GraphToFrameworkString(v)
	} else {
		data.AdminAccountPassword = priorPassword
	}
	data.AutoAdvanceSetupEnabled = convert.GraphToFrameworkBool(macProfile.GetAutoAdvanceSetupEnabled())
	data.DontAutoPopulatePrimaryAccountInfo = convert.GraphToFrameworkBool(macProfile.GetDontAutoPopulatePrimaryAccountInfo())
	data.EnableRestrictEditing = convert.GraphToFrameworkBool(macProfile.GetEnableRestrictEditing())
	data.HideAdminAccount = convert.GraphToFrameworkBool(macProfile.GetHideAdminAccount())
	data.PrimaryAccountFullName = convert.GraphToFrameworkStringNullIfEmpty(macProfile.GetPrimaryAccountFullName())
	data.PrimaryAccountUserName = convert.GraphToFrameworkStringNullIfEmpty(macProfile.GetPrimaryAccountUserName())
	data.RequestRequiresNetworkTether = convert.GraphToFrameworkBool(macProfile.GetRequestRequiresNetworkTether())
	data.SetPrimarySetupAccountAsRegularUser = convert.GraphToFrameworkBool(macProfile.GetSetPrimarySetupAccountAsRegularUser())
	data.SkipPrimarySetupAccountCreation = convert.GraphToFrameworkBool(macProfile.GetSkipPrimarySetupAccountCreation())

	// --- LAPS rotation setting nested object ---
	rotationSetting := macProfile.GetDepProfileAdminAccountPasswordRotationSetting()
	if rotationSetting != nil {
		// Map the delay sub-object
		var delayObj types.Object
		delaySetting := rotationSetting.GetDepProfileDelayAutoRotationSetting()
		if delaySetting != nil {
			delayValues := map[string]attr.Value{
				"on_retrieval_auto_rotate_password_enabled":        convert.GraphToFrameworkBool(delaySetting.GetOnRetrievalAutoRotatePasswordEnabled()),
				"on_retrieval_delay_auto_rotate_password_in_hours": convert.GraphToFrameworkInt32(delaySetting.GetOnRetrievalDelayAutoRotatePasswordInHours()),
			}
			var diags diag.Diagnostics
			delayObj, diags = types.ObjectValue(delayRotationSettingAttrTypes, delayValues)
			if diags.HasError() {
				delayObj = types.ObjectNull(delayRotationSettingAttrTypes)
			}
		} else {
			delayObj = types.ObjectNull(delayRotationSettingAttrTypes)
		}

		rotationValues := map[string]attr.Value{
			"auto_rotation_period_in_days":            convert.GraphToFrameworkInt32(rotationSetting.GetAutoRotationPeriodInDays()),
			"dep_profile_delay_auto_rotation_setting": delayObj,
		}
		rotationObj, diags := types.ObjectValue(rotationSettingAttrTypes, rotationValues)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create rotation setting object", map[string]any{
				"error": diags.Errors()[0].Detail(),
			})
			data.DepProfileAdminAccountPasswordRotationSetting = types.ObjectNull(rotationSettingAttrTypes)
		} else {
			data.DepProfileAdminAccountPasswordRotationSetting = rotationObj
		}
	} else {
		data.DepProfileAdminAccountPasswordRotationSetting = types.ObjectNull(rotationSettingAttrTypes)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// setExtendedFieldsToNull sets all DepEnrollmentBaseProfile and DepMacOSEnrollmentProfile fields
// to their null type values. This is the fallback when the API response doesn't contain the
// DepMacOSEnrollmentProfile discriminator (e.g. base EnrollmentProfile only).
func setExtendedFieldsToNull(ctx context.Context, data *DepMacOSEnrollmentProfileResourceModel, priorPassword types.String) {
	tflog.Debug(ctx, "Setting extended macOS enrollment profile fields to null")

	// DepEnrollmentBaseProfile fields
	data.ConfigurationWebUrl = types.BoolNull()
	data.DeviceNameTemplate = types.StringNull()
	data.EnrollmentTimeAzureAdGroupIds = types.SetNull(types.StringType)
	data.IsDefault = types.BoolNull()
	data.IsMandatory = types.BoolNull()
	data.ProfileRemovalDisabled = types.BoolNull()
	data.SupervisedModeEnabled = types.BoolNull()
	data.SupportDepartment = types.StringNull()
	data.SupportPhoneNumber = types.StringNull()
	data.WaitForDeviceConfiguredConfirmation = types.BoolNull()

	// Setup Assistant screen control booleans (base profile)
	data.AppleIdDisabled = types.BoolNull()
	data.ApplePayDisabled = types.BoolNull()
	data.DiagnosticsDisabled = types.BoolNull()
	data.DisplayToneSetupDisabled = types.BoolNull()
	data.LocationDisabled = types.BoolNull()
	data.PrivacyPaneDisabled = types.BoolNull()
	data.RestoreBlocked = types.BoolNull()
	data.ScreenTimeScreenDisabled = types.BoolNull()
	data.SiriDisabled = types.BoolNull()
	data.TermsAndConditionsDisabled = types.BoolNull()
	data.TouchIdDisabled = types.BoolNull()

	// Setup Assistant screen control booleans (macOS-specific)
	data.WelcomeScreenDisabled = types.BoolNull()
	data.AccessibilityScreenDisabled = types.BoolNull()
	data.AutoUnlockWithWatchDisabled = types.BoolNull()
	data.ChooseYourLockScreenDisabled = types.BoolNull()
	data.FileVaultDisabled = types.BoolNull()
	data.ICloudDiagnosticsDisabled = types.BoolNull()
	data.ICloudStorageDisabled = types.BoolNull()
	data.PassCodeDisabled = types.BoolNull()
	data.RegistrationDisabled = types.BoolNull()
	data.ZoomDisabled = types.BoolNull()

	// DepMacOSEnrollmentProfile: Account and enrollment behavior fields
	data.AdminAccountFullName = types.StringNull()
	data.AdminAccountPassword = priorPassword // Explicitly preserve from prior state
	data.AdminAccountUserName = types.StringNull()
	data.AutoAdvanceSetupEnabled = types.BoolNull()
	data.DepProfileAdminAccountPasswordRotationSetting = types.ObjectNull(rotationSettingAttrTypes)
	data.DontAutoPopulatePrimaryAccountInfo = types.BoolNull()
	data.EnableRestrictEditing = types.BoolNull()
	data.HideAdminAccount = types.BoolNull()
	data.PrimaryAccountFullName = types.StringNull()
	data.PrimaryAccountUserName = types.StringNull()
	data.RequestRequiresNetworkTether = types.BoolNull()
	data.SetPrimarySetupAccountAsRegularUser = types.BoolNull()
	data.SkipPrimarySetupAccountCreation = types.BoolNull()
}
