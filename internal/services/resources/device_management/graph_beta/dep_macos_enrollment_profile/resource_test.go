package graphBetaDepMacOSEnrollmentProfile_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaDepMacOSEnrollmentProfile "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/dep_macos_enrollment_profile"
	_ "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/dep_macos_enrollment_profile/mocks" // Import to register mocks
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	resourceType = graphBetaDepMacOSEnrollmentProfile.ResourceName
)

func setupMockEnvironment() *mocks.Mocks {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	return mockClient
}

// TestUnitResourceDepMacOSEnrollmentProfile_01_Basic tests creating a DEP macOS enrollment profile with basic configuration
func TestUnitResourceDepMacOSEnrollmentProfile_01_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_ = setupMockEnvironment()
	mocks.GlobalRegistry.ActivateMocks("dep_macos_enrollment_profile")
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_basic").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_basic").Key("display_name").HasValue("Test Basic Enrollment Profile"),
					check.That(resourceType+".test_basic").Key("requires_user_authentication").HasValue("true"),
					check.That(resourceType+".test_basic").Key("supervised_mode_enabled").HasValue("true"),
					check.That(resourceType+".test_basic").Key("admin_account_user_name").HasValue("localadmin"),
					check.That(resourceType+".test_basic").Key("configuration_endpoint_url").Exists(),
				),
			},
			{
				ResourceName:            resourceType + ".test_basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_account_password"}, // Password is write-only
			},
		},
	})
}

// TestUnitResourceDepMacOSEnrollmentProfile_02_ZeroTouchBooleans tests zero-touch deployment with boolean screen control fields
func TestUnitResourceDepMacOSEnrollmentProfile_02_ZeroTouchBooleans(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_ = setupMockEnvironment()
	mocks.GlobalRegistry.ActivateMocks("dep_macos_enrollment_profile")
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigZeroTouchBooleans(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_zero_touch").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_zero_touch").Key("display_name").HasValue("Zero-Touch Deployment"),
					check.That(resourceType+".test_zero_touch").Key("auto_advance_setup_enabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("skip_primary_setup_account_creation").HasValue("true"),

					// Verify all 20 boolean screen control fields
					// Base profile screens (11)
					check.That(resourceType+".test_zero_touch").Key("apple_id_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("apple_pay_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("diagnostics_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("display_tone_setup_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("location_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("privacy_pane_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("restore_blocked").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("screen_time_screen_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("siri_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("terms_and_conditions_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("touch_id_disabled").HasValue("true"),

					// macOS-specific screens (10, including skip-key-only welcome_screen_disabled)
					check.That(resourceType+".test_zero_touch").Key("welcome_screen_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("accessibility_screen_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("auto_unlock_with_watch_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("choose_your_lock_screen_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("file_vault_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("i_cloud_diagnostics_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("i_cloud_storage_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("pass_code_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("registration_disabled").HasValue("true"),
					check.That(resourceType+".test_zero_touch").Key("zoom_disabled").HasValue("true"),

					check.That(resourceType+".test_zero_touch").Key("enrollment_time_azure_ad_group_ids.#").HasValue("2"),
				),
			},
		},
	})
}

// TestUnitResourceDepMacOSEnrollmentProfile_03_PrivacyRegistrationBooleans tests Privacy and Registration screens
// which can ONLY be controlled via boolean fields (API rejects them as skip key strings)
func TestUnitResourceDepMacOSEnrollmentProfile_03_PrivacyRegistrationBooleans(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_ = setupMockEnvironment()
	mocks.GlobalRegistry.ActivateMocks("dep_macos_enrollment_profile")
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigPrivacyRegistration(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_privacy_registration").Key("id").Exists(),
					// These are the special boolean-only fields that work despite API rejecting skip key strings
					check.That(resourceType+".test_privacy_registration").Key("privacy_pane_disabled").HasValue("true"),
					check.That(resourceType+".test_privacy_registration").Key("registration_disabled").HasValue("true"),
				),
			},
		},
	})
}

// TestUnitResourceDepMacOSEnrollmentProfile_04_LAPS tests LAPS password rotation configuration
func TestUnitResourceDepMacOSEnrollmentProfile_04_LAPS(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_ = setupMockEnvironment()
	mocks.GlobalRegistry.ActivateMocks("dep_macos_enrollment_profile")
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigLAPS(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_laps").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_laps").Key("admin_account_user_name").HasValue("localadmin"),
					check.That(resourceType+".test_laps").Key("admin_account_full_name").HasValue("Local Admin"),
					check.That(resourceType+".test_laps").Key("hide_admin_account").HasValue("true"),
					check.That(resourceType+".test_laps").Key("dep_profile_admin_account_password_rotation_setting.auto_rotation_period_in_days").HasValue("30"),
					check.That(resourceType+".test_laps").Key("dep_profile_admin_account_password_rotation_setting.dep_profile_delay_auto_rotation_setting.on_retrieval_auto_rotate_password_enabled").HasValue("true"),
					check.That(resourceType+".test_laps").Key("dep_profile_admin_account_password_rotation_setting.dep_profile_delay_auto_rotation_setting.on_retrieval_delay_auto_rotate_password_in_hours").HasValue("24"),
				),
			},
		},
	})
}

// TestUnitResourceDepMacOSEnrollmentProfile_05_AdminAccountClearing tests clearing admin account fields to disable LAPS
func TestUnitResourceDepMacOSEnrollmentProfile_05_AdminAccountClearing(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_ = setupMockEnvironment()
	mocks.GlobalRegistry.ActivateMocks("dep_macos_enrollment_profile")
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create profile with admin account (LAPS active)
				Config: testConfigWithAdminAccount(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_admin_clearing").Key("admin_account_user_name").HasValue("localadmin"),
					check.That(resourceType+".test_admin_clearing").Key("admin_account_full_name").HasValue("Local Admin"),
					check.That(resourceType+".test_admin_clearing").Key("hide_admin_account").HasValue("true"),
				),
			},
			{
				// Step 2: Remove admin account fields from config (should clear them in API, disabling LAPS)
				Config: testConfigWithoutAdminAccount(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_admin_clearing").Key("id").Exists(),
					// Admin account fields should be empty after removal from config (Computed fields remain in state)
					check.That(resourceType+".test_admin_clearing").Key("admin_account_user_name").HasValue(""),
					check.That(resourceType+".test_admin_clearing").Key("admin_account_full_name").HasValue(""),
					check.That(resourceType+".test_admin_clearing").Key("hide_admin_account").HasValue("false"),
				),
			},
		},
	})
}

// TestUnitResourceDepMacOSEnrollmentProfile_06_Update tests updating profile configuration with boolean screen controls
func TestUnitResourceDepMacOSEnrollmentProfile_06_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_ = setupMockEnvironment()
	mocks.GlobalRegistry.ActivateMocks("dep_macos_enrollment_profile")
	defer httpmock.DeactivateAndReset()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Start with no screen control
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_basic").Key("apple_id_disabled").HasValue("false"),
					check.That(resourceType+".test_basic").Key("siri_disabled").HasValue("false"),
				),
			},
			{
				// Add screen control booleans
				Config: testConfigBasicWithBooleans(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_basic").Key("apple_id_disabled").HasValue("true"),
					check.That(resourceType+".test_basic").Key("siri_disabled").HasValue("true"),
					check.That(resourceType+".test_basic").Key("diagnostics_disabled").HasValue("true"),
					check.That(resourceType+".test_basic").Key("location_disabled").HasValue("true"),
					check.That(resourceType+".test_basic").Key("privacy_pane_disabled").HasValue("true"),
				),
			},
		},
	})
}

// Test configuration functions
func testConfigBasic() string {
	return `resource "microsoft365_graph_beta_device_management_dep_macos_enrollment_profile" "test_basic" {
  dep_onboarding_settings_id                              = "11111111-1111-1111-1111-111111111111"
  display_name                                            = "Test Basic Enrollment Profile"
  requires_user_authentication                            = true
  enable_authentication_via_company_portal                = true
  require_company_portal_on_setup_assistant_enrolled_devices = false

  supervised_mode_enabled  = true
  is_mandatory            = true

  admin_account_user_name = "localadmin"
  admin_account_full_name = "Local Admin"
  admin_account_password  = "SecureP@ss123"
  hide_admin_account      = true
}`
}

func testConfigBasicWithBooleans() string {
	return `resource "microsoft365_graph_beta_device_management_dep_macos_enrollment_profile" "test_basic" {
  dep_onboarding_settings_id                              = "11111111-1111-1111-1111-111111111111"
  display_name                                            = "Test Basic Enrollment Profile"
  requires_user_authentication                            = true
  enable_authentication_via_company_portal                = true
  require_company_portal_on_setup_assistant_enrolled_devices = false

  supervised_mode_enabled  = true
  is_mandatory            = true

  # Boolean screen controls instead of enabled_skip_keys
  apple_id_disabled       = true
  siri_disabled          = true
  diagnostics_disabled   = true
  location_disabled      = true
  privacy_pane_disabled  = true

  admin_account_user_name = "localadmin"
  admin_account_full_name = "Local Admin"
  admin_account_password  = "SecureP@ss123"
  hide_admin_account      = true
}`
}

func testConfigZeroTouchBooleans() string {
	return `resource "microsoft365_graph_beta_device_management_dep_macos_enrollment_profile" "test_zero_touch" {
  dep_onboarding_settings_id                              = "11111111-1111-1111-1111-111111111111"
  display_name                                            = "Zero-Touch Deployment"
  requires_user_authentication                            = false
  enable_authentication_via_company_portal                = false
  require_company_portal_on_setup_assistant_enrolled_devices = false

  supervised_mode_enabled      = true
  is_mandatory                = true
  auto_advance_setup_enabled  = true

  # All 20 boolean screen control fields for complete zero-touch
  # Base profile screens (11)
  apple_id_disabled              = true
  apple_pay_disabled             = true
  diagnostics_disabled           = true
  display_tone_setup_disabled    = true
  location_disabled              = true
  privacy_pane_disabled          = true  # Boolean-only field (API rejects "Privacy" skip key)
  restore_blocked                = true
  screen_time_screen_disabled    = true
  siri_disabled                  = true
  terms_and_conditions_disabled  = true
  touch_id_disabled              = true

  # macOS-specific screens (10, including skip-key-only welcome_screen_disabled)
  welcome_screen_disabled          = true  # Skip-key-only field (no API boolean property)
  accessibility_screen_disabled    = true
  auto_unlock_with_watch_disabled  = true
  choose_your_lock_screen_disabled = true
  file_vault_disabled              = true
  i_cloud_diagnostics_disabled     = true
  i_cloud_storage_disabled         = true
  pass_code_disabled               = true
  registration_disabled            = true  # Boolean-only field (API rejects "Registration" skip key)
  zoom_disabled                    = true

  enrollment_time_azure_ad_group_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222"
  ]

  skip_primary_setup_account_creation = true
}`
}

func testConfigPrivacyRegistration() string {
	return `resource "microsoft365_graph_beta_device_management_dep_macos_enrollment_profile" "test_privacy_registration" {
  dep_onboarding_settings_id                              = "11111111-1111-1111-1111-111111111111"
  display_name                                            = "Privacy/Registration Test"
  requires_user_authentication                            = false
  enable_authentication_via_company_portal                = false
  require_company_portal_on_setup_assistant_enrolled_devices = false

  supervised_mode_enabled = true

  # These two fields can ONLY be controlled via booleans
  # The Graph API rejects "Privacy" and "Registration" as skip key strings
  privacy_pane_disabled = true
  registration_disabled = true
}`
}

func testConfigLAPS() string {
	return `resource "microsoft365_graph_beta_device_management_dep_macos_enrollment_profile" "test_laps" {
  dep_onboarding_settings_id                              = "11111111-1111-1111-1111-111111111111"
  display_name                                            = "LAPS Test Profile"
  requires_user_authentication                            = true
  enable_authentication_via_company_portal                = true
  require_company_portal_on_setup_assistant_enrolled_devices = false

  supervised_mode_enabled = true

  admin_account_user_name = "localadmin"
  admin_account_full_name = "Local Admin"
  admin_account_password  = "InitialP@ss123"
  hide_admin_account      = true

  dep_profile_admin_account_password_rotation_setting = {
    auto_rotation_period_in_days = 30

    dep_profile_delay_auto_rotation_setting = {
      on_retrieval_auto_rotate_password_enabled        = true
      on_retrieval_delay_auto_rotate_password_in_hours = 24
    }
  }
}`
}

func testConfigWithAdminAccount() string {
	return `resource "microsoft365_graph_beta_device_management_dep_macos_enrollment_profile" "test_admin_clearing" {
  dep_onboarding_settings_id                              = "11111111-1111-1111-1111-111111111111"
  display_name                                            = "Admin Account Clearing Test"
  requires_user_authentication                            = false
  enable_authentication_via_company_portal                = false
  require_company_portal_on_setup_assistant_enrolled_devices = false

  supervised_mode_enabled = true

  # Admin account configured (LAPS will be active)
  admin_account_user_name = "localadmin"
  admin_account_full_name = "Local Admin"
  admin_account_password  = "InitialP@ss123"
  hide_admin_account      = true
}`
}

func testConfigWithoutAdminAccount() string {
	return `resource "microsoft365_graph_beta_device_management_dep_macos_enrollment_profile" "test_admin_clearing" {
  dep_onboarding_settings_id                              = "11111111-1111-1111-1111-111111111111"
  display_name                                            = "Admin Account Clearing Test"
  requires_user_authentication                            = false
  enable_authentication_via_company_portal                = false
  require_company_portal_on_setup_assistant_enrolled_devices = false

  supervised_mode_enabled = true

  # Explicitly set admin account fields to empty strings to clear them in API and disable LAPS
  # With Computed: true, omitting fields preserves state - must explicitly set to "" to clear
  admin_account_user_name = ""
  admin_account_full_name = ""
  hide_admin_account      = false
}`
}
