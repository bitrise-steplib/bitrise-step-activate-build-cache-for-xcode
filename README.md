# Build Cache for Xcode

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-activate-build-cache-for-xcode?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-activate-build-cache-for-xcode/releases)

Activates Bitrise Remote Build Cache add-on for subsequent Xcode builds in the workflow

<details>
<summary>Description</summary>

This Step enables Bitrise's Build Cache Add‑On for Xcode by configuring the environment with the Build Cache CLI.

After this Step runs, Xcode builds invoked via xcodebuild in subsequent workflow steps will automatically read from the remote cache and push new entries when applicable.

The Step adds an alias to ~/.zshrc and ~/.bashrc so the wrapper is available in all following steps; from that point all xcodebuild calls are wrapped to enable compilation caching.
Analytical data (command, duration, cache information, environment) is collected and sent to Bitrise and is available on the Build cache page: https://app.bitrise.io/build-cache

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `verbose` | Enable logging additional information for troubleshooting | required | `false` |
| `cache` | Whether the build uses the remote cache | required | `true` |
| `push` | Whether the build can not only read, but write new entries to the remote cache | required | `true` |
| `silent` | Whether Bitrise components should not log anything except the underlying xcodebuild output. Takes precedence over the 'Verbose logging' and 'Add timestamps' options. | required | `false` |
| `timestamps` | When enabled, the analytics wrapper adds timestamps to xcodebuild output log messages during the build. | required | `false` |
| `cache_skip_flags` | Skip passing cache flags to xcodebuild except the COMPILATION_CACHE_REMOTE_SERVICE_PATH.  Cache will have to be enabled manually in the Xcode project settings. More information can be found at the FAQ document: https://docs.bitrise.io/en/bitrise-build-cache/build-cache-for-xcode/xcode-compilation-cache-faq.html | required | `false` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-activate-build-cache-for-xcode/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-activate-build-cache-for-xcode/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
