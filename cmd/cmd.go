package cmd

import (
	"context"
	"github.com/IBAX-io/go-ibax-sdk/packages/request"
	"github.com/spf13/cobra"
	"sync"
)

var cmdList []*cobra.Command

// query
var (
	getKeyInfo = &cobra.Command{
		Use:   "getKeyInfo [Account]",
		Short: "Returns an ecosystem list containing rolues registered with the specified address",
		Long: `
Request:
	Account			(string) Account Address: "xxxx-xxxx-xxxx-xxxx-xxxx"

Returns a json object key information.
Result:
	{
		"account": "str",				(string) Account Address: "xxxx-xxxx-xxxx-xxxx-xxxx"
		"ecosystems": [
			{
				"ecosystem": "str",		(string) Ecosystem Id
				"name": "str",			(string) Ecosystem Name
				"digits": n,			(string) Digits
				"roles": [				(array) Role List
					{
						"id": "str",	(string) Role Id
						"name": "str"	(string) Role Name
					}
				]
			}
		]
	}
`,
		SuggestFor: []string{"getKeyInfo"},
		Example:    "./ibax-cli getKeyInfo 0666-7782-xxxx-xxxx-3160",
		Args:       cobra.ExactArgs(1),
		PreRun:     loadConfigPre,
		Run:        getKeysInfoCmd,
	}

	getBalance = &cobra.Command{
		Use:   "getBalance [Account] [EcosystemId]",
		Short: "Get Account Balance",
		Long: `
Request:
	Account   				(string,optional) Account Address: "xxxx-xxxx-xxxx-xxxx-xxxx",default: current account (if exist)
	EcosystemId   			(number,optional) Ecosystem Id,default: ecosystem id 1

Returns a json object Balance information.
Result:
	{
		"amount": "str",		(string) The contract account balance of the smallest unit (QIBAX).
		"digits": n,			(number) precision.
		"total": "str",			(string) he minimum unit account total balance (amount + utxo)
		"utxo": "str",			(string) The smallest unit UTXO account balance.
		"token_symbol": "str"	(string) token symbol
	}
`,
		SuggestFor: []string{"getBalance"},
		Example:    "./ibax-cli getBalance 0666-7782-xxxx-xxxx-3160",
		Args:       cobra.RangeArgs(0, 2),
		PreRun:     loadConfigPre,
		Run:        getBalanceCmd,
	}

	getVersion = &cobra.Command{
		Use:   "getVersion",
		Short: "Get Version Information.",
		Long: `
Request: 
	No parameters required

Returns String.
Result:
	1.3.0 branch.main commit.93c39396 time.2023-03-27-07:19:36(UTC)			(string) Version Information
`,
		SuggestFor: []string{"getVersion"},
		Example:    "./ibax-cli getVersion",
		PreRun:     loadConfigPre,
		Run:        getVersionCmd,
	}

	getConfig = &cobra.Command{
		Use:   "getConfig [Option]",
		Short: "Get Chain Config.",
		Long: `
Request:
	Option 			(string) centrifugo: Centrifugo server

Returns String.
Result:
	"wss://node21.ibax.io:8330"		(string) Centrifugo Server URL
`,
		SuggestFor: []string{"getConfig"},
		Example:    "./ibax-cli getConfig centrifugo",
		PreRun:     loadConfigPre,
		Args:       cobra.ExactArgs(1),
		Run:        getChainConfigCmd,
	}

	ecosystemCount = &cobra.Command{
		Use:   "ecosystemCount",
		Short: "Get Ecosystem Count",
		Long: `
Request:
	No parameters required

Returns Count.
Result:
	3				(number) - Ecosystem Count
`,
		SuggestFor: []string{"ecosystemCount"},
		Example:    "./ibax-cli ecosystemCount",
		PreRun:     loadConfigPre,
		Run:        getEcosystemCountCmd,
	}

	maxBlock = &cobra.Command{
		Use:   "maxBlock",
		Short: "Get Max Block.",
		Long: `
Request:
	No parameters required

Returns Number For Max Bolck.
Result: 
	1 			(number)  Max Block Number
`,
		SuggestFor: []string{"maxBlock"},
		Example:    "./ibax-cli maxBlock",
		PreRun:     loadConfigPre,
		Run:        getMaxBlockCmd,
	}

	txCount = &cobra.Command{
		Use:   "transactionCount",
		Short: "Get Transaction Count",
		Long: `
Request:
	No parameters required

Returns Count.
Result:
	223982					(number) - Transaction Count
`,
		SuggestFor: []string{"transactionCount"},
		Example:    "./ibax-cli transactionCount",
		PreRun:     loadConfigPre,
		Run:        getTxCountCmd,
	}

	keysCount = &cobra.Command{
		Use:   "keysCount",
		Short: "Get Account Key Count.",
		Long: `
Request:
	No parameters required

Returns Count.
Result:
	27590				(number) - Account Key Count
`,
		SuggestFor: []string{"keysCount"},
		Example:    "./ibax-cli keysCount",
		PreRun:     loadConfigPre,
		Run:        getKeysCountCmd,
	}

	honorNodesCount = &cobra.Command{
		Use:   "honorNodesCount",
		Short: "Get Honor Nodes Count.",
		Long: `
Request:
	No parameters required

Returns Count.
Result:
	1 -					(number) - Honor Nodes Count
`,
		SuggestFor: []string{"honorNodesCount"},
		Example:    "./ibax-cli honorNodesCount",
		PreRun:     loadConfigPre,
		Run:        getHonorNodesCountCmd,
	}

	detailedBlocks = &cobra.Command{
		Use:   "detailedBlocks [BlockId] [Count]",
		Short: "Get block detailed information.",
		Long: `
Request:
	BlockId				(number) The starting block height to query
	Count				(number,optional) The number of blocks, the default is 25, the maximum request is 100

Returns a json object block detail information.
Result:
	{
		"n": {								(string) Block Id
			"header": {						(object) Block Header
				"block_id": n,				(number) Block Height
				"time": n,					(number) Block generation timestamp (unit: s)
				"key_id": n,				(number) The address of the account that signed the block
				"node_position": n,			(number) The position of the node that produced the block in the honor node list
				"version": n				(string) Version
			},
			"hash": "str",					(string) Block Hash
			"node_position": n,				(number) The position of the node that produced the block in the honor node list
			"key_id": n,					(number) The address of the account that signed the block
			"time": 1681434711,				(number) Block generation timestamp
			"tx_count": 1,					(number) The block transaction count
			"size": "str",					(string) Block size
			"rollbacks_hash": "str",		(strung) Block rollback hash
			"merkle_root": "str",			(string) merkle root
			"bin_data": "str",				(string) block header, all transactions within the block, the hash of the previous block, and the private key Serialization of the node that generated the block
			"transactions": [				(array) transaction list
				{
					"hash": "str",			(string) transaction hash	
					"contract_name": "str",	(string) contract name
					"params": {},			(object) contract params or utxo params
					"key_id": n,			(number) The address of the account that signed the transaction
					"time": n,				(number) Transaction generation timestamp (unit: ms)
					"type": n,				(number) The transaction type
					"size": "str"			(string) The transaction size
				}
			]
		}
	}
`,
		SuggestFor: []string{"detailedBlocks"},
		Example:    "./ibax-cli detailedBlocks [BlockId] [Count]",
		PreRun:     loadConfigPre,
		Args:       cobra.RangeArgs(1, 2),
		Run:        detailBlocksCmd,
	}

	getBlockInfo = &cobra.Command{
		Use:   "getBlockInfo [BlockId]",
		Short: "Get block information",
		Long: `
Request:
	BlockId				(number) The block height


Returns a json object block detail information.
Result:
	{
		"hash": "str",				(string) Block Hash
		"key_id": n,				(number) The address of the account that signed the block
		"time": n,					(number) Block generation timestamp (unit: s)
		"tx_count": n,				(number) The block transaction count
		"rollbacks_hash": "str",	(string) Block rollback hash
		"node_position": n			(number) The position of the node that produced the block in the honor node list
		"consensus_mode": 1			(number) Consensus mode, parameters (1: Creator management mode 2: DAO governance mode)
	}
`,
		SuggestFor: []string{"getBlockInfo"},
		Example:    "./ibax-cli getBlockInfo [BlockId]",
		PreRun:     loadConfigPre,
		Args:       cobra.ExactArgs(1),
		Run:        getBlockInfoCmd,
	}

	blocksTxInfo = &cobra.Command{
		Use:   "blocksTxInfo [BlockId] [Count]",
		Short: "Get block information.",
		Long: `
Request:
	BlockId			(number) The starting block height to query
	Count			(number,optional) The number of blocks, the default is 25, the maximum request is 100

Returns a json object block transaction information.
Result:
	{
		"n": [								(string) Block Id
			{
				"hash": "str",					(string) transaction hash		
				"contract_name": "str",			(string) contract name
				"params": {},					(object) contract params or utxo params
				"key_id": n						(number) The address of the account that signed the transaction
			}
		]
	}
`,
		SuggestFor: []string{"blocksTxInfo"},
		Example:    "./ibax-cli blocksTxInfo [BlockId] [Count]",
		PreRun:     loadConfigPre,
		Args:       cobra.RangeArgs(1, 2),
		Run:        getBlocksTxInfoCmd,
	}

	getTableCount = &cobra.Command{
		Use:   "getTableCount [Offset] [Limit]",
		Short: "Get table count.",
		Long: `
Request:
Offset 				(number,optional) offset, default is 0
Limit  				(number,optional) The number of entries, the default is 25, and the maximum is 100.

Returns a json object Table for the Current Ecosystem Count information.
Result:
{
    "1170": [								(string) Block Id
        {
            "hash": "str",					(string) transaction hash		
            "contract_name": "str",			(string) contract name
            "params": {},					(object) contract params or utxo params
            "key_id": n						(number) The address of the account that signed the transaction
        }
    ]
}
`,
		SuggestFor: []string{"getTableCount"},
		Example:    "./ibax-cli getTableCount [Offset] [Limit]",
		PreRun:     loginPre,
		Args:       cobra.RangeArgs(0, 2),
		Run:        getTableCountCmd,
	}

	getTable = &cobra.Command{
		Use:   "getTable [TableName]",
		Short: "Get table information.",
		Long: `
Request:
	TableName				(string) table name.

Returns a json object Table information for the Current Ecosystem.
Result:
	{
		"name": "str",		(string) data table name
		"insert": "str",		(string) permissions for insert
		"new_column": "str",		(string) permissions for new column
		"update": "str",		(string) permissions for update
		"conditions": "str",		(string) Conditions for changing permissions
		"app_id": "1",		(string) application id
		"columns": [		(array) Data table field related information array
			{
				"name": "str",		(string) Field Name.
				"type": "str",		(string) Field data type.
				"perm": "str"		(string) Permission to change the value of this field.
			}
		]
	}
`,
		SuggestFor: []string{"getTable"},
		Example:    "./ibax-cli getTable [TableName]",
		PreRun:     loginPre,
		Args:       cobra.ExactArgs(1),
		Run:        getTableCmd,
	}

	getSections = &cobra.Command{
		Use:   "getSections [Language] [Offset] [Limit]",
		Short: "Go back to the tab for the current ecosystem A list of table entries.",
		Long: `
Request:
	Language				(string,optional) This field specifies a multilingual resource code or localization, for example: en, zh. 
							If the specified multilingual resource is not found, for example: en-US, then search in the multilingual resource group,default: en
	Offset				(number,optional) offset, default is 0
	Limit 				(number,optional) The number of entries, the default is 25, and the maximum is 100.

if role_access field contains a list of roles, and does not include the current role, no records are returned. title
The data in the field will be replaced by the Accept-Language language resource of the request header

Returns a json object Go back to the tab for the current ecosystem.
Result:
	{
		"count": n,				(number) Total number of tab entries
		"list": [				(array) Each element in the array contains information about all the columns in the sections table.
			{
				"ecosystem": "str",
				"id": "str",
				"page": "str",
				"roles_access": "[]",
				"status": "str",
				"title": "str",
				"urlname": "str"
			}
		]
	}
`,
		SuggestFor: []string{"getSections"},
		Example:    "./ibax-cli getSections [TableName]",
		PreRun:     loginPre,
		Args:       cobra.RangeArgs(0, 3),
		Run:        getSectionsCmd,
	}

	getPageRow = &cobra.Command{
		Use:   "getPageRow [Name]",
		Short: "Get the entries of the current ecosystem pages data table field",
		Long: `
Request:
	Name				(string) specifies the name of the entry in the table

Returns a json object Get the entries of the for the current ecosystem
Result:
	{
		"id": n,				(number) The entry ID
		"name": "str",			(string) The entry name
		"value": "str",			(string) The content
		"menu": "str",			(string) The menu
		"nodesCount": n,		(number) The number of nodes the page needs to verify
		"app_id": n,			(number) Application Id
		"conditions": "str"		(string) The permission to change the parameter
	}
`,
		SuggestFor: []string{"getPageRow"},
		Example:    `./ibax-cli getPageRow [Name]`,
		PreRun:     loginPre,
		Args:       cobra.ExactArgs(1),
		Run:        getPageRowCmd,
	}

	getMenuRow = &cobra.Command{
		Use:   "getMenuRow [Name]",
		Short: "Get the entries of the current ecosystem menu data table fields.",
		Long: `
Request:
	Name				(string) Specifies the name of the entry in the table.

Returns a json object Get the entries of the for the current ecosystem
Result:
	{
		"id": n,				(number) The entry ID
		"name": "str",			(string) The entry name
		"title": "str",			(string) The title
		"value": "str",			(string) The content
		"conditions": "str"		(string) The permission to change the parameter
	}
`,
		SuggestFor: []string{"getMenuRow"},
		Example:    "./ibax-cli getMenuRow [Name]",
		PreRun:     loginPre,
		Args:       cobra.ExactArgs(1),
		Run:        getMenuRowCmd,
	}

	getSnippetRow = &cobra.Command{
		Use:   "getSnippetRow [Name]",
		Short: "Get the entries of the current ecosystem snippet data table fields.",
		Long: `
Request:
	Name				(string) Specifies the name of the entry in the table.

Returns a json object Get the entries of the for the current ecosystem
Result:
	{
		"id": n,				(number) The entry ID
		"name": "str",			(string) The entry name
		"value": "str",			(string) The content
		"conditions": "str"		(string) The permission to change the parameter
	}
`,
		SuggestFor: []string{"getSnippetRow"},
		Example:    "./ibax-cli getSnippetRow [Name]",
		PreRun:     loginPre,
		Args:       cobra.ExactArgs(1),
		Run:        getSnippetRowCmd,
	}

	getAppContent = &cobra.Command{
		Use:   "getAppContent [Id]",
		Short: "Get Obtain application-related information (including page, snippet, menu)",
		Long: `
Request:
	Id			(number) Application Id.

Returns a json object Get application information for the current ecosystem
Result:
	{
		"snippets": [					(array) Array of code snippet information
			{
				"id": n,					(number) id
				"name": "str"				(string) the snippet name
			}
		],
		"pages": [						(array) Array of page information
			{
				"id": n,					(number) id
				"name": "str"				(string) the name of the page
			}
		],
		"contracts": [					(array) contract information array
			{
				"id": n,					(number) id
				"name": "str"				(string) contract name
			}
		]
	}
`,
		SuggestFor: []string{"getAppContent"},
		Example:    "./ibax-cli getAppContent [Id]",
		PreRun:     loginPre,
		Args:       cobra.ExactArgs(1),
		Run:        getAppContentCmd,
	}

	appParams = &cobra.Command{
		Use:   "appParams [AppId] [Names] [EcosystemId] [Offset] [Limit]",
		Short: "Returns a list of application parameters in the current or specified ecosystem",
		Long: `
Request:
	AppId 						(number) Application Id
	EcosystemId 				(number,optional) Ecosystem Id.If zero, will return the parameters of the current ecosystem
	Names 						(string,optional) Filter application parameter name. List of names separated by commas, eg: "name1,name2"
									When there are filter parameters, the offset and limit parameters are invalid
	Offset 						(number,optional) offset, default is 0
	Limit 						(number,optional) limit, the default is 10, and the maximum is 100

Returns a json object application parameters
Result:
	{
		"app_id": "str",				(string) Application ID
		"list": [
			{
				"id": "str",			(string) Parameter ID, unique
				"name": "str",			(string) The parameter name
				"value": "str",			(string) The parameter value
				"conditions": "str"		(string) The permission to change the parameter
			}
		]
	}
`,
		SuggestFor: []string{"appParams"},
		Example:    "./ibax-cli appParams [AppId] [Names] [EcosystemId] [Offset] [Limit]",
		PreRun:     loginPre,
		Args:       cobra.RangeArgs(1, 5),
		Run:        appParamsCmd,
	}

	ecosystemParams = &cobra.Command{
		Use:   "ecosystemParams [EcosystemId] [Names] [Offset] [Limit]",
		Short: "Get a list of ecosystem parameters",
		Long: `
Request:
	EcosystemId			(number,optional) Ecosystem Id.If zero, will return the parameters of the current ecosystem
	Names 				(string,optional) Filter application parameter name. List of names separated by commas, eg: "name1,name2".
							When there are filter parameters, the offset and limit parameters are invalid
	Offset 				(number,optional) offset, default is 0
	Limit 				(number,optional) limit, the default is 10, and the maximum is 100

Returns a json object ecosystem parameters
Result:
	{
		"list": [
			{
				"id": "str",		(string) Parameter id, unique.
				"name": "str",		(string) The parameter name.
				"value": "str",		(string) The parameter value
				"conditions": "str"	(string) The permission to change the parameter
			}
		]
	}
`,
		SuggestFor: []string{"ecosystemParams"},
		Example:    "./ibax-cli ecosystemParams [EcosystemId] [Names] [Offset] [Limit]",
		PreRun:     loginPre,
		Args:       cobra.RangeArgs(0, 4),
		Run:        ecosystemParamsCmd,
	}

	systemParams = &cobra.Command{
		Use:   "systemParams [Names] [Offset] [Limit]",
		Short: "Get a list of ecosystem parameters",
		Long: `
Request:
	Names 			(string,optional) Filter application parameter name. List of names separated by commas, eg: "name1,name2". 
						When there are filter parameters, the offset and limit parameters are invalid
	Offset			(number,optional) offset, default is 0
	Limit 			(number,optional) limit, the default is 10, and the maximum is 100

Returns a json object ecosystem parameters
Result:
	{
		"list": [
			{
				"id": "str",		(string) Parameter id, unique.
				"name": "str",		(string) The parameter name.
				"value": "str",		(string) The parameter value
				"conditions": "str"		(string) The permission to change the parameter
			}
		]
	}
`,
		SuggestFor: []string{"systemParams"},
		Example:    "./ibax-cli systemParams [Names] [Offset] [Limit]",
		PreRun:     loginPre,
		Args:       cobra.RangeArgs(0, 3),
		Run:        systemParamsCmd,
	}

	rowParams idObject

	getRow = &cobra.Command{
		Use:   "getRow [TableName] [Columns] [WhereColumn]",
		Short: "Returns the entries for the specified data table in the current ecosystem. You can specify which columns to return.",
		Long: `
Request:
	TableName			(string) Table name
	Columns				(string,optional) A comma-separated list of requested columns, if not specified all columns will be returned. 
							In all cases the id column is returned.
	WhereColumn			(string,optional) Find the column name (only for Number type columns)

Returns a json object entries for the specified data table in the current ecosystem
Result:
	{
		"value": {
			"id": "-5476304279945383650",				(string)  The entry ID
			...											(strung) A sequence of requested columns.
		}
	}
`,
		SuggestFor: []string{"getRow"},
		Example:    "./ibax-cli getRow [TableName] [Columns] [WhereColumn]",
		PreRun:     loginPre,
		Args:       cobra.RangeArgs(1, 3),
		Run:        getRowCmd,
	}

	getHistory = &cobra.Command{
		Use:   "getHistory [TableName] [Id]",
		Short: "Returns the change record for the entry in the specified data table in the current ecosystem",
		Long: `
Request:
	TableName 				(string) Table name
	Id						(number) The entry ID

Returns a json object entries for the specified data table in the current ecosystem
Result:
	{
		"list": [					(array) Each element in the array contains a change record for the requested item
			{
				"conditions": "str",
				"ecosystem": "str",
				"value": "str"
			}
		]
	}
`,
		SuggestFor: []string{"getHistory"},
		Example:    "./ibax-cli getHistory [TableName] [Id]",
		PreRun:     loginPre,
		Args:       cobra.ExactArgs(2),
		Run:        getHistoryCmd,
	}

	getListParams struct {
		request.GetList
		whereStr string
	}
	getList = &cobra.Command{
		Use: `getList [TableName]`,
		Short: `Returns the entries of the specified data table.
		You can specify which columns to return.
		You can set the offset and the number of entries.
		You can set query conditions.
		Do hexadecimal encoding processing for the type BYTEA (byte array, hash, bytecode array) in the data table
`,
		Long: `
Request:
	TableName		(string) Table name

Returns a json object entries for the specified data table in the current ecosystem
Result:
	{
		"count": n,						(number) Total number of entries.
		"list": [						(array) Each element in the array contains the following parameters:
			{
				"id": "str",			(string) The entry ID.
				"...": "str"			(string) additional columns of the data table.
			}
		]
	}
`,
		SuggestFor: []string{"getList"},
		Example:    `./ibax-cli getList [TableName] -w='{"id": [tableId]}' -c="amount,ecosystem" -l=3 -t=1 -r="ecosystem desc`,
		PreRun: func(cmd *cobra.Command, args []string) {
			loginPre(cmd, args)
		},
		Args: cobra.ExactArgs(1),
		Run:  getListCmd,
	}

	blockTxCount = &cobra.Command{
		Use:   "blockTxCount [BlockOrHash]",
		Short: "Get block transaction",
		Long: `
Request:
	BlockOrHash			(string || number) Block Id Or Block Hash

Returns block transaction count
Result:
	count					(number) block transaction count
`,
		SuggestFor: []string{"blockTxCount"},
		Example:    "./ibax-cli blockTxCount [BlockOrHash]",
		PreRun:     loadConfigPre,
		Args:       cobra.ExactArgs(1),
		Run:        blockTxCountCmd,
	}

	detailedBlock = &cobra.Command{
		Use:   "detailedBlock [BlockOrHash]",
		Short: "Returns a detailed list of additional information about the transactions in the block.",
		Long: `
Request:
	BlockOrHash			(string || number) Block Id Or Block Hash

Returns a json object for detailed block
Result:
	{
		"header": {							(object) Block Header
			"block_id": n,					(number) Block Height
			"time": n,						(number) Block generation timestamp (unit: s)
			"key_id": n,					(number) The address of the account that signed the block
			"node_position": n,				(number) The position of the node that produced the block in the honor node list
			"version": n					(string) Version
		},
		"hash": "str",						(string) Block Hash
		"node_position": n,					(number) The position of the node that produced the block in the honor node list
		"key_id": n,						(number) The address of the account that signed the block
		"time": 1681434711,					(number) Block generation timestamp
		"tx_count": 1,						(number) The block transaction count
		"size": "str",						(string) Block size
		"rollbacks_hash": "str",			(strung) Block rollback hash
		"merkle_root": "str",				(string) merkle root
		"bin_data": "str",					(string) block header, all transactions within the block, the hash of the previous block, and the private key Serialization of the node that generated the block
		"transactions": [					(array) transaction list
			{
				"hash": "str",				(string) transaction hash	
				"contract_name": "str",		(string) contract name
				"params": {},				(object) contract params or utxo params
				"key_id": n,				(number) The address of the account that signed the transaction
				"time": n,					(number) Transaction generation timestamp (unit: ms)
				"type": n,					(number) The transaction type
				"size": "str"				(string) The transaction size
			}
		]
	}
`,
		SuggestFor: []string{"detailedBlock"},
		Example:    "./ibax-cli detailedBlock [BlockOrHash]",
		PreRun:     loadConfigPre,
		Args:       cobra.ExactArgs(1),
		Run:        detailedBlockCmd,
	}

	ecosystemInfo = &cobra.Command{
		Use:   "ecosystemInfo [EcosystemId]",
		Short: "Get Ecosystem information",
		Long: `
Request:
	EcosystemId 				(number) Ecosystem Id

Returns a json object for Ecosystem information
Result:
	{
		"id": n,						(number) Ecosystem id
		"name": "str",					(string) Ecosystem name
		"digits": n,					(number) precision
		"token_symbol": "str",			(string) token symbol
		"token_name": "str",			(string) token name
		"total_amount": "str", 			(string) Supply Token (the first supply amount, or "0" if not supply amount)
		"is_withdraw": bool,			(boolean) can be withdrawn true: can be withdrawn false: can not be withdrawn	
		"withdraw": "str",				(string) Amount to withdraw ("0" if not withdrawn, or not withdrawn)
		"is_emission": bool,			(boolean) can be Increase supply true: can be Increase supply false: can not be Increase supply
		"emission": "str",				(string) Increase supply
		"introduction": "str",			(string) Ecosystem introduction
		"logo": n,						(number) Ecosystem logo id
		"creator": "str" 				(string) The ecosystem creator
	}
`,
		SuggestFor: []string{"ecosystemInfo"},
		Example:    "./ibax-cli ecosystemInfo [EcosystemId]",
		PreRun:     loadConfigPre,
		Args:       cobra.ExactArgs(1),
		Run:        ecosystemInfoCmd,
	}

	getMemberInfo = &cobra.Command{
		Use:   "getMemberInfo [Account] [EcosystemId]",
		Short: "Get Ecosystem information",
		Long: `
Request:
	Account 			(string) Account Address (xxxx-xxxx-xxxx-xxxx-xxxx)
	EcosystemId 		(number) Ecosystem Id

Returns a json object for Ecosystem information
Result:
	{
		"id": n,						(number) Member id
		"member_name": "str",			(string) Member Name
		"image_id": n,					(number) Image id
		"member_info": "str"			(string) Member introduction
	}
`,
		SuggestFor: []string{"getMemberInfo"},
		Example:    "./ibax-cli getMemberInfo [Account] [EcosystemId]",
		PreRun:     loadConfigPre,
		Args:       cobra.ExactArgs(2),
		Run:        getMemberInfoCmd,
	}

	binaryFileName string
	binaryVerify   = &cobra.Command{
		Use:   "binaryVerify",
		Short: "Verify export binary",
		Long: `
Request:
	BinaryId 		(number) binary file Id
	BinaryHash 		(string) binary file hash

Returns a json object for binary Verify
Result:
	{
		"name": "str",			(string) save binary file name
		"type": "str",			(string) binary file type
		"value": "str"			(string) if save binary file name is null, save result to value
	}
`,
		SuggestFor: []string{"binaryVerify"},
		Example:    "./ibax-cli binaryVerify [BinaryId] [BinaryHash]",
		PreRun:     loginPre,
		Args:       cobra.ExactArgs(2),
		Run:        binaryVerifyCmd,
	}

	exportFileName string
	export         = &cobra.Command{
		Use:   "export",
		Short: "export Application",
		Long: `
Request:
	AppId		 (number) Applications Id

Returns a json object for export
Result:
	{
		"name": "str",			(string) save binary file name
		"type": "str",			(string) binary file type
		"value": "str"			(string) if save binary file name is null, save result to value
	}
`,
		SuggestFor: []string{"export"},
		Example:    "./ibax-cli export",
		PreRun:     loginPre,
		Args:       cobra.ExactArgs(1),
		Run:        exportCmd,
	}

	importFileName string
	importUpload   = &cobra.Command{
		Use:   "import",
		Short: "import Application",
		Long: `
Request:
	AppId		 (number) Applications Id

Returns a json object for export
Result:
	{
		"block_id": n,			(number) The block id generated by the transaction
		"hash": "str",			(string) The block hash generated by the transaction
		"penalty": n,			(number) If transaction execution fails, (0: no penalty 1: penalty)
		"err": ""				(string, optional) If the execution of the transaction fails, an error text message is returned.
	}
`,
		SuggestFor: []string{"import"},
		Example:    "./ibax-cli import",
		PreRun:     loginPre,
		Run:        importCmd,
	}

	encodeFileName string
	base64Encode   = &cobra.Command{
		Use:   "base64Encode",
		Short: "base64Encode",
		Long: `
Request:
	data		 (string) need encode data

Returns Encode data
	Encode 		(string) encode data
`,
		SuggestFor: []string{"base64Encode"},
		Example:    "./ibax-cli base64Encode",
		Args:       cobra.MaximumNArgs(1),
		Run:        base64EncodeCmd,
	}

	decodeFileName string
	base64Decode   = &cobra.Command{
		Use:   "base64Decode",
		Short: "base64Decode",
		Long: `
Request:
	data		 (string) need decode data

Returns decode data
	Decode		(string) decode data
`,
		SuggestFor: []string{"base64Decode"},
		Example:    "./ibax-cli base64Decode [data]",
		Args:       cobra.MaximumNArgs(1),
		Run:        base64DecodeCmd,
	}
)

func hasErrorContext(cmd *cobra.Command) bool {
	ctx := cmd.Context()
	if val := ctx.Value("error"); val != nil {
		return true
	}
	return false
}
func cleanErrorContext(cmd *cobra.Command) {
	cmd.SetContext(context.Background())
}

type idObject struct {
	lock sync.Mutex
	Id   int64
}

func initFlags() {
	getRow.Flags().Int64VarP(&rowParams.Id, "id", "i", 0, "The entry ID")
	getRow.MarkFlagRequired("id")

	callContract.Flags().StringVarP(&contractParamsFile, "file", "f", "", "Contract Params File Name,json object,priority")

	getList.Flags().IntVarP(&getListParams.Limit, "limit", "l", 0, "the number of Omitempty entries,default 25")
	getList.Flags().IntVarP(&getListParams.Offset, "offset", "t", 0, "offset,default 0")
	getList.Flags().StringVarP(&getListParams.Order, "order", "r", "", "sorting method, default id ASC")
	getList.Flags().StringVarP(&getListParams.Columns, "columns", "c", "", "A comma-separated list of Omitempty request columns, if not specified, all columns will be returned.")
	getList.Flags().StringVarP(&getListParams.whereStr, "where", "w", "", `
Query conditions

Example: If you want to query id>2 and name = john

You can use: where:{"id":{"$gt":2},"name":{"$eq":"john"}}

For details, please refer to DBFind where syntax
`)

	export.Flags().StringVarP(&exportFileName, "file", "f", "", "Export Application file name")
	binaryVerify.Flags().StringVarP(&binaryFileName, "file", "f", "", "Save binary file name")
	importUpload.Flags().StringVarP(&importFileName, "file", "f", "", "Import Application file name")
	base64Encode.Flags().StringVarP(&encodeFileName, "file", "f", "", "need encode file name,priority")
	base64Decode.Flags().StringVarP(&decodeFileName, "file", "f", "", "decode file name,priority")

}

func initCmdList() {
	initFlags()
	//auth
	cmdList = append(cmdList, authStatus)
	cmdList = append(cmdList, refresh)

	//contract
	cmdList = append(cmdList, getContracts)
	cmdList = append(cmdList, getContractInfo)
	cmdList = append(cmdList, callContract)

	//utxo
	cmdList = append(cmdList, callUtxo)

	//query
	cmdList = append(cmdList, getKeyInfo)
	cmdList = append(cmdList, getBalance)
	cmdList = append(cmdList, getVersion)
	cmdList = append(cmdList, getConfig)
	cmdList = append(cmdList, ecosystemCount)
	cmdList = append(cmdList, maxBlock)
	cmdList = append(cmdList, txCount)
	cmdList = append(cmdList, keysCount)
	cmdList = append(cmdList, honorNodesCount)
	cmdList = append(cmdList, detailedBlocks)
	cmdList = append(cmdList, getBlockInfo)
	cmdList = append(cmdList, blocksTxInfo)
	cmdList = append(cmdList, getTableCount)
	cmdList = append(cmdList, getTable)
	cmdList = append(cmdList, getSections)
	cmdList = append(cmdList, getPageRow)
	cmdList = append(cmdList, getMenuRow)
	cmdList = append(cmdList, getSnippetRow)
	cmdList = append(cmdList, getAppContent)
	cmdList = append(cmdList, appParams)
	cmdList = append(cmdList, ecosystemParams)
	cmdList = append(cmdList, systemParams)
	cmdList = append(cmdList, getRow)
	cmdList = append(cmdList, getHistory)
	cmdList = append(cmdList, getList)
	cmdList = append(cmdList, blockTxCount)
	cmdList = append(cmdList, detailedBlock)
	cmdList = append(cmdList, ecosystemInfo)
	cmdList = append(cmdList, getMemberInfo)
	cmdList = append(cmdList, binaryVerify)
	cmdList = append(cmdList, export)
	cmdList = append(cmdList, importUpload)
	cmdList = append(cmdList, base64Encode)
	cmdList = append(cmdList, base64Decode)
}
