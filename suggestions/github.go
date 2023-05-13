package suggestions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

type IncomingIssue struct {
	owner   string
	repo    string
	number  int
	ref     string
	comment *github.IssueComment
}

type EventType struct {
	Comment      *struct{} `json:"comment,omitempty"`
	Pull_request *struct{} `json:"pull_request,omitempty"`
}

func githubClient() *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func parseEventData() error {
	eventData, err := os.ReadFile(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(eventData, &issueCommentEvent); err != nil {
		return err
	}
	return nil
}

func getIssueInfo(ctx context.Context, ghClient *github.Client, issueCommentEvent github.IssueCommentEvent) (IncomingIssue, error) {
	owner := issueCommentEvent.Repo.GetOwner().GetLogin()
	repo := issueCommentEvent.Repo.GetName()
	prNumber := issueCommentEvent.GetIssue().GetNumber()
	pull, _, err := ghClient.PullRequests.Get(ctx, owner, repo, prNumber)

	if err != nil {
		return IncomingIssue{}, err
	}

	ref := pull.GetHead().GetRef()
	commentID := issueCommentEvent.Comment.GetID()

	comment, _, err := ghClient.Issues.GetComment(ctx, owner, repo, commentID)
	if err != nil {
		return IncomingIssue{}, err
	}

	return IncomingIssue{owner: owner, repo: repo, number: prNumber, ref: ref, comment: comment}, nil
}

func listChangedFiles(ctx context.Context, ghClient *github.Client, issueData IncomingIssue) ([]*github.CommitFile, error) {
	opts := &github.ListOptions{}
	files, _, err := ghClient.PullRequests.ListFiles(ctx, issueData.owner, issueData.repo, issueData.number, opts)
	if err != nil {
		return []*github.CommitFile{}, err
	}

	return files, nil
}

func createComment(ghClient *github.Client, comment IssueSuggestionComment, issueData IncomingIssue) {
	formattedMessage := fmt.Sprintf("<details>\n<summary>Suggestion(s): %s</summary>\n\n```%s\n%s\n```\n</details>\n", comment.filename, comment.language, comment.message)
	issueComment := github.IssueComment{Body: github.String(formattedMessage)}
	if _, _, err := ghClient.Issues.CreateComment(ctx, issueData.owner, issueData.repo, issueData.number, &issueComment); err != nil {
		log.Fatal(err)
	}
}

func readFile(filename string) (string, error) {
	if isDebugMode() {
		filename = fmt.Sprintf("run-local/%s", filename)
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func getFileExtension(path string) string {
	extension := filepath.Ext(path)

	if extension != "" {
		return extension
	}

	return ""
}

func getIgnoredFileExtensions() []string {
	return strings.Split(os.Getenv("IGNORED_FILE_EXTENSIONS"), ",")
}

func isIgnoredFileExtension(extension string) bool {
	fileExtensions := getIgnoredFileExtensions()
	for _, ignoredExtension := range fileExtensions {
		if extension == ignoredExtension {
			return true
		}
	}

	return false
}

func isDebugMode() bool {
	return os.Getenv("DEBUG_MODE") == "true"
}

func supportedMarkdownLanguages() map[string]string {
	return map[string]string{
		".feature":         "cucumber",
		".abap":            "abap",
		".adb":             "ada",
		".ads":             "ada",
		".ada":             "ada",
		".ahk":             "ahk",
		".ahkl":            "ahk",
		".htaccess":        "apacheconf",
		".apache.conf":     "apacheconf",
		".apache2.conf":    "apacheconf",
		".applescript":     "applescript",
		".as":              "as",
		".asy":             "asy",
		".sh":              "bash",
		".ksh":             "bash",
		".bash":            "bash",
		".ebuild":          "bash",
		".eclass":          "bash",
		".bat":             "bat",
		".cmd":             "bat",
		".befunge":         "befunge",
		".bmx":             "blitzmax",
		".boo":             "boo",
		".bf":              "brainfuck",
		".b":               "brainfuck",
		".c":               "c",
		".h":               "c",
		".cfm":             "cfm",
		".cfml":            "cfm",
		".cfc":             "cfm",
		".tmpl":            "cheetah",
		".spt":             "cheetah",
		".cl":              "cl",
		".lisp":            "cl",
		".el":              "cl",
		".clj":             "clojure",
		".cljs":            "clojure",
		".cmake":           "cmake",
		".coffee":          "coffeescript",
		".sh-session":      "console",
		".control":         "control",
		".cpp":             "cpp",
		".hpp":             "cpp",
		".c++":             "cpp",
		".h++":             "cpp",
		".cc":              "cpp",
		".hh":              "cpp",
		".cxx":             "cpp",
		".hxx":             "cpp",
		".pde":             "cpp",
		".cs":              "csharp",
		".css":             "css",
		".pyx":             "cython",
		".pxd":             "cython",
		".pxi":             "cython",
		".d":               "d",
		".di":              "d",
		".pas":             "delphi",
		".diff":            "diff",
		".patch":           "diff",
		".dpatch":          "diff",
		".darcspatch":      "diff",
		".duel":            "duel",
		".jbst":            "duel",
		".dylan":           "dylan",
		".dyl":             "dylan",
		".erb":             "erb",
		".erl-sh":          "erl",
		".erl":             "erl",
		".hrl":             "erl",
		".evoque":          "evoque",
		".factor":          "factor",
		".flx":             "felix",
		".flxh":            "felix",
		".f":               "fortran",
		".f90":             "fortran",
		".s":               "gas",
		".S":               "gas",
		".kid":             "genshi",
		".vert":            "glsl",
		".frag":            "glsl",
		".geo":             "glsl",
		".plot":            "gnuplot",
		".plt":             "gnuplot",
		".go":              "go",
		".1234567":         "groff",
		".man":             "groff",
		".haml":            "haml",
		".hs":              "haskell",
		".html":            "html",
		".htm":             "html",
		".xhtml":           "html",
		".xslt":            "html",
		".hx":              "hx",
		".hy":              "hybris",
		".hyb":             "hybris",
		".ini":             "ini",
		".cfg":             "ini",
		".io":              "io",
		".ik":              "ioke",
		".weechatlog":      "irc",
		".jade":            "jade",
		".java":            "java",
		".js":              "js",
		".jsp":             "jsp",
		".lhs":             "lhs",
		".ll":              "llvm",
		".lgt":             "logtalk",
		".lua":             "lua",
		".wlua":            "lua",
		".mak":             "make",
		".makefile":        "make",
		".GNUmakefile":     "make",
		".mao":             "mako",
		".maql":            "maql",
		".mhtml":           "mason",
		".mc":              "mason",
		".mi":              "mason",
		".autohandler":     "mason",
		".dhandler":        "mason",
		".md":              "markdown",
		".mo":              "modelica",
		".def":             "modula2",
		".mod":             "modula2",
		".moo":             "moocode",
		".mu":              "mupad",
		".mxml":            "mxml",
		".myt":             "myghty",
		".autodelegate":    "myghty",
		".asm":             "nasm",
		".ASM":             "nasm",
		".ns2":             "newspeak",
		".objdump":         "objdump",
		".m":               "objectivec",
		".j":               "objectivej",
		".ml":              "ocaml",
		".mli":             "ocaml",
		".mll":             "ocaml",
		".mly":             "ocaml",
		".ooc":             "ooc",
		".pl":              "perl",
		".pm":              "perl",
		".php":             "php",
		".php345":          "php",
		".php3":            "php",
		".php4":            "php",
		".php5":            "php",
		".ps":              "postscript",
		".eps":             "postscript",
		".pot":             "pot",
		".po":              "pot",
		".pov":             "pov",
		".inc":             "pov",
		".prolog":          "prolog",
		".pro":             "prolog",
		".properties":      "properties",
		".proto":           "protobuf",
		".py3tb":           "py3tb",
		".pytb":            "pytb",
		".py":              "python",
		".pyw":             "python",
		".sc":              "python",
		".SConstruct":      "python",
		".SConscript":      "python",
		".tac":             "python",
		".R":               "r",
		".rb":              "ruby",
		".rbw":             "ruby",
		".Rakefile":        "ruby",
		".rake":            "ruby",
		".gemspec":         "ruby",
		".rbx":             "ruby",
		".duby":            "ruby",
		".Rout":            "rconsole",
		".r":               "rebol",
		".r3":              "rebol",
		".cw":              "redcode",
		".rhtml":           "rhtml",
		".rst":             "rst",
		".rest":            "rst",
		".sass":            "sass",
		".scala":           "scala",
		".scaml":           "scaml",
		".scm":             "scheme",
		".scss":            "scss",
		".st":              "smalltalk",
		".tpl":             "smarty",
		".sources.list":    "sourceslist",
		".sql":             "sql",
		".sqlite3-console": "sqlite3",
		".squid.conf":      "squidconf",
		".ssp":             "ssp",
		".tcl":             "tcl",
		".tcsh":            "tcsh",
		".csh":             "tcsh",
		".tex":             "tex",
		".aux":             "tex",
		".toc":             "tex",
		".txt":             "text",
		".v":               "v",
		".sv":              "v",
		".vala":            "vala",
		".vapi":            "vala",
		".vb":              "vbnet",
		".bas":             "vbnet",
		".vm":              "velocity",
		".fhtml":           "velocity",
		".vim":             "vim",
		".vimrc":           "vim",
		".xml":             "xml",
		".xsl":             "xml",
		".rss":             "xml",
		".xsd":             "xml",
		".wsdl":            "xml",
		".xqy":             "xquery",
		".xquery":          "xquery",
		".yaml":            "yaml",
		".yml":             "yaml",
	}
}
