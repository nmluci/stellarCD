package indto

type DeploymentJobs struct {
	ID              string `toml:"id"`
	Command         string `toml:"cmd"`
	WorkingDir      string `toml:"cwd"`
	Message         string `toml:"msg"`
	SignatureHeader string `toml:"signature_header"`
	SignatureVal    string `toml:"signature_value"`
	TriggerKey      string `toml:"trigger_key"`
	TriggerRegex    string `toml:"trigger_regex"`
}
