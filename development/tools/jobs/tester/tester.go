package tester

import (
	"io/ioutil"
	"os"
	"testing"

	"fmt"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/test-infra/prow/config"
)

// Preset represents a existing presets
type Preset string

const (
	// PresetDindEnabled means docker-in-docker preset
	PresetDindEnabled Preset = "preset-dind-enabled"
	// PresetGcrPush means GCR push service account
	PresetGcrPush Preset = "preset-sa-gcr-push"
	// PresetDockerPushRepo means Docker repository
	PresetDockerPushRepo Preset = "preset-docker-push-repository"
	// PresetDockerPushRepoTestInfra means Docker repository test-infra images
	PresetDockerPushRepoTestInfra Preset = "preset-docker-push-repository-test-infra"
	// PresetDockerPushRepoIncubator means Decker repository incubator images
	PresetDockerPushRepoIncubator Preset = "preset-docker-push-repository-incubator"
	// PresetBuildPr means PR environment
	PresetBuildPr Preset = "preset-build-pr"
	// PresetBuildMaster means master environment
	PresetBuildMaster Preset = "preset-build-master"
	// PresetBuildRelease means release environment
	PresetBuildRelease Preset = "preset-build-release"
	// PresetBotGithubToken means github token
	PresetBotGithubToken Preset = "preset-bot-github-token"
	// PresetBotGithubSSH means github ssh
	PresetBotGithubSSH Preset = "preset-bot-github-ssh"

	// ImageGolangBuildpackLatest means Golang buildpack image
	ImageGolangBuildpackLatest = "eu.gcr.io/kyma-project/prow/test-infra/buildpack-golang:v20181119-afd3fbd"
	// ImageNodeBuildpackLatest means Node.js buildpack image
	ImageNodeBuildpackLatest = "eu.gcr.io/kyma-project/prow/test-infra/buildpack-node:v20181130-b28250b"
	// ImageNodeChromiumBuildpackLatest means Node.js + Chromium buildpack image
	ImageNodeChromiumBuildpackLatest = "eu.gcr.io/kyma-project/prow/test-infra/buildpack-node-chromium:v20181207-d46c013"
	// ImageBootstrapLatest means Bootstrap image
	ImageBootstrapLatest = "eu.gcr.io/kyma-project/prow/test-infra/bootstrap:v20181121-f3ea5ce"
	// ImageBoostrap001 represents version 0.0.1 of bootstrap image
	ImageBoostrap001 = "eu.gcr.io/kyma-project/prow/bootstrap:0.0.1"
	// ImageBootstrapHelm20181121 represents verion of bootstrap-helm image
	ImageBootstrapHelm20181121 = "eu.gcr.io/kyma-project/prow/test-infra/bootstrap-helm:v20181121-f2f12bc"

	// KymaProjectDir means kyma project dir
	KymaProjectDir = "/home/prow/go/src/github.com/kyma-project"

	// BuildScriptDir means build script directory
	BuildScriptDir = "/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/build.sh"
	// GovernanceScriptDir means governance script directory
	GovernanceScriptDir = "/home/prow/go/src/github.com/kyma-project/test-infra/prow/scripts/governance.sh"
)

// ReadJobConfig reads job configuration from file
func ReadJobConfig(fileName string) (config.JobConfig, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return config.JobConfig{}, errors.Wrapf(err, "while opening file [%s]", fileName)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return config.JobConfig{}, errors.Wrapf(err, "while reading file [%s]", fileName)
	}
	jobConfig := config.JobConfig{}
	if err = yaml.Unmarshal(b, &jobConfig); err != nil {
		return config.JobConfig{}, errors.Wrapf(err, "while unmarshalling file [%s]", fileName)
	}
	return jobConfig, nil
}

// FindPresubmitJobByName finds presubmit job by name from provided jobs list
func FindPresubmitJobByName(jobs []config.Presubmit, name string) *config.Presubmit {
	for _, job := range jobs {
		if job.Name == name {
			return &job
		}
	}

	return nil
}

// FindPostsubmitJobByName finds postsubmit job by name from provided jobs list
func FindPostsubmitJobByName(jobs []config.Postsubmit, name string) *config.Postsubmit {
	for _, job := range jobs {
		if job.Name == name {
			return &job
		}
	}

	return nil
}

// FindPeriodicJobByName finds periodic job by name from provided jobs list
func FindPeriodicJobByName(jobs []config.Periodic, name string) *config.Periodic {
	for _, job := range jobs {
		if job.Name == name {
			return &job
		}
	}

	return nil
}

// AssertThatHasExtraRef checks if UtilityConfig has repository passed in argument defined
func AssertThatHasExtraRef(t *testing.T, in config.UtilityConfig, repository string) {
	for _, curr := range in.ExtraRefs {
		if curr.PathAlias == fmt.Sprintf("github.com/kyma-project/%s", repository) &&
			curr.Org == "kyma-project" &&
			curr.Repo == repository &&
			curr.BaseRef == "master" {
			return
		}
	}
	assert.FailNow(t, fmt.Sprintf("Job has not configured %s as a extra ref", repository))
}

// AssertThatHasExtraRefTestInfra checks if UtilityConfig has test-infra repository defined
func AssertThatHasExtraRefTestInfra(t *testing.T, in config.UtilityConfig) {
	AssertThatHasExtraRef(t, in, "test-infra")
}

// AssertThatHasExtraRefs checks if UtilityConfig has repositories passed in argument defined
func AssertThatHasExtraRefs(t *testing.T, in config.UtilityConfig, repositories []string) {
	for _, repository := range repositories {
		for _, curr := range in.ExtraRefs {
			if curr.PathAlias == fmt.Sprintf("github.com/kyma-project/%s", repository) &&
				curr.Org == "kyma-project" &&
				curr.Repo == repository &&
				curr.BaseRef == "master" {
				return
			}
		}
		assert.FailNow(t, fmt.Sprintf("Job has not configured %s as a extra ref", repository))
	}
}

// AssertThatHasPresets checks if JobBase has expected labels
func AssertThatHasPresets(t *testing.T, in config.JobBase, expected ...Preset) {
	for _, p := range expected {
		assert.Equal(t, "true", in.Labels[string(p)], "missing preset [%s]", p)
	}
}

// AssertThatJobRunIfChanged checks if Presubmit has run_if_changed parameter
func AssertThatJobRunIfChanged(t *testing.T, p config.Presubmit, changedFile string) {
	sl := []config.Presubmit{p}
	require.NoError(t, config.SetPresubmitRegexes(sl))
	assert.True(t, sl[0].RunsAgainstChanges([]string{changedFile}), "missed change [%s]", changedFile)
}

// AssertThatHasCommand checks if job has
func AssertThatHasCommand(t *testing.T, command []string) {
	assert.Equal(t, []string{BuildScriptDir}, command)
}
