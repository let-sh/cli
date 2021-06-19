#!/usr/bin/env pwsh
# Copyright 2021 Oasis Networks authors. All rights reserved. MIT license.

$ErrorActionPreference = 'Stop'

if ($v) {
  $Version = "v${v}"
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

$LetZip = "$BinDir\let.zip"
$LetExe = "$BinDir\let.exe"
$Target = 'x86_64-windows'

# GitHub requires TLS 1.2
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$LetUri = if (!$Version) {
  $Response = Invoke-WebRequest 'https://github.com/let-sh/cli/releases' -UseBasicParsing
  if ($PSVersionTable.PSEdition -eq 'Core') {
    $Response.Links |
      Where-Object { $_.href -like "/let-sh/cli/releases/download/*/lets-${Target}.zip" } |
      ForEach-Object { 'https://github.com' + $_.href } |
      Select-Object -First 1
  } else {
    $HTMLFile = New-Object -Com HTMLFile
    if ($HTMLFile.IHTMLDocument2_write) {
      $HTMLFile.IHTMLDocument2_write($Response.Content)
    } else {
      $ResponseBytes = [Text.Encoding]::Unicode.GetBytes($Response.Content)
      $HTMLFile.write($ResponseBytes)
    }
    $HTMLFile.getElementsByTagName('a') |
      Where-Object { $_.href -like "about:/let-sh/cli/releases/download/*/lets-${Target}.zip" } |
      ForEach-Object { $_.href -replace 'about:', 'https://github.com' } |
      Select-Object -First 1
  }
} else {
  "https://github.com/let-sh/cli/releases/download/${Version}/lets-${Target}.zip"
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