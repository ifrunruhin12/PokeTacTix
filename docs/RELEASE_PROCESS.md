# Release Process for PokeTacTix CLI

This document describes the process for creating and publishing new releases of PokeTacTix CLI.

## Prerequisites

- Git with commit access to the repository
- Go 1.21 or later installed
- GitHub CLI (`gh`) installed (optional, for easier releases)
- Write access to the GitHub repository

## Release Checklist

### 1. Pre-Release Preparation

- [ ] Ensure all tests pass: `go test ./...`
- [ ] Update version in relevant files if needed
- [ ] Update CHANGELOG.md with release notes
- [ ] Test CLI on all target platforms (or at least your platform)
- [ ] Review and update CLI_README.md if needed
- [ ] Commit all changes

### 2. Create Release Tag

Choose a version number following [Semantic Versioning](https://semver.org/):

- **Major** (X.0.0): Breaking changes
- **Minor** (0.X.0): New features, backwards compatible
- **Patch** (0.0.X): Bug fixes, backwards compatible

**Recommended: Use the tag-release script**

```bash
./scripts/tag-release.sh
```

The script will guide you through creating and pushing a tag.

**Alternative: Manual tag creation**

```bash
VERSION="v1.0.0"  # Change this to your version
git tag -a $VERSION -m "Release $VERSION"
git push origin $VERSION
```

### 3. Automated Build and Release

Once you push a tag, GitHub Actions will automatically:

1. Build binaries for all platforms
2. Create release archives
3. Generate checksums
4. Create a GitHub release
5. Upload all assets

Monitor the workflow at: `https://github.com/ifrunruhin12/poketactix/actions`

### 4. Manual Release (Alternative)

If you prefer to create releases manually:

```bash
# Build all binaries
make build-cli-all

# Create release package
./scripts/create-release.sh

# Follow the prompts and instructions
```

Then manually create a GitHub release:

```bash
# Using GitHub CLI
gh release create $VERSION \
  release/$VERSION/*.zip \
  release/$VERSION/*.tar.gz \
  release/$VERSION/SHA256SUMS \
  release/$VERSION/MD5SUMS \
  --title "PokeTacTix CLI $VERSION" \
  --notes-file release/$VERSION/RELEASE_NOTES.md
```

Or use the GitHub web interface:

1. Go to https://github.com/yourusername/poketactix/releases/new
2. Choose the tag you created
3. Upload the files from `release/$VERSION/`
4. Add release notes
5. Publish release

## Post-Release Tasks

### 1. Verify Release

- [ ] Download binaries from GitHub release
- [ ] Test installation on at least one platform
- [ ] Verify checksums match
- [ ] Test `--version` flag shows correct version

### 2. Update Documentation

- [ ] Update main README.md if needed
- [ ] Announce release in GitHub Discussions
- [ ] Update any external documentation

### 3. Social/Community

- [ ] Post release announcement (if applicable)
- [ ] Update any package managers (if applicable)
- [ ] Notify users of major changes

## Release Types

### Stable Release (v1.0.0)

Full production release with complete features and testing.

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Pre-Release (v1.0.0-beta.1)

Beta or release candidate for testing.

```bash
git tag -a v1.0.0-beta.1 -m "Beta release v1.0.0-beta.1"
git push origin v1.0.0-beta.1
```

Mark as "pre-release" in GitHub when creating the release.

### Hotfix Release (v1.0.1)

Quick patch for critical bugs.

```bash
# Create hotfix branch from tag
git checkout -b hotfix/v1.0.1 v1.0.0

# Make fixes and commit
git commit -am "Fix critical bug"

# Tag and push
git tag -a v1.0.1 -m "Hotfix v1.0.1"
git push origin v1.0.1
git push origin hotfix/v1.0.1

# Merge back to main
git checkout main
git merge hotfix/v1.0.1
git push origin main
```

## Troubleshooting

### Build Fails

```bash
# Clean and rebuild
rm -rf bin/
make build-cli-all
```

### GitHub Actions Fails

1. Check the Actions tab for error logs
2. Common issues:
   - Missing dependencies
   - Permission issues
   - Invalid tag format

### Checksums Don't Match

```bash
# Regenerate checksums
cd release/$VERSION
sha256sum poketactix-cli-* > SHA256SUMS
md5sum poketactix-cli-* > MD5SUMS
```

## Version Numbering Examples

- `v1.0.0` - First stable release
- `v1.1.0` - Added new features
- `v1.1.1` - Bug fixes
- `v2.0.0` - Breaking changes
- `v1.0.0-beta.1` - Beta release
- `v1.0.0-rc.1` - Release candidate

## Release Notes Template

```markdown
## What's New in v1.X.X

### New Features

- Feature 1 description
- Feature 2 description

### Improvements

- Improvement 1
- Improvement 2

### Bug Fixes

- Fix 1
- Fix 2

### Breaking Changes

- Change 1 (if any)

### Known Issues

- Issue 1 (if any)

## Installation

[Standard installation instructions]

## Checksums

[SHA256 checksums]
```

## Emergency Rollback

If a release has critical issues:

```bash
# Delete the release from GitHub
gh release delete $VERSION

# Delete the tag
git tag -d $VERSION
git push origin :refs/tags/$VERSION

# Fix issues and create new release
```

## Support

For questions about the release process:

- Open an issue on GitHub
- Contact the maintainers
- Check GitHub Actions documentation

---

**Remember**: Always test releases before publishing, and keep CHANGELOG.md up to date!
