package test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/yrn-go/yrn/module/flowmanager"
	"github.com/yrn-go/yrn/pkg/pluginhttp"
	"github.com/yrn-go/yrn/pkg/pluginmapper"
	"github.com/yrn-go/yrn/pkg/yctx"
	"golang.org/x/exp/slog"
)

func TestFlowExecutor(t *testing.T) {
	suite.Run(t, new(FlowExecutorTestSuite))
}

type FlowExecutorTestSuite struct {
	suite.Suite
	flowReaderRepositoryMock *flowmanager.FlowReaderRepositoryMock
	pluginManager            *pluginmapper.PluginManagerLocal
	statusRepoMock           *flowmanager.PluginStatusRepositoryMock
	flowExecutor             *flowmanager.FlowExecutor
}

func (s *FlowExecutorTestSuite) SetupTest() {
	s.flowReaderRepositoryMock = new(flowmanager.FlowReaderRepositoryMock)
	s.pluginManager = pluginmapper.NewPluginManagerLocal()
	s.statusRepoMock = new(flowmanager.PluginStatusRepositoryMock)
	s.flowExecutor = flowmanager.NewFlowExecutor(s.flowReaderRepositoryMock, s.pluginManager, s.statusRepoMock)
}

func (suite *FlowExecutorTestSuite) TearDownTest() {}

func (suite *FlowExecutorTestSuite) TearDownSuite() {}

func (suite *FlowExecutorTestSuite) TestExecute_WithSuccess() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/data":
			requestBody, _ := io.ReadAll(r.Body)

			var requestData any
			_ = json.Unmarshal(requestBody, &requestData)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(struct {
				Message string `json:"message"`
				Request any    `json:"request"`
			}{
				Message: "success",
				Request: requestData,
			})
			break
		default:
			http.NotFound(w, r)
		}
	}))

	defer mockServer.Close()

	var (
		ctx          = yctx.NewContext(context.Background())
		flowExecutor = flowmanager.NewFlowExecutor(
			suite.flowReaderRepositoryMock,
			suite.pluginManager,
			suite.statusRepoMock,
		)
		flowId           = "flow-test-id"
		eventRequestData = map[string]any{
			"user_email": "test@example.com",
		}
		test01SchemaBodyObj = pluginhttp.HTTPSchema{
			Request: pluginhttp.HTTPRequest{
				Method: http.MethodGet,
				URL:    mockServer.URL + "/data",
				Body: map[string]interface{}{
					"name":  "Test 01",
					"email": "{{.data.user_email}}",
				},
			},
		}
		test02SchemaBodyObj = pluginhttp.HTTPSchema{
			Request: pluginhttp.HTTPRequest{
				Method: http.MethodGet,
				URL:    mockServer.URL + "/data",
				Body: map[string]interface{}{
					"name":  "{{ .sharedForAll.test_01.request.name }}",
					"email": "{{ with .data }}{{ .request.name }}{{ end }}",
				},
			},
		}
		test01SchemaBody, _ = json.Marshal(test01SchemaBodyObj)
		test02SchemaBody, _ = json.Marshal(test02SchemaBodyObj)
		flowInfo            = &flowmanager.Flow{
			Id:               flowId,
			FirstPluginToRun: "test_01",
			Plugins: []flowmanager.FlowPlugin{
				{
					Id:                          "test_01",
					Slug:                        pluginhttp.SlugHttp,
					SchemaInput:                 string(test01SchemaBody),
					ShareResponseWithAllPlugins: true,
					NextToBeExecuted:            []string{"test_02"},
				},
				{
					Id:          "test_02",
					Slug:        pluginhttp.SlugHttp,
					SchemaInput: string(test02SchemaBody),
				},
			},
		}
	)

	suite.flowReaderRepositoryMock.
		On("GetById", mock.Anything, flowId).
		Return(flowInfo)

	// Configura as expectativas para o repositório de status
	// Cada plugin terá duas chamadas de Save (início e fim da execução)
	suite.statusRepoMock.
		On("Save", mock.Anything, mock.Anything).
		Return(nil).
		Times(4) // 2 plugins * 2 chamadas cada

	response, err := flowExecutor.Do(ctx, flowId, eventRequestData)
	suite.NoError(err)

	slog.Info("response", slog.Any("response", response))
}
