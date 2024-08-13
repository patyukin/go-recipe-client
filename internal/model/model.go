package model

type Recipe struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Instructions string `json:"instructions"`
	CreatedAt    string `json:"created_at"`
}

type RecipesResponse struct {
	Recipes []Recipe `json:"recipes"`
	Total   int      `json:"total"`
}
