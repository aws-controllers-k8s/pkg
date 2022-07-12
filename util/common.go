// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package util

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	sdkRepoURL             = "https://github.com/aws/aws-sdk-go"
	defaultGitCloneTimeout = 180 * time.Second
	defaultGitFetchTimeout = 30 * time.Second
)

func ContextWithSigterm(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	signalCh := make(chan os.Signal, 1)

	// recreate the context.CancelFunc
	cancelFunc := func() {
		signal.Stop(signalCh)
		cancel()
	}

	// notify on SIGINT or SIGTERM
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-signalCh:
			cancel()
		case <-ctx.Done():
		}
	}()

	return ctx, cancelFunc
}

// EnsureDir makes sure that a supplied directory exists and
// returns whether the directory already existed.
func EnsureDir(fp string) (bool, error) {
	fi, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return false, os.MkdirAll(fp, os.ModePerm)
		}
		return false, err
	}
	if !fi.IsDir() {
		return false, fmt.Errorf("expected %s to be a directory", fp)
	}
	if !IsDirWriteable(fp) {
		return true, fmt.Errorf("%s is not a writeable directory", fp)
	}

	return true, nil
}

// IsDirWriteable returns true if the supplied directory path is writeable,
// false otherwise
func IsDirWriteable(fp string) bool {
	testPath := filepath.Join(fp, "test")
	f, err := os.Create(testPath)
	if err != nil {
		return false
	}
	f.Close()
	os.Remove(testPath)
	return true
}

// EnsureSDKRepo ensures that we have a git clone'd copy of the aws-sdk-go
// repository, which we use model JSON files from. Upon successful return of
// this function, the sdkDir global variable will be set to the directory where
// the aws-sdk-go is found. It will also optionally fetch all the remote tags
// and checkout the given tag.
func EnsureSDKRepo(
	ctx context.Context,
	cacheDir string,
	// A boolean instructing EnsureSDKRepo whether to fetch the remote tags from
	// the upstream repository
	fetchTags bool,
	awsSDKGoVersion string,
) error {
	var err error
	srcPath := filepath.Join(cacheDir, "src")
	if err = os.MkdirAll(srcPath, os.ModePerm); err != nil {
		return err
	}

	// Clone repository if it doen't exist
	sdkDir := filepath.Join(srcPath, "aws-sdk-go")
	if _, err := os.Stat(sdkDir); os.IsNotExist(err) {

		ctx, cancel := context.WithTimeout(ctx, defaultGitCloneTimeout)
		defer cancel()
		err = CloneRepository(ctx, sdkDir, sdkRepoURL)
		if err != nil {
			return fmt.Errorf("canot clone repository: %v", err)
		}
	}

	// Fetch all tags
	if fetchTags {
		ctx, cancel := context.WithTimeout(ctx, defaultGitFetchTimeout)
		defer cancel()
		err = FetchRepositoryTags(ctx, sdkDir)
		if err != nil {
			return fmt.Errorf("cannot fetch tags: %v", err)
		}
	}

	// get sdkVersion and ensure its prefix
	sdkVersion := GetSDKVersion(awsSDKGoVersion)
	if err != nil {
		return err
	}
	sdkVersion = EnsureSemverPrefix(sdkVersion)

	repo, err := LoadRepository(sdkDir)
	if err != nil {
		return fmt.Errorf("cannot read local repository: %v", err)
	}

	// Now checkout the local repository.
	err = CheckoutRepositoryTag(repo, sdkVersion)
	if err != nil {
		return fmt.Errorf("cannot checkout tag: %v", err)
	}

	return err
}

// EnsureSemverPrefix takes a semver string and tries to append the 'v'
// prefix if it's missing.
func EnsureSemverPrefix(s string) string {
	// trim all leading 'v' runes (characters)
	s = strings.TrimLeft(s, "v")
	return fmt.Sprintf("v%s", s)
}

// GetSDKVersion returns the github.com/aws/aws-sdk-go version to use
// from the --aws-sdk-go-version flag.
func GetSDKVersion(awsSDKGoVersion string) string {
	return awsSDKGoVersion
}
