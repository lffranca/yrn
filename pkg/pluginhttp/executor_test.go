package pluginhttp

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"github.com/yrn-go/yrn/pkg/yctx"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecutor(t *testing.T) {
	suite.Run(t, &ExecutorTestSuite{})
}

type ExecutorTestSuite struct {
	suite.Suite
}

func (suite *ExecutorTestSuite) SetupTest() {}

func (suite *ExecutorTestSuite) TearDownTest() {}

func (suite *ExecutorTestSuite) TearDownSuite() {}

func (suite *ExecutorTestSuite) TestDo_WithSuccess() {
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
		userEmail              = "john.doe@yrn.com"
		telegramTextMessage    = "response_shared_for_all"
		previousPluginResponse = map[string]any{
			"token":      "token_test_TestDo_WithSuccess",
			"user_email": userEmail,
		}
		responseSharedForAll = map[string]any{
			"telegram": map[string]any{
				"text_message": telegramTextMessage,
			},
		}
		schemaBody = HTTPSchema{
			Request: HTTPRequest{
				Method: "POST",
				URL:    mockServer.URL + "/data",
				Headers: map[string]string{
					"Content-Type":  "application/json",
					"Authorization": "Bearer {{.data.token}}",
				},
				QueryParams: map[string]string{
					"limit":  "10",
					"offset": "0",
				},
				Body: map[string]interface{}{
					"name":  "John Doe",
					"email": "{{.data.user_email}}",
					"from":  "{{.sharedForAll.telegram.text_message}}",
					"age":   30,
				},
				Timeout: 5000,
			},
			Retry: &RetryConfig{
				MaxAttempts: 3,
				Delay:       1000,
			},
		}
		body, _ = json.Marshal(schemaBody)
	)

	ctx := yctx.NewContext(context.Background())

	executor := NewExecutor()

	response, err := executor.Do(ctx, string(body), previousPluginResponse, responseSharedForAll)
	suite.NoError(err)

	responseMap, responseMapOk := response.(map[string]any)
	suite.True(responseMapOk)

	requestMap, requestMapOk := responseMap["request"].(map[string]any)
	suite.True(requestMapOk)

	emailString, emailStringOk := requestMap["email"].(string)
	suite.True(emailStringOk)
	suite.Equal(userEmail, emailString)

	telegramTextMessageString, telegramTextMessageStringOk := requestMap["from"].(string)
	suite.True(telegramTextMessageStringOk)
	suite.Equal(telegramTextMessage, telegramTextMessageString)
}
