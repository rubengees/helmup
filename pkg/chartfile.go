package pkg

import (
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"path/filepath"
)

func LoadChartfile(path string) (*chart.Metadata, error) {
	_, err := chartutil.IsChartDir(path)
	if err != nil {
		return nil, err
	}

	file, err := chartutil.LoadChartfile(filepath.Join(path, chartutil.ChartfileName))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func UpdateChartfile(path string, chartfile *chart.Metadata, dependencies UpdatableDependencies) error {
	for _, dependency := range dependencies {
		for _, chartfileDependency := range chartfile.Dependencies {
			if dependency.Name == chartfileDependency.Name && dependency.Repository == chartfileDependency.Repository {
				chartfileDependency.Version = dependency.LatestVersion
				break
			}
		}
	}

	err := chartutil.SaveChartfile(filepath.Join(path, chartutil.ChartfileName), chartfile)
	if err != nil {
		return err
	}
	return nil
}
