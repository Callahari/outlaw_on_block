package databases

import (
	"outlaw_on_block/player"
	"outlaw_on_block/tiles"
)

type SaveManifest struct {
	SaveName     string
	PlayerObject *player.Player
	MapTiles     []tiles.Tile
}
