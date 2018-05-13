/*
cat input.json

{
  "limits": {
    "cpu": "200m",
    "memory": "200Mi"
  },
  "termination_grace_period": 90,
  "deadline": 300,
  "docker_user": "sumanmukherjee03"
}
*/

package k8s

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

func genImageBuilderPodConfigTemplate(input ImageBuilderTemplate) *corev1.Pod {
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.Join([]string{imageName, "image-builder-pod"}, "-"),
			Namespace: namespace,
			Labels: map[string]string{
				"app":         builderImageName,
				"app_version": builderImageTag,
			},
			Annotations: map[string]string{
				"description": fmt.Sprintf("Builds the docker image for %s", imageName),
			},
		},
		Spec: corev1.PodSpec{
			Hostname: imageName,
			Volumes: []corev1.Volume{
				corev1.Volume{
					"docker-socket",
					corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: "/var/run/docker.sock",
							Type: hostPathTypePtr(corev1.HostPathSocket),
						},
					},
				},
				corev1.Volume{
					"docker-config",
					corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName:  DOCKER_CONFIG_SECRET_NAME,
							DefaultMode: utils.Int32Ptr(420),
						},
					},
				},
				corev1.Volume{
					"builder-data",
					corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							Medium: corev1.StorageMediumMemory,
						},
					},
				},
			},
			Containers: []corev1.Container{
				corev1.Container{
					Name:            builderImageName,
					Image:           strings.Join([]string{strings.Join([]string{input.DockerRegistryDomain, input.DockerUser, builderImageName}, "/"), builderImageTag}, ":"),
					ImagePullPolicy: corev1.PullAlways,
					VolumeMounts: []corev1.VolumeMount{
						corev1.VolumeMount{
							Name:      "builder-data",
							MountPath: "/var/data/build",
							ReadOnly:  true,
						},
						corev1.VolumeMount{
							Name:      "docker-socket",
							MountPath: "/var/run/docker.sock",
							ReadOnly:  false,
						},
						corev1.VolumeMount{
							Name:      "docker-config",
							MountPath: "/root/.docker/config.json",
							ReadOnly:  true,
						},
					},
					Env: append(getDefaultEnvVars(), []corev1.EnvVar{
						corev1.EnvVar{
							Name:  DOCKER_USERNAME_ENV_VAR,
							Value: input.DockerUser,
						},
						corev1.EnvVar{
							Name:  GIT_REPO_NAME_ENV_VAR,
							Value: input.RepoName,
						},
						corev1.EnvVar{
							Name: GITHUB_USERNAME_ENV_VAR,
							ValueFrom: &corev1.EnvVarSource{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: GITHUB_TOKEN_SECRET_NAME,
									},
									Key: GITHUB_USERNAME_SECRET_KEY,
								},
							},
						},
						corev1.EnvVar{
							Name: GITHUB_TOKEN_ENV_VAR,
							ValueFrom: &corev1.EnvVarSource{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: GITHUB_TOKEN_SECRET_NAME,
									},
									Key: GITHUB_TOKEN_SECRET_KEY,
								},
							},
						},
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
						Privileged: utils.BoolPtr(true),
					},
				},
			},
			RestartPolicy:                 corev1.RestartPolicyNever,
			TerminationGracePeriodSeconds: utils.Int64Ptr(input.TerminationGracePeriod),
			ActiveDeadlineSeconds:         utils.Int64Ptr(input.Deadline),
		},
		Status: corev1.PodStatus{},
	}
}
