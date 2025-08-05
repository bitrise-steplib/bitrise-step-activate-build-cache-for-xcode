package main

import (
	"os"

	"github.com/bitrise-steplib/bitrise-step-enable-xcode-compilation-cache/step"
	"github.com/bitrise-steplib/steps-git-clone/gitclone/tracker"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/errorutil"
	"github.com/bitrise-io/go-utils/v2/exitcode"
	"github.com/bitrise-io/go-utils/v2/log"
)

func main() {
	exitCode := run()
	os.Exit(int(exitCode))
}

func run() exitcode.ExitCode {
	logger := log.NewLogger()
	envRepo := env.NewRepository()
	tracker := tracker.NewStepTracker(envRepo, logger)
	inputParser := stepconf.NewInputParser(envRepo)
	commandFactory := command.NewFactory(envRepo)

	stepInstance := step.New(
		logger,
		tracker,
		inputParser,
		commandFactory,
		envRepo,
	)

	err := stepInstance.Run()
	if err != nil {
		formattedMsg := errorutil.FormattedError(err)
		logger.Errorf("%s", formattedMsg)
		return exitcode.Failure
	}

	return exitcode.Success
}
