// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-depmacosenrollmentprofile?view=graph-rest-beta
package graphBetaDepMacOSEnrollmentProfile

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DepMacOSEnrollmentProfileResourceModel models a DEP (ABM/ASM) macOS enrollment profile
// under a DEP onboarding setting.
// Endpoint: POST /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles
// SDK type: DepMacOSEnrollmentProfile (extends DepEnrollmentBaseProfile extends EnrollmentProfile)
type DepMacOSEnrollmentProfileResourceModel struct {
	// --- EnrollmentProfile base fields ---
	ID                                                  types.String `tfsdk:"id"`
	DisplayName                                         types.String `tfsdk:"display_name"`
	Description                                         types.String `tfsdk:"description"`
	RequiresUserAuthentication                          types.Bool   `tfsdk:"requires_user_authentication"`
	EnableAuthenticationViaCompanyPortal                types.Bool   `tfsdk:"enable_authentication_via_company_portal"`
	RequireCompanyPortalOnSetupAssistantEnrolledDevices types.Bool   `tfsdk:"require_company_portal_on_setup_assistant_enrolled_devices"`
	ConfigurationEndpointUrl                            types.String `tfsdk:"configuration_endpoint_url"`
	DepOnboardingSettingsId                             types.String `tfsdk:"dep_onboarding_settings_id"`

	// --- DepEnrollmentBaseProfile fields ---
	ConfigurationWebUrl                 types.Bool   `tfsdk:"configuration_web_url"`
	DeviceNameTemplate                  types.String `tfsdk:"device_name_template"`
	EnrollmentTimeAzureAdGroupIds       types.Set    `tfsdk:"enrollment_time_azure_ad_group_ids"`
	IsDefault                           types.Bool   `tfsdk:"is_default"`
	IsMandatory                         types.Bool   `tfsdk:"is_mandatory"`
	ProfileRemovalDisabled              types.Bool   `tfsdk:"profile_removal_disabled"`
	SupervisedModeEnabled               types.Bool   `tfsdk:"supervised_mode_enabled"`
	SupportDepartment                   types.String `tfsdk:"support_department"`
	SupportPhoneNumber                  types.String `tfsdk:"support_phone_number"`
	WaitForDeviceConfiguredConfirmation types.Bool   `tfsdk:"wait_for_device_configured_confirmation"`

	// --- DepEnrollmentBaseProfile: Setup Assistant screen control booleans ---
	// These control which Setup Assistant screens to skip during enrollment.
	// The provider auto-generates the enabledSkipKeys array from these booleans
	// when constructing API requests. Users should ONLY set these boolean fields;
	// enabledSkipKeys is managed internally by the provider.
	//
	// Note: privacyPaneDisabled and registrationDisabled are exceptions - they
	// control screens that cannot be skipped via the enabledSkipKeys array due
	// to Microsoft Graph API limitations. The provider handles this automatically.
	AppleIdDisabled            types.Bool `tfsdk:"apple_id_disabled"`
	ApplePayDisabled           types.Bool `tfsdk:"apple_pay_disabled"`
	DiagnosticsDisabled        types.Bool `tfsdk:"diagnostics_disabled"`
	DisplayToneSetupDisabled   types.Bool `tfsdk:"display_tone_setup_disabled"`
	LocationDisabled           types.Bool `tfsdk:"location_disabled"`
	PrivacyPaneDisabled        types.Bool `tfsdk:"privacy_pane_disabled"`
	RestoreBlocked             types.Bool `tfsdk:"restore_blocked"`
	ScreenTimeScreenDisabled   types.Bool `tfsdk:"screen_time_screen_disabled"`
	SiriDisabled               types.Bool `tfsdk:"siri_disabled"`
	TermsAndConditionsDisabled types.Bool `tfsdk:"terms_and_conditions_disabled"`
	TouchIdDisabled            types.Bool `tfsdk:"touch_id_disabled"`

	// --- DepMacOSEnrollmentProfile: macOS-specific Setup Assistant screen control booleans ---
	WelcomeScreenDisabled        types.Bool `tfsdk:"welcome_screen_disabled"`
	AccessibilityScreenDisabled  types.Bool `tfsdk:"accessibility_screen_disabled"`
	AutoUnlockWithWatchDisabled  types.Bool `tfsdk:"auto_unlock_with_watch_disabled"`
	ChooseYourLockScreenDisabled types.Bool `tfsdk:"choose_your_lock_screen_disabled"`
	FileVaultDisabled            types.Bool `tfsdk:"file_vault_disabled"`
	ICloudDiagnosticsDisabled    types.Bool `tfsdk:"i_cloud_diagnostics_disabled"`
	ICloudStorageDisabled        types.Bool `tfsdk:"i_cloud_storage_disabled"`
	PassCodeDisabled             types.Bool `tfsdk:"pass_code_disabled"`
	RegistrationDisabled         types.Bool `tfsdk:"registration_disabled"`
	ZoomDisabled                 types.Bool `tfsdk:"zoom_disabled"`

	// --- DepMacOSEnrollmentProfile: Account and enrollment behavior fields ---
	AdminAccountFullName                          types.String `tfsdk:"admin_account_full_name"`
	AdminAccountPassword                          types.String `tfsdk:"admin_account_password"`
	AdminAccountUserName                          types.String `tfsdk:"admin_account_user_name"`
	AutoAdvanceSetupEnabled                       types.Bool   `tfsdk:"auto_advance_setup_enabled"`
	DepProfileAdminAccountPasswordRotationSetting types.Object `tfsdk:"dep_profile_admin_account_password_rotation_setting"`
	DontAutoPopulatePrimaryAccountInfo            types.Bool   `tfsdk:"dont_auto_populate_primary_account_info"`
	EnableRestrictEditing                         types.Bool   `tfsdk:"enable_restrict_editing"`
	HideAdminAccount                              types.Bool   `tfsdk:"hide_admin_account"`
	PrimaryAccountFullName                        types.String `tfsdk:"primary_account_full_name"`
	PrimaryAccountUserName                        types.String `tfsdk:"primary_account_user_name"`
	RequestRequiresNetworkTether                  types.Bool   `tfsdk:"request_requires_network_tether"`
	SetPrimarySetupAccountAsRegularUser           types.Bool   `tfsdk:"set_primary_setup_account_as_regular_user"`
	SkipPrimarySetupAccountCreation               types.Bool   `tfsdk:"skip_primary_setup_account_creation"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

// DepProfileAdminAccountPasswordRotationSettingModel models the LAPS rotation setting nested object.
type DepProfileAdminAccountPasswordRotationSettingModel struct {
	AutoRotationPeriodInDays           types.Int32  `tfsdk:"auto_rotation_period_in_days"`
	DepProfileDelayAutoRotationSetting types.Object `tfsdk:"dep_profile_delay_auto_rotation_setting"`
}

// DepProfileDelayAutoRotationSettingModel models the LAPS delay auto-rotation sub-object.
type DepProfileDelayAutoRotationSettingModel struct {
	OnRetrievalAutoRotatePasswordEnabled      types.Bool  `tfsdk:"on_retrieval_auto_rotate_password_enabled"`
	OnRetrievalDelayAutoRotatePasswordInHours types.Int32 `tfsdk:"on_retrieval_delay_auto_rotate_password_in_hours"`
}
