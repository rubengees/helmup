package pkg

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"os"
)

func Prompt(updatableDependencies UpdatableDependencies) (UpdatableDependencies, error) {
	prompt := &survey.MultiSelect{
		Message: "Choose which dependencies to update.",
		Options: updatableDependencies.Strings(),
	}

	var chosenUpdates []int
	err := survey.AskOne(prompt, &chosenUpdates)
	if err == terminal.InterruptErr {
		os.Exit(0)
	} else if err != nil {
		return nil, err
	}

	var result UpdatableDependencies
	for _, chosen := range chosenUpdates {
		result = append(result, updatableDependencies[chosen])
	}

	return result, nil
}
