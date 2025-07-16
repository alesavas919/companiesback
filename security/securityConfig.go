package security

import (
	"log"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

func ResourceSecurityData(secretRoute string) string {
	// Crear un cliente Vault
	godotenv.Load()
	config := &api.Config{
		Address: os.Getenv("SSL_ADD_PROJECT"),
	}

	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf(`No se pudo crear el cliente de Vault: ` + err.Error())
	}

	client.SetToken(os.Getenv("SSL_TOKEN_PROJECT"))

	secret, err := client.Logical().Read("secret/data/CONFIDENTIAL")
	if err != nil {
		log.Fatalf("Error al leer el secreto: %v", err)
	}

	if secret != nil && secret.Data != nil {
		secretMap, ok := secret.Data["data"].(map[string]interface{})
		if ok {
			regisInfo, ok := secretMap[secretRoute].(string)
			if ok {
				return regisInfo
			}
		}
	}
	return ""
}
