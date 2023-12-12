package authenticate

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/Appkube-awsx/awsx-metric-cli/client"
	"github.com/Appkube-awsx/awsx-metric-cli/vault"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
)

func AuthenticateData(cloudElementId string, cloudElementApiUrl string, vaultUrl string, vaultToken string, accountNo string, region string, acKey string, secKey string, crossAccountRoleArn string, externalId string) (bool, *client.Auth, error) {
	if cloudElementId != "" {
		log.Println("cloud-element-id provided. getting user credentials. cloud-element-id: " + cloudElementId)
		if cloudElementApiUrl == "" {
			log.Println("cloud-element api url missing")
			return false, nil, fmt.Errorf("cloud-element api url missing")
		}

		apiResp, statusCode, err := vault.GetUserCredential(cloudElementId, cloudElementApiUrl)
		if err != nil {
			log.Println("call to cloud-element api failed. \n", err)
			return false, nil, err
		}
		if statusCode != http.StatusOK {
			log.Println("error in calling cloud-element api. status code: "+string(statusCode)+" \n", err)
			return false, nil, err
		}
		clientAuth := client.Auth{
			CrossAccountRoleArn: apiResp.Data.CrossAccountRoleArn,
			AccessKey:           apiResp.Data.AccessKey,
			SecretKey:           apiResp.Data.SecretKey,
			ExternalId:          apiResp.Data.ExternalId,
		}
		if region != "" {
			clientAuth.Region = region
		} else {
			log.Println("region not provided. default region will be used")
			clientAuth.Region = apiResp.Data.Region
		}
		return true, &clientAuth, nil
	}

	log.Println("vault url not provided. validating provided user credentials")
	if region == "" {
		log.Println("region missing")
		return false, nil, fmt.Errorf("region missing")
	}
	if acKey == "" {
		log.Println("access key missing")
		return false, nil, fmt.Errorf("access key missing")
	}
	if secKey == "" {
		log.Println("secret key missing")
		return false, nil, fmt.Errorf("secret key missing")
	}

	if crossAccountRoleArn != "" {
		log.Println("cloud-element api url missing")

		return false, nil, fmt.Errorf("cloud-element api url missing")
	}

	if crossAccountRoleArn == "" {
		log.Println("cross account role arn missing")
		return false, nil, fmt.Errorf("cross account role arn missing")
	}
	if externalId == "" {
		log.Println("external id missing")
		return false, nil, fmt.Errorf("external id missing")
	}
	clientAuth := client.Auth{
		Region:              region,
		CrossAccountRoleArn: crossAccountRoleArn,
		AccessKey:           acKey,
		SecretKey:           secKey,
		ExternalId:          externalId,
	}
	return true, &clientAuth, nil
}

func AuthenticateDataEnv(cloudElementId string, cloudElementApiUrl string, vaultUrl string, vaultToken string, accountNo string, region string, acKey string, secKey string, crossAccountRoleArn string, externalId string) (bool, *client.Auth, error) {

	if crossAccountRoleArn == "" {
		log.Println("crossAccountRoleArn missing")
		return false, nil, fmt.Errorf("crossAccountRoleArn missing")
	}

	log.Println("corss arn provided. getting user credentials. corss arn: " + crossAccountRoleArn)
	clientAuth := client.Auth{}
	var decryptedAccessKey = ""
	var decryptedSecrtKey = ""
	vaultResp, err := vault.GetAccountDetails(vaultUrl, vaultToken, "GLOBAL_ACCESS_SECRET_AWS")
	if err != nil {
		log.Println("call to vault api failed. \n", err)
		key := []byte("qwertyuioplkjhgfdsa1234!@MNB?>P)")
		decryptedAccessKey, err = decrypt(key, "cXdlcnR5dWlvcGxramhnZg9xwxCE6a73juwgeyrAwkSWCLlY")
		if err != nil {
			log.Fatal("Error decrypting access key:", err)
		}
		decryptedSecrtKey, err = decrypt(key, "cXdlcnR5dWlvcGxramhnZioD/Dj5jPqfkOAjHXr23l9mGA8fh/0M83gZMukSd5NTYlIwSd8o24o=")
		if err != nil {
			log.Fatal("Error decrypting secret key:", err)
		}
		fmt.Println("decryptedSecrtKeydecryptedSecrtKeydecryptedSecrtKey", decryptedSecrtKey)

		if decryptedAccessKey == "" {
			decryptedAccessKey = os.Getenv("API_KEY")
		}
		if decryptedSecrtKey == "" {
			decryptedSecrtKey = os.Getenv("SECRET_KEY")
		}

		clientAuth.CrossAccountRoleArn = crossAccountRoleArn
		clientAuth.AccessKey = decryptedAccessKey
		clientAuth.SecretKey = decryptedSecrtKey
		clientAuth.CrossAccountRoleArn = crossAccountRoleArn

		if externalId != "" {
			clientAuth.ExternalId = externalId
		} else if vaultResp.Data.ExternalId != "" {
			clientAuth.ExternalId = vaultResp.Data.ExternalId
		}
		if region != "" {
			clientAuth.Region = region
		} else {
			clientAuth.Region = "us-east-1"
		}
		return true, &clientAuth, nil
	} else {
		if vaultResp.Data.AccessKey == "" || vaultResp.Data.SecretKey == "" {
			log.Println("account details not found in vault")
			return false, nil, fmt.Errorf("account details not found in vault")
		}
	}

	clientAuth.CrossAccountRoleArn = crossAccountRoleArn
	clientAuth.AccessKey = vaultResp.Data.AccessKey
	clientAuth.SecretKey = vaultResp.Data.SecretKey

	if externalId != "" {
		clientAuth.ExternalId = externalId
	} else if vaultResp.Data.ExternalId != "" {
		clientAuth.ExternalId = vaultResp.Data.ExternalId
	}
	if region != "" {
		clientAuth.Region = region
	} else {
		clientAuth.Region = "us-east-1"
	}

	return true, &clientAuth, nil

}

func CommandAuth(cmd *cobra.Command) (bool, *client.Auth, error) {
	cloudElementId, _ := cmd.PersistentFlags().GetString("cloudElementId")
	cloudElementApiUrl, _ := cmd.PersistentFlags().GetString("cloudElementApiUrl")
	vaultUrl, _ := cmd.PersistentFlags().GetString("vaultUrl")
	vaultToken, _ := cmd.PersistentFlags().GetString("vaultToken")
	accountNo, _ := cmd.PersistentFlags().GetString("accountId")
	region, _ := cmd.PersistentFlags().GetString("zone")
	acKey, _ := cmd.PersistentFlags().GetString("accessKey")
	secKey, _ := cmd.PersistentFlags().GetString("secretKey")
	crossAccountRoleArn, _ := cmd.PersistentFlags().GetString("crossAccountRoleArn")
	externalId, _ := cmd.PersistentFlags().GetString("externalId")
	authFlag, clientAuth, err := AuthenticateDataEnv(cloudElementId, cloudElementApiUrl, vaultUrl, vaultToken, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId)
	return authFlag, clientAuth, err
}

func SubCommandAuth(cmd *cobra.Command) (bool, *client.Auth, error) {
	cloudElementId, _ := cmd.Parent().PersistentFlags().GetString("cloudElementId")
	cloudElementApiUrl, _ := cmd.Parent().PersistentFlags().GetString("cloudElementApiUrl")
	vaultUrl, _ := cmd.Parent().PersistentFlags().GetString("vaultUrl")
	vaultToken, _ := cmd.Parent().PersistentFlags().GetString("vaultToken")
	accountNo, _ := cmd.Parent().PersistentFlags().GetString("accountId")
	region, _ := cmd.Parent().PersistentFlags().GetString("zone")
	acKey, _ := cmd.Parent().PersistentFlags().GetString("accessKey")
	secKey, _ := cmd.Parent().PersistentFlags().GetString("secretKey")
	crossAccountRoleArn, _ := cmd.Parent().PersistentFlags().GetString("crossAccountRoleArn")
	externalId, _ := cmd.Parent().PersistentFlags().GetString("externalId")
	authFlag, clientAuth, err := AuthenticateData(cloudElementId, cloudElementApiUrl, vaultUrl, vaultToken, accountNo, region, acKey, secKey, crossAccountRoleArn, externalId)
	return authFlag, clientAuth, err
}

func decrypt(key []byte, ciphertext string) (string, error) {
	if ciphertext != "" {
		block, err := aes.NewCipher(key)
		if err != nil {
			return "", err
		}
		decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			return "", err
		}
		if len(decodedCiphertext) < aes.BlockSize {
			return "", fmt.Errorf("ciphertext too short")
		}
		iv := decodedCiphertext[:aes.BlockSize]
		decodedCiphertext = decodedCiphertext[aes.BlockSize:]
		cipher.NewCFBDecrypter(block, iv).XORKeyStream(decodedCiphertext, decodedCiphertext)
		return string(decodedCiphertext), nil
	}
	return "", nil
}
