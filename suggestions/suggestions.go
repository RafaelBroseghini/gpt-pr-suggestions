package suggestions

import (
	"context"
	"log"

	"github.com/google/go-github/v52/github"
)

var (
	ctx               context.Context
	issueCommentEvent github.IssueCommentEvent
)

type IssueSuggestionComment struct {
	filename string
	message  string
	language string
}

func Run() error {
	ctx = context.Background()
	ghClient := githubClient()

	err := parseEventData()

	if err != nil {
		log.Fatal(err)
	}

	issueData, err := getIssueInfo(ctx, ghClient, issueCommentEvent)

	if err != nil {
		log.Fatal(err)
	}

	isValid := isValidComment(*issueData.comment)

	if !isValid {
		log.Printf("Comment: '%s' does not trigger PR suggestions.", *issueData.comment)
		return nil
	}

	changedFiles, err := listChangedFiles(ctx, ghClient, issueData)

	if err != nil {
		log.Fatal(err)
	}

	markdownLanguages := supportedMarkdownLanguages()

	oaClient := openAiClient()

	for _, changedFile := range changedFiles {
		filename := changedFile.GetFilename()

		fileExtension := getFileExtension(filename)

		if isIgnoredFileExtension(fileExtension) {
			log.Printf("Ignoring file: %s", filename)
			continue
		}

		detectedLanguage := markdownLanguages[fileExtension]

		if changedFile.GetStatus() == "deleted" {
			continue
		}

		fileContent, err := readFile(filename)

		if err != nil {
			log.Fatal(err)
		}

		chatCompletion, err := suggestImprovements(oaClient, fileContent)

		if err != nil {
			displayAvailableModels(ctx, oaClient, err)
			log.Fatal(err)
		}

		answer := chatCompletion.Choices[0].Message.Content
		suggestionComment := IssueSuggestionComment{filename: filename, message: answer, language: detectedLanguage}

		if isDebugMode() {
			log.Println(suggestionComment)
		} else {
			createComment(ghClient, suggestionComment, issueData)
		}
	}
	return nil
}
