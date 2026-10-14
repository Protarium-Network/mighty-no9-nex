# Protocol coverage

What the server answers, and what we have actually seen a console send.

| Endpoint | Protocol | Status |
|---|---|---|
| Auth | Ticket Granting | Login, LoginEx, RequestTicket. Seen on hardware. |
| Secure | Secure Connection | Register. Seen on hardware. |
| Secure | Utility | Stock handlers. |
| Secure | Ranking | UploadScore, GetRanking, GetCachedTopXRanking(s), Upload/GetCommonData, DeleteCommonData. UploadScore, GetRanking and UploadCommonData seen on hardware. |
| Secure | MatchmakeExtension | Stock handlers. GetFriendNotificationData and GetlstFriendNotificationData seen on hardware; session create/join not exercised. |
| Secure | MatchMaking, MatchMakingExt, NAT Traversal | Registered with stock handlers, not exercised. |
| Secure | DataStore | Registered (meta, prepare get/post, complete post, ratings, search). Not called by the game so far. |

Leaderboard rows are stored in `mn9_*` tables, one row per player and
category; the matchmaking schema is created by the common library.

## Open questions

- Which DataStore methods the game uses, if any (replays, ghosts).
- Whether EUR and JPN consoles request a different game server ID.
- Two-console matchmaking and the NAT traversal path.
