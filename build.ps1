# Создаем папку для сборок
$buildDir = "build"
New-Item -ItemType Directory -Path $buildDir -Force | Out-Null

# Список целевых платформ и архитектур
$targets = @(
    "darwin_amd64",
    "darwin_arm64",
    "freebsd_386",
    "freebsd_amd64",
    "freebsd_arm",
    "freebsd_arm64",
    "linux_386",
    "linux_amd64",
    "linux_arm",
    "linux_loong64",
    "linux_mips",
    "linux_mips64",
    "linux_mips64le",
    "linux_mipsle",
    "linux_riscv64",
    "linux_s390x",
    "netbsd_386",
    "netbsd_amd64",
    "netbsd_arm",
    "netbsd_arm64",
    "windows_386",
    "windows_amd64",
    "windows_arm",
    "windows_arm64"
)

foreach ($target in $targets) {
    $parts = $target.Split('_')
    $goos = $parts[0]
    $goarch = $parts[1]
    $env:GOOS = $goos
    $env:GOARCH = $goarch

    # Имя выходного файла
    $extension = if ($goos -eq "windows") { ".exe" } else { "" }
    $outputName = "$buildDir/ClientHandler_${goos}_${goarch}${extension}"

    Write-Host "▶ Building for $goos/$goarch..."

    # Сборка
    go build -o $outputName .

    # Проверка ошибки
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to build for $goos/$goarch" -ForegroundColor Red
    }
}
