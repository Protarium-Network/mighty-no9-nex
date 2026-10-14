# mighty-no9-nex

A NEX server for the Wii U version of Mighty No. 9. It handles the game's
login, leaderboards and matchmaking calls so the online features work again
now that the original servers are gone.

It is built on the Pretendo Network libraries (`nex-go`, `nex-protocols-go`,
`nex-protocols-common-go`) and started from our other Wii U servers, mostly
[dkc-tropical-freeze-nex](https://github.com/Protarium-Network/dkc-tropical-freeze-nex).

## What works

Tested on one USA console against a live deployment:

- Login: token request, `LoginEx`, ticket, secure `Register`.
- Leaderboards: the game uploads its score, reads the board back and shows the
  right player name.
- The matchmaking requests the game sends right after connecting
  (friend notification data) are answered.

Not tested yet: actually playing a multiplayer session between two consoles,
DataStore uploads (replays/ghosts, if the game uses them), and any region
other than USA.

## How the game talks to the server

Mighty No. 9 is an Unreal Engine 3 game and uses Unreal's generic Wii U online
layer, which links the whole NEX stack. Linking a protocol does not mean the
game calls it, so we started with ranking and DataStore and added matchmaking
once the console asked for it. See [RECON.md](RECON.md) for what we found in
the executable and [PROTOCOL_COVERAGE.md](PROTOCOL_COVERAGE.md) for the list of
protocols this server answers.

| | |
|---|---|
| Game server ID | `101c5a00` |
| Access key | `bb980c2e` |
| NEX version | 3.9.19 |

The ID and key come from [Kinnay's NEX database](https://kinnay.github.io/).
The game reads its server ID from a runtime field rather than a constant, so
EUR and JPN copies may ask for a different ID. If yours does, the account
server's `nex_token` route needs that ID too; open an issue with the
`game_server_id` you see.

Wire settings: structure headers on, modern PRUDPv1 CONNECT-ACK signature
(`LegacyConnectionSignature = false`).

## Running it

You need Go 1.25+ (or Docker) and PostgreSQL.

```
cp .env.example .env
cp settings.example.json settings.json
docker compose up --build
```

The auth server listens on UDP 28100 and the secure server on UDP 28101.

Two ways to authenticate players:

- **Local mode** (`PN_MN9_LOCAL_MODE=1`, the default): accounts come from
  `settings.json` and tokens are not checked. Only use this on a network you
  control.
- **Shared mode** (`PN_MN9_LOCAL_MODE=0`): set `PN_MN9_NEX_TOKEN_AES_KEY` and
  `PN_MN9_NEX_PASSWORD_SECRET` to the same values your account server uses, so
  the tokens it issues are accepted. [account-server](https://github.com/Protarium-Network/account-server)
  does this, or you can use Pretendo's.

`PN_MN9_SECURE_HOST` is the address consoles are told to connect to. Keep it
short, the game copies it into a small buffer.

DataStore uploads need an S3-compatible bucket (`PN_S3_*` in `.env.example`).
Without one, downloads still work and uploads fail.

The game uploads its own player name with `Ranking::UploadCommonData` (UTF-16)
and the server hands it back untouched.

## License

AGPL-3.0, see [LICENSE](LICENSE) and [NOTICE](NOTICE). No Nintendo or
Comcept code or assets are included.
