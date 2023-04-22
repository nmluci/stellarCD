package router

const (
	basePath       = "/v1"
	PingPath       = basePath + "/ping"
	ReflectorPath  = basePath + "/reflector"
	DeploymentPath = basePath + "/deploy/:jobID"
)
