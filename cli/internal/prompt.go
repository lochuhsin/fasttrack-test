package internal

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
)

func CreateNamePrompt() (string, error) {
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("Enter name"),
	}
	result, err := prompt.Run()
	// Handle any errors that occurred during the prompt
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	if result == "" {
		fmt.Println("Found empty username, create via uuid")
		return uuid.New().String(), nil
	}
	return result, nil
}

func RunProblemPrompt(questions []question) ([]entry, error) {
	results := make([]entry, 0, len(questions))
	for i, q := range questions {
		prompt := promptui.Select{
			Label: fmt.Sprintf("%s, Select Option", q.Q),
			Items: q.Options,
		}

		index, result, err := prompt.Run()
		if err != nil {
			return nil, err
		}
		fmt.Printf("You chose %q\n", result)
		if err != nil {
			return nil, err
		}
		results = append(results, entry{
			QuestionId: i,
			Answer:     index,
		})
	}
	return results, nil
}
