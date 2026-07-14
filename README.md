# localtld

**Give your local projects real domains on dynamic ports.** Stop chasing ports — reach your app at `panel.aaron.localtld`.

```jsonc
// package.json
{ "localtld": "panel.aaron" }
```

```bash
localtld run -- pnpm dev
#   → http://panel.aaron.localtld   (dynamic port — you never need to know it)
```

Whatever port your dev server grabs (3000, 51234, doesn't matter), you always reach it at the same clean domain. Run hundreds of projects at once: no port collisions, nothing to track.

## Zero-config, even if you've never heard of localtld

Wire it once in a project and it *just works* on any machine that has localtld set up — the developer doesn't need to know localtld exists:

```jsonc
// package.json
{
  "localtld": "panel.aaron",
  "scripts": {
    "dev": "localtld run -- next dev"
  },
  "devDependencies": {
    "localtld": "^0.1.0"
  }
}
```

Now the normal workflow does everything:

```bash
pnpm install     # localtld binary lands in node_modules/.bin
pnpm dev         # → http://panel.aaron.localtld
```

- **Machine has localtld configured** → the dev server comes up at `panel.aaron.localtld`.
- **Machine does *not* have it** → `localtld run` transparently falls back and runs the command as-is (`localhost:PORT`). The project still works; nobody is blocked.

So a teammate who clones the repo and runs `pnpm dev` either sees the pretty domain (if their machine is set up) or plain `localhost` — never an error.

## How it works

localtld doesn't reinvent anything; it orchestrates two standard tools:

```
browser → panel.aaron.localtld
   │  dnsmasq:  *.localtld → 127.0.0.1        (+ macOS /etc/resolver)
   ▼
127.0.0.1:80 → Caddy  (Host header → the right port)
   ▼
127.0.0.1:51234 ← your dev server (localtld assigned it via PORT env)
```

- **TLD is a machine-level setting** (default `.localtld`) — everyone picks their own.
- **Label is a project-level setting** (`package.json`) — one line, no TLD in it, so the repo stays portable.
- **Opt-in & graceful**: without a system setup, projects just run on `localhost:PORT`.

## Install

```bash
brew install localtld        # (coming soon — primary channel)
# or
curl -fsSL https://localtld.sh | bash

localtld setup               # pick a TLD, configure dnsmasq + Caddy (asks for sudo)
```

Requires macOS + Homebrew. `setup` installs `dnsmasq`, `caddy`, and `jq` if missing.

## Commands

| Command | What it does |
|---------|--------------|
| `localtld setup` | First-time setup: pick a TLD + configure dnsmasq/Caddy |
| `localtld run -- <cmd>` | Run a project under its domain (falls back if not set up) |
| `localtld list` | Show active projects and their domains |
| `localtld tld <new>` | Change the machine-wide TLD (e.g. `localtld tld test`) |
| `localtld doctor` | Check the health of your setup |
| `localtld uninstall` | Revert DNS/route changes |

## Choosing / changing your TLD

The default `.localtld` is **not** a real TLD, so it can never collide with a live website. You can switch to something shorter or branded:

```bash
localtld tld test         # panel.aaron.test
localtld tld localtld.sh  # panel.aaron.localtld.sh
```

> ⚠️ If you pick a **real** TLD (`.com`, `.dev`, …), every site under it resolves to
> `127.0.0.1` on your machine — you'd lose access to the real ones. `localtld` checks
> your choice against the IANA root zone and makes you confirm before doing that.

Guaranteed-safe options: `.localtld`, `.test`, `.localhost` (RFC-reserved — never real).

## License

MIT © Abdullah Harun Öztürk
