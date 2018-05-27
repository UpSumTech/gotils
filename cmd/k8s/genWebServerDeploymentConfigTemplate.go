/*
cat input.json

{
  "limits": {
    "cpu": "100m",
    "memory": "100Mi"
  },
  "termination_grace_period": 60,
  "deadline": 300,
  "docker_user": "sumanmukherjee03",
}
*/

package k8s

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

func genWebServerDeploymentTemplate(input WebServerDeployment) *appsv1.Deployment {
	deploymentName := strings.Join([]string{"deployment", imageName, imageTag}, "-")
	runLabel := strings.Join([]string{imageName, imageTag}, "-")
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
			Labels: map[string]string{
				LABEL_RUN_KEY:    runLabel,
				LABEL_SERVER_KEY: "web",
				LABEL_APP_KEY:    imageName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					LABEL_RUN_KEY:    runLabel,
					LABEL_SERVER_KEY: "web",
					LABEL_APP_KEY:    imageName,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: getIntOrStringPtr("25%"),
					MaxSurge:       getIntOrStringPtr("50%"),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						LABEL_RUN_KEY:    runLabel,
						LABEL_SERVER_KEY: "web",
						LABEL_APP_KEY:    imageName,
					},
				},
				Spec: corev1.PodSpec{
					Hostname: imageName,
					Containers: []corev1.Container{
						{
							Name:            imageName,
							Image:           strings.Join([]string{imageName, imageTag}, ":"),
							ImagePullPolicy: corev1.PullAlways,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: *utils.Int32Ptr(appPort),
								},
							},
							Env: append(getDefaultEnvVars(), []corev1.EnvVar{
								corev1.EnvVar{
									Name: BINTRAY_TOKEN_ENV_VAR,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: BINTRAY_TOKEN_SECRET_NAME,
											},
											Key: BINTRAY_TOKEN_SECRET_KEY,
										},
									},
								},
								corev1.EnvVar{
									Name: BINTRAY_USERNAME_ENV_VAR,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: BINTRAY_TOKEN_SECRET_NAME,
											},
											Key: BINTRAY_USERNAME_SECRET_KEY,
										},
									},
								},
								corev1.EnvVar{
									Name: BINTRAY_REPO_NAME_ENV_VAR,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: BINTRAY_TOKEN_SECRET_NAME,
											},
											Key: BINTRAY_REPO_NAME_SECRET_KEY,
										},
									},
								},
							}...),
							Resources: getResourceRequirements(input.Limits, input.Requests),
							SecurityContext: &corev1.SecurityContext{
								Privileged: utils.BoolPtr(false),
							},
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						corev1.LocalObjectReference{
							Name: DOCKER_CONFIG_SECRET_NAME,
						},
					},
					DNSPolicy:                     corev1.DNSClusterFirst,
					RestartPolicy:                 corev1.RestartPolicyAlways,
					TerminationGracePeriodSeconds: utils.Int64Ptr(input.TerminationGracePeriod),
				},
			},
		},
	}
}
