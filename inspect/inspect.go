package inspect

import (
	"ariga.io/atlas-go-sdk/atlasexec"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const cobraErrorPrefix = "Error: "

type InspectParams struct {
	atlasexec.SchemaInspectParams
	AtlasCliPath string
}

func Inspect(ctx context.Context, params *InspectParams) ([]Schema, error) {
	atlasClient, err := atlasexec.NewClient("", params.AtlasCliPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create atlas client: %w", err)
	}
	raw, err := atlasClient.SchemaInspect(ctx, &params.SchemaInspectParams)
	if err != nil {
		if strings.HasPrefix(err.Error(), cobraErrorPrefix) {
			// avoid printing cobra error prefix twice...
			return nil, fmt.Errorf("%s", strings.TrimPrefix(err.Error(), cobraErrorPrefix))
		}
		return nil, err
	}
	var schemas []Schema
	err = json.Unmarshal([]byte(raw), &schemas)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}
	return schemas, nil
}
