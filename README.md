# Bitrise Build Cache Add-On for Xcode

Enables the Bitrise Build Cache Add-On for Xcode.

## What this step does

- Configures the build using the Build Cache CLI so Bitrise Remote Build Cache is used.
- Ensures all subsequent Xcode builds in the workflow will read from the remote cache and push new entries.
- Adds an alias to `~/.zshrc` and `~/.bashrc`, making the CLI available in all subsequent steps.
    - From this point on, all calls to `xcodebuild` are wrapped so the underlying Xcode command has compilation caching enabled.
- Saves analytical data (command, duration, cache information, environment) to Bitrise. The data is available on the Build cache page: https://app.bitrise.io/build-cache

## Notes

- The alias persists only for subsequent steps in the current workflow run.
- Ensure your workflow uses `xcodebuild` (or compatible tooling) to benefit from the remote cache.
