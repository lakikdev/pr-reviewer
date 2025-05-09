package analyze

import (
	"fmt"
	"pr-reviewer/internal/api/apiHelper"
	"pr-reviewer/internal/api/authHelper"
	"pr-reviewer/internal/common"
	prReviewer "pr-reviewer/internal/modules/prReviever"
	aiModel "pr-reviewer/internal/modules/prReviever/ai/model"
	providerModel "pr-reviewer/internal/modules/prReviever/provider/model"

	"github.com/labstack/echo/v4"
)

type API struct {
	Permissions    authHelper.Permissions
	DataRepository *common.DataRepository

	PRReviewer *prReviewer.Reviewer
}

func SetAPI(router *echo.Group, permissions authHelper.Permissions, dataRepository *common.DataRepository) {

	// Initialize the PRReviewer instance
	reviewer, err := prReviewer.NewReviewer(&prReviewer.ReviewerOptions{
		Provider: providerModel.Bitbucket,
		AIClient: aiModel.Ollama,
	})

	if err != nil {
		// Handle error
		fmt.Println("Error initializing PRReviewer:", err)
		return
	}

	api := API{
		Permissions:    permissions,
		DataRepository: dataRepository,
		PRReviewer:     reviewer,
	}

	endpoints := []apiHelper.Endpoint{
		apiHelper.NewEndpoint("/analyze", "POST", api.Analyze, authHelper.APIKey),
	}

	for _, endpoint := range endpoints {
		router.Add(endpoint.Method, endpoint.Path, permissions.Wrap(endpoint.Func, endpoint.PermissionTypes...))
	}

}
