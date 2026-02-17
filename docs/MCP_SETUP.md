# Setting up MCP Server for PostgreSQL

To enable AI access to your PostgreSQL database, you need to configure an MCP (Model Context Protocol) server.

## Prerequisites

- **Node.js** installed (v16+).
- Your PostgreSQL database must be running (e.g., via `docker-compose up`).

## Configuration

Add the following configuration to your MCP settings file (e.g., `claude_desktop_config.json` for Claude Desktop, or your agent's configuration file).

### Windows Configuration

```json
{
  "mcpServers": {
    "postgres": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-postgres",
        "postgresql://postgres:postgres@localhost:5432/learngo"
      ]
    }
  }
}
```

## How It Works

1.  The `npx` command downloads and runs the `@modelcontextprotocol/server-postgres` package.
2.  It connects to your local PostgreSQL instance running on `localhost:5432` with the credentials provided (user: `postgres`, password: `postgres`, db: `learngo`).
3.  The AI assistant (Claude/Agent) communicates with this server to inspect schema and run readonly queries.

## Gemini AI Configuration

To use this MCP server with Gemini (e.g., via Gemini CLI or supported IDE extensions), you typically need to configure a `.gemini/settings.json` file in your project root.

I have automatically created this configuration for you at `.gemini/settings.json`.

```json
{
  "mcpServers": {
    "postgres": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-postgres",
        "postgresql://postgres:postgres@localhost:5432/learngo"
      ]
    }
  }
}
```

## Testing

Once configured and restarted, you can ask the AI to "Check the database schema" or "Show me the last 5 orders" to verify the connection.

## Panduan Bahasa Indonesia (Indonesian Guide)

Agar AI dapat mengakses database PostgreSQL Anda, server MCP (Model Context Protocol) telah dikonfigurasi.

### Cara Kerja

1.  Perintah `npx` akan mengunduh dan menjalankan paket `@modelcontextprotocol/server-postgres`.
2.  Server ini terhubung ke database PostgreSQL lokal Anda di `localhost:5432` dengan kredensial yang diberikan (user: `postgres`, password: `postgres`, db: `learngo`).
3.  Asisten AI (seperti Gemini atau Claude) berkomunikasi dengan server ini untuk memeriksa skema database dan menjalankan query (hanya baca).

### Cara Penggunaan

Setelah konfigurasi selesai dan editor/server dimulai ulang, Anda cukup meminta AI dalam bahasa natural.

**Contoh Perintah:**

*   "Cek skema database saya"
*   "Tampilkan 5 pesanan terakhir dari tabel orders"
*   "Siapa user yang paling sering melakukan pemesanan?"
*   "Buatkan query SQL untuk menghitung total pendapatan bulanan"

AI akan secara otomatis menggunakan alat (tool) `query` dari server MCP untuk melihat data yang ada di database Anda dan memberikan jawaban berdasarkan data tersebut.