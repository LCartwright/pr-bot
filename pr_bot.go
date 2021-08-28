package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v38/github"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

var (
	logger, _ = zap.NewProduction()
)

type Config struct {
	Token                   string            `required:"true" envconfig:"GH_ACCESS_TOKEN"`
	Owner                   string            `required:"true" envconfig:"GH_OWNER"`
	Repo                    string            `required:"true" envconfig:"GH_REPO"`
	URL                     string            `required:"true" envconfig:"GH_URL"`
	BaseBranch              string            `required:"true" envconfig:"GH_BASE_BRANCH"`
	GHLoginToSlackUsernames map[string]string `required:"true" envconfig:"GH_LOGIN_TO_SLACK_USERNAMES"`
}

func main() {
	cfg := mustLoadConfig()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	var allPullRequests []*github.PullRequest
	for {
		pullRequests, resp, err := client.PullRequests.List(ctx, cfg.Owner, cfg.Repo, &github.PullRequestListOptions{
			State:     "open",
			Sort:      "created",
			Direction: "asc",
			Base:      cfg.BaseBranch,
			ListOptions: github.ListOptions{
				Page:    0,
				PerPage: 10,
			},
		})
		if err != nil {
			logger.Error("error fetching pull requests", zap.Error(err))
			break
		}
	prs:
		for _, pr := range pullRequests {
			if pr.GetDraft() {
				logger.Sugar().Infof("dicarding PR [%s] as it's in draft", pr.GetTitle())
				continue prs
			}

			for _, label := range pr.Labels {
				if label.GetName() == "wip" {
					logger.Sugar().Infof("dicarding PR [%s] as it's a wip", pr.GetTitle())
					continue prs
				}
			}
			allPullRequests = append(allPullRequests, pr)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
		time.Sleep(1000 * time.Millisecond)
	}

	logger.Sugar().Infof("fetched [%d] PRs", len(allPullRequests))
	sb := strings.Builder{}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("*Pull Requests for [%s](%s", cfg.Repo, cfg.URL))
	sb.WriteString("/pulls?q=is%3Apr+is%3Aopen+-label%3Awip+review%3Arequired+sort%3Acreated-asc)*\n")
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("_Beep, boop, Happy %s! We have *%d* Open Pull Requests today!_\n", time.Now().Weekday().String(), len(allPullRequests)))
	for _, pr := range allPullRequests {
		openDuration := time.Since(pr.GetCreatedAt())
		openDurationString := openDuration.String()
		if openDuration.Hours() > 24 {
			openDurationString = fmt.Sprintf("%dd%dh", int(openDuration.Hours())/24, int(openDuration.Hours())%24)
		}

		title := strings.Replace(pr.GetTitle(), "`", "", -1)
		sb.WriteString(fmt.Sprintf(" * *[%s](%s)*\n", title, pr.GetHTMLURL()))
		sb.WriteString(fmt.Sprintf("    * Author: %s, open for: *%s*, last updated: *%s*\n", ghLoginToSlackMention(cfg, pr.GetUser()), openDurationString, pr.GetUpdatedAt().Format(time.Stamp)))
		if len(pr.RequestedReviewers) > 0 {
			sb.WriteString("    * Requested reviewer(s): ")
			prefix := ""
			for _, user := range pr.RequestedReviewers {
				sb.WriteString(fmt.Sprintf("%s%s", prefix, ghLoginToSlackMention(cfg, user)))
				prefix = ", "
			}
			sb.WriteString(".\n")
		} else {
			sb.WriteString("    * *No pending reviewers!*\n")
		}
	}
	logger.Info("Output string")
	fmt.Print(sb.String())
}

func ghLoginToSlackMention(cfg *Config, user *github.User) string {
	slackMention, ok := cfg.GHLoginToSlackUsernames[user.GetLogin()]
	if !ok {
		return fmt.Sprintf("[%s](%s)", user.GetLogin(), user.GetHTMLURL())
	}
	return fmt.Sprintf("@%s ", slackMention)
}

func mustLoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Info("failed to load .env, continuing", zap.Error(err))
	}
	var cfg Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		logger.Fatal("error binding config from env", zap.Error(err))
	}
	return &cfg
}
