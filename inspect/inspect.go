package inspect

import (
	"ariga.io/atlas-go-sdk/atlasexec"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

const cobraErrorPrefix = "Error: "
const jsonFormat = "json"

type Params struct {
	atlasexec.SchemaInspectParams
	FromFilePath string
	AtlasCliPath string
}

func Inspect(ctx context.Context, params *Params) (*Data, error) {
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
	return unmarshal([]byte(raw))
}

func LoadFromFile(fpath string) (*Data, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file at '%s': %w", fpath, err)
	}
	defer f.Close()
	raw, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file at '%s': %w", fpath, err)
	}
	return unmarshal(raw)
}

func unmarshal(raw []byte) (*Data, error) {
	var data Data
	err := json.Unmarshal(raw, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}
	return &data, nil
}
