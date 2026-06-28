param(
    [Parameter(Position = 0)]
    [string]$Version = "latest"
)

$ErrorActionPreference = 'Stop'

# cnext - Windows installer
# Usage: irm https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.ps1 | iex
# Or:    irm .../install.ps1 -OutFile install.ps1; .\install.ps1 v1.0.0

$BinaryName = "cnext"
$Repo = "neko233-com/cnext"
$InstallDir = Join-Path $env:LOCALAPPDATA $BinaryName
$Asset = "${BinaryName}-windows-amd64.exe"

function Get-NormalizedVersion([string]$Value) {
    $v = $Value.Trim()
    while ($v.StartsWith('v') -or $v.StartsWith('V')) { $v = $v.Substring(1) }
    return $v
}

function Test-PathInUserPath([string]$Dir) {
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ([string]::IsNullOrWhiteSpace($userPath)) { return $false }

    $normalizedDir = (Resolve-Path -LiteralPath $Dir).Path.TrimEnd('\')
    foreach ($entry in $userPath -split ';') {
        if ([string]::IsNullOrWhiteSpace($entry)) { continue }
        try {
            $normalizedEntry = (Resolve-Path -LiteralPath $entry -ErrorAction Stop).Path.TrimEnd('\')
            if ($normalizedEntry -ieq $normalizedDir) { return $true }
        } catch {
            if ($entry.TrimEnd('\') -ieq $normalizedDir) { return $true }
        }
    }
    return $false
}

function Add-BinaryLink([string]$Source, [string]$TargetDir) {
    $linkPath = Join-Path $TargetDir "$BinaryName.exe"
    if (Test-Path -LiteralPath $linkPath) {
        Remove-Item -LiteralPath $linkPath -Force
    }

    try {
        New-Item -ItemType HardLink -Path $linkPath -Target $Source -Force | Out-Null
        return $linkPath
    } catch {
        Copy-Item -LiteralPath $Source -Destination $linkPath -Force
        return $linkPath
    }
}

function Add-ToUserPath([string]$Dir) {
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    $normalizedDir = (Resolve-Path -LiteralPath $Dir).Path

    if (Test-PathInUserPath $normalizedDir) { return $false }

    $newPath = if ([string]::IsNullOrWhiteSpace($userPath)) {
        $normalizedDir
    } else {
        "$normalizedDir;$userPath"
    }

    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    return $true
}

function Notify-PathChanged {
    $signature = @'
[DllImport("user32.dll", SetLastError = true, CharSet = CharSet.Auto)]
public static extern IntPtr SendMessageTimeout(
    IntPtr hWnd, uint Msg, UIntPtr wParam, string lParam,
    uint fuFlags, uint uTimeout, out UIntPtr lpdwResult);
'@
    try {
        Add-Type -MemberDefinition $signature -Name NativeMethods -Namespace Win32 -ErrorAction Stop
        $null = [UIntPtr]::Zero
        [Win32.NativeMethods]::SendMessageTimeout(
            [IntPtr]0xffff, 0x1A, [UIntPtr]::Zero, "Environment", 2, 5000, [ref]$null) | Out-Null
    } catch {
        # Best-effort only; a new terminal still picks up registry PATH.
    }
}

function Test-Compiler {
    $hasGcc = $false
    $hasClang = $false
    $hasMsvc = $false

    try { Get-Command g++ -ErrorAction Stop | Out-Null; $hasGcc = $true } catch {}
    try { Get-Command clang++ -ErrorAction Stop | Out-Null; $hasClang = $true } catch {}
    try { Get-Command cl -ErrorAction Stop | Out-Null; $hasMsvc = $true } catch {}

    if ($hasGcc -or $hasClang -or $hasMsvc) {
        return $true
    }

    Write-Host ""
    Write-Host "⚠ No C/C++ compiler found!" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "cnext requires a compiler. Install one of:"
    Write-Host ""
    Write-Host "  Option 1: Visual Studio Build Tools"
    Write-Host "    https://visualstudio.microsoft.com/visual-cpp-build-tools/"
    Write-Host ""
    Write-Host "  Option 2: LLVM/Clang"
    Write-Host "    winget install LLVM.LLVM"
    Write-Host ""
    Write-Host "  Option 3: MinGW-w64 (GCC)"
    Write-Host "    winget install MSYS2.MSYS2"
    Write-Host "    Then: pacman -S mingw-w64-x86_64-gcc"
    Write-Host ""
    return $false
}

Write-Host "========================================="
Write-Host "  cnext installer"
Write-Host "========================================="
Write-Host ""

if ($Version -eq "latest" -or [string]::IsNullOrWhiteSpace($Version)) {
    $url = "https://github.com/$Repo/releases/latest/download/$Asset"
    $versionLabel = "latest"
} else {
    $Version = Get-NormalizedVersion $Version
    $url = "https://github.com/$Repo/releases/download/v$Version/$Asset"
    $versionLabel = "v$Version"
}

$dest = Join-Path $InstallDir "$BinaryName.exe"

Write-Host "Detected: windows/amd64"
Write-Host "Version:  $versionLabel"
Write-Host ""

Write-Host "Downloading $url..."
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
Invoke-WebRequest -Uri $url -OutFile $dest
Write-Host "✓ Downloaded to $dest"

$pathCandidates = @(
    (Join-Path $env:USERPROFILE ".local\bin"),
    (Join-Path $env:LOCALAPPDATA "Microsoft\WinGet\Links"),
    (Join-Path $env:USERPROFILE "go\bin")
)

$linked = $false
foreach ($candidate in $pathCandidates) {
    if (-not (Test-Path -LiteralPath $candidate)) { continue }
    if (-not (Test-PathInUserPath $candidate)) { continue }

    $linkPath = Add-BinaryLink -Source $dest -TargetDir $candidate
    Write-Host "✓ Linked $linkPath -> $dest"
    $linked = $true
    break
}

if (-not $linked) {
    if (Add-ToUserPath $InstallDir) {
        Write-Host "✓ Added $InstallDir to user PATH"
    } else {
        Write-Host "✓ $InstallDir is already in user PATH"
    }
}

Notify-PathChanged
$env:Path = [Environment]::GetEnvironmentVariable("Path", "Machine") + ";" +
            [Environment]::GetEnvironmentVariable("Path", "User")

Write-Host ""
Write-Host "========================================="
Write-Host "  Verifying installation..."
Write-Host "========================================="

if (Get-Command $BinaryName -ErrorAction SilentlyContinue) {
    Write-Host "✓ $BinaryName is in PATH"
    & $BinaryName --help | Select-Object -First 5
} else {
    Write-Host ""
    Write-Host "⚠ $BinaryName installed but not in current PATH."
    Write-Host "Restart your terminal to use it."
}

Test-Compiler | Out-Null

Write-Host ""
Write-Host "========================================="
Write-Host "  Quick Start"
Write-Host "========================================="
Write-Host ""
Write-Host "  cnext init my-app    # Create new project"
Write-Host "  cd my-app"
Write-Host "  cnext build          # Build project"
Write-Host "  cnext run            # Run executable"
Write-Host ""
