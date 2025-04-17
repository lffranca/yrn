package plugingdrive

//import (
//	"context"
//	"encoding/json"
//	"github.com/stretchr/testify/suite"
//	"github.com/yrn-go/yrn/pkg/yctx"
//	"log/slog"
//
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//type ExecutorTestSuite struct {
//	suite.Suite
//	mockServer *httptest.Server
//}
//
//func TestExecutor(t *testing.T) {
//	suite.Run(t, new(ExecutorTestSuite))
//}
//
//func (suite *ExecutorTestSuite) SetupTest() {
//	suite.mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		switch r.URL.Path {
//		case "/drive/v3/files":
//			// Mock da listagem de arquivos
//			files := map[string]any{
//				"files": []map[string]any{
//					{
//						"id":           "file-id-1",
//						"name":         "test-file.txt",
//						"mimeType":     "text/plain",
//						"modifiedTime": "2023-01-01T00:00:00Z",
//					},
//				},
//			}
//			w.Header().Set("Content-Type", "application/json")
//			_ = json.NewEncoder(w).Encode(files)
//
//		case "/drive/v3/files/file-id-1":
//			// Mock do conteúdo do arquivo
//			w.WriteHeader(http.StatusOK)
//			_, _ = w.Write([]byte("conteudo de teste"))
//
//		default:
//			http.NotFound(w, r)
//		}
//	}))
//}
//
//func (suite *ExecutorTestSuite) TearDownTest() {
//	suite.mockServer.Close()
//}
//
//func (suite *ExecutorTestSuite) TestDo_WithSuccess() {
//	// Montando um JSON fake com client_id/secret só para preencher o campo
//	// No teste, não vai ser usado porque estamos injetando o mock server
//	credentialsTemplate := `{
//	}`
//
//	//credentials := fmt.Sprintf(credentialsTemplate,
//	//	suite.mockServer.URL,
//	//	suite.mockServer.URL,
//	//	suite.mockServer.URL,
//	//	suite.mockServer.URL,
//	//)
//
//	input := DriveSchema{
//		Credentials:   credentialsTemplate,
//		FolderID:      "1MRMbqkgFeHxX-GRovtoGcM5nGpV4De_2",
//		SharedDriveID: "0ANV-Pf5HP2CEUk9PVA",
//	}
//	body, err := json.Marshal(input)
//	suite.NoError(err)
//
//	ctx := yctx.NewContext(context.Background())
//
//	executor := NewExecutor()
//
//	output, err := executor.Do(ctx, string(body), nil, nil)
//	suite.NoError(err)
//
//	slog.Info(
//		"output",
//		slog.Any("output", output))
//
//	//files, ok := output.([]map[string]any)
//	//suite.True(ok)
//	//suite.Len(files, 1)
//	//
//	//file := files[0]
//	//suite.Equal("test-file.txt", file["name"])
//	//suite.Equal("text/plain", file["mimeType"])
//	//suite.Equal("conteudo de teste", file["content"])
//}
