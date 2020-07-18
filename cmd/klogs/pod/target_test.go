package pod

import (
	"reflect"
	"testing"
)

func TestNewLogThread(t *testing.T) {
	type args struct {
		context   string
		namespace string
		pod       string
		container string
	}
	tests := []struct {
		name string
		args args
		want *Target
	}{
		{
			name: "all non empty",
			args: args{
				context:   "testCtx",
				namespace: "testNs",
				pod:       "testPod",
				container: "testContainer",
			},
			want: &Target{
				Context:   "testCtx",
				Namespace: "testNs",
				Pod:       "testPod",
				Container: "testContainer",
			},
		},
		{
			name: "all empty",
			args: args{
				context:   "",
				namespace: "",
				pod:       "",
				container: "",
			},
			want: &Target{
				Context:   "",
				Namespace: "",
				Pod:       "",
				Container: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTarget(tt.args.context, tt.args.namespace, tt.args.pod, tt.args.container, 0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}
