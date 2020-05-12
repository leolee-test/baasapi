param (
  [string]$platform,
  [string]$arch
)

$ErrorActionPreference = "Stop";

$binary = "baasapi.exe"
$go_path = "$($(Get-ITEM -Path env:AGENT_HOMEDIRECTORY).Value)\go"

Set-Item env:GOPATH "$go_path"

New-Item -Name dist -Path "." -ItemType Directory -Force | Out-Null
New-Item -Name baasapi -Path "$go_path\src\github.com\baasapi" -ItemType Directory -Force | Out-Null

Copy-Item -Path "api" -Destination "$go_path\src\github.com\baasapi\baasapi\api" -Recurse -Force -ErrorAction:SilentlyContinue

Set-Location -Path "api\cmd\baasapi"

go.exe get -t -d -v ./...
go.exe build -v

Move-Item -Path "$($binary)" -Destination "dist"
