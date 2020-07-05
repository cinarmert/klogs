package ui

import (
	"github.com/cinarmert/klogs/cmd/klogs/pod"
	"github.com/rivo/tview"
	"gotest.tools/assert"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"testing"
)

func TestLayoutManager_createTextView(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantTitle string
	}{
		{
			name:      "empty title",
			title:     "",
			wantTitle: "  ",
		},
		{
			name:      "regular title",
			title:     "test title",
			wantTitle: " test title ",
		},
		{
			name:      "special chars in title",
			title:     "!/#test\"",
			wantTitle: " !/#test\" ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lm := &LayoutManager{
				App:  tview.NewApplication(),
				Grid: tview.NewGrid(),
			}
			got := lm.createTextView(tt.title)
			assert.Equal(t, got.GetTitle(), tt.wantTitle)
		})
	}
}

func TestNewManager(t *testing.T) {
	type args struct {
		targets []*pod.Target
		cs      v1.CoreV1Interface
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid request",
			args: args{
				targets: []*pod.Target{{Context: "test"}},
				cs:      &v1.CoreV1Client{},
			},
			wantErr: false,
		},
		{
			name: "empty targets",
			args: args{
				cs: &v1.CoreV1Client{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewManager(tt.args.targets, tt.args.cs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
