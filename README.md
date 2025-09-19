# Wooper Bot

A Discord bot built in Go that sends random images from organized local directories.

## Features

- **Flexible Image Commands**: Create commands for any image category by organizing images in folders
- **Local Image Storage**: Fast, reliable image serving from local directories
- **Dynamic Command Discovery**: Automatically creates commands based on available image folders
- **Help System**: Built-in help command to list available image categories
- **Comprehensive Logging**: Structured logging with Zap for command tracking, user metrics, and performance monitoring
- **Clean Architecture**: Modular design with separate packages for config, services, handlers, and bot logic
- **Environment Configuration**: Support for `.env` files and environment variables
- **Graceful Shutdown**: Proper signal handling for clean shutdowns

## Commands

- `!<category>` - Sends a random image from the specified category (e.g., `!wooper`, `!cats`, `!dogs`)
- `!help` or `!list` - Shows all available image categories and image counts

## Image Organization

Organize your images in the `img/` directory structure:

```
img/
├── wooper/          # Images for !wooper command
│   ├── wooper1.jpg
│   ├── wooper2.jpg
│   └── Wooper_anime.webp
├── cats/            # Images for !cats command
│   ├── cat1.jpg
│   └── cat2.png
└── dogs/            # Images for !dogs command
    ├── dog1.jpg
    └── dog2.gif
```

Supported image formats: `.png`, `.jpg`, `.jpeg`, `.gif`, `.webp`

## Logging

The bot includes comprehensive structured logging using Zap. Logs include:

- **Command Tracking**: Every command execution with user info, timing, and results
- **Performance Metrics**: Response times for image operations
- **User Analytics**: Who is using which commands and when
- **Error Tracking**: Detailed error information for debugging
- **System Events**: Bot startup, image service initialization, etc.

### Log Levels

Set the log level using the `LOG_LEVEL` environment variable:
- `debug` - All messages including debug info
- `info` - General information (default)
- `warn` - Warnings and errors
- `error` - Errors only

### Example Log Output

```json
{
  "timestamp": "2024-01-20T10:30:45Z",
  "level": "info",
  "message": "Command received",
  "command": "!wooper",
  "category": "wooper",
  "user": "username",
  "user_id": "123456789",
  "channel_id": "987654321",
  "guild_id": "111222333"
}
```

### Adding New Image Categories

To add a new image category (e.g., `dogs`):

1. Create a new directory: `mkdir img/dogs`
2. Add your images to the directory: `cp your-dog-images/* img/dogs/`
3. Restart the bot
4. Use the new command: `!dogs`

The bot will automatically detect the new category and make it available as a command.

## Setup

### Prerequisites

- Go 1.19 or later
- A Discord bot token

### Installation

1. **Clone the repository**
   ```bash
   git clone git@github.com:hlaclau/wooper-discord-bot.git
   cd wooper-bot
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure environment**

   **Option A: Using .env file (recommended)**
   ```bash
   cp .env.example .env
   # Edit .env and add your Discord bot token
   ```

   **Option B: Using environment variables**
   ```bash
   export DISCORD_BOT_TOKEN=your_bot_token_here
   ```

4. **Run the bot**
   ```bash
   go run ./...
   ```

## Creating a Discord Bot

1. Go to the [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application" and give it a name
3. Go to the "Bot" section
4. Click "Add Bot"
5. Copy the bot token (this is your `DISCORD_BOT_TOKEN`)
6. Under "Privileged Gateway Intents", enable "Message Content Intent"
7. Go to the "OAuth2" > "URL Generator" section
8. Select "bot" scope and "Send Messages" permission
9. Use the generated URL to invite your bot to a server

## Project Structure

```
wooper-bot/
├── .env.example          # Environment template
├── .env                  # Your environment variables (gitignored)
├── go.mod               # Go module file
├── go.sum               # Go module checksums
├── main.go              # Application entry point
├── README.md            # This file
├── img/                 # Image directories
│   ├── wooper/          # Wooper images
│   │   ├── wooper1.jpg
│   │   ├── wooper2.jpg
│   │   └── Wooper_anime.webp
│   └── cats/            # Cat images (example)
│       └── README.txt
└── internal/            # Internal packages
    ├── bot/             # Discord bot wrapper
    │   └── bot.go
    ├── config/          # Configuration management
    │   └── config.go
    ├── handlers/        # Message event handlers
    │   └── messages.go
    ├── logger/          # Structured logging with Zap
    │   └── logger.go
    └── services/        # Business logic services
        └── image.go
```

## Architecture

The bot follows a clean, layered architecture:

- **`internal/config`**: Environment variable loading with `.env` support
- **`internal/logger`**: Structured logging configuration and initialization
- **`internal/services`**: Business logic for local image management and category discovery
- **`internal/handlers`**: Discord message event processing with dynamic command support and comprehensive logging
- **`internal/bot`**: Discord session management and lifecycle
- **`main.go`**: Dependency injection and application startup

## Dependencies

- **discordgo**: Discord API client for Go
- **godotenv**: Environment variable loading from `.env` files
- **zap**: High-performance structured logging

## Development

### Building

```bash
go build -o wooper-bot .
```

### Running Tests

```bash
go test ./...
```

### Code Quality

The project follows Go best practices:
- Clean architecture with separated concerns
- Interface-based design for testability
- Proper error handling and context usage
- Graceful shutdown handling

## Troubleshooting

### Bot doesn't respond to commands

1. Ensure the bot has "Message Content Intent" enabled in Discord Developer Portal
2. Check that the bot has "Send Messages" permission in the server
3. Verify the bot token is correct in your `.env` file
4. Use `!help` to see available image categories

### No images found for a category

- Check that the category folder exists in the `img/` directory
- Ensure the folder contains supported image files (`.png`, `.jpg`, `.jpeg`, `.gif`, `.webp`)
- Verify file permissions allow the bot to read the images

### Bot shows "no image categories available"

- Ensure the `img/` directory exists and contains at least one subdirectory with images
- Check that image files have supported extensions
- Verify the bot has read permissions for the `img/` directory

## License

This project is open source and available under the MIT License.
