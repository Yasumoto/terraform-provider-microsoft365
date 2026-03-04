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
