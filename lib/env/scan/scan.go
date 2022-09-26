package scan

import (
	"SMEI/lib/env/project"
	"SMEI/lib/env/ue"
	"SMEI/lib/env/vs"
)

type EnvInfo struct {
	UE      *ue.Info
	VS      *vs.Info
	Project *project.Info
}

func Scan() (EnvInfo, error) {
	var info EnvInfo
	return info, nil
}
