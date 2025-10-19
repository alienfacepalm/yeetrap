# Quick Start Guide for YeeTrap

This guide will help you get started with YeeTrap in just a few minutes.

## Step 1: Prerequisites

Before you begin, make sure you have:

1. ✅ **Go** installed (version 1.24 or later)
2. ✅ **yt-dlp** installed and available in PATH
3. ✅ **Google Cloud Project** with YouTube Data API v3 enabled

### Install yt-dlp (if not already installed)

**Windows:**

```powershell
# Using winget
winget install yt-dlp

# Or using Chocolatey
choco install yt-dlp

# Or download from: https://github.com/yt-dlp/yt-dlp/releases
```

**macOS:**

```bash
brew install yt-dlp
```

**Linux:**

```bash
sudo apt install yt-dlp
# or
pip install yt-dlp
```

## Step 2: Set Up Google Cloud Credentials

### 2.1 Create a Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project (or use an existing one)
3. Note down your project name

### 2.2 Enable YouTube Data API v3

1. In Google Cloud Console, go to **"APIs & Services"** → **"Library"**
2. Search for **"YouTube Data API v3"**
3. Click on it and press **"Enable"**

### 2.3 Create OAuth2 Credentials

1. Go to **"APIs & Services"** → **"Credentials"**
2. Click **"Create Credentials"** → **"OAuth client ID"**
3. If prompted, configure the OAuth consent screen:
   - Choose "External" user type
   - Fill in the required fields (app name, support email)
   - Add your email as a test user
   - Save and continue
4. Back to Create OAuth client ID:
   - Select **"Desktop app"** as application type
   - Give it a name (e.g., "YeeTrap Desktop")
   - Click **"Create"**
5. Download the credentials JSON file

### 2.4 Place Credentials File

Create the YeeTrap config directory and place your credentials:

**Windows (PowerShell):**

```powershell
# Create directory
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.yeetrap"

# Move your downloaded credentials file to:
# C:\Users\YourUsername\.yeetrap\credentials.json
```

**macOS/Linux:**

```bash
# Create directory
mkdir -p ~/.yeetrap

# Move your downloaded credentials file
mv ~/Downloads/client_secret_*.json ~/.yeetrap/credentials.json
```

## Step 3: Build YeeTrap

```bash
# Navigate to the YeeTrap directory
cd YeeTrap

# Build the application
go build -o yeetrap.exe

# Verify it works
./yeetrap.exe --help
```

## Step 4: Authenticate with YouTube

Run the authentication command:

```bash
./yeetrap.exe auth
```

This will:

1. Display a URL in your terminal
2. Open your browser to that URL (or copy-paste it)
3. Ask you to log in with your Google account
4. Request permission to access your YouTube data
5. Give you an authorization code
6. Paste the code back into the terminal

After successful authentication, you'll see: ✓ Authentication successful!

## Step 5: List Your Videos

```bash
# List videos from your channel
./yeetrap.exe list

# List with more options
./yeetrap.exe list --max 100
```

## Step 6: Download Your Videos

```bash
# Download videos
./yeetrap.exe download

# Download with custom settings
./yeetrap.exe download --max 10 --quality 1080p --output ./my-backups
```

## Common Commands

### Authentication

```bash
yeetrap auth
```

### List Videos

```bash
# Your channel (default)
yeetrap list

# Specific channel
yeetrap list --channel UCxxxxxxxxxxxxxx

# Limit results
yeetrap list --max 20
```

### Download Videos

```bash
# Download all (up to 50)
yeetrap download

# Custom quality
yeetrap download --quality 720p

# Custom output directory
yeetrap download --output D:\YouTube-Backup

# More concurrent downloads (faster but uses more bandwidth)
yeetrap download --concurrent 5

# Download only 10 videos
yeetrap download --max 10
```

## Troubleshooting

### Error: "credentials.json not found"

- Make sure the file is at `~/.yeetrap/credentials.json` (or `%USERPROFILE%\.yeetrap\credentials.json` on Windows)
- Check that the file name is exactly `credentials.json` (not `credentials (1).json` or similar)

### Error: "yt-dlp is not installed"

- Install yt-dlp using the instructions in Step 1
- Make sure it's in your PATH by running: `yt-dlp --version`

### Error: "Please run 'yeetrap auth' first"

- You need to authenticate before listing or downloading
- Run: `yeetrap auth`

### Error: API quota exceeded

- The YouTube Data API has daily quotas
- Wait 24 hours or request a quota increase in Google Cloud Console

### Downloads are slow

- Increase concurrent downloads: `--concurrent 5`
- Check your internet connection
- Some videos may be large files

## Tips & Best Practices

1. **Start Small**: Test with `--max 5` before downloading your entire channel
2. **Quality vs Size**: Use `--quality 720p` if you want to save disk space
3. **Organize Downloads**: Use custom output directories for different channels
4. **Regular Backups**: Set up a scheduled task to run downloads regularly
5. **Check Disk Space**: Ensure you have enough space before large downloads

## Next Steps

- Check the full [README.md](README.md) for detailed documentation
- Review command options with `yeetrap [command] --help`
- Star the repository if you find it useful!

## Getting Help

If you encounter issues:

1. Check the [Troubleshooting](#troubleshooting) section above
2. Review the [README.md](README.md)
3. [Open an issue](https://github.com/AlienFacepalm/YeeTrap/issues) on GitHub

Happy backing up! 🎬

