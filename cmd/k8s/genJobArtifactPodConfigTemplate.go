/*
cat input.json

{
  "limits": {
    "cpu": "200m",
    "memory": "200Mi"
  },
  "termination_grace_period": 90,
  "deadline": 300
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

func genJobArtifactTemplate(input JobArtifactTemplate) *corev1.Pod {
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
			Containers: []corev1.Container{
				corev1.Container{
					Name:            imageName,
					Image:           strings.Join([]string{strings.Join([]string{"sumanmukherjee03", imageName}, "/"), imageTag}, ":"),
					ImagePullPolicy: corev1.PullAlways,
					VolumeMounts: []corev1.VolumeMount{
						corev1.VolumeMount{
							Name:      "artifact-data",
							MountPath: "/var/data/build",
						},
						corev1.VolumeMount{
							Name:      "dshm",
							MountPath: "/dev/shm",
						},
					},
					Command: []string{
						"ls",
						"-lah",
						"/var/data",
					},
					Lifecycle: &corev1.Lifecycle{
						PreStop: &corev1.Handler{
							Exec: &corev1.ExecAction{
								Command: []string{
									"/usr/bin/env",
									"bash",
									"-c",
									"test ! -z $(ls -A /var/data/build)",
								},
							},
						},
					},
					Resources: getResourceRequirements(input.Limits, input.Requests),
				},
				corev1.Container{
					Name:            volumeImageName,
					Image:           strings.Join([]string{strings.Join([]string{"sumanmukherjee03", volumeImageName}, "/"), volumeImageTag}, ":"),
					ImagePullPolicy: corev1.PullAlways,
					VolumeMounts: []corev1.VolumeMount{
						corev1.VolumeMount{
							Name:      "artifact-data",
							MountPath: "/var/data/build",
						},
						corev1.VolumeMount{
							Name:      "dshm",
							MountPath: "/dev/shm",
						},
					},
					Ports: []corev1.ContainerPort{
						corev1.ContainerPort{
							ContainerPort: int32(volumeContainerPort),
							Protocol:      corev1.ProtocolTCP,
						},
					},
					Lifecycle: &corev1.Lifecycle{
						PostStart: &corev1.Handler{
							Exec: &corev1.ExecAction{
								Command: []string{
									"/usr/bin/env",
									"bash",
									"-c",
									"test ! -z $(ls -A /var/data/build)",
								},
							},
						},
					},
					Resources: getResourceRequirements(input.Limits, input.Requests),
				},
			},
			RestartPolicy:                 corev1.RestartPolicyNever,
			TerminationGracePeriodSeconds: utils.Int64Ptr(input.TerminationGracePeriod),
			ActiveDeadlineSeconds:         utils.Int64Ptr(input.Deadline),
		},
		Status: corev1.PodStatus{},
	}
}
