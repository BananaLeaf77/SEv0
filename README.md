SEv0 is a simple Go backend that connects to WhatsApp (via WhatsMeow) and provides an HTTP API to broadcast messages.
It was built for testing and automation — **use responsibly**. Sending unsolicited messages (spam) violates WhatsApp’s Terms of Service and may be illegal in many places. You are responsible for how you use this tool.

> ⚠️ **Warning:** Do not use this for unsolicited mass messaging. Always have recipient consent and respect rate limits, anti-abuse rules, and local laws.

---

## Features

* Connects to WhatsApp using WhatsMeow and stores session data in PostgreSQL.
* Provides a `/send` HTTP endpoint to send plain or formatted WhatsApp messages to one or many recipients.
* Supports sending multiple messages and repeating messages (`repeater`).
* Prints a QR code in the terminal for initial login (scan with your WhatsApp app).

---

## Prerequisites

* Go toolchain (recommended Go 1.20+)
* PostgreSQL database (can be local or hosted)
* (Optional) Docker if you plan to containerize
* A terminal that supports ANSI colors for nice logs

---

## Environment variables

Set these environment variables before running:

```
DB_NAME=your_db_name
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=your_db_host
DB_PORT=5432
DB_SSLMODE=disable            # use `disable` for local DB; use secure mode for production
ALLOW_ORIGINS=*               # CORS origins allowed (or your frontend origin)
API_SECRET=supersecrettoken   # token used to authenticate API requests
```

### DB connection string

Your app expects a Postgres URL like:

```
postgres://DB_USER:DB_PASSWORD@DB_HOST:DB_PORT/DB_NAME?sslmode=disable
```

Your app code likely builds that from the `DB_*` vars.

---

## Quick start (local)

1. Clone repo:

   ```bash
   git clone https://github.com/youruser/SEv0.git
   cd SEv0
   ```

2. Create a `.env` file with values from **Environment variables** above (or set them in your shell).

3. Build & run:

   ```bash
   go build -o app ./...
   ./app
   ```

   * On first run the program will print a QR code in the terminal. Open WhatsApp on your phone → Settings → Linked Devices → Link a device → scan the QR code.
   * After successful login the session is stored in your PostgreSQL DB and you won’t need to rescan.

4. Test server:

   ```
   GET http://localhost:3000/ping
   ```

---

## API usage

Endpoint: `POST /send` (protected by `Authorization: Bearer <API_SECRET>` header)

### Payload options

You can send either a single repeated message or multiple messages. Example (multiple messages + repeater):

```json
{
  "to": ["6281234567890"],
  "messages": [
    "Hello *World*! 👋",
    "_This is italic_",
    "Check this link: https://example.com"
  ],
  "repeater": 2
}
```

* `to`: array of phone numbers in international format (no `+`), e.g. `6281234567890`.
* `messages`: list of message texts. WhatsApp inline formatting works: `*bold*`, `_italic_`, `~strike~`, `` `mono` ``.
* `repeater`: how many times to send each message (optional; default `1`).

### Example `curl`

```bash
curl -X POST "http://localhost:3000/send" \
  -H "Authorization: Bearer supersecrettoken" \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["6281234567890"],
    "messages": ["Hello *World*!", "Second message _italic_"],
    "repeater": 1
  }'
```

### Timeout & failure

* The API supports a timeout (configured in code) — if sending takes too long, the request returns 500 with an error message.
* The app logs colorized statuses (200 green, 400 yellow, 500 red) for easier debugging.

---

## Security & best practices

* Never hardcode `API_SECRET` in code — use environment variables.
* Rate-limit your requests and implement per-client API keys if you expose the endpoint to multiple users.
* Keep your WhatsApp client usage within WhatsApp’s policies; frequent/spammy behavior can get your number banned.
* Store secrets (DB password, API secret) securely and rotate keys when needed.

---

## Deployment notes

* The app is suitable for platforms that support long-running processes (the WhatsApp client must stay connected). Good hosts: Fly.io, Koyeb, Railway (requires credit card), or Deta (no card, but DB options differ).
* When deploying, make sure the web server listens on `0.0.0.0:$PORT` (many platforms require `$PORT` from environment).
* Add all env vars in your platform’s dashboard (Render/Koyeb/Fly/Deta) — they don’t read local `.env` files automatically.

---

## Docker (coming soon)

You mentioned packaging into a Docker container — good idea for portability. Typical steps:

* Add `Dockerfile` that builds the Go binary and exposes `$PORT`.
* Use environment variables for DB and API secret.
* Push the image to a registry and deploy to your hosting provider.

---

## Legal & Ethical notice (important)

This project can be used to send many messages quickly. That capability is easily abused for spam or harassment. The author is not responsible for misuse. Use this tool only for legitimate communication with people who consent to be contacted. Check local laws and WhatsApp’s Terms before using.

---

## Troubleshooting

* `pq: SSL is not enabled on the server` → add `?sslmode=disable` to your Postgres URL for local DB.
* `sql: unknown driver "postgres"` → make sure your code imports a Postgres driver (e.g. `_ "github.com/lib/pq"` or `_ "github.com/jackc/pgx/v5/stdlib"`).
* QR not scanning → ensure your terminal supports ANSI block characters (or use `qrterminal.Generate` vs `GenerateHalfBlock` to adjust size).

---

If you want, I can:

* Turn this into a ready-to-copy `README.md` file.
* Add a sample `docker-compose.yml` for local testing (Go + Postgres).
* Create a small `deploy.md` with steps for Koyeb or Fly.io (no-credit-card required options).

Which one do you want next?

