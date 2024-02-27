package daemon

import (
	"encoding/json"
	"os"
)

var stdout = json.NewEncoder(os.Stdout)
