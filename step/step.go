package step

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/bitrise-io/go-steputils/tools"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"

	"github.com/bitrise-steplib/steps-git-clone/gitclone"
	"github.com/bitrise-steplib/steps-git-clone/gitclone/bitriseapi"
	"github.com/bitrise-steplib/steps-git-clone/gitclone/tracker"
	"github.com/bitrise-steplib/steps-git-clone/transport"
)

const (
	cacheToolsGitURL = "git@github.com:bitrise-io/xcode-cache-tools.git"
	toolsPath        = "xcode-cache-tools"
	proxyPath        = "/cmd/proxy"
	proxyBinaryName  = "xcode_cache_proxy"

	infoStartedProxy = "Started detached process with PID %d"

	errFailedToExportEnable         = "failed to export BITRISE_XCODE_COMPILATION_CACHE_ENABLED: %s"
	errFailedToExportCompCacheArgs  = "failed to export BITRISE_XCODE_COMPILATION_CACHE_ARGS: %s"
	errFailedToExportAdditionalArgs = "failed to export BITRISE_XCODE_ADDITIONAL_ARGS: %s"
	errFailedToExportProxyPID       = "failed to export BITRISE_XCODE_COMPILATION_CACHE_PROXY_PID: %s"
	errFailedToStartProxy           = "failed to start proxy: %s"
	errFailedToCompileProxy         = "failed to compile proxy: %s"
	errProxyBuildFailed             = "proxy build failed: %s"
)

type Input struct {
	Verbose bool `env:"verbose,required"`

	GitHTTPUsername string `env:"git_http_username"`
	GitHTTPPassword string `env:"git_http_password"`
	BuildURL        string `env:"build_url"`
	BuildAPIToken   string `env:"build_api_token"`
}

type Step struct {
	logger         log.Logger
	tracker        tracker.StepTracker
	inputParser    stepconf.InputParser
	commandFactory command.Factory
	envRepo        env.Repository
}

func New(
	logger log.Logger,
	tracker tracker.StepTracker,
	inputParser stepconf.InputParser,
	commandFactory command.Factory,
	envRepo env.Repository,
) Step {
	return Step{
		logger:         logger,
		tracker:        tracker,
		inputParser:    inputParser,
		commandFactory: commandFactory,
		envRepo:        envRepo,
	}
}

func (step Step) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return err
	}
	stepconf.Print(input)
	step.logger.Println()
	step.logger.EnableDebugLog(input.Verbose)

	step.clone(input, cacheToolsGitURL, toolsPath)

	pathToBinary, err := step.compileProxy(toolsPath, proxyPath)
	if err != nil {
		return fmt.Errorf(errFailedToCompileProxy, err)
	}

	err = step.startProxy(pathToBinary)
	if err != nil {
		return fmt.Errorf(errFailedToStartProxy, err)
	}

	if err := tools.ExportEnvironmentWithEnvman("BITRISE_XCODE_COMPILATION_CACHE_ENABLED", "true"); err != nil {
		return fmt.Errorf(errFailedToExportEnable, err)
	}

	var compilationCacheArgs string = ""
	compilationCacheArgs += "SWIFT_ENABLE_EXPLICIT_MODULES=YES "
	compilationCacheArgs += "COMPILATION_CACHE_ENABLE_CACHING=YES "
	compilationCacheArgs += "SWIFT_ENABLE_COMPILE_CACHE=1 "
	compilationCacheArgs += "COMPILATION_CACHE_REMOTE_SERVICE_PATH=/tmp/llvmproxy.sock "
	compilationCacheArgs += "COMPILATION_CACHE_ENABLE_PLUGIN=1"

	if err := tools.ExportEnvironmentWithEnvman("BITRISE_XCODE_COMPILATION_CACHE_ARGS", compilationCacheArgs); err != nil {
		return fmt.Errorf(errFailedToExportCompCacheArgs, err)
	}

	var additionalXcodeArgs string = os.Getenv("BITRISE_XCODE_ADDITIONAL_ARGS")
	if additionalXcodeArgs != "" {
		additionalXcodeArgs += " "
	}
	additionalXcodeArgs += compilationCacheArgs

	if err := tools.ExportEnvironmentWithEnvman("BITRISE_XCODE_ADDITIONAL_ARGS", additionalXcodeArgs); err != nil {
		return fmt.Errorf(errFailedToExportAdditionalArgs, err)
	}

	return nil
}

func (s Step) clone(input Input, url string, dir string) (gitclone.CheckoutStateResult, error) {
	if err := transport.Setup(transport.Config{
		URL:          url,
		HTTPUsername: input.GitHTTPUsername,
		HTTPPassword: input.GitHTTPPassword,
	}); err != nil {
		return gitclone.CheckoutStateResult{}, err
	}

	gitCloneCfg := convertConfig(url, dir)
	patchSource := bitriseapi.NewPatchSource(input.BuildURL, input.BuildAPIToken)
	mergeRefChecker := bitriseapi.NewMergeRefChecker(input.BuildURL, input.BuildAPIToken, retry.NewHTTPClient(), s.logger, s.tracker)
	cloner := gitclone.NewGitCloner(s.logger, s.tracker, s.commandFactory, patchSource, mergeRefChecker, false)
	return cloner.CheckoutState(gitCloneCfg)
}

func convertConfig(url string, dir string) gitclone.Config {
	return gitclone.Config{
		CloneIntoDir:         dir,
		CloneDepth:           1,
		SubmoduleUpdateDepth: 1,
		RepositoryURL:        url,
		Branch:               "main",
	}
}

func (s Step) compileProxy(projectPath string, modulePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("go", "build", "-o", proxyBinaryName, filepath.Join(cwd, projectPath, modulePath))
	cmd.Dir = filepath.Join(cwd, projectPath)

	// (TODO) Connect command output to your program's stdout/stderr for debugging
	cmd.Stdout = os.Stdout // or os.Stdout
	cmd.Stderr = os.Stderr // or os.Stderr

	// Run the command
	err = cmd.Run()
	if err != nil {
		s.logger.Errorf(errProxyBuildFailed, err.Error())
		return "", err
	}

	s.logger.Donef("Build succeeded!")

	return filepath.Join(cwd, projectPath, proxyBinaryName), nil
}

func (s Step) startProxy(pathToBinary string) error {
	cmd := exec.Command(pathToBinary)

	// Detach from terminal and run in its own session
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	// Redirect standard file descriptors to /dev/null
	devNull, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if err != nil {
		s.logger.Errorf("Error opening /dev/null:", err)
		return err
	}
	defer devNull.Close()

	cmd.Stdin = devNull
	cmd.Stdout = devNull
	cmd.Stderr = devNull

	env := os.Environ()
	env = append(
		env,
		"REMOTE_CACHE_TOKEN="+s.envRepo.Get("REMOTE_CACHE_TOKEN"),
		"APP_ID="+s.envRepo.Get("APP_ID"),
		"ORG_ID="+s.envRepo.Get("ORG_ID"),
		"LLVM_PROXY_TARGET_ADDR="+s.envRepo.Get("LLVM_PROXY_TARGET_ADDR"),
	)
	cmd.Env = env

	// Start asynchronously
	if err := cmd.Start(); err != nil {
		return fmt.Errorf(errFailedToStartProxy, err)
	}

	s.logger.Infof(infoStartedProxy, cmd.Process.Pid)
	if err := tools.ExportEnvironmentWithEnvman("BITRISE_XCODE_COMPILATION_CACHE_PROXY_PID", fmt.Sprint(cmd.Process.Pid)); err != nil {
		return fmt.Errorf(errFailedToExportProxyPID, err)
	}

	return nil
}
