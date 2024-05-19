package usecase

import (
	"fmt"

	"github.com/blcvn/lib-golang-test/security/hsm/util"
)

func LoadConfig_HSM() {
	// lib_path := "/usr/lib/softhsm/libsofthsm2.so"
	lib_path := "/opt/homebrew/lib/softhsm/libsofthsm2.so"
	slotID := uint(2058063310)
	pin := "8764329"

	sessionKeyEncFile := "test/hsm/config/session.enc"
	configDataEncFile := "test/hsm/config/config.enc"

	cfg, err := util.LoadConfigUsingHSM(lib_path, slotID, pin, sessionKeyEncFile, configDataEncFile)

	if err != nil {
		fmt.Printf("LoadConfigUsingFileKey failed: %s \n", err.Error())
		return
	}
	fmt.Println(cfg.AllKeys())
	configData := util.YamlStringSettings(cfg)
	fmt.Printf("configData: %+v \n", string(configData))

}

func LoadConfig_FileKey() {
	privKey := "test/hsm/key/api.pem"
	sessionKeyEncFile := "test/hsm/config/session.enc"
	configDataEncFile := "test/hsm/config/config.enc"

	cfg, err := util.LoadConfigUsingFileKey(privKey, sessionKeyEncFile, configDataEncFile)
	if err != nil {
		fmt.Printf("LoadConfigUsingFileKey failed: %s \n", err.Error())
		return
	}
	fmt.Println(cfg.AllKeys())
	configData := util.YamlStringSettings(cfg)
	fmt.Printf("configData: %+v \n", string(configData))
}
