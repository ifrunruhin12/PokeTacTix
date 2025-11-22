# Platform Compatibility Test Script for Windows
# Tests cross-platform features of PokeTacTix CLI

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "PokeTacTix CLI Platform Compatibility Test" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Detect platform
$OS = [System.Environment]::OSVersion.Platform
$OSVersion = [System.Environment]::OSVersion.VersionString
$Arch = [System.Environment]::GetEnvironmentVariable("PROCESSOR_ARCHITECTURE")

Write-Host "Platform Information:" -ForegroundColor Yellow
Write-Host "  OS: $OS"
Write-Host "  Version: $OSVersion"
Write-Host "  Architecture: $Arch"
Write-Host "  PowerShell Version: $($PSVersionTable.PSVersion)"
Write-Host ""

# Detect terminal
Write-Host "Terminal Information:" -ForegroundColor Yellow
Write-Host "  WT_SESSION: $env:WT_SESSION"
Write-Host "  ConEmuANSI: $env:ConEmuANSI"
Write-Host "  TERM: $env:TERM"
Write-Host "  COLORTERM: $env:COLORTERM"

# Check if running in Windows Terminal
if ($env:WT_SESSION) {
    Write-Host "  ✓ Running in Windows Terminal" -ForegroundColor Green
} elseif ($env:ConEmuANSI -eq "ON") {
    Write-Host "  ✓ Running in ConEmu with ANSI support" -ForegroundColor Green
} else {
    Write-Host "  ℹ Not running in Windows Terminal (colors may be limited)" -ForegroundColor Yellow
}
Write-Host ""

# Test color support
Write-Host "Color Support Tests:" -ForegroundColor Yellow
Write-Host "  NO_COLOR: $env:NO_COLOR"
Write-Host "  CLICOLOR: $env:CLICOLOR"
Write-Host "  CLICOLOR_FORCE: $env:CLICOLOR_FORCE"

# Test ANSI colors (works in PowerShell 5.1+ and Windows Terminal)
if (-not $env:NO_COLOR) {
    Write-Host "  " -NoNewline
    Write-Host "Red " -ForegroundColor Red -NoNewline
    Write-Host "Green " -ForegroundColor Green -NoNewline
    Write-Host "Yellow " -ForegroundColor Yellow -NoNewline
    Write-Host "Blue " -ForegroundColor Blue -NoNewline
    Write-Host "- Colors working!" -ForegroundColor White
} else {
    Write-Host "  Colors disabled (NO_COLOR set)"
}
Write-Host ""

# Test file paths
Write-Host "File Path Tests:" -ForegroundColor Yellow
$HomeDir = $env:USERPROFILE
Write-Host "  Home Directory: $HomeDir"

$SaveDir = Join-Path $HomeDir ".poketactix"
Write-Host "  Save Directory: $SaveDir"

if (Test-Path $SaveDir) {
    Write-Host "  ✓ Save directory exists" -ForegroundColor Green
    
    # List contents
    $Files = Get-ChildItem $SaveDir -ErrorAction SilentlyContinue
    if ($Files) {
        Write-Host "  Files in save directory:"
        foreach ($File in $Files) {
            $Size = if ($File.PSIsContainer) { "DIR" } else { "{0:N2} KB" -f ($File.Length / 1KB) }
            Write-Host "    - $($File.Name) ($Size)"
        }
    }
} else {
    Write-Host "  ℹ Save directory does not exist yet (will be created on first run)" -ForegroundColor Yellow
}
Write-Host ""

# Test Go environment
Write-Host "Go Environment:" -ForegroundColor Yellow
$GoPath = Get-Command go -ErrorAction SilentlyContinue
if ($GoPath) {
    $GoVersion = go version
    Write-Host "  ✓ Go installed: $GoVersion" -ForegroundColor Green
    
    # Test build for Windows
    Write-Host ""
    Write-Host "Testing build for Windows..." -ForegroundColor Yellow
    
    $BuildOutput = "bin\poketactix-cli-test.exe"
    
    try {
        $BuildResult = go build -o $BuildOutput cmd\cli\main.go 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✓ Build successful: $BuildOutput" -ForegroundColor Green
            
            # Check binary size
            if (Test-Path $BuildOutput) {
                $Size = (Get-Item $BuildOutput).Length
                $SizeMB = "{0:N2} MB" -f ($Size / 1MB)
                Write-Host "  Binary size: $SizeMB"
                
                # Clean up test binary
                Remove-Item $BuildOutput -ErrorAction SilentlyContinue
            }
        } else {
            Write-Host "  ✗ Build failed" -ForegroundColor Red
            Write-Host $BuildResult
            exit 1
        }
    } catch {
        Write-Host "  ✗ Build failed: $_" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "  ✗ Go not installed" -ForegroundColor Red
    Write-Host "  Download from: https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}
Write-Host ""

# Run Go tests
Write-Host "Running Go Tests:" -ForegroundColor Yellow

Write-Host "  Testing UI package..." -ForegroundColor Cyan
$TestResult = go test -v ./internal/cli/ui/... 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ UI tests passed" -ForegroundColor Green
} else {
    Write-Host "  ✗ UI tests failed" -ForegroundColor Red
    Write-Host $TestResult
}

Write-Host ""
Write-Host "  Testing storage package..." -ForegroundColor Cyan
$TestResult = go test -v ./internal/cli/storage/... 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Storage tests passed" -ForegroundColor Green
} else {
    Write-Host "  ✗ Storage tests failed" -ForegroundColor Red
    Write-Host $TestResult
}

Write-Host ""
Write-Host "  Testing commands package..." -ForegroundColor Cyan
$TestResult = go test -v ./internal/cli/commands/... 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "  ✓ Commands tests passed" -ForegroundColor Green
} else {
    Write-Host "  ✗ Commands tests failed" -ForegroundColor Red
    Write-Host $TestResult
}

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Platform Compatibility Test Complete" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Windows-specific notes
Write-Host "Windows Notes:" -ForegroundColor Yellow
Write-Host "  - Save directory: $SaveDir"
Write-Host "  - Recommended terminals: Windows Terminal, PowerShell"
Write-Host "  - cmd.exe has limited color support"
Write-Host "  - For best experience, use Windows Terminal from Microsoft Store"
Write-Host ""
Write-Host "For detailed platform testing guide, see: docs\cli-platform-testing.md" -ForegroundColor Cyan
