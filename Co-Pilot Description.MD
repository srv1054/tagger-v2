# Tagger-v2 Slack Bot Summary

## Overview

**Tagger-v2** is a Slack bot written in Go that automatically adds emoji reactions ("spray cans") to Slack messages based on keyword triggers. It replaces the original RTM-based bots-tagger, now using Slack's websocket API for improved reliability and interactivity. Tagger-v2 allows dynamic management of emoji tags and keywords directly from Slack, with fine-grained control over who can modify or delete tags.

---

## Key Features

- **Automated Emoji Tagging:**  
  Monitors messages in Slack channels and adds emoji reactions when messages contain configured keywords.

- **Dynamic Tag & Keyword Management:**  
  Users can add or remove emoji tags ("spray cans") and their associated keywords from within Slack using bot commands.

- **Configurable Permissions:**  
  Deletion of tags/keywords can be restricted to a specific channel, as set in `config.json`.

- **Live Reloading:**  
  The bot can reload its tag configuration (`tags.json`) on command, without needing a restart.

- **Logging:**  
  Logs bot activity to daily log files and optionally to a Slack channel.

- **Custom Emoji Support:**  
  Includes sample emoji images for use as bot avatars or Slack emoji.

---

## How It Works

1. **Startup:**  
   Loads configuration from `config.json` and tag definitions from `tags.json`. Connects to Slack via websocket.

2. **Message Monitoring:**  
   Listens for all messages in channels where the bot is present.

3. **Tagging:**  
   When a message contains a keyword from any spray can, the bot reacts with the corresponding emoji.

4. **Slack Commands:**  
   Users can interact with the bot using commands like:
   - `@tagger help` — Get help and command list
   - `@tagger list spray cans` — List all tags and their keywords
   - `@tagger add spray can <emoji>` — Add a new emoji tag
   - `@tagger add word <spray can> "<keyword>"` — Add a keyword to a tag
   - `@tagger delete spray can <emoji>` — Delete a tag (restricted by channel)
   - `@tagger delete word <spray can> "<keyword>"` — Delete a keyword (restricted by channel)
   - `@tagger reload spray cans` — Reload tags.json

5. **Logging:**  
   All actions and errors are logged to a daily log file and, if configured, to a Slack channel.

---

## Configuration

- **config.json**  
  Contains Slack tokens, bot name, log channel, path to tags.json, debug mode, and channel restrictions for deletes.

- **tags.json**  
  Defines spray cans (emoji tags) and their associated keywords. Example:
  ```json
  [
    {
      "spray": "businesscat",
      "words": ["synergy", "low-hanging fruit", "win-win"]
    },
    {
      "spray": "drevil",
      "words": ["laser", "million dollars"]
    }
  ]
  ```

- **Emoji Directory**  
  Contains images for bot avatar and custom emoji.

---

## Command Line Options

- `-v` Show version and exit
- `-c` Path to configuration file
- `-j` Path to tags.json file

---

## Example Workflow

1. User posts "synergy" in a channel.
2. Bot reacts with :businesscat: emoji.
3. Admin adds a new spray can or keyword via Slack command.
4. Bot reloads configuration on command, applying changes instantly.

---

## References

- [README.md](x:/src/tagger-v2/README.md)
- [config.go](x:/src/tagger-v2/config.go)
- [main.go](x:/src/tagger-v2/main.go)
- [slack.go](x:/src/tagger-v2/slack.go)
- [spray.go](x:/src/tagger-v2/spray.go)
- [tags.json.example](x:/src/tagger-v2/configs/tags.json.example)
- [config.json.example](x:/src/tagger-v2/configs/config.json.example)

---

**Note:**  
Be careful with keywords to avoid excessive tagging. For best results, use