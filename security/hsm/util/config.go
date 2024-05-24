package util

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/cloudflare/cfssl/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func YamlStringSettings(cfg *viper.Viper) string {
	c := cfg.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("unable to marshal config to YAML: %v", err)
	}
	return string(bs)
}

func LoadConfigUsingHSM(lib_path string, pin string, token, label string,
	sessionKeyEncFile, configDataEncFile string) (*viper.Viper, error) {
	sessionKeyEnc, err := LoadFile(sessionKeyEncFile)
	if err != nil {
		fmt.Printf("LoadFile %s failed: %s \n", sessionKeyEncFile, err.Error())
		return nil, err
	}
	sDec, err := base64.StdEncoding.DecodeString(string(sessionKeyEnc))
	if err != nil {
		fmt.Printf("DecodeString %s failed: %s \n", sessionKeyEncFile, err.Error())
		return nil, err
	}
	sessionKey, err := HSM_AES_CBC_Decrypt(lib_path, pin, token, label, sDec)
	if err != nil {
		fmt.Printf("HSM_Decrypt failed: %s \n", err.Error())
		return nil, err
	}
	configDataEnc, err := LoadFile(configDataEncFile)
	if err != nil {
		fmt.Printf("LoadFile %s failed: %s \n", configDataEncFile, err.Error())
		return nil, err
	}
	configData, err := AES_decrypt(string(configDataEnc), string(sessionKey))
	if err != nil {
		fmt.Printf("AES_decrypt  failed: %s \n", err.Error())
		return nil, err
	}
	cfg := viper.New()
	cfg.SetConfigType("json")
	cfg.ReadConfig(bytes.NewBuffer(configData))
	return cfg, nil
}
func LoadConfigUsingFileKey(privKey, sessionKeyEncFile, configDataEncFile string) (*viper.Viper, error) {
	privateKey, err := LoadPrivateKeyFile_PKCS8(privKey)
	if err != nil {
		fmt.Printf("LoadPrivateKeyFile failed: %s \n", err.Error())
		return nil, err
	}
	// fmt.Printf("PrivateKey: %+v ", privateKey)

	sessionKeyEnc, err := LoadFile(sessionKeyEncFile)
	if err != nil {
		fmt.Printf("LoadFile %s failed: %s \n", sessionKeyEncFile, err.Error())
		return nil, err
	}
	sessionKey, err := DecryptWithPrivateKey(string(sessionKeyEnc), privateKey)
	if err != nil {
		fmt.Printf("DecryptWithPrivateKey  failed: %s \n", err.Error())
		return nil, err
	}
	fmt.Printf("sessionKey: %s \n", string(sessionKey))

	configDataEnc, err := LoadFile(configDataEncFile)
	if err != nil {
		fmt.Printf("LoadFile %s failed: %s \n", configDataEncFile, err.Error())
		return nil, err
	}
	configData, err := AES_decrypt(string(configDataEnc), sessionKey)
	if err != nil {
		fmt.Printf("AES_decrypt  failed: %s \n", err.Error())
		return nil, err
	}
	cfg := viper.New()
	cfg.SetConfigType("json")
	cfg.ReadConfig(bytes.NewBuffer(configData))
	return cfg, nil
}
func LoadFile(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func LoadPrivateKeyFile_PKCS8(file string) (*rsa.PrivateKey, error) {
	privBytes, err := ioutil.ReadFile(file) // This is fine with Encryption
	if err != nil {
		fmt.Printf("LoadKeyFile: open file failed %s\n", err)
		return nil, err
	}
	return BytesToPrivateKey_PKCS8(privBytes)
}
