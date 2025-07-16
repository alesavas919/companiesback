package security

import (
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

func ResourceSecurityData(secretRoute string) (string, error) {
	// Crear un cliente Vault
	godotenv.Load()
	config := &api.Config{
		Address: os.Getenv("SSL_ADD_PROJECT"),
	}

	client, err := api.NewClient(config)
	if err != nil {
		return "", err
	}

	client.SetToken(os.Getenv("SSL_TOKEN_PROJECT"))

	secret, err := client.Logical().Read("secret/data/CONFIDENTIAL")
	if err != nil {
		return "", err
	}

	if secret != nil && secret.Data != nil {
		secretMap, ok := secret.Data["data"].(map[string]interface{})
		if ok {
			regisInfo, ok := secretMap[secretRoute].(string)
			if ok {
				return regisInfo, nil
			}
		}
	}
	return "", nil
}
