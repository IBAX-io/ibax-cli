package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/converter"
	"github.com/IBAX-io/go-ibax/packages/common/crypto"
	"github.com/IBAX-io/ibax-cli/conf"
	"github.com/IBAX-io/ibax-cli/models"
	"github.com/IBAX-io/ibax-cli/packages/consts"
	"github.com/IBAX-io/ibax-cli/packages/parameter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const fileMode = 0600

var (
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "Manage accounts",
		Long: `
Manage accounts, list all existing accounts, create a new account,  
address convert keyId, keyId convert address
`,
	}
)

var (
	accountNewCmd = &cobra.Command{
		Use:   "new",
		Short: "create a new account",
		Long: `
Creates a new account and prints the address and keyId and publicKey.
`,
		PreRun:     loadConfigPre,
		Run:        accountNew,
		SuggestFor: []string{"new"},
		Example:    `./ibax-cli account new`,
	}

	accountListCmd = &cobra.Command{
		Use:   "list",
		Short: "Print information of existing accounts",
		Long: `
Print a short information of all accounts
`,
		PreRun:     loadConfigPre,
		Run:        accountList,
		SuggestFor: []string{"list"},
		Example:    `./ibax-cli account list`,
	}

	addressToKeyIdCmd = &cobra.Command{
		Use:   "addressToKeyId [address]",
		Short: "Account Address convert to Key Id",
		Long: `
Request:
	Account			(string) Account Address: "xxxx-xxxx-xxxx-xxxx-xxxx"

Returns converted address
Result:
	KeyId			(string) Key Id
`,

		PreRun:     loadConfigPre,
		Run:        addressToKeyId,
		SuggestFor: []string{"addressToKeyId"},
		Example:    `./ibax-cli account addressToKeyId [address]`,
	}

	keyId             int64
	keyIdToAddressCmd = &cobra.Command{
		Use:   "keyIdToAddress",
		Short: "Key Id convert to Account Address",
		Long: `
Request:
	No parameters required

Returns converted address
Result:
	Address			(string) Account Address: "xxxx-xxxx-xxxx-xxxx-xxxx"
`,
		PreRun: func(cmd *cobra.Command, args []string) {
			keyId = 0
			loadConfigPre(cmd, args)
		},
		Run:        keyIdToAddress,
		SuggestFor: []string{"keyIdToAddress --keyId="},
		Example:    `./ibax-cli account keyIdToAddress --keyId=[KeyId]`,
	}

	publicKeyToAddressCmd = &cobra.Command{
		Use:   "publicKeyToAddress [publicKey]",
		Short: "publicKey convert to Account Address",
		Long: `
Request:
	publicKey			(string) Public Key

Returns converted address
Result:
	Address			(string) Account Address: "xxxx-xxxx-xxxx-xxxx-xxxx"
`,
		PreRun:     loadConfigPre,
		Run:        publicKeyToAddress,
		SuggestFor: []string{"publicKeyToAddress"},
		Example:    `./ibax-cli account publicKeyToAddress [publicKey]`,
	}

	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "current account information",
		Long: `
Request:
	No parameters required

Returns a json object for the current account information
Result:
	{
		"public_key": "str",		(string) Public Key
		"ecosystem_id": n,			(string) Login Ecosystem Id
		"key_id": n,				(number) Key Id
		"account": "str"			(string) Account Address: "xxxx-xxxx-xxxx-xxxx-xxxx"
	}
`,
		PreRun:     loadConfigPre,
		Run:        accountInfo,
		SuggestFor: []string{"info"},
		Example:    `./ibax-cli account info`,
	}
)

const newAccountWarning = `
Warning
You can share your public key or address or keyId with anyone. Others need it to interact with you.
You must NEVER share the secret key with anyone! The key controls access to your funds!
You must BACKUP your key file! Without the key, it's impossible to access account funds!
`

func init() {
	accountListCmd.Flags().StringVar(&conf.Config.DirPathConf.KeysDir, "keysDir", "", "Keys Directory")
	accountNewCmd.Flags().StringVar(&conf.Config.DirPathConf.KeysDir, "keysDir", "", "Keys Directory")

	keyIdToAddressCmd.Flags().Int64Var(&keyId, "keyId", 0, "Key Id")
	keyIdToAddressCmd.MarkFlagRequired("keyId")
}

func createFile(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0775)
		if err != nil {
			return errors.Wrapf(err, "creating dir %s", dir)
		}
	}

	return os.WriteFile(filename, data, fileMode)
}

func createKeyPair(privFilename, pubFilename string) (priv, pub []byte, err error) {
	priv, pub, err = crypto.GenKeyPair()
	if err != nil {
		log.WithError(err).Error("generate keys")
		return
	}

	err = createFile(privFilename, []byte(hex.EncodeToString(priv)))
	if err != nil {
		log.WithFields(log.Fields{"error": err, "path": privFilename}).Error("creating private key")
		return
	}

	err = createFile(pubFilename, []byte(crypto.PubToHex(pub)))
	if err != nil {
		log.WithFields(log.Fields{"error": err, "path": pubFilename}).Error("creating public key")
		return
	}
	return
}

func toISO8601(t time.Time) string {
	var tz string
	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}

func accountNew(cmd *cobra.Command, params []string) {
	err := cobra.NoArgs(cmd, params)
	if err != nil {
		log.Infof("no parameters required: %s", err.Error())
		return
	}
	nowTimeStr := toISO8601(time.Now())
	privateKeyName := filepath.Join(conf.Config.DirPathConf.KeysDir, fmt.Sprintf("%s-UTC-%s", nowTimeStr, consts.PrivateKeyFilename))
	publicKeyName := filepath.Join(conf.Config.DirPathConf.KeysDir, fmt.Sprintf("%s-UTC-%s", nowTimeStr, consts.PublicKeyFilename))
	_, publicKey, err := createKeyPair(
		privateKeyName,
		publicKeyName,
	)
	if err != nil {
		log.Infof("account new failed: %s", err.Error())
		return
	}
	keyId := crypto.Address(publicKey)
	address := converter.AddressToString(keyId)

	fmt.Printf("Path of the private key file: %s\n", privateKeyName)
	fmt.Printf("Path of the public key file: %s\n", publicKeyName)
	fmt.Printf("Public Key: %s\n", crypto.PubToHex(publicKey))
	fmt.Printf("KeyId: %d\n", keyId)
	fmt.Printf("address: %s\n", address)
	fmt.Printf("%s\n", newAccountWarning)
}

func accountList(cmd *cobra.Command, params []string) {
	err := cobra.NoArgs(cmd, params)
	if err != nil {
		log.Infof("no parameters required: %s", err.Error())
		return
	}
	dir, err := os.ReadDir(filepath.Join(conf.Config.DirPathConf.KeysDir))
	if err != nil {
		log.Infof("get keys dir failed:%s", err.Error())
		return
	}
	type accountInfo struct {
		Path struct {
			PrivateKey string
			PublicKey  string
		}
		PublicKey string
		keyId     int64
		Address   string
	}
	var accountIndex int = 1
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		fileName := entry.Name()
		if strings.HasSuffix(fileName, consts.PublicKeyFilename) {
			publicKeyName := filepath.Join(conf.Config.DirPathConf.KeysDir, fileName)
			utcTimePre := fileName[:len(fileName)-len(consts.PublicKeyFilename)]
			privateKeyName := filepath.Join(conf.Config.DirPathConf.KeysDir, utcTimePre+consts.PublicKeyFilename)
			publicKey, err := os.ReadFile(publicKeyName)
			if err != nil {
				fmt.Printf("readFile failed:%s\n", err.Error())
				continue
			}
			pub, err := hex.DecodeString(string(publicKey))
			if err != nil {
				fmt.Printf("publicKey[%s] decode failed:%s\n", publicKeyName, err.Error())
				continue
			}

			keyId := crypto.Address(pub)
			address := converter.AddressToString(keyId)
			var info accountInfo
			info.Path.PublicKey = publicKeyName
			info.Path.PrivateKey = privateKeyName
			info.PublicKey = crypto.PubToHex(pub)
			info.keyId = keyId
			info.Address = address
			fmt.Printf("account #%d:%+v\n", accountIndex, info)
			accountIndex += 1
		}
	}
	fmt.Println()
}

func publicKeyToAddress(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	publicKeyStr, err := args.Set(0, true).String()
	if err != nil {
		log.Infof("publicKey invalid:%s", err.Error())
		return
	}
	pub, err := hex.DecodeString(publicKeyStr)
	if err != nil {
		log.Infof("publicKey decode failed:%s", err.Error())
		return
	}
	keyId := crypto.Address(pub)
	address := converter.AddressToString(keyId)

	fmt.Printf("\nkeyId: %d,account address: %s\n", keyId, address)
}

func addressToKeyId(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	address, err := args.Set(0, true).String()
	if err != nil {
		log.Infof("address invalid:%s", err.Error())
		return
	}
	keyId := converter.StringToAddress(address)
	if keyId == 0 {
		log.Infof("address invalid:%s", address)
		return
	}

	fmt.Printf("\n%d\n", keyId)
}

func keyIdToAddress(cmd *cobra.Command, params []string) {
	err := cobra.NoArgs(cmd, params)
	if err != nil {
		log.Infof("no parameters required: %s", err.Error())
		return
	}
	address := converter.AddressToString(keyId)
	if address == "" {
		fmt.Printf("invalid Key Id:%d\n", keyId)
		return
	}
	if converter.StringToAddress(address) == 0 {
		log.Infof("KeyId invalid:%d", keyId)
		return
	}

	fmt.Printf("\n%s\n", address)
}

func accountInfo(cmd *cobra.Command, params []string) {
	err := cobra.NoArgs(cmd, params)
	if err != nil {
		log.Infof("no parameters required: %s", err.Error())
		return
	}
	type accountInfo struct {
		PublicKey   string `json:"public_key"`
		EcosystemId int64  `json:"ecosystem_id"`
		KeyId       int64  `json:"key_id"` // key id
		Account     string `json:"account"`
	}
	var info accountInfo
	cnf := models.Client.GetConfig()
	info.PublicKey = hex.EncodeToString(cnf.PublicKey)
	info.EcosystemId = cnf.Ecosystem
	info.KeyId = cnf.KeyId
	info.Account = cnf.Account

	str, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}
