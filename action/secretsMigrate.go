package action

import (
	"fmt"

	"github.com/carsdotcom/swarm-tool/types"
)

func MigrateSecrets(source, target *types.SwarmConnection) {
	fmt.Printf("Migrating secrets from [%v] to [%v]\n", source, target)
}
