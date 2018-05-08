package k8s

import (
	"github.com/sumanmukherjee03/gotils/cmd/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	CPU_LIMIT              = "500m"
	MEM_LIMIT              = "500Mi"
	CPU_REQUEST            = "100m"
	MEM_REQUEST            = "100Mi"
	TERMINATION_LIMIT_SECS = 120
	DEADLINE_LIMIT_SECS    = 600
)

type JsonInput interface {
	readInput() error
	validateInput() error
}

type JsonOutput interface {
	jsonOutput() (string, error)
}

type K8sTemplate interface {
	JsonInput
	JsonOutput
}

type ResourceLimitConfig struct {
	Cpu    string `json:"cpu" validate:"required"`
	Memory string `json:"memory" validate:"required"`
}

type ResourceRequestConfig struct {
	Cpu    string `json:"cpu" validate:"omitempty"`
	Memory string `json:"memory" validate:"omitempty"`
}

type JobArtifactTemplate struct {
	Limits                 ResourceLimitConfig   `json:"limits" validate:"required,dive,required"`
	Requests               ResourceRequestConfig `json:"requests"`
	TerminationGracePeriod int                   `json:"termination_grace_period" validate:"required"`
	Commands               []string              `json:"commands"`
	Deadline               int                   `json:"deadline" validate:"required"`
}

func (i *JobArtifactTemplate) readInput() error     { return readJson(i, src) }
func (i *JobArtifactTemplate) validateInput() error { return validateJsonInput(i) }
func (i *JobArtifactTemplate) build() *corev1.Pod   { return genJobArtifactTemplate(*i) }
func (i *JobArtifactTemplate) jsonOutput() (string, error) {
	var data string

	err := i.readInput()
	if err != nil {
		return data, err
	}

	err = i.validateInput()
	if err != nil {
		return data, err
	}

	o := i.build()
	data, err = utils.ToJson(o)
	if err != nil {
		return data, err
	}

	return data, nil
}

type WebServerDeployment struct {
	Limits                 ResourceLimitConfig   `json:"requests" validate:"required,dive,required"`
	Requests               ResourceRequestConfig `json:"requests"`
	TerminationGracePeriod int                   `json:"termination_grace_period" validate:"required"`
	Commands               []string              `json:"commands"`
}

func (i *WebServerDeployment) readInput() error          { return readJson(i, src) }
func (i *WebServerDeployment) validateInput() error      { return validateJsonInput(i) }
func (i *WebServerDeployment) build() *appsv1.Deployment { return genWebServerDeploymentTemplate(*i) }
func (i *WebServerDeployment) jsonOutput() (string, error) {
	var data string

	err := i.readInput()
	if err != nil {
		return data, err
	}

	err = i.validateInput()
	if err != nil {
		return data, err
	}

	o := i.build()
	data, err = utils.ToJson(o)
	if err != nil {
		return data, err
	}

	return data, nil
}
