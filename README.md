# InfluxDB to Telegram Notifier

This application provides an HTTP endpoint to receive notifications from InfluxDB and forward them to a Telegram channel.

## Configuration

The following environment variables must be set:

- `TELEGRAM_BOT_TOKEN`: Your Telegram bot token.
- `TELEGRAM_CHANNEL_ID`: The ID or username of your Telegram channel.

## Running the Application

### With Docker

1.  Build the Docker image:
    ```sh
    docker build -t telegram-notifier .
    ```

2.  Run the Docker container:
    ```sh
    docker run -d -p 8080:8080 \
      -e TELEGRAM_BOT_TOKEN="your_bot_token" \
      -e TELEGRAM_CHANNEL_ID="your_channel_id" \
      --name telegram-notifier \
      telegram-notifier
    ```

### Locally

1.  Install dependencies:
    ```sh
    go mod tidy
    ```

2.  Run the application:
    ```sh
    export TELEGRAM_BOT_TOKEN="your_bot_token"
    export TELEGRAM_CHANNEL_ID="your_channel_id"
    go run main.go
    ```

## InfluxDB Configuration

1.  In your InfluxDB instance, go to "Alerts" -> "Notification Endpoints".
2.  Create a new "HTTP" endpoint.
3.  Set the URL to `http://<your_server_ip>:8080/notify`.
4.  Configure your alert rules to use this notification endpoint.
