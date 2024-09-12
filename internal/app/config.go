package app

import (
	"github.com/spf13/viper"
)

type Config struct {
	AWSProfile       string
	AWSTargetProfile string
	AWSRegion        string
}

var cfg Config

// Get new config with precedence flags > environment variabls > configuration file > default value
func NewConfig() *Config {
	awsProfile := viper.GetString("aws-profile")
	// if !ok {
	// 	panic("could not set aws-profile")
	// }

	awsTargetProfile, ok := viper.Get("aws-target-profile").(string)
	if !ok {
		panic("could not set aws-target-profile")
	}

	awsRegion, ok := viper.Get("aws-region").(string)
	if !ok {
		panic("could not set aws-region")
	}

	return &Config{
		AWSProfile:       awsProfile,
		AWSTargetProfile: awsTargetProfile,
		AWSRegion:        awsRegion,
	}
}
