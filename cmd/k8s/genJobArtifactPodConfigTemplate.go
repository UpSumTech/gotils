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

func genArtifactJobTemplate(input JobArtifactTemplate) *corev1.Pod {
	if len(volumeImageName) == 0 {
		utils.CheckErr("You need to provide a name for the volume image that holds the artifact to download")
	}
	if len(volumeImageTag) == 0 {
		utils.CheckErr("You need to provide a tag for the volume image that holds the artifact to download")
	}
	if volumeContainerPort == 0 {
		utils.CheckErr("You need to provide a port for the volume image that holds the artifact to download")
	}

	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.Join([]string{imageName, "pod"}, "-"),
			Namespace: namespace,
			Labels: map[string]string{
				"app":         imageName,
				"app_version": imageTag,
			},
			Annotations: map[string]string{
				"description": fmt.Sprintf("Builds the artifact for %s", imageName),
			},
		},
		Spec: corev1.PodSpec{
			Hostname: imageName,
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    utils.Int64Ptr(userId),
				RunAsNonRoot: utils.BoolPtr(true),
				FSGroup:      utils.Int64Ptr(groupId),
				SupplementalGroups: []int64{
					int64(groupId),
				},
			},
			Volumes: []corev1.Volume{
				corev1.Volume{
					"artifact-data",
					corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							Medium: corev1.StorageMediumDefault,
						},
					},
				},
				corev1.Volume{
					"dshm",
					corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							Medium: corev1.StorageMediumMemory,
						},
					},
				},
			},
			InitContainers: []corev1.Container{
				corev1.Container{
					Name:            imageName,
					Image:           strings.Join([]string{strings.Join([]string{input.DockerRegistryDomain, input.DockerUser, imageName}, "/"), imageTag}, ":"),
					ImagePullPolicy: corev1.PullAlways,
					VolumeMounts: []corev1.VolumeMount{
						corev1.VolumeMount{
							Name:      "artifact-data",
							MountPath: "/var/data/build",
							ReadOnly:  false,
						},
						corev1.VolumeMount{
							Name:      "dshm",
							MountPath: "/dev/shm",
							ReadOnly:  false,
						},
					},
					Env: append(getDefaultEnvVars(), []corev1.EnvVar{
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
					}...),
					Resources: getResourceRequirements(input.Limits, input.Requests),
					SecurityContext: &corev1.SecurityContext{
						AllowPrivilegeEscalation: utils.BoolPtr(false),
					},
				},
			},
			Containers: []corev1.Container{
				corev1.Container{
					Name:            volumeImageName,
					Image:           strings.Join([]string{strings.Join([]string{input.DockerRegistryDomain, input.DockerUser, volumeImageName}, "/"), volumeImageTag}, ":"),
					ImagePullPolicy: corev1.PullAlways,
					VolumeMounts: []corev1.VolumeMount{
						corev1.VolumeMount{
							Name:      "artifact-data",
							MountPath: "/var/data/build",
							ReadOnly:  false,
						},
						corev1.VolumeMount{
							Name:      "dshm",
							MountPath: "/dev/shm",
							ReadOnly:  false,
						},
					},
					Env: append(getDefaultEnvVars(), []corev1.EnvVar{
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
					}...),
					Ports: []corev1.ContainerPort{
						corev1.ContainerPort{
							ContainerPort: int32(volumeContainerPort),
							Protocol:      corev1.ProtocolTCP,
						},
					},
					Resources: getResourceRequirements(input.Limits, input.Requests),
					SecurityContext: &corev1.SecurityContext{
						AllowPrivilegeEscalation: utils.BoolPtr(false),
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
