package utils

import "github.com/spf13/viper"

func ViperGetStrOr(key, defaultVal string) string {
	conf := viper.GetString(key)
	if conf == "" {
		conf = defaultVal
	}
	return conf
}

func ViperGetIntOr(key string, defaultVal int) int {
	if viper.InConfig(key) {
		return viper.GetInt(key)
	}
	return defaultVal
}

func ViperGetBoolOr(key string, defaultVal bool) bool {
	if viper.InConfig(key) {
		return viper.GetBool(key)
	}
	return defaultVal
}
