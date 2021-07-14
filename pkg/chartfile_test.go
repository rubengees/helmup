package pkg

import (
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadChartfile(t *testing.T) {
	type args struct {
		path string
	}

	expectedChartfile, err := chartutil.LoadChartfile(filepath.Join("testdata", "project", "Chart.yaml"))
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name    string
		args    args
		want    *chart.Metadata
		wantErr bool
	}{
		{
			name:    "Existing Chart.yaml",
			args:    args{path: filepath.Join("testdata", "project")},
			want:    expectedChartfile,
			wantErr: false,
		},
		{
			name:    "Invalid path",
			args:    args{path: filepath.Join("testdata")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadChartfile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadChartfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadChartfile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateChartfile(t *testing.T) {
	type args struct {
		path         string
		chartfile    *chart.Metadata
		dependencies UpdatableDependencies
	}

	chartfile, err := chartutil.LoadChartfile(filepath.Join("testdata", "project", "Chart.yaml"))
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Update version",
			args: args{
				path:      t.TempDir(),
				chartfile: chartfile,
				dependencies: UpdatableDependencies{
					&UpdatableDependency{
						Name:           "helmup-test",
						Repository:     "file://testdata/repository",
						CurrentVersion: "10.0.0",
						LatestVersion:  "11.0.0",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateChartfile(tt.args.path, tt.args.chartfile, tt.args.dependencies)
			if err != nil {
				t.Errorf("LoadChartfile() error = %v", err)
				return
			}

			file, err := os.ReadFile(filepath.Join(tt.args.path, "Chart.yaml"))
			if err != nil {
				t.Errorf("LoadChartfile() expected file Chart.yaml could not be read: %v", err)
				return
			}

			fileString := string(file)
			if !strings.Contains(fileString, "11.0.0") {
				t.Error("LoadChartfile() saved chartfile does not contain expected string '11.0.0'")
				return
			}
		})
	}
}
