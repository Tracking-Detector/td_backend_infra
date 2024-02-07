package config

import "os"

func EnvDomain() string {
	return os.Getenv("DOMAIN")
}

func EnvApiServiceDomain() string {
	return os.Getenv("API_SERVICE_DOMAIN")
}

func EnvDatasetServiceDomain() string {
	return os.Getenv("DATASET_SERVICE_DOMAIN")
}

func EnvDispatchServiceDomain() string {
	return os.Getenv("DISPATCH_SERVICE_DOMAIN")
}

func EnvDownloadServiceDomain() string {
	return os.Getenv("DOWNLOAD_SERVICE_DOMAIN")
}

func EnvExportServiceDomain() string {
	return os.Getenv("EXPORT_SERVICE_DOMAIN")
}

func EnvModelServiceDomain() string {
	return os.Getenv("MODEL_SERVICE_DOMAIN")
}

func EnvRequestServiceDomain() string {
	return os.Getenv("REQUEST_SERVICE_DOMAIN")
}

func EnvUserServiceDomain() string {
	return os.Getenv("USER_SERVICE_DOMAIN")
}
