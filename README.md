# LinkedIn Automation Tool (Go/Rod)

A comprehensive technical proof-of-concept demonstrating advanced browser automation, anti-detection techniques, and clean Go architecture. This tool sends connection requests and messages while mimicking human behavior to avoid bot detection.

## 🚀 Features

-   **Modular Architecture**: Clean separation of concerns (Auth, Search, Connection, Stealth).
-   **Stealth Mode**:
    -   **Bezier Mouse**: Human-like mouse curvature and speed implementation.
    -   **Fingerprint Masking**: Automatic masking of WebDriver flags.
    -   **Randomized Timing**: Dynamic sleep intervals to simulate "think time".
-   **Robost Authentication**:
    -   Env-based credentials.
    -   Cookie persistence for session reuse.
    -   Basic checkpoint detection.
-   **Smart Targeting**:
    -   Scrapes search results securely.
    -   Handles pagination and "Connect" vs "More" button variants.

## 🛠️ Setup

### Prerequisites
-   [Go 1.21+](https://go.dev/dl/)
-   Google Chrome installed.

### Installation

1.  Clone the repository:
    ```bash
    git clone https://github.com/yourusername/linkedin-automation.git
    cd linkedin-automation
    ```

2.  Install dependencies:
    ```bash
    go mod tidy
    ```

3.  Configure Environment:
    Copy the example file and fill in your credentials:
    ```bash
    cp .env.example .env
    ```
    *Note: You can also configure settings in `config.yaml`.*

## 🏃 Usage

Build the binary:
```bash
go build -o linkedin-bot.exe ./cmd/bot
```

### 1. Search and Connect
Search for a job title and send connection requests with a limit.
```bash
./linkedin-bot.exe -task search-connect -keywords "Software Engineer" -limit 5
```

### 2. Dry Run / Debug
Verify browser behavior without running actions (requires editing config to Headless: false).
```bash
./linkedin-bot.exe
```

## 🎥 Demonstration

[Insert Video Link Here]

*This video demonstrates the tool's setup, configuration, and execution flow, highlighting the human-like mouse movements and stealth navigation.*

## ⚠️ Disclaimer
This tool is for educational and testing purposes only. Automated interaction with LinkedIn violates their Terms of Service. Use responsibly and at your own risk.
