# Recon — Mighty No. 9 (Wii U)

Source: USA `MN9Game.rpx` (zlib sections decompressed with `recon/unrpx.py`) +
kinnay's `nexwiiu.json`.

- Unreal Engine 3 build. NEX statically linked; `nn::act::AcquireNexServiceToken`
  with log string `GameID:%08x Success:%d AcquireNexServiceTokenResult:%08x`.
- UE3 Wii U online layer: `UOnlineWiiUSystem` → `OnlineWiiUNexSystemWiiU`
  (Bind{Ranking,DataStore,Matchmake}Service, GetRankingList, DeleteRankingScore,
  FindDataStore/DownloadDataStore, Create/FindMatchmakeSession, Login/Logout
  MatchmakeCommunity, GetFriendList, ConnectFP/ConnectNGS, GetMiiNickName).
- The `AcquireNexServiceToken` call site passes `r4 = *(this+0xEC)`: the game
  server ID is a field of the online-system object, set elsewhere (likely from
  a config), **not** a code literal — so region variants are unresolved.
- Kinnay: id `0x101c5a00`, key `bb980c2e`, build `3_9_19_2005_0`.

`recon/*.bin`, `strings.txt` are derived from the proprietary executable and
are git-ignored; regenerate with `python recon/unrpx.py <rpx> out.bin`.

## What a real console sent (USA, hardware)

Right after `Register`: MatchmakeExtension `GetFriendNotificationData` /
`GetlstFriendNotificationData`, then Ranking `UploadCommonData` (the player's
Mii nickname as UTF-16BE, padded), `GetRanking` (mode 4, category 100) and
`GetCachedTopXRanking`. Finishing a stage produced a Ranking `UploadScore`
(category 203, ascending order). The access token's title ID was
`00050000101C5A00`.
