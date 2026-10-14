package globals

// NEX configuration for "Mighty No. 9" (Wii U, Comcept / Deep Silver, 2016).
//
// GameServerID and AccessKey come from kinnay.github.io's public Wii U NEX
// database (id 270,... = 0x101c5a00). The retail executable (MN9Game.rpx, an
// Unreal Engine 3 build using Nintendo's OnlineWiiUNexSystem) reads the game
// server ID from a field of its online-system object, not from a literal, so
// whether other regions use a different ID is NOT settled - watch the token
// service logs for other game_server_id values.
const (
	GameServerID = "101c5a00"
	AccessKey    = "bb980c2e"

	// kinnay build tag "build:3_9_19_2005_0" (same NEX build as Rio 2016 /
	// Mario Tennis: Ultra Smash / Fast Racing NEO).
	NEXMajor = 3
	NEXMinor = 9
	NEXPatch = 19
)
