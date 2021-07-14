package pkg

import (
	"fmt"
	"github.com/mgutz/ansi"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/repo"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type UpdatableDependency struct {
	Name           string
	Repository     string
	CurrentVersion string
	LatestVersion  string
}

func (dep UpdatableDependency) String() string {
	return fmt.Sprintf(
		"%s (%s): %s -> %s",
		ansi.Color(dep.Name, "black+b"),
		dep.Repository,
		ansi.Color(dep.CurrentVersion, "blue"),
		ansi.Color(dep.LatestVersion, "green"),
	)
}

type UpdatableDependencies []*UpdatableDependency

func (deps UpdatableDependencies) Strings() []string {
	var result []string

	for _, dependency := range deps {
		result = append(result, dependency.String())
	}

	return result
}

type ChartDependencies []*chart.Dependency

type ResolverSettings struct {
	Env     *cli.EnvSettings
	Getters getter.Providers
}

func ResolveUpdates(chartfile *chart.Metadata, settings *ResolverSettings) (UpdatableDependencies, error) {
	repofile, err := repo.LoadFile(settings.Env.RepositoryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load repositories: %w", err)
	}

	groupedDependencies := groupByRepository(repofile, chartfile)

	var result UpdatableDependencies
	for resolvedDependency := range resolveParallel(groupedDependencies, settings) {
		result = append(result, resolvedDependency)
	}

	sort.Slice(result, func(i, j int) bool {
		return strings.Compare(result[i].Name, result[j].Name) < 0
	})

	return result, nil
}

func groupByRepository(repofile *repo.File, chartfile *chart.Metadata) map[*repo.Entry]ChartDependencies {
	var result = make(map[*repo.Entry]ChartDependencies)

	for _, dependency := range chartfile.Dependencies {
		repository := findRepoByURL(repofile.Repositories, dependency.Repository)
		if repository == nil {
			fmt.Println(
				fmt.Sprintf(
					"Repository %s for %s not found. Ignoring...\n",
					dependency.Repository,
					ansi.Color(dependency.Name, "black+b"),
				),
			)
			continue
		}

		if result[repository] != nil {
			result[repository] = append(result[repository], dependency)
		} else {
			result[repository] = ChartDependencies{dependency}
		}
	}

	return result
}

func findRepoByURL(repositories []*repo.Entry, url string) *repo.Entry {
	for _, repository := range repositories {
		if repository.URL == url {
			return repository
		}
	}

	return nil
}

func resolveParallel(groupedDependencies map[*repo.Entry]ChartDependencies, settings *ResolverSettings) <-chan *UpdatableDependency {
	out := make(chan *UpdatableDependency)

	go func() {
		defer close(out)

		var wg sync.WaitGroup

		for repository, dependencies := range groupedDependencies {
			wg.Add(1)

			go func(repository *repo.Entry, dependencies ChartDependencies) {
				for updatableDependency := range resolveForRepository(repository, dependencies, settings) {
					out <- updatableDependency
				}

				wg.Done()
			}(repository, dependencies)
		}

		wg.Wait()
	}()

	return out
}

func resolveForRepository(repository *repo.Entry, dependencies ChartDependencies, settings *ResolverSettings) chan *UpdatableDependency {
	out := make(chan *UpdatableDependency)

	go func() {
		defer close(out)

		chartRepository, err := repo.NewChartRepository(repository, settings.Getters)
		if err != nil {
			fmt.Println(
				fmt.Sprintf("Failed to update repository %s: %v. Skipping...", repository.Name, err),
			)
			return
		}

		_, err = chartRepository.DownloadIndexFile()
		if err != nil {
			fmt.Println(
				fmt.Sprintf("Failed to update repository %s: %v. Skipping...", repository.Name, err),
			)
			return
		}

		repositoryIndex := filepath.Join(settings.Env.RepositoryCache, helmpath.CacheIndexFile(repository.Name))
		loadedIndex, err := repo.LoadIndexFile(repositoryIndex)
		if err != nil {
			fmt.Println(
				fmt.Sprintf("Failed to load index for repository %s: %v. Skipping...", repository.Name, err),
			)
			return
		}

		for _, dependency := range dependencies {
			latestVersion, err := loadedIndex.Get(dependency.Name, "")
			if err != nil {
				fmt.Println(
					fmt.Sprintf(
						"Failed to resolve latest version for dependency %s: %v. Skipping...",
						ansi.Color(repository.Name, "black+b"),
						err,
					),
				)
				continue
			}

			if dependency.Version != latestVersion.Version {
				out <- &UpdatableDependency{
					Name:           dependency.Name,
					Repository:     dependency.Repository,
					CurrentVersion: dependency.Version,
					LatestVersion:  latestVersion.Version,
				}
			}
		}
	}()

	return out
}
