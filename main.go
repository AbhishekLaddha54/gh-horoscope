package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/cli/go-gh/pkg/tableprinter"
	"github.com/spf13/cobra"
)

type Issue struct {
	Title       string           `json:"title"`
	Number      int              `json:"number"`
	State       string           `json:"state"`
	PullRequest *PullRequestLink `json:"pull_request,omitempty"` // Nullable
}
type PullRequestLink struct {
	URL string `json:"url"`
}

type PullRequest struct {
	Title  string `json:"title"`
	Number int    `json:"number"`
	State  string `json:"state"`
	Draft  bool   `json:"draft"`
}

var (
	repoFlag string
	zodiac   = []string{"♈ Aries", "♉ Taurus", "♊ Gemini", "♋ Cancer", "♌ Leo", "♍ Virgo", "♎ Libra", "♏ Scorpio", "♐ Sagittarius", "♑ Capricorn", "♒ Aquarius", "♓ Pisces"}
	fortunes = []string{
		"Merge conflicts will test your patience – breathe 🧘",
		"A green CI run is in the stars today 🟢",
		"Someone will finally review your 3-week-old PR 👀",
		"Your next bug fix will earn you 10 new internet points 🏆",
		"A retrograde issue will re-open; handle with cosmic calm 🌌",
	}
)

func main() {
	cmd := &cobra.Command{
		Use:   "horoscope",
		Short: "Daily horoscope powered by your GitHub issues & PRs",
		RunE:  run,
	}

	cmd.Flags().StringVarP(&repoFlag, "repo", "r", "", "GitHub repository in owner/name format (e.g., cli/cli)")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	client, err := gh.RESTClient(nil)
	if err != nil {
		return err
	}

	// Determine repository context
	var repo repository.Repository
	if repoFlag != "" {
		repo, err = repository.Parse(repoFlag)
		if err != nil {
			return fmt.Errorf("invalid repo format: %v", err)
		}
	} else {
		repo, err = gh.CurrentRepository()
		if err != nil {
			return fmt.Errorf("could not determine current repository. Use --repo to specify one")
		}
	}

	owner := repo.Owner()
	name := repo.Name()

	// Current user
	var user struct{ Login string }
	if err := client.Get("user", &user); err != nil {
		return err
	}

	// Issues
	var issues []Issue
	issuesURL := fmt.Sprintf("repos/%s/%s/issues?state=open&per_page=100", owner, name)
	if err := client.Get(issuesURL, &issues); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch issues: %v\n", err)
		return err
	}

	// Filter out PRs from issues (GitHub's /issues includes both issues and PRs)
	filteredIssues := make([]Issue, 0)
	for _, i := range issues {
		if i.PullRequest == nil { // This works if you include a PullRequest field
			filteredIssues = append(filteredIssues, i)
		}
	}

	// Pull requests
	var prs []PullRequest
	prsURL := fmt.Sprintf("repos/%s/%s/pulls?state=open&per_page=100", owner, name)
	if err := client.Get(prsURL, &prs); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch pull requests: %v\n", err)
		return err
	}

	// Generate horoscope
	sign := zodiac[int(time.Now().Unix()/86400)%len(zodiac)]
	fortune := fortunes[rand.Intn(len(fortunes))]

	fmt.Printf("🌠 Daily Horoscope for @%s (%s)\n", user.Login, sign)
	fmt.Printf("🔗 Repository: %s/%s\n\n", owner, name)
	fmt.Printf("📊 %d open issues, %d open PRs\n", len(filteredIssues), len(prs))
	fmt.Printf("🔮 %s\n\n", fortune)

	// Render table
	tp := tableprinter.New(os.Stdout, true, 120)
	tp.AddField("Type")
	tp.AddField("#")
	tp.AddField("Title")
	tp.EndRow()
	for _, i := range filteredIssues {
		tp.AddField("Issue")
		tp.AddField(fmt.Sprintf("#%d", i.Number))
		tp.AddField(i.Title)
		tp.EndRow()
	}
	for _, p := range prs {
		typ := "PR"
		if p.Draft {
			typ = "Draft"
		}
		tp.AddField(typ)
		tp.AddField(fmt.Sprintf("#%d", p.Number))
		tp.AddField(p.Title)
		tp.EndRow()
	}
	tp.Render()

	return nil
}
