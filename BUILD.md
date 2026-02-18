# Build Guide — devtools

Komplett-Anleitung zum Kompilieren für jedes Betriebssystem und jede CPU-Architektur.

> **Wichtig:** Das SQLite-Package (`glebarez/sqlite`) ist **Pure Go** — kein CGO nötig.
> Cross-Compilation funktioniert daher ohne C-Compiler oder Docker.

---

## Voraussetzungen

| Tool | Mindestversion | Download |
|------|---------------|----------|
| Go | 1.21+ | https://go.dev/dl |
| Node.js | 20.19+ oder 22.12+ | https://nodejs.org |
| npm | mitgeliefert mit Node | — |

---

## Schritt 1 — Frontend bauen (einmalig, plattformunabhängig)

Das Frontend-Build (`dist/`) ist für alle Plattformen gleich.

```bash
cd frontend-vite
npm install
npm run build
cd ..
```

Ergebnis: `frontend-vite/dist/` — wird vom Go-Binary zur Laufzeit eingelesen.

---

## Schritt 2 — Go-Backend kompilieren

```bash
$env:CGO_ENABLED="0"; $env:GOOS="linux"; $env:GOARCH="amd64"; go build -o devtools-linux-amd64
```

### Syntax

```bash
GOOS=<os> GOARCH=<arch> go build -o <output-name> .
```

Unter **Windows** (PowerShell):

```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o devtools-linux-amd64 .
```

---

## Komplette Build-Matrix

### Windows

| Architektur | Befehl | Binary |
|-------------|--------|--------|
| x64 (Intel/AMD) | `GOOS=windows GOARCH=amd64 go build -o devtools-windows-amd64.exe .` | `devtools-windows-amd64.exe` |
| ARM64 (Surface Pro X, Snapdragon) | `GOOS=windows GOARCH=arm64 go build -o devtools-windows-arm64.exe .` | `devtools-windows-arm64.exe` |
| x86 (32-bit) | `GOOS=windows GOARCH=386 go build -o devtools-windows-386.exe .` | `devtools-windows-386.exe` |

### Linux

| Architektur | Befehl | Binary |
|-------------|--------|--------|
| x64 | `GOOS=linux GOARCH=amd64 go build -o devtools-linux-amd64 .` | `devtools-linux-amd64` |
| ARM64 (Raspberry Pi 4/5, AWS Graviton, Oracle Ampere) | `GOOS=linux GOARCH=arm64 go build -o devtools-linux-arm64 .` | `devtools-linux-arm64` |
| ARMv7 (Raspberry Pi 2/3) | `GOOS=linux GOARCH=arm GOARM=7 go build -o devtools-linux-armv7 .` | `devtools-linux-armv7` |
| ARMv6 (Raspberry Pi 1/Zero) | `GOOS=linux GOARCH=arm GOARM=6 go build -o devtools-linux-armv6 .` | `devtools-linux-armv6` |
| x86 (32-bit) | `GOOS=linux GOARCH=386 go build -o devtools-linux-386 .` | `devtools-linux-386` |
| RISC-V 64 | `GOOS=linux GOARCH=riscv64 go build -o devtools-linux-riscv64 .` | `devtools-linux-riscv64` |
| PowerPC 64 LE (IBM) | `GOOS=linux GOARCH=ppc64le go build -o devtools-linux-ppc64le .` | `devtools-linux-ppc64le` |
| PowerPC 64 BE | `GOOS=linux GOARCH=ppc64 go build -o devtools-linux-ppc64 .` | `devtools-linux-ppc64` |
| IBM S390x (Mainframe) | `GOOS=linux GOARCH=s390x go build -o devtools-linux-s390x .` | `devtools-linux-s390x` |
| MIPS 64 LE | `GOOS=linux GOARCH=mips64le go build -o devtools-linux-mips64le .` | `devtools-linux-mips64le` |
| MIPS 64 BE | `GOOS=linux GOARCH=mips64 go build -o devtools-linux-mips64 .` | `devtools-linux-mips64` |
| MIPS LE (Router) | `GOOS=linux GOARCH=mipsle go build -o devtools-linux-mipsle .` | `devtools-linux-mipsle` |
| MIPS BE (Router) | `GOOS=linux GOARCH=mips go build -o devtools-linux-mips .` | `devtools-linux-mips` |
| LoongArch 64 | `GOOS=linux GOARCH=loong64 go build -o devtools-linux-loong64 .` | `devtools-linux-loong64` |

### macOS

| Architektur | Befehl | Binary |
|-------------|--------|--------|
| Apple Silicon (M1/M2/M3/M4) | `GOOS=darwin GOARCH=arm64 go build -o devtools-darwin-arm64 .` | `devtools-darwin-arm64` |
| Intel Mac | `GOOS=darwin GOARCH=amd64 go build -o devtools-darwin-amd64 .` | `devtools-darwin-amd64` |
| Universal Binary (beide) | siehe unten | `devtools-darwin-universal` |

**macOS Universal Binary** (läuft auf Intel und Apple Silicon):

```bash
GOOS=darwin GOARCH=amd64 go build -o devtools-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o devtools-darwin-arm64 .
lipo -create -output devtools-darwin-universal devtools-darwin-amd64 devtools-darwin-arm64
```

> `lipo` ist nur auf macOS verfügbar.

### FreeBSD

| Architektur | Befehl | Binary |
|-------------|--------|--------|
| x64 | `GOOS=freebsd GOARCH=amd64 go build -o devtools-freebsd-amd64 .` | `devtools-freebsd-amd64` |
| ARM64 | `GOOS=freebsd GOARCH=arm64 go build -o devtools-freebsd-arm64 .` | `devtools-freebsd-arm64` |
| ARMv7 | `GOOS=freebsd GOARCH=arm GOARM=7 go build -o devtools-freebsd-armv7 .` | `devtools-freebsd-armv7` |
| x86 (32-bit) | `GOOS=freebsd GOARCH=386 go build -o devtools-freebsd-386 .` | `devtools-freebsd-386` |

### OpenBSD / NetBSD

```bash
# OpenBSD
GOOS=openbsd GOARCH=amd64 go build -o devtools-openbsd-amd64 .
GOOS=openbsd GOARCH=arm64 go build -o devtools-openbsd-arm64 .

# NetBSD
GOOS=netbsd GOARCH=amd64 go build -o devtools-netbsd-amd64 .
GOOS=netbsd GOARCH=arm64 go build -o devtools-netbsd-arm64 .
```

---

## Schritt 3 — yt-dlp Binary beschaffen

Das Binary muss zur kompilierten App im **gleichen Ordner** liegen.
Download: https://github.com/yt-dlp/yt-dlp/releases/latest

| Plattform | Dateiname auf GitHub | Umbenennen zu |
|-----------|---------------------|---------------|
| Windows x64 | `yt-dlp.exe` | `yt-dlp.exe` |
| Windows ARM64 | `yt-dlp_x86.exe` (Emulation) oder nativ aus dem Source build | `yt-dlp.exe` |
| Linux x64 | `yt-dlp_linux` | `yt-dlp.exe`* |
| Linux ARM64 | `yt-dlp_linux_aarch64` | `yt-dlp.exe`* |
| Linux ARMv7 | `yt-dlp_linux_armv7l` | `yt-dlp.exe`* |
| macOS (Universal) | `yt-dlp_macos` | `yt-dlp.exe`* |
| macOS Intel | `yt-dlp_macos_legacy` | `yt-dlp.exe`* |

> \* Der Code referenziert aktuell `./yt-dlp.exe` hardcoded. Entweder die Datei so benennen
> **oder** in `tools/ytdl/service.go` Zeile 142 anpassen:
>
> ```go
> // Aktuell:
> cmd := exec.Command("./yt-dlp.exe", args...)
>
> // Cross-platform:
> ytdlpBin := "./yt-dlp.exe"
> if runtime.GOOS != "windows" {
>     ytdlpBin = "./yt-dlp"
> }
> cmd := exec.Command(ytdlpBin, args...)
> ```
> Dann `"runtime"` zu den Imports hinzufügen.

Auf Linux/macOS muss das Binary außerdem ausführbar sein:

```bash
chmod +x ./yt-dlp.exe   # oder ./yt-dlp
```

---

## Schritt 4 — Deployment-Struktur

Nach dem Build braucht die App diese Dateien im selben Verzeichnis:

```
deploy/
├── devtools-linux-amd64    ← kompiliertes Go-Binary
├── yt-dlp.exe              ← yt-dlp Binary (oder yt-dlp auf Linux/macOS)
├── ytdl.db                 ← wird automatisch erstellt
├── downloads/              ← wird automatisch erstellt
└── frontend-vite/
    └── dist/               ← gebautes Vue-Frontend
        ├── index.html
        └── assets/
```

---

## Alles auf einmal — Build-Scripts

### Bash (Linux/macOS) — alle Plattformen auf einmal

```bash
#!/bin/bash
# build-all.sh

set -e

echo "Building frontend..."
cd frontend-vite && npm install && npm run build && cd ..

mkdir -p dist

targets=(
    "windows/amd64/devtools-windows-amd64.exe"
    "windows/arm64/devtools-windows-arm64.exe"
    "windows/386/devtools-windows-386.exe"
    "linux/amd64/devtools-linux-amd64"
    "linux/arm64/devtools-linux-arm64"
    "linux/arm/devtools-linux-armv7"
    "linux/386/devtools-linux-386"
    "linux/riscv64/devtools-linux-riscv64"
    "linux/ppc64le/devtools-linux-ppc64le"
    "linux/s390x/devtools-linux-s390x"
    "linux/mips64le/devtools-linux-mips64le"
    "linux/mipsle/devtools-linux-mipsle"
    "darwin/amd64/devtools-darwin-amd64"
    "darwin/arm64/devtools-darwin-arm64"
    "freebsd/amd64/devtools-freebsd-amd64"
    "freebsd/arm64/devtools-freebsd-arm64"
)

for target in "${targets[@]}"; do
    IFS='/' read -r os arch name <<< "$target"
    echo "Building $name (GOOS=$os GOARCH=$arch)..."
    GOOS=$os GOARCH=$arch go build -ldflags="-s -w" -o "dist/$name" .
done

echo "Done! Binaries are in dist/"
```

Ausführen:
```bash
chmod +x build-all.sh
./build-all.sh
```

### PowerShell (Windows) — alle Plattformen auf einmal

```powershell
# build-all.ps1

Write-Host "Building frontend..." -ForegroundColor Cyan
Set-Location frontend-vite
npm install
npm run build
Set-Location ..

New-Item -ItemType Directory -Force -Path dist | Out-Null

$targets = @(
    @{os="windows"; arch="amd64"; name="devtools-windows-amd64.exe"},
    @{os="windows"; arch="arm64"; name="devtools-windows-arm64.exe"},
    @{os="windows"; arch="386";   name="devtools-windows-386.exe"},
    @{os="linux";   arch="amd64"; name="devtools-linux-amd64"},
    @{os="linux";   arch="arm64"; name="devtools-linux-arm64"},
    @{os="linux";   arch="arm";   name="devtools-linux-armv7"},
    @{os="linux";   arch="386";   name="devtools-linux-386"},
    @{os="linux";   arch="riscv64"; name="devtools-linux-riscv64"},
    @{os="linux";   arch="ppc64le"; name="devtools-linux-ppc64le"},
    @{os="linux";   arch="s390x"; name="devtools-linux-s390x"},
    @{os="linux";   arch="mips64le"; name="devtools-linux-mips64le"},
    @{os="linux";   arch="mipsle"; name="devtools-linux-mipsle"},
    @{os="darwin";  arch="amd64"; name="devtools-darwin-amd64"},
    @{os="darwin";  arch="arm64"; name="devtools-darwin-arm64"},
    @{os="freebsd"; arch="amd64"; name="devtools-freebsd-amd64"},
    @{os="freebsd"; arch="arm64"; name="devtools-freebsd-arm64"}
)

foreach ($t in $targets) {
    Write-Host "Building $($t.name)..." -ForegroundColor Green
    $env:GOOS = $t.os
    $env:GOARCH = $t.arch
    go build -ldflags="-s -w" -o "dist\$($t.name)" .
}

Remove-Item Env:GOOS
Remove-Item Env:GOARCH

Write-Host "Done! Binaries are in dist/" -ForegroundColor Cyan
```

Ausführen:
```powershell
.\build-all.ps1
```

---

## Build-Flags (optional, empfohlen)

```bash
go build -ldflags="-s -w" -o output .
```

| Flag | Effekt |
|------|--------|
| `-s` | Symbol-Tabelle entfernen (kleinere Binary) |
| `-w` | DWARF Debug-Info entfernen (kleinere Binary) |

Mit beiden Flags wird die Binary ca. 25–30% kleiner.

Noch kleiner mit `upx` (Packer):
```bash
upx --best dist/devtools-linux-amd64
```

---

## Alle unterstützten GOOS/GOARCH Kombinationen anzeigen

```bash
go tool dist list
```

Gibt alle offiziell unterstützten Kombinationen aus (aktuell 45+).

---

## Proxmox LXC / Debian Server (Empfehlung)

Für einen Proxmox LXC Container mit Debian x64:

```bash
# Build (von Windows/Mac/Linux)
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o devtools-linux-amd64 .

# Auf den Server kopieren
scp devtools-linux-amd64 user@server:/opt/devtools/devtools
scp yt-dlp_linux user@server:/opt/devtools/yt-dlp.exe

# Auf dem Server
chmod +x /opt/devtools/devtools
chmod +x /opt/devtools/yt-dlp.exe
cd /opt/devtools && ./devtools
```

Systemd Service (`/etc/systemd/system/devtools.service`):

```ini
[Unit]
Description=devtools
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/devtools
ExecStart=/opt/devtools/devtools
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
systemctl enable --now devtools
```
