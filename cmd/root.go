package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	repogithub "repoready/internal/github"
	"repoready/internal/models"
	"repoready/internal/runner"
	"repoready/internal/scanner"
	"repoready/internal/system"
	"repoready/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:               "repoready [github-url]",
	Short:             "Scan a GitHub repository and compare it with your system",
	Args:              cobra.MaximumNArgs(1),
	SilenceErrors:     true,
	SilenceUsage:      true,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	RunE:              runScan,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		ui.PrintError(err.Error(), "")
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Println(ui.HelpScreen())
	})
	rootCmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Println(ui.HelpScreen())
		return nil
	})
}

func runScan(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	ui.PrintBanner()

	rawURL := ""
	if len(args) == 1 {
		rawURL = args[0]
	} else {
		value, err := ui.PromptRepositoryURL()
		if err != nil {
			return err
		}
		rawURL = value
	}

	repo, err := repogithub.ParseRepositoryURL(rawURL)
	if err != nil {
		return friendlyError(err)
	}

	client := repogithub.NewClient()

	tempDir, err := os.MkdirTemp("", "repoready-scan-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	target := filepath.Join(tempDir, repo.Name)
	var project models.ProjectInfo
	var report models.SystemReport
	if err := ui.RunScanSequence("Live Scan", []ui.ScanTask{
		{
			Label: "validate GitHub repository",
			Run: func() error {
				return client.Validate(ctx, repo)
			},
		},
		{
			Label: "download repository files",
			Run: func() error {
				return client.Clone(ctx, repo, target)
			},
		},
		{
			Label: "inspect project files",
			Run: func() error {
				var analyzeErr error
				project, analyzeErr = scanner.Analyze(target, repo.URL, repo.Owner, repo.Name)
				return analyzeErr
			},
		},
		{
			Label: "compare with your system",
			Run: func() error {
				report = system.NewChecker().CheckRequired(ctx, project.RequiredTools)
				return nil
			},
		},
	}); err != nil {
		return friendlyError(err)
	}
	fmt.Println()

	ui.PrintProjectAnalysis(project)
	ui.PrintSystemCheck(report)
	missing := ui.PrintDiagnosis(project, report)
	if len(missing) == 0 && project.MainLanguage != "Unknown" {
		if err := maybeDownloadRepository(ctx, client, repo); err != nil {
			return err
		}
	} else if len(missing) > 0 {
		shouldFix, err := ui.Confirm("Should I fix this problem for you now?", false)
		if err != nil {
			return err
		}
		if shouldFix {
			steps := system.BuildFixPlan(missing)
			ui.PrintFixPlan(steps)
			if err := runFixPlan(ctx, steps); err != nil {
				return err
			}
			report = system.NewChecker().CheckRequired(ctx, project.RequiredTools)
			ui.PrintFixResult(report)
			if len(system.MissingTools(report)) == 0 {
				if err := maybeDownloadRepository(ctx, client, repo); err != nil {
					return err
				}
			}
		}
	}
	ui.PrintFooter()
	return nil
}

func runFixPlan(ctx context.Context, steps []system.FixStep) error {
	for _, step := range steps {
		if !step.CanRun || step.RequiresAdmin {
			continue
		}

		fmt.Println()
		fmt.Println(ui.InfoBox("RepoReady can run:\n\n" + step.Command))
		confirmed, err := ui.Confirm("Run this command now?", false)
		if err != nil {
			return err
		}
		if !confirmed {
			continue
		}

		var output string
		if err := ui.RunSpinner("installing "+step.Tool.Name, func() error {
			var runErr error
			output, runErr = runner.RunShell(ctx, step.Command)
			return runErr
		}); err != nil {
			if strings.TrimSpace(output) != "" {
				fmt.Println(ui.ErrorBox(output))
			}
			return err
		}
		fmt.Println(ui.SuccessStyle.Render("✓ installed " + step.Tool.Name))
	}
	return nil
}

func maybeDownloadRepository(ctx context.Context, client *repogithub.Client, repo repogithub.Repository) error {
	target, err := repogithub.SafeClonePath(defaultDownloadRoot(), repo.Owner, repo.Name)
	if err != nil {
		return err
	}

	ui.PrintDownloadOffer(target)
	shouldDownload, err := ui.Confirm("Do you want me to download this project for you now?", false)
	if err != nil {
		return err
	}
	if !shouldDownload {
		return nil
	}

	if err := ui.RunSpinner("downloading project", func() error {
		return client.Clone(ctx, repo, target)
	}); err != nil {
		return friendlyError(err)
	}
	ui.PrintDownloadSuccess(target)
	return nil
}

func defaultDownloadRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "RepoReady")
	}
	return filepath.Join(home, "RepoReady")
}
