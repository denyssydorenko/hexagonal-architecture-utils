package health

import "hexagonal-architexture-utils/config"

type Health struct {
	Version           string             `json:"version"`
	Healthy           bool               `json:"healthy"`
	Host              string             `json:"host"`
	ApplicationConfig config.AppConfig   `json:"application"`
	InfraConfig       config.InfraConfig `json:"infrastructure"`
}
