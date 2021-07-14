package pkg

import (
	"bytes"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

type FilesystemGetter struct {
	repoUrls []string
}

func (g *FilesystemGetter) Get(href string, _ ...getter.Option) (*bytes.Buffer, error) {
	url := strings.TrimPrefix(href, "file://")

	file, err := os.ReadFile(url)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(file), nil
}

func TestResolveUpdates(t *testing.T) {
	type args struct {
		chartfile *chart.Metadata
		settings  *ResolverSettings
	}

	helmCli := cli.New()
	helmCli.RepositoryConfig = filepath.Join("testdata", "repositories.yaml")
	settings := &ResolverSettings{
		Env: helmCli,
		Getters: getter.Providers{getter.Provider{
			Schemes: []string{"file"},
			New: func(options ...getter.Option) (getter.Getter, error) {
				return &FilesystemGetter{}, nil
			},
		}},
	}

	chartfile, err := chartutil.LoadChartfile(filepath.Join("testdata", "project", "Chart.yaml"))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		args    args
		want    UpdatableDependencies
		wantErr bool
	}{
		{
			name: "Test",
			args: args{chartfile: chartfile, settings: settings},
			want: []*UpdatableDependency{
				{
					Name:           "helmup-test",
					Repository:     "file://testdata/repository",
					CurrentVersion: "10.0.0",
					LatestVersion:  "11.0.0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveUpdates(tt.args.chartfile, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveUpdates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveUpdates() got = %v, want %v", got, tt.want)
			}
		})
	}
}
