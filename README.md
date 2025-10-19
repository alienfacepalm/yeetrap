# YeeTrap 🎬

YeeTrap is a command-line tool for YouTube content creators to download their videos in bulk for backup purposes. It uses OAuth2 authentication to securely access your YouTube channel and downloads videos using yt-dlp.

## Features

- 🔐 **Professional OAuth2 App** - Seamless Google authentication with automatic browser opening
- 📺 **Smart Video Listing** - List all videos from your channel or any specified channel
- ⬇️ **Bulk Downloads** - Download multiple videos with quality selection and concurrent processing
- ⚡ **High Performance** - Concurrent downloads for faster backups
- 📝 **Complete Metadata** - Saves video description, thumbnail, and info JSON
- 🎯 **Flexible Access** - Download from any channel you have access to
- 🛠️ **Easy Setup** - Built-in setup wizard with step-by-step OAuth2 configuration
- 🔒 **Secure Storage** - Local token storage with automatic refresh

## Prerequisites

1. **Go** (1.24 or later) - [Download](https://golang.org/dl/)
2. **yt-dlp** - [Installation instructions](https://github.com/yt-dlp/yt-dlp#installation)
3. **Google Cloud Project with YouTube Data API v3 enabled**

## Setup

### 1. Install yt-dlp

#### Windows (via Chocolatey)

```bash
choco install yt-dlp
```

#### Windows (manual)

Download from [GitHub releases](https://github.com/yt-dlp/yt-dlp/releases) and add to PATH.

#### macOS

```bash
brew install yt-dlp
```

#### Linux

```bash
pip install yt-dlp
```

### 2. Set up Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the **YouTube Data API v3**:
   - Go to "APIs & Services" > "Library"
   - Search for "YouTube Data API v3"
   - Click "Enable"
4. Create OAuth2 credentials:
   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth client ID"
   - Choose "Desktop app" as the application type
   - Download the credentials JSON file

### 3. Install YeeTrap

```bash
# Clone or download this repository
git clone https://github.com/AlienFacepalm/YeeTrap.git
cd YeeTrap

# Build the application
go build -o yeetrap

# (Optional) Install to your PATH
# Windows: copy yeetrap.exe to a folder in your PATH
# macOS/Linux: sudo mv yeetrap /usr/local/bin/
```

### 4. Configure OAuth2 Credentials

**Easy Setup (Recommended):**

```bash
yeetrap setup
```

**Manual Setup:**
Place your downloaded credentials file at:

- **Windows**: `%USERPROFILE%\.yeetrap\credentials.json`
- **macOS/Linux**: `~/.yeetrap/credentials.json`

## Usage

### Quick Start

```bash
# 1. Set up OAuth2 (one-time)
yeetrap setup

# 2. Authenticate with YouTube (one-time)
yeetrap auth

# 3. List your videos
yeetrap list --max 10

# 4. Download videos
yeetrap download --max 5 --quality 720p
```

### Authenticate with YouTube

```bash
yeetrap auth
```

This will automatically open your browser for Google OAuth2 login. The authentication token is saved locally and reused automatically.

### List Videos from Your Channel

```bash
# List videos from your authenticated channel
yeetrap list

# List with a maximum number of videos
yeetrap list --max 100

# List from a specific channel ID
yeetrap list --channel UC_x5XG1OV2P6uZZ5FSM9Ttw
```

### Download Videos

```bash
# Download all videos from your channel
yeetrap download

# Download with custom settings
yeetrap download --max 50 --quality 1080p --output ./my-backups --concurrent 5

# Download from a specific channel
yeetrap download --channel UC_x5XG1OV2P6uZZ5FSM9Ttw
```

#### Download Options

- `--channel`, `-c`: YouTube channel ID (leave empty for your own channel)
- `--max`, `-m`: Maximum number of videos to download (default: 50)
- `--output`, `-o`: Output directory (default: ./downloads)
- `--quality`, `-q`: Video quality - `best`, `1080p`, `720p`, `480p` (default: best)
- `--concurrent`, `-j`: Number of concurrent downloads (default: 3)

## Configuration

Configuration is stored at:

- **Windows**: `%USERPROFILE%\.yeetrap\config.json`
- **macOS/Linux**: `~/.yeetrap/config.json`

Example configuration:

```json
{
  "default_channel_id": "",
  "default_quality": "best",
  "output_dir": "./downloads",
  "max_concurrent": 3
}
```

## Project Structure

```
YeeTrap/
├── main.go                    # Application entry point
├── cmd/                       # CLI commands
│   ├── root.go               # Root command setup
│   ├── auth.go               # Authentication command
│   ├── list.go               # List videos command
│   └── download.go           # Download videos command
├── internal/
│   ├── auth/                 # OAuth2 authentication
│   │   └── auth.go
│   ├── youtube/              # YouTube API integration
│   │   └── service.go
│   ├── downloader/           # Video download logic
│   │   └── downloader.go
│   └── config/               # Configuration management
│       └── config.go
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
└── README.md                 # This file
```

## Troubleshooting

### "credentials.json not found"

Make sure you've placed your OAuth2 credentials file at `~/.yeetrap/credentials.json` (or the Windows equivalent).

### "Please run 'yeetrap auth' first"

You need to authenticate before using other commands. Run `yeetrap auth` and follow the prompts.

### "yt-dlp is not installed"

Install yt-dlp following the instructions in the Prerequisites section.

### API Quota Exceeded

The YouTube Data API has daily quota limits. If you hit the limit, you'll need to wait until the next day or request a quota increase from Google Cloud Console.

## Legal & Ethical Considerations

⚠️ **Important**: This tool is designed for content creators to backup their own videos. Please ensure you:

- Only download videos you own or have permission to download
- Respect YouTube's Terms of Service
- Comply with copyright laws in your jurisdiction
- Use this tool responsibly and ethically

## License

MIT License - feel free to use and modify as needed.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/AlienFacepalm/YeeTrap/issues) on GitHub.
