package pkg

import (
	"errors"
	"fmt"
	"helm.sh/helm/v3/pkg/chartutil"
	"os"
	"path/filepath"
)

func GetProjectPath(args []string) (string, error) {
	var result string

	if len(args) == 1 {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return "", err
		}

		stat, err := os.Stat(path)
		if err != nil {
			return "", err
		}

		if !stat.IsDir() {
			return "", errors.New(
				fmt.Sprintf(
					"Given path is not a directory. Do not pass the %s file but the containing directory to helmup.",
					chartutil.ChartfileName,
				),
			)
		}

		result = path
	} else {
		var err error
		result, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	return result, nil
}
