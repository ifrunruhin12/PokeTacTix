# Release Scripts

## Creating a Release

To create a new release of PokeTacTix CLI:

### 1. Use the tag-release script (Recommended)

```bash
./scripts/tag-release.sh
```

This interactive script will:
- Show you recent tags
- Prompt for a new version number (e.g., v1.0.0)
- Validate the version format
- Prompt for release notes
- Create an annotated git tag
- Push the tag to GitHub
- Display links to monitor the release

### 2. What happens next?

Once the tag is pushed, GitHub Actions automatically:
1. Builds binaries for all platforms (Windows, macOS, Linux)
2. Creates release archives (.zip for Windows, .tar.gz for Unix)
3. Generates SHA256 and MD5 checksums
4. Creates a GitHub Release
5. Uploads all binaries and archives

### 3. Monitor the release

The script will show you links to:
- GitHub Actions workflow (to monitor build progress)
- GitHub Release page (to view the published release)

## Manual Tag Creation

If you prefer to create tags manually:

```bash
# Create an annotated tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Push the tag to trigger the workflow
git push origin v1.0.0
```

## Version Format

Follow [Semantic Versioning](https://semver.org/):

- `v1.0.0` - Stable release
- `v1.0.0-beta.1` - Beta release
- `v1.0.0-rc.1` - Release candidate

Examples:
- `v1.0.0` - First stable release
- `v1.1.0` - New features added
- `v1.1.1` - Bug fixes
- `v2.0.0` - Breaking changes

## Troubleshooting

### Workflow doesn't run

- Make sure you pushed a tag starting with `v` (e.g., `v1.0.0`)
- Check that the workflow file exists: `.github/workflows/release-cli.yml`
- Verify GitHub Actions is enabled in your repository settings

### Build fails

- Check the Actions tab for error logs
- Ensure all tests pass locally: `go test ./...`
- Verify binaries build locally: `make build-cli-all`

### Tag already exists

```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin :refs/tags/v1.0.0

# Create new tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## Installation Scripts

The repository includes installation scripts for end users:

- `install.sh` - Unix/Linux/macOS installation
- `install.ps1` - Windows PowerShell installation

These scripts download the latest release from GitHub and install it to the appropriate location.

## See Also

- [RELEASE_PROCESS.md](../docs/RELEASE_PROCESS.md) - Detailed release process
- [BUILD_AND_RELEASE.md](../docs/BUILD_AND_RELEASE.md) - Build and release guide
- [CHANGELOG.md](../CHANGELOG.md) - Version history
