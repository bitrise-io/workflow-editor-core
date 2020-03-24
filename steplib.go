package core

import (
	"encoding/json"
	"fmt"

	"github.com/bitrise-io/depman/pathutil"
	"github.com/bitrise-io/go-utils/fileutil"
	stepmanModels "github.com/bitrise-io/stepman/models"
)

// LoadSpec ...
func LoadSpec(specPth string) (stepmanModels.StepCollectionModel, error) {
	if exist, err := pathutil.IsPathExists(specPth); err != nil {
		return stepmanModels.StepCollectionModel{}, fmt.Errorf("failed to check if spec exists at: %s, error: %s", specPth, err)
	} else if !exist {
		return stepmanModels.StepCollectionModel{}, fmt.Errorf("spec not exists at: %s", specPth)
	}

	specBytes, err := fileutil.ReadBytesFromFile(specPth)
	if err != nil {
		return stepmanModels.StepCollectionModel{}, fmt.Errorf("failed to read spec at: %s, error: %s", specPth, err)
	}

	var spec stepmanModels.StepCollectionModel
	if err := json.Unmarshal(specBytes, &spec); err != nil {
		return stepmanModels.StepCollectionModel{}, fmt.Errorf("failed to serialize spec, error: %s", err)
	}

	return spec, nil
}

// IsLibrarySetup ...
func IsLibrarySetup(libraryURI string, libraryInfos []stepmanModels.SteplibInfoModel) bool {
	for _, libraryInfo := range libraryInfos {
		if libraryInfo.URI == libraryURI {
			return true
		}
	}
	return false
}

// LibraryInfo ...
func LibraryInfo(URI string, infos []stepmanModels.SteplibInfoModel) (stepmanModels.SteplibInfoModel, bool) {
	for _, info := range infos {
		if info.URI == URI {
			return info, true
		}
	}
	return stepmanModels.SteplibInfoModel{}, false
}
