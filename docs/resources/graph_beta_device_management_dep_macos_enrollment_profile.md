---
page_title: "microsoft365_graph_beta_device_management_dep_macos_enrollment_profile Resource - microsoft365"
subcategory: "Device Management"

description: |-
  Manages macOS DEP enrollment profiles using the /deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles endpoint. This resource configures zero-touch macOS deployment with full Setup Assistant control via the depMacOSEnrollmentProfile Graph API type. See Microsoft Graph API documentation https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-depmacosenrollmentprofile?view=graph-rest-beta for details.
---

# microsoft365_graph_beta_device_management_dep_macos_enrollment_profile (Resource)

Manages macOS DEP enrollment profiles using the `/deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles` endpoint. This resource configures zero-touch macOS deployment with full Setup Assistant control via the `depMacOSEnrollmentProfile` Graph API type. See [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-depmacosenrollmentprofile?view=graph-rest-beta) for details.

## Microsoft Documentation

- [depMacOSEnrollmentProfile resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-depmacosenrollmentprofile?view=graph-rest-beta)
- [Create enrollmentProfile](https://learn.microsoft.com/en-us/graph/api/intune-enrollment-deponboardingsetting-post-enrollmentprofiles?view=graph-rest-beta)
- [Get depMacOSEnrollmentProfile](https://learn.microsoft.com/en-us/graph/api/intune-enrollment-enrollmentprofile-get?view=graph-rest-beta)
- [Update depMacOSEnrollmentProfile](https://learn.microsoft.com/en-us/graph/api/intune-enrollment-enrollmentprofile-update?view=graph-rest-beta)
- [Delete enrollmentProfile](https://learn.microsoft.com/en-us/graph/api/intune-enrollment-enrollmentprofile-delete?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `DeviceManagementServiceConfig.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.46.0-alpha | Experimental | Initial release |

## Example Usage

### Basic Configuration

```terraform
# Basic macOS DEP Enrollment Profile
# Creates a standard enrollment profile with required settings and admin account creation
resource "microsoft365_graph_beta_device_management_dep_macos_enrollment_profile" "basic" {
  display_name                                            = "Standard macOS Enrollment"
  description                                             = "Standard DEP enrollment profile for corporate Macs"
  requires_user_authentication                            = true
  enable_authentication_via_company_portal                = true
  require_company_portal_on_setup_assistant_enrolled_devices = false

  # Behavior settings
  supervised_mode_enabled                   = true
  is_mandatory                              = true
  wait_for_device_configured_confirmation  = true
  auto_advance_setup_enabled               = false

  # Admin account creation
  admin_account_user_name = "localadmin"
  admin_account_full_name = "Local Administrator"
  admin_account_password  = "SecureP@ssw0rd!"  # Use secure secret management in production
  hide_admin_account      = true

  # Primary user account settings
  skip_primary_setup_account_creation          = false
  set_primary_setup_account_as_regular_user   = true
  dont_auto_populate_primary_account_info     = false

  # Support information
  support_department    = "IT Support"
  support_phone_number  = "+1-555-0123"
  device_name_template  = "MAC-%SERIAL%"

  # Profile management
  profile_removal_disabled = true
  enable_restrict_editing  = false

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
```

### Zero-Touch Deployment with LAPS

```terraform
# Zero-Touch macOS DEP Enrollment Profile
# Creates an automated enrollment profile that skips Setup Assistant screens
# for a streamlined deployment experience
resource "microsoft365_graph_beta_device_management_dep_macos_enrollment_profile" "zero_touch" {
  display_name                                            = "Zero-Touch macOS Deployment"
  description                                             = "Automated DEP enrollment with minimal user interaction"
  requires_user_authentication                            = true
  enable_authentication_via_company_portal                = true
  require_company_portal_on_setup_assistant_enrolled_devices = false

  # Behavior settings for automated deployment
  supervised_mode_enabled                   = true
  is_mandatory                              = true
  wait_for_device_configured_confirmation  = true
  auto_advance_setup_enabled               = true  # Automatically advance through Setup Assistant

  # Skip all non-essential Setup Assistant screens for zero-touch experience
  # Configure individual boolean fields to control which screens are skipped
  apple_id_disabled              = true  # Skip Apple ID setup
  apple_pay_disabled             = true  # Skip Apple Pay setup
  diagnostics_disabled           = true  # Skip diagnostics & usage data
  display_tone_setup_disabled    = true  # Skip True Tone display setup
  location_disabled              = true  # Skip location services
  privacy_pane_disabled          = true  # Skip privacy information
  restore_blocked                = true  # Block restore from backup
  screen_time_screen_disabled    = true  # Skip Screen Time setup
  siri_disabled                  = true  # Skip Siri setup
  terms_and_conditions_disabled  = true  # Skip Terms and Conditions
  touch_id_disabled              = true  # Skip Touch ID setup

  # macOS-specific screens
  accessibility_screen_disabled    = true  # Skip Accessibility settings
  auto_unlock_with_watch_disabled  = true  # Skip Apple Watch unlock
  choose_your_lock_screen_disabled = true  # Skip wallpaper selection
  file_vault_disabled              = true  # Skip FileVault setup
  i_cloud_diagnostics_disabled     = true  # Skip iCloud Analytics
  i_cloud_storage_disabled         = true  # Skip iCloud Storage setup
  pass_code_disabled               = true  # Skip Passcode setup
  registration_disabled            = true  # Skip device registration
  zoom_disabled                    = true  # Skip Zoom setup

  # Admin account with LAPS (Local Administrator Password Solution)
  admin_account_user_name = "localadmin"
  admin_account_full_name = "Local Administrator"
  admin_account_password  = "Initial!P@ssw0rd123"  # Use secure secret management
  hide_admin_account      = true

  # LAPS configuration for automatic password rotation
  dep_profile_admin_account_password_rotation_setting = {
    auto_rotation_period_in_days = 30  # Rotate password every 30 days

    dep_profile_delay_auto_rotation_setting = {
      on_retrieval_auto_rotate_password_enabled        = true
      on_retrieval_delay_auto_rotate_password_in_hours = 24  # Rotate 24 hours after retrieval
    }
  }

  # Primary user account settings - skip account creation for SSO enrollment
  skip_primary_setup_account_creation          = true
  set_primary_setup_account_as_regular_user   = false
  dont_auto_populate_primary_account_info     = true

  # Enrollment-time Azure AD group assignment
  enrollment_time_azure_ad_group_ids = [
    "11111111-1111-1111-1111-111111111111",  # All macOS Devices group
    "22222222-2222-2222-2222-222222222222",  # Corporate Managed Devices group
  ]

  # Support information
  support_department    = "Enterprise IT"
  support_phone_number  = "+1-555-0100"
  device_name_template  = "CORP-MAC-%SERIAL%"

  # Security settings
  profile_removal_disabled = true
  enable_restrict_editing  = true
  request_requires_network_tether = false

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The name of the macOS DEP enrollment profile displayed in Intune.
- `enable_authentication_via_company_portal` (Boolean) How users first sign in to authenticate with Intune. If your organization uses multi-factor authentication, set this to true; the app will then automatically install on devices at time of enrollment.
- `require_company_portal_on_setup_assistant_enrolled_devices` (Boolean) How users first sign in to authenticate with Intune. Setup assistant as a standalone authentication method has been superseded by setup assistant with modern authentication.
- `requires_user_authentication` (Boolean) Indicates whether the user must authenticate during Apple device setup.

### Optional

- `accessibility_screen_disabled` (Boolean) When true, skips the Accessibility screen during enrollment.
- `admin_account_full_name` (String) Full name of the admin account to be created during enrollment. When removed from config, this field is cleared in the API.
- `admin_account_password` (String, Sensitive) Password for the admin account created during enrollment. This is a write-only field; the API does not return this value. When removed from config, this field is cleared in the API.
- `admin_account_user_name` (String) User name of the admin account to be created during enrollment. When removed from config, this field is cleared in the API.
- `apple_id_disabled` (Boolean) When true, skips the Apple ID setup screen during enrollment.
- `apple_pay_disabled` (Boolean) When true, skips the Apple Pay setup screen during enrollment.
- `auto_advance_setup_enabled` (Boolean) Indicates if Setup Assistant will automatically advance through its screen.
- `auto_unlock_with_watch_disabled` (Boolean) When true, skips the Auto Unlock with Apple Watch screen during enrollment.
- `choose_your_lock_screen_disabled` (Boolean) When true, skips the Choose Your Lock Screen (Wallpaper) screen during enrollment.
- `configuration_web_url` (Boolean) Enables the setup assistant login web URL. **Note:** Despite the field name suggesting a URL string, the Graph Beta API defines this as a boolean on `depEnrollmentBaseProfile`.
- `dep_onboarding_settings_id` (String) Identifier of the DEP onboarding setting (ABM/ASM token) to create this profile under. If not specified, falls back to the intuneAccountId from /deviceManagement. Use `dep_onboarding_settings_id` when your tenant has multiple DEP tokens.
- `dep_profile_admin_account_password_rotation_setting` (Attributes) Settings for local admin account password automatic rotation (LAPS). (see [below for nested schema](#nestedatt--dep_profile_admin_account_password_rotation_setting))
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `device_name_template` (String) Sets a literal or name pattern for the device name. Supports variables like %SERIAL% (device serial number) and %DEVICETYPE% (device model/type). Examples: 'CORP-MAC-%SERIAL%', '%DEVICETYPE%-%SERIAL%'. Limited to ~63 characters and only applied to supervised devices.
- `diagnostics_disabled` (Boolean) When true, skips the Diagnostics & Usage screen during enrollment.
- `display_tone_setup_disabled` (Boolean) When true, skips the Display Tone setup screen during enrollment.
- `dont_auto_populate_primary_account_info` (Boolean) Indicates whether Setup Assistant will auto populate the primary account information.
- `enable_restrict_editing` (Boolean) Indicates whether the user will enable restricting editing.
- `enrollment_time_azure_ad_group_ids` (Set of String) List of Azure AD group IDs to associate with the profile at enrollment time.
- `file_vault_disabled` (Boolean) When true, skips the FileVault encryption screen during enrollment.
- `hide_admin_account` (Boolean) Indicates whether the admin account should be hidden or not. Defaults to false when absent from config.
- `i_cloud_diagnostics_disabled` (Boolean) When true, skips the iCloud Diagnostics screen during enrollment.
- `i_cloud_storage_disabled` (Boolean) When true, skips the iCloud Storage screen during enrollment.
- `is_mandatory` (Boolean) Indicates if the profile is mandatory.
- `location_disabled` (Boolean) When true, skips the Location Services screen during enrollment.
- `pass_code_disabled` (Boolean) When true, skips the Passcode setup screen during enrollment.
- `primary_account_full_name` (String) Full name of the primary account to be created during enrollment.
- `primary_account_user_name` (String) User name of the primary account to be created during enrollment.
- `privacy_pane_disabled` (Boolean) When true, skips the Privacy consent screen during enrollment. **Note:** This screen can ONLY be controlled via this boolean field. The Microsoft Graph API rejects 'Privacy' as a skip key string, even though the boolean property works correctly. The provider handles this limitation automatically.
- `profile_removal_disabled` (Boolean) Indicates if the profile removal option is disabled.
- `registration_disabled` (Boolean) When true, skips the device registration screen during enrollment. **Note:** This screen can ONLY be controlled via this boolean field. The Microsoft Graph API rejects 'Registration' as a skip key string, even though the boolean property works correctly. The provider handles this limitation automatically.
- `request_requires_network_tether` (Boolean) Indicates if the device is network tethered to run the command.
- `restore_blocked` (Boolean) When true, blocks the Restore from iCloud/Mac screen during enrollment.
- `screen_time_screen_disabled` (Boolean) When true, skips the Screen Time setup screen during enrollment.
- `set_primary_setup_account_as_regular_user` (Boolean) Indicates whether Setup Assistant will set the account as a regular user.
- `siri_disabled` (Boolean) When true, skips the Siri setup screen during enrollment.
- `skip_primary_setup_account_creation` (Boolean) Indicates whether Setup Assistant will skip the user interface for primary account setup.
- `supervised_mode_enabled` (Boolean) Supervised mode. If true, the device enters supervised mode during enrollment.
- `support_department` (String) Support department information.
- `support_phone_number` (String) Support phone number.
- `terms_and_conditions_disabled` (Boolean) When true, skips the Terms and Conditions screen during enrollment.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `touch_id_disabled` (Boolean) When true, skips the Touch ID setup screen during enrollment.
- `wait_for_device_configured_confirmation` (Boolean) Indicates if the device will need to wait for configured confirmation.
- `welcome_screen_disabled` (Boolean) When true, skips the Welcome ("Hello" multilingual) screen during enrollment. Requires macOS Sonoma or later. **Note:** This screen can ONLY be controlled via the "Welcome" skip key string. The Microsoft Graph API does not provide a boolean property for this screen.
- `zoom_disabled` (Boolean) When true, skips the Zoom setup screen during enrollment.

### Read-Only

- `configuration_endpoint_url` (String) Apple Configurator enrollment configuration endpoint URL generated by Intune.
- `id` (String) The unique identifier of the enrollment profile.
- `is_default` (Boolean) Indicates if this is the default profile. This is managed server-side by Intune and cannot be set via Terraform.

<a id="nestedatt--dep_profile_admin_account_password_rotation_setting"></a>
### Nested Schema for `dep_profile_admin_account_password_rotation_setting`

Required:

- `auto_rotation_period_in_days` (Number) Number of days between 1 and 180 since the last rotation after which to auto-rotate the local admin password.

Optional:

- `dep_profile_delay_auto_rotation_setting` (Attributes) Settings for delaying auto-rotation of the local admin password after retrieval. (see [below for nested schema](#nestedatt--dep_profile_admin_account_password_rotation_setting--dep_profile_delay_auto_rotation_setting))

<a id="nestedatt--dep_profile_admin_account_password_rotation_setting--dep_profile_delay_auto_rotation_setting"></a>
### Nested Schema for `dep_profile_admin_account_password_rotation_setting.dep_profile_delay_auto_rotation_setting`

Optional:

- `on_retrieval_auto_rotate_password_enabled` (Boolean) Indicates whether to auto-rotate the password after retrieval.
- `on_retrieval_delay_auto_rotate_password_in_hours` (Number) Number of hours to delay auto-rotation of the password after retrieval.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Apple Business Manager / School Manager**: This resource requires that your organization has set up Apple Business Manager (ABM) or Apple School Manager (ASM) and synced tokens with Microsoft Intune.
- **DEP Onboarding Settings**: The `dep_onboarding_settings_id` is optional - if not provided, the provider will automatically resolve it from the Intune account ID.
- **Setup Assistant Customization**: Use boolean fields to control which Setup Assistant screens are shown during enrollment. The provider automatically converts these to the `enabledSkipKeys` array in the API.
- **LAPS Support**: Configure Local Administrator Password Solution (LAPS) for automatic password rotation using the `dep_profile_admin_account_password_rotation_setting` nested block.
- **Admin Account Creation**: When `admin_account_user_name`, `admin_account_full_name`, and `admin_account_password` are provided, a local administrator account is created on enrolled devices.
- **Write-Only Password**: The `admin_account_password` field is write-only - the API never returns this value. Terraform preserves it from the configuration for proper state management.
- **Privacy and Registration Screens**: The `privacy_pane_disabled` and `registration_disabled` fields control screens that cannot be skipped via the `enabledSkipKeys` array due to Microsoft Graph API limitations. These work via their boolean properties only.
- **Zero-Touch Enrollment**: Set `skip_primary_setup_account_creation` to `true` and configure admin accounts for fully automated zero-touch deployment scenarios.
- **Supervised Mode**: Enable `supervised_mode_enabled` to unlock additional device management capabilities including advanced restrictions and configuration options.
- **Device Naming**: Use the `device_name_template` field with variables like `%SERIAL%`, `%RAND:X%`, or `%DEVICETYPE%` to automatically name devices during enrollment.
- **Mandatory Enrollment**: Set `is_mandatory` to `true` to prevent users from skipping enrollment, ensuring all devices are properly managed.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Import using the enrollment profile ID.
# After import, set dep_onboarding_settings_id in your config if your tenant
# has multiple DEP tokens (ABM/ASM). Otherwise, the provider will auto-resolve
# it from the Intune account ID, which may select the wrong token.
#
# The admin_account_password field is write-only and cannot be recovered after
# import. You must re-supply it in your Terraform config and run terraform apply
# to reconcile state.

# {resource_id}
terraform import microsoft365_graph_beta_device_management_dep_macos_enrollment_profile.example 00000000-0000-0000-0000-000000000000
```
