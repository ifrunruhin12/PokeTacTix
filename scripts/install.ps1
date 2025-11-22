# PokeTacTix CLI Installation Script for Windows
# Run with: powershell -ExecutionPolicy Bypass -File install.ps1

$ErrorActionPreference = "Stop"

# Configuration
$REPO = "yourusername/poketactix"  # Update this with actual GitHub repo
$BINARY_NAME = "poketactix-cli.exe"
$INSTALL_DIR = "$env:LOCALAPPDATA\PokeTacTix"

Write-Host ""
Write-Host "╔════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║                                                        ║" -ForegroundColor Cyan
Write-Host "║         PokeTacTix CLI Installation Script            ║" -ForegroundColor Cyan
Write-Host "║                                                        ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Detect architecture
$ARCH = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
$PLATFORM = "windows"
$BINARY_FILE = "poketactix-cli-$PLATFORM-$ARCH.exe"

Write-Host "Detected platform: $PLATFORM/$ARCH" -ForegroundColor Blue
Write-Host ""

# Create installation directory
if (-not (Test-Path $INSTALL_DIR)) {
    Write-Host "Creating directory: $INSTALL_DIR" -ForegroundColor Blue
    New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
}

# Fetch latest release
Write-Host "Fetching latest release..." -ForegroundColor Blue
try {
    $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
    $LATEST_RELEASE = $response.tag_name
    Write-Host "Latest version: $LATEST_RELEASE" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "✗ Failed to fetch latest release" -ForegroundColor Red
    Write-Host "Please check your internet connection and repository name" -ForegroundColor Yellow
    exit 1
}

# Download binary
$DOWNLOAD_URL = "https://github.com/$REPO/releases/download/$LATEST_RELEASE/$BINARY_FILE"
$TARGET_PATH = Join-Path $INSTALL_DIR $BINARY_NAME

Write-Host "Downloading from: $DOWNLOAD_URL" -ForegroundColor Blue
try {
    Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $TARGET_PATH
    Write-Host "✓ Download complete" -ForegroundColor Green
} catch {
    Write-Host "✗ Download failed" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

# Add to PATH if not already present
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$INSTALL_DIR*") {
    Write-Host ""
    Write-Host "Adding $INSTALL_DIR to PATH..." -ForegroundColor Blue
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$currentPath;$INSTALL_DIR",
        "User"
    )
    Write-Host "✓ PATH updated" -ForegroundColor Green
    Write-Host ""
    Write-Host "⚠ Please restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
} else {
    Write-Host ""
    Write-Host "✓ Installation directory already in PATH" -ForegroundColor Green
}

# Display success message
Write-Host ""
Write-Host "╔════════════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║                                                        ║" -ForegroundColor Green
Write-Host "║         ✓ PokeTacTix CLI installed successfully!      ║" -ForegroundColor Green
Write-Host "║                                                        ║" -ForegroundColor Green
Write-Host "╚════════════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Host "Installation location: $TARGET_PATH" -ForegroundColor Cyan
Write-Host ""
Write-Host "Run the game with: $BINARY_NAME" -ForegroundColor Green
Write-Host "Check version with: $BINARY_NAME --version" -ForegroundColor Green
Write-Host ""
Write-Host "Note: You may need to restart your terminal or run:" -ForegroundColor Yellow
Write-Host "  `$env:Path = [System.Environment]::GetEnvironmentVariable('Path','User')" -ForegroundColor Yellow
Write-Host ""
