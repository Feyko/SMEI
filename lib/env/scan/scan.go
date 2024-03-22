package scan

import (
	"github.com/satisfactorymodding/SMEI/lib/env/project"
	"github.com/satisfactorymodding/SMEI/lib/env/ue"
	"github.com/satisfactorymodding/SMEI/lib/env/vs"
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
