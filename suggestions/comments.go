package suggestions

import "github.com/google/go-github/v52/github"

func isValidComment(issueComment github.IssueComment) bool {
	body := issueComment.GetBody()
	return body == "/suggest"
}
