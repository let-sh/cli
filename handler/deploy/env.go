package deploy

import (
	"github.com/joho/godotenv"
)

func (c *DeployContext) LoadEnvFiles() error {
	if FileExists(".env") {
		err := godotenv.Load(".env")
		if err != nil {
			return err
		}
		c.LetConfig.Env, err = godotenv.Read()
		if err != nil {
			return err
		}
	}
	return nil
}
