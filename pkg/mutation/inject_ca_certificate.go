package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// injectEnv is a container for the mutation injecting environment vars
type injectCaCertificate struct {
	Logger logrus.FieldLogger
}

// injectEnv implements the podMutator interface
var _ podMutator = (*injectCaCertificate)(nil)

// Name returns the struct name
func (se injectCaCertificate) Name() string {
	return "inject_ca_certificate"
}

// Mutate returns a new mutated pod according to set env rules
func (se injectCaCertificate) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {
	se.Logger = se.Logger.WithField("mutation", se.Name()).WithField("configMap_name", "ca-certificates")
	mpod := pod.DeepCopy()

	volume := corev1.Volume{
		Name: "ca-certificates",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: "ca-certificates"},
			},
		},
	}

	volumeMount := corev1.VolumeMount{
		Name:      volume.Name,
		MountPath: "/etc/ssl/certs/ca-certificates.crt",
		SubPath:   "ca-certificates.crt",
		ReadOnly:  true,
	}

	if HasVolume(*pod, volume) {
		se.Logger.Info("pod has already a volume for the configmap")
		return pod, nil
	}

	injectVolume(mpod, volume)

	injectVolumeMount(mpod, volumeMount)
	se.Logger.Info("added volume and volumeMount")

	return mpod, nil
}

func injectVolume(pod *corev1.Pod, volume corev1.Volume) {
	pod.Spec.Volumes = append(pod.Spec.Volumes, volume)
}

func HasVolume(pod corev1.Pod, checkVolume corev1.Volume) bool {
	for _, volume := range pod.Spec.Volumes {
		if volume.Name == checkVolume.Name {
			return true
		}
	}
	return false
}

func injectVolumeMount(pod *corev1.Pod, volumeMount corev1.VolumeMount) {
	for i, container := range pod.Spec.Containers {
		if !HasVolumeMount(container, volumeMount) {
			pod.Spec.Containers[i].VolumeMounts = append(container.VolumeMounts, volumeMount)
		}
	}
	for i, container := range pod.Spec.InitContainers {
		if !HasVolumeMount(container, volumeMount) {
			pod.Spec.InitContainers[i].VolumeMounts = append(container.VolumeMounts, volumeMount)
		}
	}
}

func HasVolumeMount(container corev1.Container, checkVolumeMount corev1.VolumeMount) bool {
	for _, volumeMount := range container.VolumeMounts {
		if volumeMount.Name == checkVolumeMount.Name {
			return true
		}
	}
	return false
}
