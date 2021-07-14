package pkg

import (
	"helm.sh/helm/v3/pkg/downloader"
	"os"
)

func UpdateCharts(path string, settings *ResolverSettings) error {
	helmDownloader := downloader.Manager{
		Out:              os.Stdout,
		ChartPath:        path,
		Verify:           downloader.VerifyNever,
		Debug:            false,
		SkipUpdate:       true,
		Getters:          settings.Getters,
		RepositoryConfig: settings.Env.RepositoryConfig,
		RepositoryCache:  settings.Env.RepositoryCache,
	}

	err := helmDownloader.Update()
	if err != nil {
		return err
	}

	return nil
}
