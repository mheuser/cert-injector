package mutation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInjectCaCertificateMutate(t *testing.T) {
	want := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{{
				Name: "ca-certificates",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "ca-certificates",
						},
					},
				},
			}},
			Containers: []corev1.Container{{
				Name: "test",
				VolumeMounts: []corev1.VolumeMount{{
					Name:      "ca-certificates",
					MountPath: "/etc/ssl/certs/ca-certificates.crt",
					SubPath:   "ca-certificates.crt",
					ReadOnly:  true,
				}},
			}},
			InitContainers: []corev1.Container{{
				Name: "inittest",
				VolumeMounts: []corev1.VolumeMount{{
					Name:      "ca-certificates",
					MountPath: "/etc/ssl/certs/ca-certificates.crt",
					SubPath:   "ca-certificates.crt",
					ReadOnly:  true,
				}},
			}},
		},
	}

	pod := &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "test",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name: "test",
			}},
			InitContainers: []corev1.Container{{
				Name: "inittest",
			}},
		},
	}

	got, err := injectCaCertificate{Logger: logger()}.Mutate(pod)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, want, got)
}

func TestHasVolumeMount(t *testing.T) {
	yes := corev1.VolumeMount{
		Name:      "yes",
		MountPath: "/bar",
	}

	no := corev1.VolumeMount{
		Name:      "no",
		MountPath: "/foo",
	}

	c := corev1.Container{
		Name:         "test",
		VolumeMounts: []corev1.VolumeMount{yes},
	}

	assert.True(t, HasVolumeMount(c, yes))
	assert.False(t, HasVolumeMount(c, no))
}
