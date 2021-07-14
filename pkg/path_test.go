package pkg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetProjectPath(t *testing.T) {
	type args struct {
		args []string
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Path from args",
			args:    args{args: []string{filepath.Join("testdata", "project")}},
			want:    filepath.Join(workingDirectory, "testdata", "project"),
			wantErr: false,
		},
		{
			name:    "Not existing path from args",
			args:    args{args: []string{"invalid"}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Path not directory from args",
			args:    args{args: []string{filepath.Join("testdata", "project", "Chart.yaml")}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Empty args",
			args:    args{args: []string{}},
			want:    workingDirectory,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProjectPath(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetProjectPath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
