package test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/yrn-go/yrn/module/flowmanager"
	"github.com/yrn-go/yrn/pkg/pluginhttp"
	"github.com/yrn-go/yrn/pkg/pluginmapper"
	"github.com/yrn-go/yrn/pkg/yctx"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlowExecutor(t *testing.T) {
	suite.Run(t, new(FlowExecutorTestSuite))
}

type FlowExecutorTestSuite struct {
	suite.Suite
	flowReaderRepository *flowmanager.FlowReaderRepositoryMock
}

func (suite *FlowExecutorTestSuite) SetupTest() {
	suite.flowReaderRepository = new(flowmanager.FlowReaderRepositoryMock)
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
			suite.flowReaderRepository,
			pluginmapper.NewPluginManagerLocal(),
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
					//"from":  "{{.sharedForAll.telegram.text_message}}",
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
					//"from":  "{{.sharedForAll.telegram.text_message}}",
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

	suite.flowReaderRepository.
		On("GetById", mock.Anything, flowId).
		Return(flowInfo)

	response, err := flowExecutor.Do(ctx, flowId, eventRequestData)
	suite.NoError(err)

	slog.Info("response", slog.Any("response", response))
}
