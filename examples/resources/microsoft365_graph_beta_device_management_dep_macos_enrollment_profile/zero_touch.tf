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
