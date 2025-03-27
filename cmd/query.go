package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/converter"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/ibax-cli/conf"
	"github.com/IBAX-io/ibax-cli/models"
	"github.com/IBAX-io/ibax-cli/packages/parameter"
	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
)

func getKeysInfoCmd(cmd *cobra.Command, args []string) {
	if args[0] == "" {
		log.Infof("Account Address Can't Not Be Empty")
		return
	}

	result, err := models.Client.GetKeyInfo(args[0])
	if err != nil {
		log.Infof("Get Key Info Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("Get Key Info Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getBalanceCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	account, err := args.Set(0, false).String()
	if err != nil {
		log.Infof("account invalid:%s", err.Error())
		return
	}

	ecosystemId, err := args.Set(1, false).NumberInt64()
	if err != nil {
		log.Infof("ecosystem id invalid:%s", err.Error())
		return
	}
	if account == "" {
		cnf := models.Client.GetConfig()
		if cnf.Account != "" {
			account = cnf.Account
			fmt.Printf("current account:%s\n", cnf.Account)
		} else {
			fmt.Printf("current account not exist,Please set in the configuration file:%s or specify an account\n", conf.Config.ConfigPath)
			return
		}
	}

	result, err := models.Client.Balance(account, ecosystemId)
	if err != nil {
		log.Infof("Get Balance Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("Get Balance Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getVersionCmd(cmd *cobra.Command, args []string) {
	err := cobra.NoArgs(cmd, args)
	if err != nil {
		log.Infof("no parameters required: %s", err.Error())
		return
	}

	result, err := models.Client.GetVersion()
	if err != nil {
		log.Infof("Get Version Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("Get Version Result Empty")
		return
	}
	fmt.Printf("%s\n", *result)
}

func getChainConfigCmd(cmd *cobra.Command, args []string) {
	if args[0] == "" {
		log.Infof("Option Can't Not Be Empty")
		return
	}

	result, err := models.Client.GetIBAXConfig(args[0])
	if err != nil {
		log.Infof("Get Chain Config Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("Get Chain Config Result Empty")
		return
	}
	fmt.Printf("%s\n", *result)
}

func getEcosystemCountCmd(cmd *cobra.Command, args []string) {
	err := cobra.NoArgs(cmd, args)
	if err != nil {
		log.Infof("no parameters required: %s", err.Error())
		return
	}

	count, err := models.Client.EcosystemCount()
	if err != nil {
		log.Infof("Get Version Failed: %s", err.Error())
		return
	}

	fmt.Printf("\n%d\n", count)
}

func getMaxBlockCmd(cmd *cobra.Command, args []string) {
	err := cobra.NoArgs(cmd, args)
	if err != nil {
		log.Infof("no parameters required: %s", err.Error())
		return
	}

	count, err := models.Client.GetMaxBlockID()
	if err != nil {
		log.Infof("Get Max Block Failed: %s", err.Error())
		return
	}

	fmt.Printf("\n%d\n", count)
}

func getTxCountCmd(cmd *cobra.Command, args []string) {
	err := cobra.NoArgs(cmd, args)
	if err != nil {
		log.Infof("no parameters required: %s", err.Error())
		return
	}

	count, err := models.Client.TransactionsCount()
	if err != nil {
		log.Infof("Get Transaction Count Failed: %s", err.Error())
		return
	}

	fmt.Printf("\n%d\n", count)
}

func getKeysCountCmd(cmd *cobra.Command, args []string) {
	err := cobra.NoArgs(cmd, args)
	if err != nil {
		log.Infof("no parameters required: %s", err.Error())
		return
	}

	count, err := models.Client.KeysCount()
	if err != nil {
		log.Infof("Get keys Count Failed: %s", err.Error())
		return
	}

	fmt.Printf("\n%d\n", count)
}

func getHonorNodesCountCmd(cmd *cobra.Command, args []string) {
	err := cobra.NoArgs(cmd, args)
	if err != nil {
		log.Infof("no parameters required: %s", err.Error())
		return
	}

	count, err := models.Client.HonorNodesCount()
	if err != nil {
		log.Infof("Get Honor Nodes Count Failed: %s", err.Error())
		return
	}

	fmt.Printf("\n%d\n", count)
}

func detailBlocksCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	blockId, err := args.Set(0, true).NumberInt64()
	if err != nil {
		log.Infof("BlockId invalid:%s", err.Error())
		return
	}

	count, err := args.Set(1, false).NumberInt64()
	if err != nil {
		log.Infof("Count invalid:%s", err.Error())
		return
	}

	result, err := models.Client.DetailedBlocks(blockId, count)
	if err != nil {
		log.Infof("Detailed Blocks Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("Detailed Blocks Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getBlockInfoCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	blockId, err := args.Set(0, true).NumberInt64()
	if err != nil {
		log.Infof("blockId invalid:%s", err.Error())
		return
	}

	result, err := models.Client.GetBlockInfo(blockId)
	if err != nil {
		log.Infof("get Block Info Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("Block Info Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getBlocksTxInfoCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	blockId, err := args.Set(0, true).NumberInt64()
	if err != nil {
		log.Infof("blockId invalid:%s", err.Error())
		return
	}
	count, err := args.Set(1, false).NumberInt64()
	if err != nil {
		log.Infof("count invalid:%s", err.Error())
		return
	}

	result, err := models.Client.BlocksTxInfo(blockId, count)
	if err != nil {
		log.Infof("get blocks tx info Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get blocks tx info Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getTableCountCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	offset, err := args.Set(0, false).NumberInt()
	if err != nil {
		log.Infof("offset invalid:%s", err.Error())
		return
	}
	limit, err := args.Set(1, false).NumberInt()
	if err != nil {
		log.Infof("limit invalid:%s", err.Error())
		return
	}

	result, err := models.Client.GetTableCount(offset, limit)
	if err != nil {
		log.Infof("get table count Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get table count Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getTableCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	tableName, err := args.Set(0, true).String()
	if err != nil {
		log.Infof("offset invalid:%s", err.Error())
		return
	}

	result, err := models.Client.GetTable(tableName)
	if err != nil {
		log.Infof("get table Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get table Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getSectionsCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	language, err := args.Set(0, false).String()
	if err != nil {
		log.Infof("language invalid:%s", err.Error())
		return
	}
	offset, err := args.Set(1, false).NumberInt()
	if err != nil {
		log.Infof("offset invalid:%s", err.Error())
		return
	}
	limit, err := args.Set(2, false).NumberInt()
	if err != nil {
		log.Infof("limit invalid:%s", err.Error())
		return
	}

	result, err := models.Client.GetSections(language, offset, limit)
	if err != nil {
		log.Infof("get sections Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get sections Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getPageRowCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	name, err := args.Set(0, true).String()
	if err != nil {
		log.Infof("name invalid:%s", err.Error())
		return
	}

	result, err := models.Client.GetPageRow(name)
	if err != nil {
		log.Infof("get page row Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get page row Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getMenuRowCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	name, err := args.Set(0, true).String()
	if err != nil {
		log.Infof("name invalid:%s", err.Error())
		return
	}

	result, err := models.Client.GetMenuRow(name)
	if err != nil {
		log.Infof("get menu row Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get menu row Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getSnippetRowCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	name, err := args.Set(0, true).String()
	if err != nil {
		log.Infof("name invalid:%s", err.Error())
		return
	}

	result, err := models.Client.GetSnippetRow(name)
	if err != nil {
		log.Infof("get snippet row Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get snippet row Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getAppContentCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	appId, err := args.Set(0, true).NumberInt64()
	if err != nil {
		log.Infof("appId invalid:%s", err.Error())
		return
	}

	result, err := models.Client.GetAppContent(appId)
	if err != nil {
		log.Infof("get app content Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get app content Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func appParamsCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	appId, err := args.Set(0, true).NumberInt64()
	if err != nil {
		log.Infof("appId invalid:%s", err.Error())
		return
	}
	names, err := args.Set(1, false).String()
	if err != nil {
		log.Infof("names invalid:%s", err.Error())
		return
	}
	ecosystemId, err := args.Set(2, false).NumberInt64()
	if err != nil {
		log.Infof("ecosystemId invalid:%s", err.Error())
		return
	}
	offset, err := args.Set(3, false).NumberInt()
	if err != nil {
		log.Infof("offset invalid:%s", err.Error())
		return
	}
	limit, err := args.Set(4, false).NumberInt()
	if err != nil {
		log.Infof("limit invalid:%s", err.Error())
		return
	}

	result, err := models.Client.AppParams(appId, names, ecosystemId, offset, limit)
	if err != nil {
		log.Infof("get app params Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get app params Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func ecosystemParamsCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	ecosystemId, err := args.Set(0, false).NumberInt64()
	if err != nil {
		log.Infof("ecosystemId invalid:%s", err.Error())
		return
	}
	names, err := args.Set(1, false).String()
	if err != nil {
		log.Infof("names invalid:%s", err.Error())
		return
	}
	offset, err := args.Set(2, false).NumberInt()
	if err != nil {
		log.Infof("offset invalid:%s", err.Error())
		return
	}
	limit, err := args.Set(3, false).NumberInt()
	if err != nil {
		log.Infof("limit invalid:%s", err.Error())
		return
	}

	result, err := models.Client.EcosystemParams(ecosystemId, names, offset, limit)
	if err != nil {
		log.Infof("get ecosystem params Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get ecosystem params Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func systemParamsCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	names, err := args.Set(0, false).String()
	if err != nil {
		log.Infof("names invalid:%s", err.Error())
		return
	}
	offset, err := args.Set(1, false).NumberInt()
	if err != nil {
		log.Infof("offset invalid:%s", err.Error())
		return
	}
	limit, err := args.Set(2, false).NumberInt()
	if err != nil {
		log.Infof("limit invalid:%s", err.Error())
		return
	}

	result, err := models.Client.SystemParams(names, offset, limit)
	if err != nil {
		log.Infof("get system params Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get system params Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getRowCmd(cmd *cobra.Command, params []string) {
	rowParams.lock.Lock()
	defer rowParams.lock.Unlock()
	args := parameter.New(params)
	tableName, err := args.Set(0, true).String()
	if err != nil {
		log.Infof("tableName invalid:%s", err.Error())
		return
	}
	columns, err := args.Set(1, false).String()
	if err != nil {
		log.Infof("columns invalid:%s", err.Error())
		return
	}
	WhereColumn, err := args.Set(2, false).String()
	if err != nil {
		log.Infof("WhereColumn invalid:%s", err.Error())
		return
	}

	result, err := models.Client.GetRow(tableName, rowParams.Id, columns, WhereColumn)
	if err != nil {
		log.Infof("get row Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get row Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getHistoryCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	tableName, err := args.Set(0, true).String()
	if err != nil {
		log.Infof("tableName invalid:%s", err.Error())
		return
	}

	id, err := args.Set(1, true).NumberUint64()
	if err != nil {
		log.Infof("tableName invalid:%s", err.Error())
		return
	}

	result, err := models.Client.GetHistory(tableName, id)
	if err != nil {
		log.Infof("get history Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get history Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getListCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	tableName, err := args.Set(0, true).String()
	if err != nil {
		log.Infof("tableName invalid:%s", err.Error())
		return
	}

	getListParams.Name = tableName
	if getListParams.whereStr != "" {
		getListParams.Where = getListParams.whereStr
	}
	if getListParams.orderStr != "" {
		err = json.Unmarshal([]byte(getListParams.orderStr), &getListParams.Order)
		if err != nil {
			log.Infof("order invalid:%s", err.Error())
			return
		}
	}
	defer func() {
		getListParams.Name = ""
		getListParams.Where = nil
	}()

	result, err := models.Client.GetList(getListParams.GetList)
	if err != nil {
		log.Infof("get list Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("get list Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func blockTxCountCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	var bh request.BlockIdOrHash
	var err error
	bh.Id, err = args.Set(0, true).NumberInt64()
	if err != nil {
		bh.Hash, err = args.Set(0, true).String()
		if err != nil {
			log.Infof("block id or block hash invalid:%s", err.Error())
			return
		}
	}
	count, err := models.Client.BlockTxCount(bh)
	if err != nil {
		log.Infof("block tx count Failed: %s", err.Error())
		return
	}
	fmt.Printf("\n%d\n", count)
}

func detailedBlockCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	var bh request.BlockIdOrHash
	var err error
	bh.Id, err = args.Set(0, true).NumberInt64()
	if err != nil {
		bh.Hash, err = args.Set(0, true).String()
		if err != nil {
			log.Infof("block id or block hash invalid:%s", err.Error())
			return
		}
	}

	result, err := models.Client.DetailedBlock(bh)
	if err != nil {
		log.Infof("detailed Block Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("detailed Block Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func ecosystemInfoCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	ecosystemId, err := args.Set(0, true).NumberInt64()
	if err != nil {
		log.Infof("ecosystem Id invalid:%s", err.Error())
		return
	}

	result, err := models.Client.EcosystemInfo(ecosystemId)
	if err != nil {
		log.Infof("Ecosystem Info Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("Ecosystem Info Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getMemberInfoCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	account, err := args.Set(0, true).String()
	if err != nil {
		log.Infof("account invalid:%s", err.Error())
		return
	}
	keyId := converter.StringToAddress(account)
	if keyId == 0 {
		log.Infof("account[%s] invalid:%s", account, "format not supported")
		return
	}

	ecosystemId, err := args.Set(1, true).NumberInt64()
	if err != nil {
		log.Infof("ecosystem Id invalid:%s", err.Error())
		return
	}

	result, err := models.Client.GetMemberInfo(account, ecosystemId)
	if err != nil {
		log.Infof("Get Member Info Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("Get Member Info Result Empty")
		return
	}
	str, err := json.MarshalIndent(*result, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func binaryVerifyCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	binaryId, err := args.Set(0, true).NumberInt64()
	if err != nil {
		log.Infof("binary id invalid:%s", err.Error())
		return
	}
	binaryHash, err := args.Set(1, true).String()
	if err != nil {
		log.Infof("binary hash invalid:%s", err.Error())
		return
	}

	fileInfo, err := models.Client.BinaryVerify(binaryId, binaryHash, binaryFileName)
	if err != nil {
		fmt.Printf("binary verify failed: %s", err.Error())
		return
	}
	str, err := json.MarshalIndent(fileInfo, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func exportCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	appId, err := args.Set(0, true).NumberInt64()
	if err != nil {
		log.Infof("account invalid:%s", err.Error())
		return
	}
	cnf := models.Client.GetConfig()

	account := cnf.Account
	ecosystem := cnf.Ecosystem
	contractParams := make(request.MapParams)
	contractParams["ApplicationId"] = appId
	result, err := models.Client.AutoCallContract("@1ExportNewApp", &contractParams, "")
	if err != nil {
		log.Infof("call @1ExportNewApp Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("call @1ExportNewApp Result Empty")
		return
	}
	if result.BlockId == 0 || result.Hash == "" || result.Penalty == 1 || result.Err != "" {
		str, err := json.MarshalIndent(*result, "", "    ")
		if err != nil {
			fmt.Printf("call @1ExportNewApp Result marshall Failed:%s\n", err.Error())
			return
		}
		fmt.Printf("export process @1ExportNewApp failed: \n%+v\n", string(str))
		return
	}

	result, err = models.Client.AutoCallContract("@1Export", nil, "")
	if err != nil {
		log.Infof("call @1Export Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("call @1Export Result Empty")
		return
	}
	if result.BlockId == 0 || result.Hash == "" || result.Penalty == 1 || result.Err != "" {
		str, err := json.MarshalIndent(*result, "", "    ")
		if err != nil {
			fmt.Printf("call @1Export Result marshall Failed:%s\n", err.Error())
			return
		}
		fmt.Printf("export process @1Export failed: \n%+v\n", string(str))
		return
	}
	var getListParams request.GetList
	getListParams.Name = "@1binaries"
	getListParams.Limit = 1
	getListParams.Columns = "id,hash"
	getListParams.Where = fmt.Sprintf(`{"name": "export", "account": "%s", "ecosystem": %d, "app_id": %d}`, account, ecosystem, appId)
	listResult, err := models.Client.GetList(getListParams)
	if err != nil {
		fmt.Printf("export process GetList failed: %s", err.Error())
		return
	}
	if listResult == nil {
		log.Info("export process GetList Result Empty")
		return
	}

	var binaryId int64
	var binaryHash string
	if listResult.Count == 1 {
		value, ok := listResult.List[0]["id"]
		if ok {
			binaryId, _ = strconv.ParseInt(value, 10, 64)
		}

		value, ok = listResult.List[0]["hash"]
		if ok {
			binaryHash = value
		}
	}
	if binaryHash == "" || binaryId == 0 {
		str, err := json.MarshalIndent(*listResult, "", "    ")
		if err != nil {
			fmt.Printf("GetList Result marshall Failed:%s\n", err.Error())
			return
		}
		fmt.Printf("export process GetList failed binaryHash or binaryId empty: \n%+v\n", string(str))
		return
	}

	fileInfo, err := models.Client.BinaryVerify(binaryId, binaryHash, exportFileName)
	if err != nil {
		fmt.Printf("export process GetBinary failed: %s", err.Error())
		return
	}
	str, err := json.MarshalIndent(fileInfo, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func importCmd(cmd *cobra.Command, params []string) {
	err := cobra.NoArgs(cmd, params)
	if err != nil {
		log.Infof("no parameters required: %s", err.Error())
		return
	}
	if importFileName == "" {
		log.Info("import file name can't not be empty")
		return
	}

	data, err := os.ReadFile(importFileName)
	if err != nil {
		log.Infof("filename: [%s] readfile err: %s", importFileName, err.Error())
		return
	}
	mimeType, _, err := getMimeType(importFileName)
	if err != nil {
		log.Infof("filename: [%s] getMimeType err: %s", importFileName, err.Error())
		return
	}

	importInfo := make(map[string]any)
	importInfo["Name"] = filepath.Base(importFileName)
	importInfo["MimeType"] = mimeType
	importInfo["Body"] = data
	contractParams := make(request.MapParams)
	contractParams["Data"] = importInfo
	result, err := models.Client.AutoCallContract("@1ImportUpload", &contractParams, "")
	if err != nil {
		log.Infof("call @1ImportUpload Failed: %s", err.Error())
		return
	}
	if result == nil {
		log.Info("call @1ImportUpload Result Empty")
		return
	}
	if result.BlockId == 0 || result.Hash == "" || result.Penalty == 1 || result.Err != "" {
		str, err := json.MarshalIndent(*result, "", "    ")
		if err != nil {
			fmt.Printf("call @1ImportUpload Result marshall Failed:%s\n", err.Error())
			return
		}
		fmt.Printf("import process @1ImportUpload failed: \n%+v\n", string(str))
		return
	}

	cnf := models.Client.GetConfig()
	var getListParams request.GetList
	getListParams.Name = "@1buffer_data"
	getListParams.Limit = 1
	getListParams.Columns = "value->'data'"
	getListParams.Where = fmt.Sprintf(`{"key": "import", "account": "%s", "ecosystem": %d}`, cnf.Account, cnf.Ecosystem)
	listResult, err := models.Client.GetList(getListParams)
	if err != nil {
		fmt.Printf("import process GetList failed: %s", err.Error())
		return
	}
	if listResult == nil {
		log.Info("import process GetList Result Empty")
		return
	}
	var bufferData string
	var ret = make([]any, 0)
	if listResult.Count == 1 {
		value, ok := listResult.List[0]["value.data"]
		if ok {
			var d []map[string]any
			err := json.Unmarshal([]byte(value), &d)
			if err != nil {
				return
			}
			for _, i2 := range d {
				a, ok := i2["Data"]
				if ok {
					var b []any
					err := json.Unmarshal([]byte(a.(string)), &b)
					if err != nil {
						return
					}
					ret = append(ret, b...)
				}
			}
			bytes, _ := json.Marshal(ret)
			bufferData = string(bytes)
		}
	}
	if bufferData == "" {
		str, err := json.MarshalIndent(*listResult, "", "    ")
		if err != nil {
			fmt.Printf("GetList Result marshall Failed:%s\n", err.Error())
			return
		}
		fmt.Printf("import process GetList failed bufferData empty: \n%+v\n", string(str))
		return
	}

	contractParams["Data"] = bufferData
	importResult, err := models.Client.AutoCallContract("@1Import", &contractParams, "")
	if err != nil {
		log.Infof("call @1Import Failed: %s", err.Error())
		return
	}
	if importResult == nil {
		log.Info("call @1Import Result Empty")
		return
	}

	str, err := json.MarshalIndent(importResult, "", "    ")
	if err != nil {
		fmt.Printf("Result marshall Failed:%s\n", err.Error())
		return
	}
	fmt.Printf("\n%+v\n", string(str))
}

func getMimeType(fileName string) (string, string, error) {
	mType, err := mimetype.DetectFile(fileName)
	if err != nil {
		fmt.Printf("DetectFile err:%s\n", err.Error())
		return "", "", err
	}
	return mType.String(), mType.Extension(), nil
}

func base64DecodeCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	decodeData, err := args.Set(0, false).String()
	if err != nil {
		log.Infof("data invalid:%s", err.Error())
		return
	}
	if decodeData == "" && decodeFileName == "" {
		log.Info("decode data and file name must have one")
		return
	}

	data, err := base64.StdEncoding.DecodeString(decodeData)
	if err != nil {
		log.Infof("data Decode failed:%s", err.Error())
		return
	}

	if decodeFileName != "" {
		err := os.WriteFile(decodeFileName, data, 0644)
		if err != nil {
			log.Info(err)
			return
		}
	} else {
		fmt.Printf("\nDecode:%s\n", string(data))
	}
}

func base64EncodeCmd(cmd *cobra.Command, params []string) {
	args := parameter.New(params)
	encodeData, err := args.Set(0, false).String()
	if err != nil {
		log.Infof("data invalid:%s", err.Error())
		return
	}
	if encodeData == "" && encodeFileName == "" {
		log.Info("encode data and file name must have one")
		return
	}

	var data []byte
	if encodeFileName != "" {
		data, err = os.ReadFile(encodeFileName)
		if err != nil {
			log.Infof("ReadFile failed: %s", err.Error())
			return
		}
	} else {
		data = []byte(encodeData)
	}
	fmt.Printf("\nEncode:%s\n", base64.StdEncoding.EncodeToString(data))
}
