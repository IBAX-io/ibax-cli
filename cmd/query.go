package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-ibax-sdk/packages/pkg/converter"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/IBAX-io/ibax-cli/conf"
	"github.com/IBAX-io/ibax-cli/models"
	"github.com/IBAX-io/ibax-cli/packages/parameter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	defer func() {
		//clean buffer
		getListParams.Name = ""
		getListParams.Offset = 0
		getListParams.Limit = 0
		getListParams.Where = nil
		getListParams.Columns = ""
		getListParams.Order = ""
	}()

	getListParams.Name = tableName
	if getListParams.whereStr != "" {
		getListParams.Where = getListParams.whereStr
	}

	result, err := models.Client.GetList(getListParams.GetList)
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
