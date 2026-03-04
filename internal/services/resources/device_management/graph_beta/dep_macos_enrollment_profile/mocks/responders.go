package mocks

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
)

var mockState struct {
	sync.Mutex
	enrollmentProfiles map[string]map[string]any
}

func init() {
	mockState.enrollmentProfiles = make(map[string]map[string]any)
	httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, `{"error":{"code":"ResourceNotFound","message":"Resource not found"}}`))
	mocks.GlobalRegistry.Register("dep_macos_enrollment_profile", &DepMacOSEnrollmentProfileMock{})
}

type DepMacOSEnrollmentProfileMock struct{}

var _ mocks.MockRegistrar = (*DepMacOSEnrollmentProfileMock)(nil)

func (m *DepMacOSEnrollmentProfileMock) RegisterMocks() {
	mockState.Lock()
	mockState.enrollmentProfiles = make(map[string]map[string]any)
	mockState.Unlock()

	// 0. GET /deviceManagement - needed for import to resolve intuneAccountId
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement(\?.*)?$`, func(req *http.Request) (*http.Response, error) {
		responseObj := map[string]any{
			"@odata.context":   "https://graph.microsoft.com/beta/$metadata#deviceManagement",
			"id":               "00000000-0000-0000-0000-000000000000",
			"intuneAccountId":  "11111111-1111-1111-1111-111111111111",
			"subscriptionState": "active",
		}
		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 1. Group validation - called during validateRequest for enrollment_time_azure_ad_group_ids
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/groups/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		groupId := parts[len(parts)-1]

		responseObj := map[string]any{
			"id":          groupId,
			"displayName": "Test Group " + groupId[:8],
		}

		return httpmock.NewJsonResponse(200, responseObj)
	})

	// 2. Create DEP macOS Enrollment Profile - POST /deviceManagement/depOnboardingSettings/{id}/enrollmentProfiles
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/enrollmentProfiles$`, func(req *http.Request) (*http.Response, error) {
		var requestBody map[string]any
		json.NewDecoder(req.Body).Decode(&requestBody)

		profileId := uuid.New().String()

		// Choose JSON file based on config characteristics
		var jsonFile string
		if _, hasLAPS := requestBody["depProfileAdminAccountPasswordRotationSetting"]; hasLAPS {
			jsonFile = "../tests/responses/validate_create/post_dep_macos_enrollment_profile_laps.json"
		} else if admin, hasAdmin := requestBody["adminAccountUserName"]; hasAdmin && admin != nil {
			jsonFile = "../tests/responses/validate_create/post_dep_macos_enrollment_profile_with_admin.json"
		} else {
			jsonFile = "../tests/responses/validate_create/post_dep_macos_enrollment_profile_zero_touch.json"
		}

		jsonStr, _ := helpers.ParseJSONFile(jsonFile)
		var responseObj map[string]any
		json.Unmarshal([]byte(jsonStr), &responseObj)

		// Set dynamic fields from request
		responseObj["id"] = profileId

		// Copy all scalar fields from request to response (except password which is write-only)
		scalarFields := []string{
			"displayName", "description", "requiresUserAuthentication",
			"enableAuthenticationViaCompanyPortal", "requireCompanyPortalOnSetupAssistantEnrolledDevices",
			"supervisedModeEnabled", "isMandatory", "profileRemovalDisabled",
			"waitForDeviceConfiguredConfirmation", "autoAdvanceSetupEnabled",
			"skipPrimarySetupAccountCreation", "dontAutoPopulatePrimaryAccountInfo",
			"requestRequiresNetworkTether", "setPrimarySetupAccountAsRegularUser",
			"enableRestrictEditing", "setAutoAdvanceSetupEnabled",
		}
		for _, field := range scalarFields {
			if val, ok := requestBody[field]; ok {
				responseObj[field] = val
			}
		}

		// Copy array fields
		if groupIds, ok := requestBody["enrollmentTimeAzureAdGroupIds"].([]interface{}); ok {
			responseObj["enrollmentTimeAzureAdGroupIds"] = groupIds
		}
		if skipKeys, ok := requestBody["enabledSkipKeys"].([]interface{}); ok {
			responseObj["enabledSkipKeys"] = skipKeys
		}

		// Copy all boolean screen control fields from request to response
		booleanFields := []string{
			"appleIdDisabled", "applePayDisabled", "diagnosticsDisabled",
			"displayToneSetupDisabled", "locationDisabled", "privacyPaneDisabled",
			"restoreBlocked", "screenTimeScreenDisabled", "siriDisabled",
			"termsAndConditionsDisabled", "touchIdDisabled",
			"accessibilityScreenDisabled", "autoUnlockWithWatchDisabled",
			"chooseYourLockScreenDisabled", "fileVaultDisabled",
			"iCloudDiagnosticsDisabled", "iCloudStorageDisabled",
			"passCodeDisabled", "registrationDisabled", "zoomDisabled",
		}
		for _, field := range booleanFields {
			if val, ok := requestBody[field]; ok {
				responseObj[field] = val
			}
		}

		// Copy admin account fields (but not password - it's write-only)
		// Handle nil values for admin account clearing
		if val, ok := requestBody["adminAccountUserName"]; ok {
			if val == nil {
				responseObj["adminAccountUserName"] = nil
			} else {
				responseObj["adminAccountUserName"] = val
			}
		}
		if val, ok := requestBody["adminAccountFullName"]; ok {
			if val == nil {
				responseObj["adminAccountFullName"] = nil
			} else {
				responseObj["adminAccountFullName"] = val
			}
		}
		if val, ok := requestBody["hideAdminAccount"]; ok {
			responseObj["hideAdminAccount"] = val
		}

		// Password is write-only - never return it in response
		responseObj["adminAccountPassword"] = nil

		// LAPS rotation settings: only include if present in request
		if lapsSettings, ok := requestBody["depProfileAdminAccountPasswordRotationSetting"].(map[string]any); ok {
			responseObj["depProfileAdminAccountPasswordRotationSetting"] = lapsSettings
		} else {
			// Ensure no LAPS settings in response if not in request
			delete(responseObj, "depProfileAdminAccountPasswordRotationSetting")
		}

		mockState.Lock()
		mockState.enrollmentProfiles[profileId] = responseObj
		mockState.Unlock()

		return httpmock.NewJsonResponse(201, responseObj)
	})

	// 3. Get DEP macOS Enrollment Profile - GET /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles/{profileId}
	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/enrollmentProfiles/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		profileId := parts[len(parts)-1]

		mockState.Lock()
		profile, exists := mockState.enrollmentProfiles[profileId]
		mockState.Unlock()

		if !exists {
			errorObj := map[string]any{
				"error": map[string]any{
					"code":    "ResourceNotFound",
					"message": "DEP enrollment profile not found",
				},
			}
			return httpmock.NewJsonResponse(404, errorObj)
		}

		return httpmock.NewJsonResponse(200, profile)
	})

	// 4. Update DEP macOS Enrollment Profile - PATCH /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles/{profileId}
	httpmock.RegisterResponder("PATCH", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/enrollmentProfiles/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		profileId := parts[len(parts)-1]

		var requestBody map[string]any
		json.NewDecoder(req.Body).Decode(&requestBody)

		mockState.Lock()
		profile, exists := mockState.enrollmentProfiles[profileId]
		if exists {
			// Check if admin account is being cleared (both username and fullname are nil)
			adminClearing := false
			if userVal, hasUser := requestBody["adminAccountUserName"]; hasUser && userVal == nil {
				if fullVal, hasFull := requestBody["adminAccountFullName"]; hasFull && fullVal == nil {
					adminClearing = true
				}
			}

			// Update profile with request body fields
			for k, v := range requestBody {
				// Special handling for admin account clearing
				if (k == "adminAccountUserName" || k == "adminAccountFullName") && v == nil {
					profile[k] = nil
				} else if k != "adminAccountPassword" {
					// Don't store password (write-only)
					profile[k] = v
				}
			}

			// When clearing admin account, also set hideAdminAccount to false
			if adminClearing {
				profile["hideAdminAccount"] = false
			}

			// Password is always nil in response (write-only)
			profile["adminAccountPassword"] = nil
			mockState.enrollmentProfiles[profileId] = profile
		}
		mockState.Unlock()

		if !exists {
			errorObj := map[string]any{
				"error": map[string]any{
					"code":    "ResourceNotFound",
					"message": "DEP enrollment profile not found",
				},
			}
			return httpmock.NewJsonResponse(404, errorObj)
		}

		return httpmock.NewJsonResponse(200, profile)
	})

	// 5. Delete DEP macOS Enrollment Profile - DELETE /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles/{profileId}
	httpmock.RegisterResponder("DELETE", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/enrollmentProfiles/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		parts := strings.Split(req.URL.Path, "/")
		profileId := parts[len(parts)-1]

		mockState.Lock()
		delete(mockState.enrollmentProfiles, profileId)
		mockState.Unlock()

		return httpmock.NewStringResponse(204, ""), nil
	})
}

func (m *DepMacOSEnrollmentProfileMock) RegisterErrorMocks() {
	mockState.Lock()
	mockState.enrollmentProfiles = make(map[string]map[string]any)
	mockState.Unlock()

	// Return errors for all operations
	httpmock.RegisterResponder("POST", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/enrollmentProfiles$`, func(req *http.Request) (*http.Response, error) {
		jsonStr, _ := helpers.ParseJSONFile("../tests/responses/validate_create/post_dep_macos_enrollment_profile_error.json")
		var errorObj map[string]any
		json.Unmarshal([]byte(jsonStr), &errorObj)
		return httpmock.NewJsonResponse(400, errorObj)
	})

	httpmock.RegisterResponder("GET", `=~^https://graph\.microsoft\.com/beta/deviceManagement/depOnboardingSettings/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/enrollmentProfiles/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`, func(req *http.Request) (*http.Response, error) {
		errorObj := map[string]any{
			"error": map[string]any{
				"code":    "ResourceNotFound",
				"message": "DEP enrollment profile not found",
			},
		}
		return httpmock.NewJsonResponse(404, errorObj)
	})
}

func (m *DepMacOSEnrollmentProfileMock) CleanupMockState() {
	mockState.Lock()
	mockState.enrollmentProfiles = make(map[string]map[string]any)
	mockState.Unlock()
}
