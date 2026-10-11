package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yundream/goslide/internal/i18n"
	"github.com/yundream/goslide/internal/model"
	"github.com/yundream/goslide/internal/updater"
	"github.com/yundream/goslide/pkg/goslide"
)

var (
	updateCheckFlag  bool
	updateDryRunFlag bool
	updateForceFlag  bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Automatically update Goslide binary to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		execPath, _ := os.Executable()
		if realPath, err := filepath.EvalSymlinks(execPath); err == nil && realPath != "" {
			execPath = realPath
		}

		u := updater.NewUpdater(nil, nil)

		opts := updater.UpdateOptions{
			CurrentVersion: goslide.Version,
			TargetPath:     execPath,
			DryRun:         updateDryRunFlag || updateCheckFlag,
			Force:          updateForceFlag,
			OnProgress: func(stage, detail string) {
				switch stage {
				case "checking":
					fmt.Println("[goslide] 🔍 " + i18n.T("update.checking"))
				case "downloading":
					fmt.Printf("[goslide] ⬇️  "+i18n.T("update.downloading")+"\n", detail)
				case "verifying":
					fmt.Println("[goslide] 🔒 " + i18n.T("update.verifying"))
				case "replacing":
					fmt.Printf("[goslide] 🔄 "+i18n.T("update.replacing")+"\n", detail)
				}
			},
		}

		res, err := u.Execute(ctx, opts)
		if err != nil {
			if errors.Is(err, updater.ErrPermissionDenied) {
				return newCLIError(model.ExitGeneralError, fmt.Errorf("%s\n(Detail: %w)", i18n.T("update.err_permission"), err))
			}
			return newCLIError(model.ExitGeneralError, fmt.Errorf("failed to update goslide: %w", err))
		}

		if res.AlreadyLatest {
			fmt.Printf("[goslide] ✨ "+i18n.T("update.already_latest")+"\n", res.UpdatedVersion)
			return nil
		}

		if updateCheckFlag {
			fmt.Printf("[goslide] 💡 "+i18n.T("update.found_new")+"\n", res.UpdatedVersion)
			fmt.Printf("[goslide]    "+i18n.T("update.current_and_latest")+"\n", res.CurrentVersion, res.UpdatedVersion)
			if res.ReleaseURL != "" {
				fmt.Printf(i18n.T("update.release_notes_link")+"\n", res.ReleaseURL)
			}
			return nil
		}

		if updateDryRunFlag {
			fmt.Printf("[goslide] 🧪 Dry-run complete. Would upgrade from %s to %s.\n", res.CurrentVersion, res.UpdatedVersion)
			return nil
		}

		fmt.Printf("\n[goslide] 🎉 "+i18n.T("update.success")+"\n", res.UpdatedVersion)

		if len(res.Highlights) > 0 {
			fmt.Printf("\n"+i18n.T("update.whats_new")+"\n", res.UpdatedVersion)
			for _, h := range res.Highlights {
				fmt.Printf("  • %s\n", h)
			}
		}

		if res.ReleaseURL != "" {
			fmt.Printf("\n"+i18n.T("update.release_notes_link")+"\n", res.ReleaseURL)
		}

		return nil
	},
}

func initUpdateFlags() {
	if updateCmd.Flags().Lookup("check") == nil {
		updateCmd.Flags().BoolVar(&updateCheckFlag, "check", false, i18n.T("cli.update.flag.check"))
	}
	if updateCmd.Flags().Lookup("dry-run") == nil {
		updateCmd.Flags().BoolVar(&updateDryRunFlag, "dry-run", false, i18n.T("cli.update.flag.dry_run"))
	}
	if updateCmd.Flags().Lookup("force") == nil {
		updateCmd.Flags().BoolVar(&updateForceFlag, "force", false, i18n.T("cli.update.flag.force"))
	}
}
