package pod

import (
	"github.com/cinarmert/klogs/cmd/klogs/config"
	v1 "k8s.io/api/core/v1"
	"reflect"
	"regexp"
	"testing"
)

type fakePod struct {
	pod *v1.Pod
}

func (fp *fakePod) withName(name string) *fakePod {
	fp.pod.Name = name
	return fp
}

func (fp *fakePod) withNamespace(namespace string) *fakePod {
	fp.pod.Namespace = namespace
	return fp
}

func (fp *fakePod) withInitContainerStatuses(st []v1.ContainerStatus) *fakePod {
	fp.pod.Status.InitContainerStatuses = st
	return fp
}

func (fp *fakePod) withContainerStatuses(st []v1.ContainerStatus) *fakePod {
	fp.pod.Status.ContainerStatuses = st
	return fp
}

func TestManager_GetTargets(t *testing.T) {
	matchAll, _ := regexp.Compile(".*")
	pm := &Manager{
		Pods: nil,
		Config: &config.Config{
			PodRegex:       matchAll,
			ContainerRegex: matchAll,
		},
	}

	cs := []v1.ContainerStatus{{Name: "testContainer1"}, {Name: "testContainer2"}}
	fp := (&fakePod{pod: &v1.Pod{}}).
		withNamespace("testNamespace").
		withName("testPod").
		withContainerStatuses(cs[:1]).
		withInitContainerStatuses(cs[1:])

	pm.Pods = &v1.PodList{
		Items: []v1.Pod{*fp.pod},
	}

	tests := []struct {
		name string
		want []*Target
	}{
		{
			name: "one init, one regular container",
			want: []*Target{
				{
					Context:   "",
					Namespace: "testNamespace",
					Pod:       "testPod",
					Container: "testContainer2",
				},
				{
					Context:   "",
					Namespace: "testNamespace",
					Pod:       "testPod",
					Container: "testContainer1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pm.GetTargets(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTargets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_getTargetsFromPod(t *testing.T) {
	matchAll, _ := regexp.Compile(".*")
	pm := &Manager{
		Config: &config.Config{
			Context:        "testCtx",
			PodRegex:       matchAll,
			ContainerRegex: matchAll,
		},
	}

	cs := []v1.ContainerStatus{{Name: "testContainer1"}, {Name: "testContainer2"}}
	fp := (&fakePod{pod: &v1.Pod{}}).
		withNamespace("testNamespace").
		withName("testPod").
		withContainerStatuses(cs[:1]).
		withInitContainerStatuses(cs[1:])

	tests := []struct {
		name string
		want []*Target
	}{
		{
			name: "one init, one regular container",
			want: []*Target{
				{
					Context:   "testCtx",
					Namespace: "testNamespace",
					Pod:       "testPod",
					Container: "testContainer2",
				},
				{
					Context:   "testCtx",
					Namespace: "testNamespace",
					Pod:       "testPod",
					Container: "testContainer1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pm.getTargetsFromPod(fp.pod); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTargetsFromPod() = %v, want %v", got, tt.want)
			}
		})
	}
}
