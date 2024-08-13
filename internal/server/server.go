package server

import (
	"app/internal/config"
	"app/internal/model"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

type Server struct {
	HostClient     *fasthttp.HostClient
	PipelineClient *fasthttp.PipelineClient
	cfg            *config.Config
}

func New(cfg *config.Config, hostClient *fasthttp.HostClient, pipelineClient *fasthttp.PipelineClient) *Server {
	return &Server{
		HostClient:     hostClient,
		PipelineClient: pipelineClient,
		cfg:            cfg,
	}
}

func (s *Server) GetAllRecipes(ctx *fasthttp.RequestCtx) {
	initialRequest := fasthttp.AcquireRequest()
	initialRequest.SetRequestURI(fmt.Sprintf("http://%s:%d/recipes?limit=1&offset=0", s.cfg.RemoteHost, s.cfg.RemotePort))
	initialRequest.Header.SetMethod(fasthttp.MethodGet)

	initialResponse := fasthttp.AcquireResponse()
	err := s.HostClient.Do(initialRequest, initialResponse)
	if err != nil {
		writeErrorResponse(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf("Failed to make initial request: %v", err))
		return
	}
	defer fasthttp.ReleaseRequest(initialRequest)
	defer fasthttp.ReleaseResponse(initialResponse)

	var initialRecipesResponse model.RecipesResponse
	if err = json.Unmarshal(initialResponse.Body(), &initialRecipesResponse); err != nil {
		writeErrorResponse(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf("Failed to parse initial response: %v", err))
		return
	}

	totalRecipes := initialRecipesResponse.Total
	fmt.Printf("Total recipes available: %d\n", totalRecipes)

	var allRecipes []model.Recipe
	limit := 10
	for offset := 0; offset < totalRecipes; offset += limit {
		request := fasthttp.AcquireRequest()
		request.SetRequestURI(fmt.Sprintf(fmt.Sprintf("http://%s:%d/recipes?limit=%d&offset=%d", s.cfg.RemoteHost, s.cfg.RemotePort, limit, offset)))
		request.Header.SetMethod(fasthttp.MethodGet)

		response := fasthttp.AcquireResponse()
		err = s.PipelineClient.Do(request, response)
		if err != nil {
			writeErrorResponse(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf("Failed to make paginated request: %v", err))
			return
		}

		// Парсинг рецептов
		var paginatedRecipesResponse model.RecipesResponse
		if err := json.Unmarshal(response.Body(), &paginatedRecipesResponse); err != nil {
			writeErrorResponse(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf("Failed to parse paginated response: %v", err))
			return
		}

		allRecipes = append(allRecipes, paginatedRecipesResponse.Recipes...)
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}

	writeJsonResponse(ctx, fasthttp.StatusOK, allRecipes)
}

func writeJsonResponse(ctx *fasthttp.RequestCtx, statusCode int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		writeErrorResponse(ctx, fasthttp.StatusInternalServerError, "Failed to marshal response")
		return
	}

	ctx.SetStatusCode(statusCode)
	ctx.SetContentType("application/json")
	ctx.SetBody(response)
}

func writeErrorResponse(ctx *fasthttp.RequestCtx, statusCode int, message string) {
	errorResponse := map[string]string{
		"error": message,
	}
	writeJsonResponse(ctx, statusCode, errorResponse)
}
