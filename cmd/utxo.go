package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/ibax-cli/models"
	"github.com/IBAX-io/ibax-cli/packages/parameter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var callUtxo = &cobra.Command{
	Use:   "callUtxo [type] [params] [expedite]",
	Short: "Call UTXO",
	Long: `
Request:
	type   					(string) Transfer || ContractToUTXO (contract account to utxo account) || UTXOToContract (utxo account to contract account)
	params  				(json object,optional) uxto params
		recipient  			(string,optional) recipient address: "xxxx-xxxx-xxxx-xxxx-xxxx"
		comment  			(string,optional) transaction comment
		amount 				(string) amount
	expedite 				(string,optional) expedite unit: QIBAX

Returns a json object transaction status information.
Result:
	{
		"block_id": n,			(number) The block id generated by the transaction
		"hash": "str",			(string) The block hash generated by the transaction
		"penalty": n,			(number) If transaction execution fails, (0: no penalty 1: penalty)
		"err": ""				(string, optional) If the execution of the transaction fails, an error text message is returned.
	}
`,
	SuggestFor: []string{"callUtxo " + TypeUTXOToContract, "callUtxo " + TypeTransfer, "callUtxo " + TypeContractToUTXO},
	Example: `
./ibax-cli callUtxo Transfer '{"recipient": "0666-7782-xxxx-xxxx-3160", "amount": "1", "comment": ""}' '1'
./ibax-cli callUtxo ContractToUTXO '{"amount": "1"}' '1'
./ibax-cli callUtxo UTXOToContract '{"amount": "1"}' '1'
`,
	Args:   cobra.RangeArgs(2, 3),
	PreRun: loginPre,
	Run:    callUtxoCmd,
}

const (
	TypeTransfer       = "TypeTransfer"
	TypeContractToUTXO = "TypeContractToUTXO"
	TypeUTXOToContract = "TypeUTXOToContract"
)

func callUtxoCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	utxoTypeStr, err := args.Set(0, true).String()
	if err != nil {
		log.Infof("Type invalid:%s", err.Error())
		return
	}
	paramsStr, err := args.Set(1, true).String()
	if err != nil {
		log.Infof("Params invalid:%s", err.Error())
		return
	}
	expedite, err := args.Set(2, false).String()
	if err != nil {
		log.Infof("Params invalid:%s", err.Error())
		return
	}
	var utxoParams request.MapParams
	err = json.Unmarshal([]byte(paramsStr), &utxoParams)
	if err != nil {
		log.Infof("Params JSON Parsing Failed: %s", err.Error())
		return
	}
	var utxoType request.UtxoType
	switch utxoTypeStr {
	case TypeTransfer:
		utxoType = request.TypeTransfer
	case TypeContractToUTXO:
		utxoType = request.TypeContractToUTXO
	case TypeUTXOToContract:
		utxoType = request.TypeUTXOToContract
	}

	result, err := models.Client.AutoCallUtxo(utxoType, &utxoParams, expedite)
	if err != nil {
		log.Infof("Call UTXO Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("Call UTXO Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}