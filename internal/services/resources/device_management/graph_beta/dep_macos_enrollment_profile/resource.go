package graphBetaDepMacOSEnrollmentProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	validate "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_dep_macos_enrollment_profile"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &DepMacOSEnrollmentProfileResource{}
	_ resource.ResourceWithConfigure   = &DepMacOSEnrollmentProfileResource{}
	_ resource.ResourceWithImportState = &DepMacOSEnrollmentProfileResource{}
)

func NewDepMacOSEnrollmentProfileResource() resource.Resource {
	return &DepMacOSEnrollmentProfileResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
		ResourcePath: "deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles",
	}
}

type DepMacOSEnrollmentProfileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *DepMacOSEnrollmentProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *DepMacOSEnrollmentProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *DepMacOSEnrollmentProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DepMacOSEnrollmentProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages macOS DEP enrollment profiles using the `/deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles` endpoint. " +
			"This resource configures zero-touch macOS deployment with full Setup Assistant control via the `depMacOSEnrollmentProfile` Graph API type. " +
			"See [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-depmacosenrollmentprofile?view=graph-rest-beta) for details.",
		Attributes: map[string]schema.Attribute{
			// --- EnrollmentProfile base fields ---
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier of the enrollment profile.",
			},
			"dep_onboarding_settings_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Identifier of the DEP onboarding setting (ABM/ASM token) to create this profile under. " +
					"If not specified, falls back to the intuneAccountId from /deviceManagement. " +
					"Use `dep_onboarding_settings_id` when your tenant has multiple DEP tokens.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the macOS DEP enrollment profile displayed in Intune.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"requires_user_authentication": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Indicates whether the user must authenticate during Apple device setup.",
			},
			"enable_authentication_via_company_portal": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "How users first sign in to authenticate with Intune. If your organization uses multi-factor authentication, set this to true; the app will then automatically install on devices at time of enrollment.",
				Validators: []validator.Bool{
					validate.MutuallyExclusiveBool("require_company_portal_on_setup_assistant_enrolled_devices", "enable_authentication_via_company_portal and require_company_portal_on_setup_assistant_enrolled_devices cannot both be set to true"),
				},
			},
			"require_company_portal_on_setup_assistant_enrolled_devices": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "How users first sign in to authenticate with Intune. Setup assistant as a standalone authentication method has been superseded by setup assistant with modern authentication.",
				Validators: []validator.Bool{
					validate.MutuallyExclusiveBool("enable_authentication_via_company_portal", "enable_authentication_via_company_portal and require_company_portal_on_setup_assistant_enrolled_devices cannot both be set to true"),
				},
			},
			"configuration_endpoint_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Apple Configurator enrollment configuration endpoint URL generated by Intune.",
			},

			// --- DepEnrollmentBaseProfile fields ---
			"configuration_web_url": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Enables the setup assistant login web URL. **Note:** Despite the field name suggesting a URL string, " +
					"the Graph Beta API defines this as a boolean on `depEnrollmentBaseProfile`.",
				PlanModifiers: []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"device_name_template": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Sets a literal or name pattern for the device name. Supports variables like %SERIAL% (device serial number) and %DEVICETYPE% (device model/type). Examples: 'CORP-MAC-%SERIAL%', '%DEVICETYPE%-%SERIAL%'. Limited to ~63 characters and only applied to supervised devices.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},

			// --- Setup Assistant Screen Control Booleans ---
			// Design: This provider exposes individual boolean fields for each Setup Assistant
			// screen instead of the enabledSkipKeys array. The provider auto-generates the
			// enabledSkipKeys array internally when constructing API requests.
			//
			// Rationale: Microsoft Graph API has an inconsistency where most screens can be
			// controlled via both enabledSkipKeys AND boolean properties, but Privacy and
			// Registration screens ONLY work via booleans (API rejects "Privacy" and
			// "Registration" skip key strings). Using booleans exclusively provides:
			// 1. Complete coverage - all screens can be controlled
			// 2. Consistent UX - same mechanism for all screens
			// 3. Explicit configuration - no guessing about valid skip key strings
			// 4. Automatic exception handling - Privacy/Registration handled correctly

			// Base profile Setup Assistant screens (from DepEnrollmentBaseProfile)
			"apple_id_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Apple ID setup screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"apple_pay_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Apple Pay setup screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"diagnostics_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Diagnostics & Usage screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"display_tone_setup_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Display Tone setup screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"location_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Location Services screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"privacy_pane_disabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "When true, skips the Privacy consent screen during enrollment. " +
					"**Note:** This screen can ONLY be controlled via this boolean field. The Microsoft Graph API " +
					"rejects 'Privacy' as a skip key string, even though the boolean property works correctly. " +
					"The provider handles this limitation automatically.",
				PlanModifiers: []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"restore_blocked": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, blocks the Restore from iCloud/Mac screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"screen_time_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Screen Time setup screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"siri_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Siri setup screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"terms_and_conditions_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Terms and Conditions screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"touch_id_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Touch ID setup screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},

			// macOS-specific Setup Assistant screens (from DepMacOSEnrollmentProfile)
			"welcome_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Welcome (\"Hello\" multilingual) screen during enrollment. Requires macOS Sonoma or later. **Note:** This screen can ONLY be controlled via the \"Welcome\" skip key string. The Microsoft Graph API does not provide a boolean property for this screen.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"accessibility_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Accessibility screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"auto_unlock_with_watch_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Auto Unlock with Apple Watch screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"choose_your_lock_screen_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Choose Your Lock Screen (Wallpaper) screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"file_vault_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the FileVault encryption screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"i_cloud_diagnostics_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the iCloud Diagnostics screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"i_cloud_storage_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the iCloud Storage screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"pass_code_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Passcode setup screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"registration_disabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "When true, skips the device registration screen during enrollment. " +
					"**Note:** This screen can ONLY be controlled via this boolean field. The Microsoft Graph API " +
					"rejects 'Registration' as a skip key string, even though the boolean property works correctly. " +
					"The provider handles this limitation automatically.",
				PlanModifiers: []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"zoom_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, skips the Zoom setup screen during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},

			"enrollment_time_azure_ad_group_ids": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "List of Azure AD group IDs to associate with the profile at enrollment time.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.UseStateForUnknownSet(),
				},
			},
			"is_default": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates if this is the default profile. This is managed server-side by Intune and cannot be set via Terraform.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.UseStateForUnknownBool()},
			},
			"is_mandatory": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the profile is mandatory.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"profile_removal_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the profile removal option is disabled.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"supervised_mode_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Supervised mode. If true, the device enters supervised mode during enrollment.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"support_department": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Support department information.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"support_phone_number": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Support phone number.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"wait_for_device_configured_confirmation": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the device will need to wait for configured confirmation.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},

			// --- DepMacOSEnrollmentProfile fields ---
			"admin_account_full_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Full name of the admin account to be created during enrollment. When removed from config, this field is cleared in the API.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"admin_account_password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				MarkdownDescription: "Password for the admin account created during enrollment. This is a write-only field; the API does not return this value. " +
					"When removed from config, this field is cleared in the API.",
			},
			"admin_account_user_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "User name of the admin account to be created during enrollment. When removed from config, this field is cleared in the API.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"auto_advance_setup_enabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if Setup Assistant will automatically advance through its screen.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"dep_profile_admin_account_password_rotation_setting": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Settings for local admin account password automatic rotation (LAPS).",
				Attributes: map[string]schema.Attribute{
					"auto_rotation_period_in_days": schema.Int32Attribute{
						Required:            true,
						MarkdownDescription: "Number of days between 1 and 180 since the last rotation after which to auto-rotate the local admin password.",
						Validators: []validator.Int32{
							int32validator.Between(1, 180),
						},
					},
					"dep_profile_delay_auto_rotation_setting": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Settings for delaying auto-rotation of the local admin password after retrieval.",
						Attributes: map[string]schema.Attribute{
							"on_retrieval_auto_rotate_password_enabled": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Indicates whether to auto-rotate the password after retrieval.",
								PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
							},
							"on_retrieval_delay_auto_rotate_password_in_hours": schema.Int32Attribute{
								Optional:            true,
								MarkdownDescription: "Number of hours to delay auto-rotation of the password after retrieval.",
							},
						},
					},
				},
			},
			"dont_auto_populate_primary_account_info": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether Setup Assistant will auto populate the primary account information.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"enable_restrict_editing": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether the user will enable restricting editing.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"hide_admin_account": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether the admin account should be hidden or not. Defaults to false when absent from config.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"primary_account_full_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Full name of the primary account to be created during enrollment.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"primary_account_user_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "User name of the primary account to be created during enrollment.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"request_requires_network_tether": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates if the device is network tethered to run the command.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"set_primary_setup_account_as_regular_user": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether Setup Assistant will set the account as a regular user.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},
			"skip_primary_setup_account_creation": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Indicates whether Setup Assistant will skip the user interface for primary account setup.",
				PlanModifiers:       []planmodifier.Bool{planmodifiers.BoolDefaultValue(false)},
			},

			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
