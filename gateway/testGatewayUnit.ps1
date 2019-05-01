$CurDir = pwd

$CurTime=(New-TimeSpan -Start (Get-Date "01/01/1970") -End (Get-Date)).TotalSeconds
$Coverage = "$Env:TEMP\coverage_$CurTime"
mkdir $Coverage
$CoverFileName = "cover.out"
$CoverFilePath="$Coverage\$CoverFileName"
Set-Location -Path .\gateway\
go test `
-covermode=atomic `
-coverprofile $CoverFileName `
-outputdir $Coverage `
-tags "unit" `
.\...

$CoverHtmlPath="$Coverage\coverage.html"
go tool cover -html="$CoverFilePath" -o $CoverHtmlPath
start $CoverHtmlPath
Write-Host "Coverage information can be found here: $Coverage"
Set-Location -Path $CurDir