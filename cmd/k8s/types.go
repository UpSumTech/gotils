package k8s

const (
	CPU_LIMIT                       = "500m"
	MEM_LIMIT                       = "500Mi"
	CPU_REQUEST                     = "100m"
	MEM_REQUEST                     = "100Mi"
	TERMINATION_LIMIT_SECS          = 120
	DEADLINE_LIMIT_SECS             = 600
	DOCKER_USERNAME_ENV_VAR         = "DOCKERHUB_USERNAME"
	DOCKER_CONFIG_SECRET_NAME       = "docker-config"
	DOCKER_REGISTRY                 = "https://index.docker.io/v1"
	DOCKER_REGISTRY_DOMAIN          = "index.docker.io"
	GITHUB_TOKEN_SECRET_NAME        = "github-token"
	GITHUB_TOKEN_SECRET_KEY         = "token"
	GITHUB_USERNAME_SECRET_KEY      = "user"
	GITHUB_EMAIL_SECRET_KEY         = "email"
	GITHUB_USER_FULLNAME_SECRET_KEY = "fullname"
	GITHUB_USERNAME_ENV_VAR         = "GITHUB_USERNAME"
	GITHUB_USER_FULLNAME_ENV_VAR    = "GIT_USER"
	GITHUB_EMAIL_ENV_VAR            = "GIT_EMAIL"
	GITHUB_TOKEN_ENV_VAR            = "DEPLOY_GITHUB_TOKEN"
	BINTRAY_TOKEN_SECRET_NAME       = "bintray-token"
	BINTRAY_TOKEN_SECRET_KEY        = "token"
	BINTRAY_USERNAME_SECRET_KEY     = "user"
	BINTRAY_REPO_NAME_SECRET_KEY    = "repo"
	BINTRAY_USERNAME_ENV_VAR        = "BINTRAY_USERNAME"
	BINTRAY_TOKEN_ENV_VAR           = "BINTRAY_API_KEY"
	BINTRAY_REPO_NAME_ENV_VAR       = "BINTRAY_REPO_NAME"
	GIT_REPO_NAME_ENV_VAR           = "GIT_REPO_NAME"
	DEFAULT_GIT_BRANCH              = "master"
	DEFAULT_RELEASE_VERSION         = "patch"
	BUILDER_DATA_DIR_ENV_VAR        = "BUILDER_DATA_DIR"
	DEFAULT_BUILDER_DATA_DIR        = "/var/data/build"
	LABEL_RUN_KEY                   = "run"
	LABEL_SERVER_KEY                = "server"
	LABEL_APP_KEY                   = "app"
	CLUSTER_IP_DEFAULT              = "None"
)

type jsonInput interface {
	readInput() error
	validateInput() error
}

type jsonOutput interface {
	jsonOutput() (string, error)
}

type k8sTemplate interface {
	jsonInput
	build() interface{}
	jsonOutput
}

type ResourceLimitConfig struct {
	Cpu    string `json:"cpu" validate:"required"`
	Memory string `json:"memory" validate:"required"`
}

type ResourceRequestConfig struct {
	Cpu    string `json:"cpu" validate:"omitempty"`
	Memory string `json:"memory" validate:"omitempty"`
}

type DockerConfigSecretTemplate struct {
	DockerConfigFile     string `json:"dockercfg_file" validate:"required"`
	DockerUser           string `json:"docker_user" validate:"required"`
	DockerRegistry       string `json:"docker_registry" validate:"required"`
	DockerRegistryDomain string `json:"docker_registry_domain" validate:"required"`
}

func (i *DockerConfigSecretTemplate) readInput() error            { return readJson(i, src) }
func (i *DockerConfigSecretTemplate) validateInput() error        { return validateJsonInput(i) }
func (i *DockerConfigSecretTemplate) build() interface{}          { return genSecretDockerConfigTemplate(*i) }
func (i *DockerConfigSecretTemplate) jsonOutput() (string, error) { return getJsonTemplateOutput(i) }

type GithubTokenSecretTemplate struct {
	GithubUser         string `json:"github_user" validate:"required"`
	GithubUserFullname string `json:"github_user_fullname" validate:"required"`
	GithubEmail        string `json:"github_email" validate:"required"`
	GithubToken        string `json:"github_token" validate:"required"`
}

func (i *GithubTokenSecretTemplate) readInput() error            { return readJson(i, src) }
func (i *GithubTokenSecretTemplate) validateInput() error        { return validateJsonInput(i) }
func (i *GithubTokenSecretTemplate) build() interface{}          { return genSecretGithubTokenTemplate(*i) }
func (i *GithubTokenSecretTemplate) jsonOutput() (string, error) { return getJsonTemplateOutput(i) }

type BintrayTokenSecretTemplate struct {
	BintrayUser  string `json:"bintray_user" validate:"required"`
	BintrayRepo  string `json:"bintray_repo" validate:"required"`
	BintrayToken string `json:"bintray_token" validate:"required"`
}

func (i *BintrayTokenSecretTemplate) readInput() error            { return readJson(i, src) }
func (i *BintrayTokenSecretTemplate) validateInput() error        { return validateJsonInput(i) }
func (i *BintrayTokenSecretTemplate) build() interface{}          { return genSecretBintrayTokenTemplate(*i) }
func (i *BintrayTokenSecretTemplate) jsonOutput() (string, error) { return getJsonTemplateOutput(i) }

type JobArtifactTemplate struct {
	Limits                 ResourceLimitConfig   `json:"limits" validate:"required,dive,required"`
	Requests               ResourceRequestConfig `json:"requests"`
	TerminationGracePeriod int                   `json:"termination_grace_period" validate:"required"`
	Commands               []string              `json:"commands"`
	Deadline               int                   `json:"deadline" validate:"required"`
	DockerUser             string                `json:"docker_user" validate:"required"`
	DockerRegistry         string                `json:"docker_registry" validate:"required"`
	DockerRegistryDomain   string                `json:"docker_registry_domain" validate:"required"`
}

func (i *JobArtifactTemplate) readInput() error            { return readJson(i, src) }
func (i *JobArtifactTemplate) validateInput() error        { return validateJsonInput(i) }
func (i *JobArtifactTemplate) build() interface{}          { return genArtifactJobTemplate(*i) }
func (i *JobArtifactTemplate) jsonOutput() (string, error) { return getJsonTemplateOutput(i) }

type WebServerDeployment struct {
	Limits                 ResourceLimitConfig   `json:"requests" validate:"required,dive,required"`
	Requests               ResourceRequestConfig `json:"requests"`
	TerminationGracePeriod int                   `json:"termination_grace_period" validate:"required"`
	Commands               []string              `json:"commands"`
	DockerUser             string                `json:"docker_user" validate:"required"`
	DockerRegistry         string                `json:"docker_registry" validate:"required"`
	DockerRegistryDomain   string                `json:"docker_registry_domain" validate:"required"`
}

func (i *WebServerDeployment) readInput() error            { return readJson(i, src) }
func (i *WebServerDeployment) validateInput() error        { return validateJsonInput(i) }
func (i *WebServerDeployment) build() interface{}          { return genWebServerDeploymentTemplate(*i) }
func (i *WebServerDeployment) jsonOutput() (string, error) { return getJsonTemplateOutput(i) }

type ImageBuilderTemplate struct {
	Limits                 ResourceLimitConfig   `json:"limits" validate:"required,dive,required"`
	Requests               ResourceRequestConfig `json:"requests"`
	TerminationGracePeriod int                   `json:"termination_grace_period" validate:"required"`
	Commands               []string              `json:"commands"`
	Deadline               int                   `json:"deadline" validate:"required"`
	DockerUser             string                `json:"docker_user" validate:"required"`
	DockerRegistry         string                `json:"docker_registry" validate:"required"`
	DockerRegistryDomain   string                `json:"docker_registry_domain" validate:"required"`
	GitRepoUrl             string                `json:"git_repo_url" validate:"required"`
	GitBranch              string                `json:"git_branch" validate:"required"`
	ReleaseVersion         string                `json:"release_version" validate:"required"`
}

func (i *ImageBuilderTemplate) readInput() error            { return readJson(i, src) }
func (i *ImageBuilderTemplate) validateInput() error        { return validateJsonInput(i) }
func (i *ImageBuilderTemplate) build() interface{}          { return genImageBuilderPodConfigTemplate(*i) }
func (i *ImageBuilderTemplate) jsonOutput() (string, error) { return getJsonTemplateOutput(i) }

type ServiceCategory string

const (
	SERVICE_TYPE_INTERNAL ServiceCategory = "Internal"
	SERVICE_TYPE_EXTERNAL ServiceCategory = "External"
)

type ServiceTemplate struct {
	Category ServiceCategory `json:"category" validate:"required"`
}

func (i *ServiceTemplate) readInput() error            { return readJson(i, src) }
func (i *ServiceTemplate) validateInput() error        { return validateJsonInput(i) }
func (i *ServiceTemplate) build() interface{}          { return genServiceTemplate(*i) }
func (i *ServiceTemplate) jsonOutput() (string, error) { return getJsonTemplateOutput(i) }
