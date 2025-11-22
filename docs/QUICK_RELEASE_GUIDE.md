# Quick Release Guide

## How to Create a Release

### Step 1: Prepare for Release

Make sure everything is ready:

```bash
# Ensure all changes are committed
git status

# Run tests
go test ./...

# Test build
make build-cli
./bin/poketactix-cli --version
```

### Step 2: Create and Push a Tag

**Option A: Use the helper script (easiest)**

```bash
./scripts/tag-release.sh
```

Follow the prompts:
1. Enter version (e.g., `v1.0.0`)
2. Enter release notes
3. Confirm

**Option B: Manual tag creation**

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Step 3: Wait for GitHub Actions

The workflow will automatically:
- Build binaries for all platforms
- Create release archives
- Generate checksums
- Create GitHub Release
- Upload all assets

Monitor at: `https://github.com/ifrunruhin12/poketactix/actions`

### Step 4: Verify Release

Once complete, check:
- GitHub Releases page has the new release
- All binaries are uploaded
- Checksums are present
- Download and test one binary

## That's It!

The release is now live and users can install it with:

```bash
# macOS/Linux
curl -L https://github.com/YOUR_USERNAME/poketactix/raw/main/scripts/install.sh | bash

# Windows
powershell -ExecutionPolicy Bypass -Command "iwr https://github.com/YOUR_USERNAME/poketactix/raw/main/scripts/install.ps1 | iex"
```

## Why the Workflow Didn't Run Before

The GitHub Actions workflow only runs when you **push a tag** that starts with `v` (like `v1.0.0`).

It does NOT run on:
- Regular commits
- Branch pushes
- Pull requests

To trigger it, you MUST push a version tag:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0  # This triggers the workflow
```

## Common Issues

### "Workflow didn't run"
- Did you push a tag? (not just commit)
- Does the tag start with `v`?
- Is GitHub Actions enabled in repo settings?

### "Build failed"
- Check Actions tab for logs
- Test locally: `make build-cli-all`
- Ensure all tests pass

### "Tag already exists"
```bash
# Delete and recreate
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- `v1.0.0` - First stable release
- `v1.1.0` - New features (backwards compatible)
- `v1.1.1` - Bug fixes
- `v2.0.0` - Breaking changes
- `v1.0.0-beta.1` - Pre-release

## Next Steps

After releasing:
1. Update CHANGELOG.md
2. Announce the release
3. Test installation on different platforms
4. Monitor for issues

---

**Pro Tip**: Use `./scripts/tag-release.sh` - it handles everything for you!
