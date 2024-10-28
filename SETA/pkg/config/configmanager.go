package config

import "os"

type ConfigManager struct {
	configModel ConfigModel
}

type ConfigModel struct {
	GatewayAEndpoint string
	GatewayBEndpoint string
	DatabaseDSN      string
}

func GetConfigManager() *ConfigManager {
	return &ConfigManager{
		ConfigModel{
			GatewayAEndpoint: os.Getenv("GATEWAY_A_ENDPOINT"),
			GatewayBEndpoint: os.Getenv("GATEWAY_B_ENDPOINT"),
			DatabaseDSN:      os.Getenv("DATABASE_DSN"),
		},
	}
}

func (cm *ConfigManager) GetGatewayAEndpoint() string {
	return cm.configModel.GatewayAEndpoint
}

func (cm *ConfigManager) GetGatewayBEndpoint() string {
	return cm.configModel.GatewayBEndpoint
}

func (cm *ConfigManager) GetDatabaseDSN() string {
	return cm.configModel.DatabaseDSN
}
