package inspect

import (
	"ariga.io/atlas-go-sdk/atlasexec"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const cobraErrorPrefix = "Error: "
const jsonFormat = "json"

type InspectParams struct {
	atlasexec.SchemaInspectParams
	AtlasCliPath string
}

func Inspect(ctx context.Context, params *InspectParams) (*Data, error) {
	atlasClient, err := atlasexec.NewClient("", params.AtlasCliPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create atlas client: %w", err)
	}
	params.SchemaInspectParams.Format = jsonFormat
	raw, err := atlasClient.SchemaInspect(ctx, &params.SchemaInspectParams)
	if err != nil {
		if strings.HasPrefix(err.Error(), cobraErrorPrefix) {
			// avoid printing cobra error prefix twice...
			return nil, fmt.Errorf("%s", strings.TrimPrefix(err.Error(), cobraErrorPrefix))
		}
		return nil, err
	}
	var data Data
	err = json.Unmarshal([]byte(raw), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}
	return &data, nil
}
