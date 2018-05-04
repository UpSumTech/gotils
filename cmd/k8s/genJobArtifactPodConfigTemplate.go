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
					"build-artifact-data",
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
				{
					Name:            imageName,
					Image:           strings.Join([]string{strings.Join([]string{"sumanmukherjee03", imageName}, "/"), imageTag}, ":"),
					ImagePullPolicy: corev1.PullAlways,
					VolumeMounts: []corev1.VolumeMount{
						corev1.VolumeMount{
							Name:      "build-artifact-data",
							MountPath: "/var/data",
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
									"test ! -z $(ls -A /var/data)",
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
