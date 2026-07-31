#Requires -Version 5.1
# Installs ctx to C:\Users\<username>\tools and adds it to the user PATH.

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$installDir = Join-Path $env:USERPROFILE "tools"
$exePath = Join-Path $installDir "ctx.exe"
$mingwDir = Join-Path $repoRoot ".tools\mingw64"
$mingwBin = Join-Path $mingwDir "bin"
$gccPath = Join-Path $mingwBin "gcc.exe"

function Find-GCC {
    $inPath = (Get-Command gcc.exe -ErrorAction SilentlyContinue)
    if ($inPath) {
        return $inPath.Source
    }
    if (Test-Path $gccPath) {
        return $gccPath
    }
    return $null
}

# Ensure install directory exists
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
    Write-Host "Created $installDir"
}

# Locate or download a C compiler for CGO
$cc = Find-GCC
if (!$cc) {
    Write-Host "No gcc found. Downloading portable MinGW..."
    if (!(Test-Path $mingwDir)) {
        New-Item -ItemType Directory -Path $mingwDir -Force | Out-Null
    }
    $zip = Join-Path $repoRoot ".tools\mingw.zip"
    $url = "https://github.com/brechtsanders/winlibs_mingw/releases/download/16.1.0posix-14.0.0-msvcrt-r4/winlibs-x86_64-posix-seh-gcc-16.1.0-mingw-w64msvcrt-14.0.0-r4.zip"
    Invoke-WebRequest -Uri $url -OutFile $zip -UseBasicParsing
    Expand-Archive -Path $zip -DestinationPath (Join-Path $repoRoot ".tools") -Force
    Remove-Item $zip
    $cc = $gccPath
}

$env:CC = $cc
$env:CXX = Join-Path (Split-Path $cc) "g++.exe"

# Build the binary
Write-Host "Building ctx with CC=$env:CC ..."
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
