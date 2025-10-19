# YeeTrap Project Information

## Overview

**YeeTrap** is a command-line application written in Go that allows YouTube content creators to download their videos in bulk for backup purposes. It uses OAuth2 authentication to securely access YouTube channels and leverages yt-dlp for video downloading.

## Project Status

✅ **Complete and Ready to Use**

## Features Implemented

### Core Features

- ✅ OAuth2 authentication with Google/YouTube
- ✅ List videos from authenticated user's channel
- ✅ List videos from any specified channel ID
- ✅ Bulk video download with quality selection
- ✅ Concurrent downloads for efficiency
- ✅ Progress tracking during downloads
- ✅ Video metadata preservation (description, thumbnail, info JSON)

### CLI Commands

- ✅ `auth` - Authenticate with YouTube
- ✅ `list` - List videos from a channel
- ✅ `download` - Download videos in bulk
- ✅ `version` - Display version information

### Configuration

- ✅ Configuration file support (~/.yeetrap/config.json)
- ✅ Command-line flags for all options
- ✅ Sensible defaults

### Additional Features

- ✅ Cross-platform support (Windows, macOS, Linux)
- ✅ Setup scripts for easy installation
- ✅ Comprehensive documentation
- ✅ MIT License
- ✅ Makefile for building

## Architecture

### Directory Structure

```
YeeTrap/
├── main.go                     # Entry point
├── cmd/                        # CLI commands
│   ├── root.go                # Root command & setup
│   ├── auth.go                # Authentication command
│   ├── list.go                # List videos command
│   ├── download.go            # Download videos command
│   └── version.go             # Version command
├── internal/                   # Internal packages
│   ├── auth/                  # OAuth2 authentication
│   │   └── auth.go
│   ├── youtube/               # YouTube API integration
│   │   └── service.go
│   ├── downloader/            # Video download logic
│   │   └── downloader.go
│   └── config/                # Configuration management
│       └── config.go
├── README.md                   # Main documentation
├── QUICKSTART.md              # Quick start guide
├── LICENSE                     # MIT License
├── Makefile                    # Build automation
├── setup.ps1                   # Windows setup script
└── setup.sh                    # Unix setup script
```

### Key Technologies

#### Language & Framework

- **Go 1.24+** - Primary language
- **Cobra** - CLI framework (github.com/spf13/cobra)

#### Authentication & API

- **golang.org/x/oauth2** - OAuth2 implementation
- **google.golang.org/api/youtube/v3** - YouTube Data API v3 client

#### Video Download

- **yt-dlp** - External tool for video downloading (Python-based)

### Data Flow

1. **Authentication Flow**

   ```
   User runs 'auth' → Opens OAuth2 URL → User authorizes
   → Receives auth code → Exchanges for token → Saves token locally
   ```

2. **List Videos Flow**

   ```
   User runs 'list' → Loads saved token → Authenticates with YouTube API
   → Retrieves channel ID (if not provided) → Fetches uploads playlist
   → Iterates through videos → Displays results
   ```

3. **Download Flow**
   ```
   User runs 'download' → Lists videos (same as above)
   → Creates output directory → Spawns concurrent workers
   → Each worker calls yt-dlp with video URL → Downloads video + metadata
   → Reports progress → Completes
   ```

## Configuration Files

### User Configuration Directory

- **Windows**: `%USERPROFILE%\.yeetrap\`
- **macOS/Linux**: `~/.yeetrap/`

### Files Stored

1. **credentials.json** - OAuth2 client credentials (from Google Cloud Console)
2. **token.json** - OAuth2 access/refresh tokens (generated after auth)
3. **config.json** - User preferences (optional)

## API Usage & Quotas

### YouTube Data API v3 Quotas

- Default quota: 10,000 units per day
- List channels: ~1-3 units
- List playlist items (50): ~1 unit
- Listing 500 videos ≈ 10-15 units

### yt-dlp

- No API used (downloads directly from YouTube)
- Subject to YouTube's rate limiting
- Recommended: 3-5 concurrent downloads

## Dependencies

### Go Modules

```
github.com/spf13/cobra v1.10.1
golang.org/x/oauth2 v0.32.0
google.golang.org/api v0.252.0
```

### External Tools

- **yt-dlp** - Must be installed separately and available in PATH

## Building & Installation

### Quick Build

```bash
go build -o yeetrap.exe
```

### Using Makefile

```bash
make build        # Build for current platform
make build-all    # Build for all platforms
make clean        # Remove build artifacts
make deps         # Install dependencies
```

### Using Setup Scripts

```bash
# Windows
.\setup.ps1

# macOS/Linux
chmod +x setup.sh
./setup.sh
```

## Usage Examples

### Basic Workflow

```bash
# 1. Authenticate
yeetrap auth

# 2. List your videos
yeetrap list --max 10

# 3. Download videos
yeetrap download --max 10 --quality 1080p --output ./backups
```

### Advanced Usage

```bash
# Download from specific channel
yeetrap download --channel UC_x5XG1OV2P6uZZ5FSM9Ttw --max 100

# Fast downloads with high concurrency
yeetrap download --concurrent 10 --quality 720p

# Download all videos (no limit)
yeetrap download --max 0
```

## Security Considerations

### OAuth2 Token Storage

- Tokens stored in `~/.yeetrap/token.json`
- File permissions: 0600 (user read/write only)
- Refresh tokens used for long-term access
- No passwords stored

### API Credentials

- Client ID and secret in `credentials.json`
- Should be kept private
- Desktop app credentials (not for web use)

### Recommendations

1. Don't commit credentials or tokens to version control
2. Keep `.yeetrap` directory permissions restricted
3. Revoke access via Google Account settings when not needed

## Limitations & Known Issues

### Limitations

1. Requires manual OAuth2 flow (no headless auth)
2. Subject to YouTube API quotas
3. yt-dlp must be installed separately
4. Download speed depends on internet connection
5. Private/unlisted videos require appropriate channel access

### Potential Issues

1. **Large channels**: Listing thousands of videos may take time
2. **Network errors**: Downloads may fail and need retry
3. **Disk space**: High-quality videos are large files
4. **API changes**: YouTube API updates may require code changes

## Future Enhancements (Not Implemented)

Potential features for future versions:

- Resume interrupted downloads
- Download queue persistence
- Filtering by date/views/etc
- Playlist support
- Multiple channel batch download
- Progress bar improvements
- Retry logic for failed downloads
- Notification on completion
- Docker container support
- Web UI option

## Testing

Currently, the project has:

- ✅ Manual testing performed
- ❌ No automated unit tests
- ❌ No integration tests

To add tests in the future:

```bash
go test ./...
```

## Contributing

This is a personal/organizational project. Contributions welcome via:

1. GitHub issues for bug reports
2. Pull requests for enhancements
3. Documentation improvements

## License

MIT License - See LICENSE file for details

## Support & Documentation

- **README.md** - Comprehensive documentation
- **QUICKSTART.md** - Step-by-step getting started guide
- **Command help** - Run `yeetrap [command] --help`
- **GitHub Issues** - For bug reports and feature requests

## Version History

### v1.0.0 (Current)

- Initial release
- Core authentication, listing, and download features
- Cross-platform support
- Comprehensive documentation

## Contact & Links

- **Repository**: https://github.com/AlienFacepalm/YeeTrap
- **Author**: AlienFacepalm
- **License**: MIT

---

**Last Updated**: October 19, 2025

