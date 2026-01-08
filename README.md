# gorig

```
 ╔═══════════════════════════╗
 ║   *╭───╮*╭───╮*╭───╮*     ║
 ║  .*│▓▓▓│*│▓▓▓│*│▓▓▓│*.    ║
 ║   *╰─┬─╯*╰─┬─╯*╰─┬─╯*     ║
 ║  ╭───────────────────╮    ║
 ║  │                   │    ║
 ║  │   ◉ ╲       ╱ ◉   │    ║
 ║  │                   │    ║
 ║  │      ╭─────╮      │    ║
 ║  │                   │    ║
 ║  ╰───────────────────╯    ║
 ║                           ║
 ║   ◉ ◉ ◉   DISTORT    ◉    ║
 ╚═══════════════════════════╝
```

A pet project — simple audio bypass with acceptable latency and the ability to write guitar pedals and effects in pure Go without any middleware.

## Current State

- Fully functional audio bypass
- TUI interface with preset management
- Write your own presets and use them on the fly
- Effects written in Go — no DSLs, no intermediate layers
- Switch between any available system input/output audio devices on the fly

## Requirements

PortAudio is required for audio I/O.

### macOS

```bash
brew install portaudio
```

### Ubuntu/Debian

```bash
sudo apt install libportaudio2 libportaudio-dev
```

### Windows

```powershell
choco install portaudio
```

Or download prebuilt binaries from [portaudio.com](https://files.portaudio.com/download.html).

## Building

```bash
go build -o gorig ./cmd/ui
```

## Usage

```bash
./gorig
```

## Writing Effects

Effects are plain Go code. See `effects/` directory for examples.