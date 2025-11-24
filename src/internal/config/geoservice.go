package config

import (
	"github.com/spf13/viper"
	"googlemaps.github.io/maps"
)

type GeoService struct {
	Client *maps.Client
}

func NewGeoService(viper *viper.Viper) (*GeoService, error) {
	client, err := maps.NewClient(maps.WithAPIKey(viper.GetString("thirdparty.google.api_key")))
	if err != nil {
		return nil, err
	}
	return &GeoService{Client: client}, nil
}
