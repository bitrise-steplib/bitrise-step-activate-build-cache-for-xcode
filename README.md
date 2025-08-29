# Bitrise Build Cache Add-On for Xcode

Enables the Bitrise Build Cache Add-On for Xcode.

## What this step does

- Configures the build using the Build Cache CLI so Bitrise Remote Build Cache is used.
- Ensures all subsequent Xcode builds in the workflow will read from the remote cache and push new entries.
- Adds a PATH entry to make the CLI available in all subsequent steps.
    - From this point on, all calls to `xcodebuild` are wrapped so the underlying Xcode command has compilation caching enabled.
    - The wrapper script in in `~/.bitrise-xcelerate/bin` and it's added as `xcodebuild` to the path in `~/.zshrc` and `~/.bashrc`,
    - The PATH is also exported through `envman` so it's available in subsequent steps.
- Saves analytical data (command, duration, cache information, environment) to Bitrise. The data is available on the Build cache page: https://app.bitrise.io/build-cache

## Notes

- Xcode 26 Beta 4 is required to use the remote cache.
- The wrapper in the PATH persists only for subsequent steps in the current workflow run.
- Ensure your workflow uses `xcodebuild` (or compatible tooling) to benefit from the remote cache.
