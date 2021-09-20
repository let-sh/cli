#!/usr/bin/env pwsh
# Copyright 2021 Oasis Networks authors. All rights reserved. MIT license.

$ErrorActionPreference = 'Stop'

if ($v) {
  $Version = "${v}"
}
if ($args.Length -eq 1) {
  $Version = $args.Get(0)
}

$LetInstall = $env:LET_INSTALL
$BinDir = if ($LetInstall) {
  "$LetInstall\bin"
} else {
  "$Home\.let\bin"
}

$LetZip = "$BinDir\lets.zip"
$LetExe = "$BinDir\lets.exe"
$Target = 'windows_amd64'

# GitHub requires TLS 1.2
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$LetUri = if (!$Version) {
  $Response = Invoke-WebRequest 'https://install.let-sh.com/version' -UseBasicParsing

  $Version = $Response.Content.Split([Environment]::NewLine) |
  Where-Object { $_ -match "latest:(.*)" } |
  Select-String 'latest:(.*)' |
  ForEach-Object {$_.matches[0].Groups[1].value} |
  Select-Object -First 1

  "https://install.let-sh.com/cli_${Version}_${Target}.zip"
} else {
  "https://install.let-sh.com/cli_${Version}_${Target}.zip"
}

if (!(Test-Path $BinDir)) {
  New-Item $BinDir -ItemType Directory | Out-Null
}

Invoke-WebRequest $LetUri -OutFile $LetZip -UseBasicParsing

if (Get-Command Expand-Archive -ErrorAction SilentlyContinue) {
  Expand-Archive $LetZip -Destination $BinDir -Force
} else {
  if (Test-Path $LetExe) {
    Remove-Item $LetExe
  }
  Add-Type -AssemblyName System.IO.Compression.FileSystem
  [IO.Compression.ZipFile]::ExtractToDirectory($LetZip, $BinDir)
}

Remove-Item $LetZip

$User = [EnvironmentVariableTarget]::User
$Path = [Environment]::GetEnvironmentVariable('Path', $User)
if (!(";$Path;".ToLower() -like "*;$BinDir;*".ToLower())) {
  [Environment]::SetEnvironmentVariable('Path', "$Path;$BinDir", $User)
  $Env:Path += ";$BinDir"
}

Write-Output "let.sh was installed successfully to $LetExe"
Write-Output "Run 'lets --help' to get started"