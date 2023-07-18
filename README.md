# gpt-pr-suggestions

## What is `gpt-pr-suggestions`?

This project leverages the power of OpenAI's models to review and suggest improvements on code changes in a GitHub pull request. 

When an user comments `/suggest` on a pull request, the model then returns suggestions and improvements, which the bot automatically comments on the pull request.

## Installation

Before running the bot, you need to configure a few environment variables:

- `GITHUB_TOKEN`: Your GitHub API token.
- `OPEN_AI_TOKEN`: Your OpenAI API key to use the ChatGPT service.
- `OPEN_AI_MODEL`: OpenAI model, e.g., `gpt-3.5-turbo`.

1. Clone the repository: `git clone https://github.com/RafaelBroseghini/gpt-pr-suggestions`
2. Install the required dependencies: `go mod download`

### Testing Locally

To test this bot locally, you'll need to set up the `GITHUB_EVENT_PATH` environment variable. This variable should point to a file containing a JSON representation of the GitHub event.

One of the easiest ways to get a realistic JSON event is to extract it from a real GitHub workflow.

To get this event payload:

1. Set up a workflow on your repository that gets triggered on comments to a PR.
2. In the workflow run, add a step that prints the contents of the `GITHUB_EVENT_PATH` file (`cat $GITHUB_EVENT_PATH`). 
3. Run the workflow that includes the event you want to test (e.g., a comment on issue event).

Here's an example of such a step in a workflow:

```yaml
on:
  issue_comment

jobs:
  show-event:    
    name: Print event
    runs-on: ubuntu-latest
    steps:
    - name: Dump GitHub event
      cat $GITHUB_EVENT_PATH
```

4. Go to a PR in the target repository and comment `/suggest`.
5. Check the workflow logs to find the printed JSON event data.
6. Copy the JSON data into a local file, and set the `GITHUB_EVENT_PATH` environment variable to the path of this file before running your bot locally.

For more detailed instructions, you can reference: [Stack Overflow question](https://stackoverflow.com/questions/63803136/how-to-get-my-own-github-events-payload-json-for-testing-github-actions-locally).

Once the environment variables are set, you can run the bot with:

```bash
go run main.go
```
