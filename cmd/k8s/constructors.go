package k8s

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

func NewBintrayTokenSecretTemplate() *BintrayTokenSecretTemplate {
	i := BintrayTokenSecretTemplate{}
	i.BintrayUser = utils.GetBintrayUser()
	i.BintrayRepo = utils.GetBintrayRepo()
	i.BintrayToken = utils.GetBintrayToken()
	return &i
}

func NewGithubTokenSecretTemplate() *GithubTokenSecretTemplate {
	i := GithubTokenSecretTemplate{}
	i.GithubUser = utils.GetGithubUser()
	i.GithubToken = utils.GetGithubToken()
	return &i
}

func NewDockerConfigSecretTemplate() *DockerConfigSecretTemplate {
	i := DockerConfigSecretTemplate{}
	home, err := homedir.Dir()
	if err != nil {
		utils.CheckErr(err.Error())
	}
	i.DockerConfigFile = filepath.Join(home, ".docker", "config.json")
	i.DockerUser = utils.GetDockerhubUser()
	i.DockerRegistry = DOCKER_REGISTRY
	i.DockerRegistryDomain = DOCKER_REGISTRY_DOMAIN
	return &i
}

func NewJobArtifactTemplate() *JobArtifactTemplate {
	i := JobArtifactTemplate{}
	i.Commands = []string{}
	i.Limits.Cpu = CPU_LIMIT
	i.Limits.Memory = MEM_LIMIT
	i.Requests.Cpu = CPU_REQUEST
	i.Requests.Memory = MEM_REQUEST
	i.TerminationGracePeriod = TERMINATION_LIMIT_SECS
	i.Deadline = DEADLINE_LIMIT_SECS
	i.DockerUser = utils.GetDockerhubUser()
	i.DockerRegistry = DOCKER_REGISTRY
	i.DockerRegistryDomain = DOCKER_REGISTRY_DOMAIN
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
	i.DockerUser = utils.GetDockerhubUser()
	i.DockerRegistry = DOCKER_REGISTRY
	i.DockerRegistryDomain = DOCKER_REGISTRY_DOMAIN
	return &i
}

func NewImageBuilderTemplate() *ImageBuilderTemplate {
	i := ImageBuilderTemplate{}
	i.Commands = []string{}
	i.Limits.Cpu = CPU_LIMIT
	i.Limits.Memory = MEM_LIMIT
	i.Requests.Cpu = CPU_REQUEST
	i.Requests.Memory = MEM_REQUEST
	i.TerminationGracePeriod = TERMINATION_LIMIT_SECS
	i.Deadline = DEADLINE_LIMIT_SECS
	i.DockerUser = utils.GetDockerhubUser()
	i.DockerRegistry = DOCKER_REGISTRY
	i.DockerRegistryDomain = DOCKER_REGISTRY_DOMAIN
	i.GitBranch = DEFAULT_GIT_BRANCH
	i.ReleaseVersion = DEFAULT_RELEASE_VERSION
	return &i
}
