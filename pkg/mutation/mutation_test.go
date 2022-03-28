package mutation

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMutatePodPatch(t *testing.T) {
	m := NewMutator(logger())
	got, err := m.MutatePodPatch(pod())
	if err != nil {
		t.Fatal(err)
	}

	p := patch()
	g := string(got)
	assert.Equal(t, p, g)
}

func BenchmarkMutatePodPatch(b *testing.B) {
	m := NewMutator(logger())
	pod := pod()

	for i := 0; i < b.N; i++ {
		_, err := m.MutatePodPatch(pod)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func pod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "lifespan",
			Labels: map[string]string{
				"acme.com/lifespan-requested": "7",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "lifespan",
				Image: "busybox",
			}},
		},
	}
}

func patch() string {
	patch := `
[
  {
    "op": "add",
    "path": "/spec/containers/0/volumeMounts",
    "value": [
      {
        "mountPath": "/etc/ssl/certs/ca-certificates.crt",
        "name": "ca-certificates",
        "readOnly": true,
        "subPath": "ca-certificates.crt"
      }
    ]
  },
  {
    "op": "add",
    "path": "/spec/volumes",
    "value": [
      {
        "configMap": {
          "name": "ca-certificates"
        },
        "name": "ca-certificates"
      }
    ]
  }
]
`

	patch = strings.ReplaceAll(patch, "\n", "")
	patch = strings.ReplaceAll(patch, "\t", "")
	patch = strings.ReplaceAll(patch, " ", "")

	return patch
}

func logger() *logrus.Entry {
	mute := logrus.StandardLogger()
	mute.Out = ioutil.Discard
	return mute.WithField("logger", "test")
}
