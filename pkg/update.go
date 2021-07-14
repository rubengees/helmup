package pkg

import (
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"os"
)

func UpdateCharts(path string) error {
	helmCli := cli.New()
	helmDownloader := downloader.Manager{
		Out:              os.Stdout,
		ChartPath:        path,
		Verify:           downloader.VerifyNever,
		Debug:            false,
		SkipUpdate:       true,
		Getters:          getter.All(helmCli),
		RepositoryConfig: helmCli.RepositoryConfig,
		RepositoryCache:  helmCli.RepositoryCache,
	}

	err := helmDownloader.Update()
	if err != nil {
		return err
	}

	return nil
}
