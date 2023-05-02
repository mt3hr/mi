package miapp

type GetApplicationConfigResponse struct {
	Errors            []string          `json:"errors"`
	ApplicationConfig ApplicationConfig `json:"application_config"`
}
