# YeeTrap OAuth2 App Setup Guide

This guide will help you set up YeeTrap as a proper Google OAuth2 application for seamless YouTube authentication.

## 🚀 Quick Setup

The easiest way to set up OAuth2 is using YeeTrap's built-in setup command:

```bash
yeetrap setup
```

This will guide you through the entire process step by step.

## 📋 Manual Setup Steps

### Step 1: Create Google Cloud Project

1. **Go to [Google Cloud Console](https://console.cloud.google.com/)**
2. **Create a new project**:
   - Click "Select a project" → "New Project"
   - Project name: `YeeTrap` (or any name you prefer)
   - Click "Create"

### Step 2: Enable YouTube Data API v3

1. **Navigate to APIs & Services**:
   - Go to "APIs & Services" → "Library"
2. **Search for YouTube Data API**:
   - Search for "YouTube Data API v3"
   - Click on the result
3. **Enable the API**:
   - Click "Enable"

### Step 3: Configure OAuth Consent Screen

1. **Go to OAuth consent screen**:
   - Navigate to "APIs & Services" → "OAuth consent screen"
2. **Choose user type**:
   - Select "External" (unless you have a Google Workspace account)
   - Click "Create"
3. **Fill in app information**:
   - **App name**: `YeeTrap`
   - **User support email**: Your email address
   - **Developer contact information**: Your email address
   - Click "Save and Continue"
4. **Add scopes** (optional):
   - Click "Add or Remove Scopes"
   - Search for "YouTube Data API v3"
   - Select the scope (usually added automatically)
   - Click "Update" → "Save and Continue"
5. **Add test users**:
   - Add your email address as a test user
   - Click "Save and Continue"
6. **Review and submit**:
   - Review the summary
   - Click "Back to Dashboard"

### Step 4: Create OAuth2 Credentials

1. **Go to Credentials**:
   - Navigate to "APIs & Services" → "Credentials"
2. **Create OAuth client ID**:
   - Click "Create Credentials" → "OAuth client ID"
3. **Configure the OAuth client**:
   - **Application type**: Desktop application
   - **Name**: `YeeTrap Desktop`
   - Click "Create"
4. **Download credentials**:
   - Click the download button (⬇️) next to your new OAuth client
   - Save the JSON file

### Step 5: Place Credentials File

1. **Rename the file**:
   - Rename the downloaded file to `credentials.json`
2. **Create YeeTrap config directory**:

   ```bash
   # Windows
   mkdir %USERPROFILE%\.yeetrap

   # macOS/Linux
   mkdir -p ~/.yeetrap
   ```

3. **Move the file**:

   - Move `credentials.json` to the config directory:

   ```
   # Windows
   C:\Users\YourUsername\.yeetrap\credentials.json

   # macOS/Linux
   ~/.yeetrap/credentials.json
   ```

### Step 6: Test the Setup

1. **Run the setup test**:
   ```bash
   yeetrap setup
   ```
2. **Authenticate**:
   ```bash
   yeetrap auth
   ```
3. **Test listing videos**:
   ```bash
   yeetrap list --max 5
   ```

## 🔧 OAuth2 App Configuration Details

### App Information

- **App Name**: YeeTrap
- **App Type**: Desktop Application
- **Scopes**: YouTube Data API v3 (read-only)
- **Redirect URI**: Not required for desktop apps

### Security Features

- **State Parameter**: Random state token for CSRF protection
- **Offline Access**: Refresh tokens for long-term access
- **Local Token Storage**: Encrypted token storage in user's home directory

### Permissions Requested

- **YouTube Data API v3**: Read access to your YouTube channel and videos
- **Scope**: `https://www.googleapis.com/auth/youtube.readonly`

## 🎯 Enhanced Features

### Automatic Browser Opening

YeeTrap automatically opens your default browser to the Google OAuth2 login page.

### Cross-Platform Support

- **Windows**: Uses `rundll32 url.dll,FileProtocolHandler`
- **macOS**: Uses `open` command
- **Linux**: Uses `xdg-open` command

### Token Management

- **Automatic Refresh**: Tokens are automatically refreshed when needed
- **Secure Storage**: Tokens stored with restricted file permissions (600)
- **Persistent Sessions**: No need to re-authenticate frequently

## 🛠️ Troubleshooting

### Common Issues

1. **"credentials.json not found"**

   - Ensure the file is named exactly `credentials.json`
   - Check the file is in the correct directory
   - Run `yeetrap setup` to verify the path

2. **"Invalid client"**

   - Verify the OAuth2 client is configured as "Desktop application"
   - Check that the JSON file is not corrupted
   - Re-download the credentials from Google Cloud Console

3. **"Access blocked"**

   - Ensure your email is added as a test user in OAuth consent screen
   - Check that the app is in "Testing" mode (not "In production")
   - Verify the YouTube Data API v3 is enabled

4. **"Browser won't open"**
   - YeeTrap will show the URL to copy-paste manually
   - Ensure you have a default browser set
   - Try running from a different terminal

### API Quota Issues

- **Daily quota**: YouTube Data API has a 10,000 unit daily quota
- **List operations**: Each list operation uses 1-3 units
- **Quota exceeded**: Wait 24 hours or request quota increase

### Token Issues

- **Expired token**: Run `yeetrap auth` again
- **Invalid token**: Delete `~/.yeetrap/token.json` and re-authenticate
- **Permission denied**: Check file permissions on the config directory

## 🔒 Security Best Practices

1. **Keep credentials private**: Never share your `credentials.json` file
2. **Regular token refresh**: Tokens are automatically refreshed
3. **Revoke access**: You can revoke access in your Google Account settings
4. **Secure storage**: Config directory has restricted permissions
5. **Use test users**: Add only trusted email addresses as test users

## 📚 Additional Resources

- [Google OAuth2 Documentation](https://developers.google.com/identity/protocols/oauth2)
- [YouTube Data API v3 Documentation](https://developers.google.com/youtube/v3)
- [Google Cloud Console](https://console.cloud.google.com/)
- [OAuth2 Consent Screen Guide](https://developers.google.com/identity/protocols/oauth2/openid-connect#consentscreen)

## 🎉 Success!

Once you've completed these steps, YeeTrap will work as a proper OAuth2 application with:

- ✅ Automatic browser opening
- ✅ Secure token management
- ✅ Cross-platform compatibility
- ✅ Professional user experience
- ✅ Proper error handling and guidance

You're now ready to use YeeTrap to backup your YouTube videos! 🎬
