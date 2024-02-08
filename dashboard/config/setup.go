package config

func GetServicesHealth() map[string]string {
	return map[string]string{
		"Api Service":      EnvApiServiceDomain() + "/api/health",
		"Dataset Service":  EnvDatasetServiceDomain() + "/datasets/health",
		"Dispatch Service": EnvDispatchServiceDomain() + "/dispatch/health",
		"Download Service": EnvDownloadServiceDomain() + "/transfer/health",
		"Export Service":   EnvExportServiceDomain() + "/export/health",
		"Model Service":    EnvModelServiceDomain() + "/models/health",
		"Request Service":  EnvRequestServiceDomain() + "/requests/health",
		"User Service":     EnvUserServiceDomain() + "/users/health",
	}
}
