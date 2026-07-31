#Requires -Version 5.1
# Installs ctx to C:\Users\<username>\tools and adds it to the user PATH.

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$installDir = Join-Path $env:USERPROFILE "tools"
$exePath = Join-Path $installDir "ctx.exe"

# Ensure install directory exists
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
    Write-Host "Created $installDir"
}

# Build the binary
Write-Host "Building ctx..."
Set-Location $repoRoot
go build -o ctx.exe ./cmd/ctx
if (!(Test-Path (Join-Path $repoRoot "ctx.exe"))) {
    throw "Build failed: ctx.exe not found"
}

# Install the binary
Write-Host "Installing ctx.exe to $installDir..."
Copy-Item -Path (Join-Path $repoRoot "ctx.exe") -Destination $exePath -Force

# Add to user PATH if not already present
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$installDir*") {
    Write-Host "Adding $installDir to user PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$installDir", "User")
} else {
    Write-Host "$installDir is already in user PATH"
}

Write-Host ""
Write-Host "Installation complete." -ForegroundColor Green
Write-Host "Restart your terminal and run: ctx --help"
