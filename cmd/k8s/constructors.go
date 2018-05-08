package k8s

func NewJobArtifactTemplate() *JobArtifactTemplate {
	i := JobArtifactTemplate{}
	i.Commands = []string{}
	i.Limits.Cpu = CPU_LIMIT
	i.Limits.Memory = MEM_LIMIT
	i.Requests.Cpu = CPU_REQUEST
	i.Requests.Memory = MEM_REQUEST
	i.TerminationGracePeriod = TERMINATION_LIMIT_SECS
	i.Deadline = DEADLINE_LIMIT_SECS
	return &i
}

func NewWebServerDeployment() *WebServerDeployment {
	i := WebServerDeployment{}
	i.Commands = []string{}
	i.Limits.Cpu = CPU_LIMIT
	i.Limits.Memory = MEM_LIMIT
	i.Requests.Cpu = CPU_LIMIT
	i.Requests.Memory = MEM_LIMIT
	i.TerminationGracePeriod = TERMINATION_LIMIT_SECS
	return &i
}
