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
