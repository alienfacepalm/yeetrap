# YeeTrap Setup Script for Windows

Write-Host "========================================"
Write-Host "   YeeTrap Setup Script"
Write-Host "========================================"
Write-Host ""

# Check Go
Write-Host "Checking for Go..." -ForegroundColor Yellow
$goCheck = Get-Command go -ErrorAction SilentlyContinue
if ($goCheck) {
    $goVersion = go version
    Write-Host "OK: $goVersion" -ForegroundColor Green
} else {
    Write-Host "ERROR: Go is not installed!" -ForegroundColor Red
    Write-Host "Install from: https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

# Check yt-dlp
Write-Host "Checking for yt-dlp..." -ForegroundColor Yellow
$ytdlpCheck = Get-Command yt-dlp -ErrorAction SilentlyContinue
if ($ytdlpCheck) {
    $ytdlpVersion = yt-dlp --version
    Write-Host "OK: yt-dlp v$ytdlpVersion" -ForegroundColor Green
} else {
    Write-Host "WARNING: yt-dlp is not installed!" -ForegroundColor Yellow
    Write-Host "Install with: winget install yt-dlp" -ForegroundColor Cyan
}

# Create config dir
Write-Host ""
Write-Host "Creating config directory..." -ForegroundColor Yellow
$configDir = "$env:USERPROFILE\.yeetrap"
if (!(Test-Path $configDir)) {
    New-Item -ItemType Directory -Force -Path $configDir | Out-Null
    Write-Host "OK: Created $configDir" -ForegroundColor Green
} else {
    Write-Host "OK: Directory exists" -ForegroundColor Green
}

# Check credentials
Write-Host ""
Write-Host "Checking for credentials..." -ForegroundColor Yellow
$credPath = Join-Path $configDir "credentials.json"
if (Test-Path $credPath) {
    Write-Host "OK: credentials.json found" -ForegroundColor Green
} else {
    Write-Host "WARNING: credentials.json not found!" -ForegroundColor Yellow
    Write-Host "Place it at: $credPath" -ForegroundColor Cyan
}

# Build
Write-Host ""
Write-Host "Building YeeTrap..." -ForegroundColor Yellow
go build -o yeetrap.exe
if ($LASTEXITCODE -eq 0) {
    Write-Host "OK: Build successful" -ForegroundColor Green
} else {
    Write-Host "ERROR: Build failed" -ForegroundColor Red
    exit 1
}

# Test
Write-Host ""
Write-Host "Testing..." -ForegroundColor Yellow
.\yeetrap.exe version

# Summary
Write-Host ""
Write-Host "========================================"
Write-Host "   Setup Complete!"
Write-Host "========================================"
Write-Host ""
Write-Host "Next steps:"
Write-Host "1. Run: .\yeetrap.exe auth"
Write-Host "2. Run: .\yeetrap.exe list"
Write-Host "3. Run: .\yeetrap.exe download"
Write-Host ""
