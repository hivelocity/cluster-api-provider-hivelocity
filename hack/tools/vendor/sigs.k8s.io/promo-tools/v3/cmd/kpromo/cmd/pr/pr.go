/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pr

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"sigs.k8s.io/promo-tools/v3/image"
	reg "sigs.k8s.io/promo-tools/v3/legacy/dockerregistry"
	"sigs.k8s.io/release-sdk/git"
	"sigs.k8s.io/release-sdk/github"
	"sigs.k8s.io/release-utils/util"
)

const (
	k8sioRepo             = "k8s.io"
	k8sioDefaultBranch    = "main"
	promotionBranchSuffix = "-image-promotion"
	defaultProject        = image.StagingRepoSuffix
	defaultReviewers      = "@kubernetes/release-engineering"
)

// PRCmd is the kpromo subcommand to promote container images
var PRCmd = &cobra.Command{
	Use:   "pr",
	Short: "Starts an image promotion for a given image tag",
	Long: `kpromo pr

This command updates image promoter manifests and creates a PR in
kubernetes/k8s.io`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Run the PR creation function
		return runPromote(promoteOpts)
	},
}

type promoteOptions struct {
	project         string
	userFork        string
	tags            []string
	reviewers       string
	interactiveMode bool
}

func (o *promoteOptions) Validate() error {
	if len(o.tags) == 0 {
		return errors.New("cannot start promotion --tag is required")
	}
	if o.userFork == "" {
		return errors.New("cannot start promotion --fork is required")
	}

	// Check the fork slug
	if _, _, err := git.ParseRepoSlug(o.userFork); err != nil {
		return errors.Wrap(err, "checking user's fork")
	}

	// Verify we got a valid tag
	for _, tag := range o.tags {
		if _, err := util.TagStringToSemver(tag); err != nil {
			return errors.Wrapf(err, "verifying tag: %s", tag)
		}
	}

	// Check that the GitHub token is set
	token, isSet := os.LookupEnv(github.TokenEnvKey)
	if !isSet || token == "" {
		return fmt.Errorf("cannot promote images if GitHub token env var %s is not set", github.TokenEnvKey)
	}
	return nil
}

var promoteOpts = &promoteOptions{}

func init() {
	PRCmd.PersistentFlags().StringVar(
		&promoteOpts.project,
		"project",
		defaultProject,
		"the name of the project to promote images for",
	)

	PRCmd.PersistentFlags().StringSliceVarP(
		&promoteOpts.tags,
		"tag",
		"t",
		[]string{},
		"version tag of the images we will promote",
	)

	PRCmd.PersistentFlags().StringVar(
		&promoteOpts.userFork,
		"fork",
		"",
		"the user's fork of kubernetes/k8s.io",
	)

	PRCmd.PersistentFlags().StringVar(
		&promoteOpts.reviewers,
		"reviewers",
		defaultReviewers,
		"the list of GitHub users or teams to assign to the PR",
	)

	PRCmd.PersistentFlags().BoolVarP(
		&promoteOpts.interactiveMode,
		"interactive",
		"i",
		false,
		"interactive mode, asks before every step",
	)

	for _, flagName := range []string{"tag", "fork"} {
		if err := PRCmd.MarkPersistentFlagRequired(flagName); err != nil {
			logrus.Error(errors.Wrapf(err, "marking tag %s as required", flagName))
		}
	}
}

func runPromote(opts *promoteOptions) error {
	// Check the cmd line opts
	if err := opts.Validate(); err != nil {
		return errors.Wrap(err, "checking command line options")
	}

	ctx := context.Background()

	// Validate options
	branchname := opts.project + "-" + opts.tags[0] + promotionBranchSuffix

	// Get the github org and repo from the fork slug
	userForkOrg, userForkRepo, err := git.ParseRepoSlug(opts.userFork)
	if err != nil {
		return errors.Wrap(err, "parsing user's fork")
	}
	if userForkRepo == "" {
		userForkRepo = k8sioRepo
	}

	// Check Environment
	gh := github.New()

	// Verify the repository is a fork of k8s.io
	if err = github.VerifyFork(
		branchname, userForkOrg, userForkRepo, git.DefaultGithubOrg, k8sioRepo,
	); err != nil {
		return errors.Wrapf(err, "while checking fork of %s/%s ", git.DefaultGithubOrg, k8sioRepo)
	}

	// Clone k8s.io
	repo, err := github.PrepareFork(branchname, git.DefaultGithubOrg, k8sioRepo, userForkOrg, userForkRepo)
	if err != nil {
		return errors.Wrap(err, "while preparing k/k8s.io fork")
	}

	defer func() {
		if mustRun(opts, "Clean fork directory?") {
			err = repo.Cleanup()
		} else {
			logrus.Infof("All modified files will be left untouched in %s", repo.Dir())
		}
	}()

	// Path to the promoter image list
	imagesListPath := filepath.Join(
		image.ProdRegistry,
		"images",
		filepath.Base(image.StagingRepoPrefix)+opts.project,
		"images.yaml",
	)

	// Read the current manifest to check later if new images come up
	oldlist := make([]byte, 0)

	// Run the promoter manifest grower
	if mustRun(opts, "Grow the manifests to add the new tags?") {
		if util.Exists(filepath.Join(repo.Dir(), imagesListPath)) {
			logrus.Debug("Reading the current image promoter manifest (image list)")
			oldlist, err = os.ReadFile(filepath.Join(repo.Dir(), imagesListPath))
			if err != nil {
				return errors.Wrap(err, "while reading the current promoter image list")
			}
		}

		for _, tag := range opts.tags {
			opt := reg.GrowManifestOptions{}
			if err := opt.Populate(
				filepath.Join(repo.Dir(), image.ProdRegistry),
				image.StagingRepoPrefix+opts.project, "", "", tag); err != nil {
				return errors.Wrapf(err, "populating image promoter options for tag %s", tag)
			}

			if err := opt.Validate(); err != nil {
				return errors.Wrapf(err, "validate promoter options for tag %s", tag)
			}

			logrus.Infof("Growing manifests with images matching tag %s", tag)
			if err := reg.GrowManifest(ctx, &opt); err != nil {
				return errors.Wrapf(err, "Growing manifest with tag %s", tag)
			}
		}
	}

	// Re-write the image list without the mock images
	rawImageList, err := image.NewManifestListFromFile(filepath.Join(repo.Dir(), imagesListPath))
	if err != nil {
		return errors.Wrap(err, "parsing the current manifest")
	}

	// Create a new imagelist to copy the non-mock images
	newImageList := &image.ManifestList{}

	// Copy all non mock-images:
	for _, imageData := range *rawImageList {
		if !strings.Contains(imageData.Name, "mock/") {
			*newImageList = append(*newImageList, imageData)
		}
	}

	// Write the modified manifest
	if err := newImageList.Write(filepath.Join(repo.Dir(), imagesListPath)); err != nil {
		return errors.Wrap(err, "while writing the promoter image list")
	}

	// Check if the image list was modified
	if len(oldlist) > 0 {
		logrus.Debug("Checking if the image list was modified")
		// read the newly modified manifest
		newlist, err := os.ReadFile(filepath.Join(repo.Dir(), imagesListPath))
		if err != nil {
			return errors.Wrap(err, "while reading the modified manifest images list")
		}

		// If the manifest was not modified, exit now
		if bytes.Equal(newlist, oldlist) {
			logrus.Info("No changes detected in the promoter images list, exiting without changes")
			return nil
		}
	}

	// add the modified manifest to staging
	logrus.Debugf("Adding %s to staging area", imagesListPath)
	if err := repo.Add(imagesListPath); err != nil {
		return errors.Wrap(err, "adding image manifest to staging area")
	}

	commitMessage := "Image promotion for " + opts.project + " " + strings.Join(opts.tags, " / ")
	if opts.project == image.StagingRepoSuffix {
		commitMessage = "releng: " + commitMessage
	}

	// Commit files
	logrus.Debug("Creating commit")
	if err := repo.UserCommit(commitMessage); err != nil {
		return errors.Wrapf(err, "Error creating commit in %s/%s", git.DefaultGithubOrg, k8sioRepo)
	}

	// Push to fork
	if mustRun(opts, fmt.Sprintf("Push changes to user's fork at %s/%s?", userForkOrg, userForkRepo)) {
		logrus.Infof("Pushing manifest changes to %s/%s", userForkOrg, userForkRepo)
		if err := repo.PushToRemote(github.UserForkName, branchname); err != nil {
			return errors.Wrapf(err, "pushing %s to %s/%s", github.UserForkName, userForkOrg, userForkRepo)
		}
	} else {
		// Exit if no push was made

		logrus.Infof("Exiting without creating a PR since changes were not pushed to %s/%s", userForkOrg, userForkRepo)
		return nil
	}

	// Create the Pull Request
	if mustRun(opts, "Create pull request?") {
		pr, err := gh.CreatePullRequest(
			git.DefaultGithubOrg, k8sioRepo, k8sioDefaultBranch,
			fmt.Sprintf("%s:%s", userForkOrg, branchname),
			commitMessage, generatePRBody(opts),
		)
		if err != nil {
			return errors.Wrap(err, "creating the pull request in k/k8s.io")
		}
		logrus.Infof(
			"Successfully created PR: %s%s/%s/pull/%d",
			github.GitHubURL, git.DefaultGithubOrg, k8sioRepo, pr.GetNumber(),
		)
	}

	// Success!
	return nil
}

// mustRun avoids running when a users chooses n in interactive mode
func mustRun(opts *promoteOptions, question string) bool {
	if !opts.interactiveMode {
		return true
	}
	_, success, err := util.Ask(fmt.Sprintf("%s (Y/n)", question), "y:Y:yes|n:N:no|y", 10)
	if err != nil {
		logrus.Error(err)
		if err.(util.UserInputError).IsCtrlC() {
			os.Exit(1)
		}
		return false
	}
	if success {
		return true
	}
	return false
}

// generatePRBody creates the body of the Image Promotion Pull Request
func generatePRBody(opts *promoteOptions) string {
	args := fmt.Sprintf("--fork %s", opts.userFork)
	if opts.interactiveMode {
		args += " --interactive"
	}

	if opts.project != defaultProject {
		args += " --project " + opts.project
	}

	if opts.reviewers != defaultReviewers {
		args += " --reviewers \"" + opts.reviewers + "\""
	}

	for _, tag := range opts.tags {
		args += " --tag " + tag
	}

	prBody := fmt.Sprintf("Image promotion for %s %s\n", opts.project, strings.Join(opts.tags, " / "))
	prBody += "This is an automated PR generated from `kpromo`\n"
	prBody += fmt.Sprintf("```\nkpromo pr %s\n```\n\n", args)
	prBody += fmt.Sprintf("/hold\ncc: %s\n", opts.reviewers)

	return prBody
}
