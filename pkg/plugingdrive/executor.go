package plugingdrive

import (
	"fmt"
	"io"
	"log"

	_ "embed"

	"github.com/yrn-go/yrn/module/flowmanager"
	"github.com/yrn-go/yrn/pkg/plugincore"
	"github.com/yrn-go/yrn/pkg/yctx"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	Slug = "google-drive"
)

var (
	_ flowmanager.PluginExecutor = (*Executor)(nil)

	//go:embed schema.json
	Schema []byte
)

type DriveSchema struct {
	Credentials   string `json:"credentials"`
	FolderID      string `json:"folderId"`
	SharedDriveID string `json:"sharedDriveId"`
}

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Do(ctx *yctx.Context, schemaInputs string, previousPluginResponse any, responseSharedForAll map[string]any) (output any, err error) {
	var (
		requestData *DriveSchema
	)

	// Valida e carrega o schema
	requestData, err = plugincore.ValidateAndGetRequestBody[DriveSchema](
		Schema,
		schemaInputs,
		previousPluginResponse,
		responseSharedForAll,
	)
	if err != nil {
		return
	}

	// Cria serviço do Google Drive
	driveService, err := drive.NewService(
		ctx.Context(),
		option.WithCredentialsJSON([]byte(requestData.Credentials)),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar o drive service: %w", err)
	}

	query := fmt.Sprintf("'%s' in parents and trashed = false", requestData.FolderID)

	// Lista arquivos da pasta
	resp, err := driveService.Files.List().
		Q(query).
		Corpora("drive").
		DriveId(requestData.SharedDriveID).
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		Fields("files(id, name, mimeType, modifiedTime)").
		Do()
	if err != nil {
		return nil, fmt.Errorf("erro ao listar arquivos: %w", err)
	}

	var result []map[string]any

	for _, file := range resp.Files {
		log.Printf("Baixando: %s", file.Name)

		fileData, err := driveService.Files.Get(file.Id).SupportsAllDrives(true).Download()
		if err != nil {
			log.Printf("erro ao baixar o arquivo %s: %v", file.Name, err)
			continue
		}

		defer fileData.Body.Close()

		content, err := io.ReadAll(fileData.Body)
		if err != nil {
			log.Printf("erro ao ler o conteúdo de %s: %v", file.Name, err)
			continue
		}

		result = append(result, map[string]any{
			"id":       file.Id,
			"name":     file.Name,
			"mimeType": file.MimeType,
			"content":  string(content), // pode mudar pra base64 se quiser
		})
	}

	output = result
	return
}
