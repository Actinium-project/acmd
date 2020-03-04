// Copyright (c) 2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package acmjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/Actinium-project/acmd/acmjson"
	"github.com/Actinium-project/acmd/wire"
)

// TestChainSvrCmds tests all of the chain server commands marshal and unmarshal
// into valid results include handling of optional fields being omitted in the
// marshalled command, while optional fields with defaults have the default
// assigned on unmarshalled commands.
func TestChainSvrCmds(t *testing.T) {
	t.Parallel()

	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "addnode",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("addnode", "127.0.0.1", acmjson.ANRemove)
			},
			staticCmd: func() interface{} {
				return acmjson.NewAddNodeCmd("127.0.0.1", acmjson.ANRemove)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"addnode","params":["127.0.0.1","remove"],"id":1}`,
			unmarshalled: &acmjson.AddNodeCmd{Addr: "127.0.0.1", SubCmd: acmjson.ANRemove},
		},
		{
			name: "createrawtransaction",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`)
			},
			staticCmd: func() interface{} {
				txInputs := []acmjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return acmjson.NewCreateRawTransactionCmd(txInputs, amounts, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1}],{"456":0.0123}],"id":1}`,
			unmarshalled: &acmjson.CreateRawTransactionCmd{
				Inputs:  []acmjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts: map[string]float64{"456": .0123},
			},
		},
		{
			name: "createrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("createrawtransaction", `[{"txid":"123","vout":1}]`,
					`{"456":0.0123}`, int64(12312333333))
			},
			staticCmd: func() interface{} {
				txInputs := []acmjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				amounts := map[string]float64{"456": .0123}
				return acmjson.NewCreateRawTransactionCmd(txInputs, amounts, acmjson.Int64(12312333333))
			},
			marshalled: `{"jsonrpc":"1.0","method":"createrawtransaction","params":[[{"txid":"123","vout":1}],{"456":0.0123},12312333333],"id":1}`,
			unmarshalled: &acmjson.CreateRawTransactionCmd{
				Inputs:   []acmjson.TransactionInput{{Txid: "123", Vout: 1}},
				Amounts:  map[string]float64{"456": .0123},
				LockTime: acmjson.Int64(12312333333),
			},
		},

		{
			name: "decoderawtransaction",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("decoderawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return acmjson.NewDecodeRawTransactionCmd("123")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decoderawtransaction","params":["123"],"id":1}`,
			unmarshalled: &acmjson.DecodeRawTransactionCmd{HexTx: "123"},
		},
		{
			name: "decodescript",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("decodescript", "00")
			},
			staticCmd: func() interface{} {
				return acmjson.NewDecodeScriptCmd("00")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"decodescript","params":["00"],"id":1}`,
			unmarshalled: &acmjson.DecodeScriptCmd{HexScript: "00"},
		},
		{
			name: "getaddednodeinfo",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getaddednodeinfo", true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetAddedNodeInfoCmd(true, nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true],"id":1}`,
			unmarshalled: &acmjson.GetAddedNodeInfoCmd{DNS: true, Node: nil},
		},
		{
			name: "getaddednodeinfo optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getaddednodeinfo", true, "127.0.0.1")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetAddedNodeInfoCmd(true, acmjson.String("127.0.0.1"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddednodeinfo","params":[true,"127.0.0.1"],"id":1}`,
			unmarshalled: &acmjson.GetAddedNodeInfoCmd{
				DNS:  true,
				Node: acmjson.String("127.0.0.1"),
			},
		},
		{
			name: "getbestblockhash",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getbestblockhash")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBestBlockHashCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getbestblockhash","params":[],"id":1}`,
			unmarshalled: &acmjson.GetBestBlockHashCmd{},
		},
		{
			name: "getblock",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getblock", "123")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBlockCmd("123", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123"],"id":1}`,
			unmarshalled: &acmjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   acmjson.Bool(true),
				VerboseTx: acmjson.Bool(false),
			},
		},
		{
			name: "getblock required optional1",
			newCmd: func() (interface{}, error) {
				// Intentionally use a source param that is
				// more pointers than the destination to
				// exercise that path.
				verbosePtr := acmjson.Bool(true)
				return acmjson.NewCmd("getblock", "123", &verbosePtr)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBlockCmd("123", acmjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true],"id":1}`,
			unmarshalled: &acmjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   acmjson.Bool(true),
				VerboseTx: acmjson.Bool(false),
			},
		},
		{
			name: "getblock required optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getblock", "123", true, true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBlockCmd("123", acmjson.Bool(true), acmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblock","params":["123",true,true],"id":1}`,
			unmarshalled: &acmjson.GetBlockCmd{
				Hash:      "123",
				Verbose:   acmjson.Bool(true),
				VerboseTx: acmjson.Bool(true),
			},
		},
		{
			name: "getblockchaininfo",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getblockchaininfo")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBlockChainInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockchaininfo","params":[],"id":1}`,
			unmarshalled: &acmjson.GetBlockChainInfoCmd{},
		},
		{
			name: "getblockcount",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getblockcount")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBlockCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockcount","params":[],"id":1}`,
			unmarshalled: &acmjson.GetBlockCountCmd{},
		},
		{
			name: "getblockhash",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getblockhash", 123)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBlockHashCmd(123)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblockhash","params":[123],"id":1}`,
			unmarshalled: &acmjson.GetBlockHashCmd{Index: 123},
		},
		{
			name: "getblockheader",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getblockheader", "123")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBlockHeaderCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblockheader","params":["123"],"id":1}`,
			unmarshalled: &acmjson.GetBlockHeaderCmd{
				Hash:    "123",
				Verbose: acmjson.Bool(true),
			},
		},
		{
			name: "getblocktemplate",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getblocktemplate")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBlockTemplateCmd(nil)
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getblocktemplate","params":[],"id":1}`,
			unmarshalled: &acmjson.GetBlockTemplateCmd{Request: nil},
		},
		{
			name: "getblocktemplate optional - template request",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"]}`)
			},
			staticCmd: func() interface{} {
				template := acmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				}
				return acmjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"]}],"id":1}`,
			unmarshalled: &acmjson.GetBlockTemplateCmd{
				Request: &acmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := acmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   500,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return acmjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":500,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &acmjson.GetBlockTemplateCmd{
				Request: &acmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   int64(500),
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getblocktemplate optional - template request with tweaks 2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getblocktemplate", `{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}`)
			},
			staticCmd: func() interface{} {
				template := acmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    100000000,
					MaxVersion:   2,
				}
				return acmjson.NewGetBlockTemplateCmd(&template)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getblocktemplate","params":[{"mode":"template","capabilities":["longpoll","coinbasetxn"],"sigoplimit":true,"sizelimit":100000000,"maxversion":2}],"id":1}`,
			unmarshalled: &acmjson.GetBlockTemplateCmd{
				Request: &acmjson.TemplateRequest{
					Mode:         "template",
					Capabilities: []string{"longpoll", "coinbasetxn"},
					SigOpLimit:   true,
					SizeLimit:    int64(100000000),
					MaxVersion:   2,
				},
			},
		},
		{
			name: "getcfilter",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getcfilter", "123",
					wire.GCSFilterRegular)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetCFilterCmd("123",
					wire.GCSFilterRegular)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilter","params":["123",0],"id":1}`,
			unmarshalled: &acmjson.GetCFilterCmd{
				Hash:       "123",
				FilterType: wire.GCSFilterRegular,
			},
		},
		{
			name: "getcfilterheader",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getcfilterheader", "123",
					wire.GCSFilterRegular)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetCFilterHeaderCmd("123",
					wire.GCSFilterRegular)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilterheader","params":["123",0],"id":1}`,
			unmarshalled: &acmjson.GetCFilterHeaderCmd{
				Hash:       "123",
				FilterType: wire.GCSFilterRegular,
			},
		},
		{
			name: "getcfilter",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getcfilter", "123", wire.GCSFilterRegular)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetCFilterCmd("123", wire.GCSFilterRegular)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilter","params":["123",0],"id":1}`,
			unmarshalled: &acmjson.GetCFilterCmd{
				Hash:    "123",
				FilterType: wire.GCSFilterRegular,
			},
		},
		{
			name: "getcfilterheader",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getcfilterheader", "123", wire.GCSFilterRegular)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetCFilterHeaderCmd("123", wire.GCSFilterRegular)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getcfilterheader","params":["123",0],"id":1}`,
			unmarshalled: &acmjson.GetCFilterHeaderCmd{
				Hash:    "123",
				FilterType: wire.GCSFilterRegular,
			},
		},
		{
			name: "getchaintips",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getchaintips")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetChainTipsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getchaintips","params":[],"id":1}`,
			unmarshalled: &acmjson.GetChainTipsCmd{},
		},
		{
			name: "getconnectioncount",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getconnectioncount")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetConnectionCountCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getconnectioncount","params":[],"id":1}`,
			unmarshalled: &acmjson.GetConnectionCountCmd{},
		},
		{
			name: "getdifficulty",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getdifficulty")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetDifficultyCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getdifficulty","params":[],"id":1}`,
			unmarshalled: &acmjson.GetDifficultyCmd{},
		},
		{
			name: "getgenerate",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getgenerate")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetGenerateCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getgenerate","params":[],"id":1}`,
			unmarshalled: &acmjson.GetGenerateCmd{},
		},
		{
			name: "gethashespersec",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("gethashespersec")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetHashesPerSecCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gethashespersec","params":[],"id":1}`,
			unmarshalled: &acmjson.GetHashesPerSecCmd{},
		},
		{
			name: "getinfo",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getinfo")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getinfo","params":[],"id":1}`,
			unmarshalled: &acmjson.GetInfoCmd{},
		},
		{
			name: "getmempoolentry",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getmempoolentry", "txhash")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetMempoolEntryCmd("txhash")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getmempoolentry","params":["txhash"],"id":1}`,
			unmarshalled: &acmjson.GetMempoolEntryCmd{
				TxID: "txhash",
			},
		},
		{
			name: "getmempoolinfo",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getmempoolinfo")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetMempoolInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmempoolinfo","params":[],"id":1}`,
			unmarshalled: &acmjson.GetMempoolInfoCmd{},
		},
		{
			name: "getmininginfo",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getmininginfo")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetMiningInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getmininginfo","params":[],"id":1}`,
			unmarshalled: &acmjson.GetMiningInfoCmd{},
		},
		{
			name: "getnetworkinfo",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getnetworkinfo")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetNetworkInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnetworkinfo","params":[],"id":1}`,
			unmarshalled: &acmjson.GetNetworkInfoCmd{},
		},
		{
			name: "getnettotals",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getnettotals")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetNetTotalsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getnettotals","params":[],"id":1}`,
			unmarshalled: &acmjson.GetNetTotalsCmd{},
		},
		{
			name: "getnetworkhashps",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getnetworkhashps")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetNetworkHashPSCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[],"id":1}`,
			unmarshalled: &acmjson.GetNetworkHashPSCmd{
				Blocks: acmjson.Int(120),
				Height: acmjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getnetworkhashps", 200)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetNetworkHashPSCmd(acmjson.Int(200), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200],"id":1}`,
			unmarshalled: &acmjson.GetNetworkHashPSCmd{
				Blocks: acmjson.Int(200),
				Height: acmjson.Int(-1),
			},
		},
		{
			name: "getnetworkhashps optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getnetworkhashps", 200, 123)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetNetworkHashPSCmd(acmjson.Int(200), acmjson.Int(123))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnetworkhashps","params":[200,123],"id":1}`,
			unmarshalled: &acmjson.GetNetworkHashPSCmd{
				Blocks: acmjson.Int(200),
				Height: acmjson.Int(123),
			},
		},
		{
			name: "getpeerinfo",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getpeerinfo")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetPeerInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getpeerinfo","params":[],"id":1}`,
			unmarshalled: &acmjson.GetPeerInfoCmd{},
		},
		{
			name: "getrawmempool",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getrawmempool")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetRawMempoolCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[],"id":1}`,
			unmarshalled: &acmjson.GetRawMempoolCmd{
				Verbose: acmjson.Bool(false),
			},
		},
		{
			name: "getrawmempool optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getrawmempool", false)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetRawMempoolCmd(acmjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawmempool","params":[false],"id":1}`,
			unmarshalled: &acmjson.GetRawMempoolCmd{
				Verbose: acmjson.Bool(false),
			},
		},
		{
			name: "getrawtransaction",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getrawtransaction", "123")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetRawTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123"],"id":1}`,
			unmarshalled: &acmjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: acmjson.Int(0),
			},
		},
		{
			name: "getrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getrawtransaction", "123", 1)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetRawTransactionCmd("123", acmjson.Int(1))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawtransaction","params":["123",1],"id":1}`,
			unmarshalled: &acmjson.GetRawTransactionCmd{
				Txid:    "123",
				Verbose: acmjson.Int(1),
			},
		},
		{
			name: "gettxout",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("gettxout", "123", 1)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetTxOutCmd("123", 1, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1],"id":1}`,
			unmarshalled: &acmjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: acmjson.Bool(true),
			},
		},
		{
			name: "gettxout optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("gettxout", "123", 1, true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetTxOutCmd("123", 1, acmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxout","params":["123",1,true],"id":1}`,
			unmarshalled: &acmjson.GetTxOutCmd{
				Txid:           "123",
				Vout:           1,
				IncludeMempool: acmjson.Bool(true),
			},
		},
		{
			name: "gettxoutproof",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("gettxoutproof", []string{"123", "456"})
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetTxOutProofCmd([]string{"123", "456"}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxoutproof","params":[["123","456"]],"id":1}`,
			unmarshalled: &acmjson.GetTxOutProofCmd{
				TxIDs: []string{"123", "456"},
			},
		},
		{
			name: "gettxoutproof optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("gettxoutproof", []string{"123", "456"},
					acmjson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"))
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetTxOutProofCmd([]string{"123", "456"},
					acmjson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettxoutproof","params":[["123","456"],` +
				`"000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"],"id":1}`,
			unmarshalled: &acmjson.GetTxOutProofCmd{
				TxIDs:     []string{"123", "456"},
				BlockHash: acmjson.String("000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf"),
			},
		},
		{
			name: "gettxoutsetinfo",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("gettxoutsetinfo")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetTxOutSetInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"gettxoutsetinfo","params":[],"id":1}`,
			unmarshalled: &acmjson.GetTxOutSetInfoCmd{},
		},
		{
			name: "getwork",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getwork")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetWorkCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":[],"id":1}`,
			unmarshalled: &acmjson.GetWorkCmd{
				Data: nil,
			},
		},
		{
			name: "getwork optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getwork", "00112233")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetWorkCmd(acmjson.String("00112233"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getwork","params":["00112233"],"id":1}`,
			unmarshalled: &acmjson.GetWorkCmd{
				Data: acmjson.String("00112233"),
			},
		},
		{
			name: "help",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("help")
			},
			staticCmd: func() interface{} {
				return acmjson.NewHelpCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":[],"id":1}`,
			unmarshalled: &acmjson.HelpCmd{
				Command: nil,
			},
		},
		{
			name: "help optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("help", "getblock")
			},
			staticCmd: func() interface{} {
				return acmjson.NewHelpCmd(acmjson.String("getblock"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"help","params":["getblock"],"id":1}`,
			unmarshalled: &acmjson.HelpCmd{
				Command: acmjson.String("getblock"),
			},
		},
		{
			name: "invalidateblock",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("invalidateblock", "123")
			},
			staticCmd: func() interface{} {
				return acmjson.NewInvalidateBlockCmd("123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"invalidateblock","params":["123"],"id":1}`,
			unmarshalled: &acmjson.InvalidateBlockCmd{
				BlockHash: "123",
			},
		},
		{
			name: "ping",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("ping")
			},
			staticCmd: func() interface{} {
				return acmjson.NewPingCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"ping","params":[],"id":1}`,
			unmarshalled: &acmjson.PingCmd{},
		},
		{
			name: "preciousblock",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("preciousblock", "0123")
			},
			staticCmd: func() interface{} {
				return acmjson.NewPreciousBlockCmd("0123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"preciousblock","params":["0123"],"id":1}`,
			unmarshalled: &acmjson.PreciousBlockCmd{
				BlockHash: "0123",
			},
		},
		{
			name: "reconsiderblock",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("reconsiderblock", "123")
			},
			staticCmd: func() interface{} {
				return acmjson.NewReconsiderBlockCmd("123")
			},
			marshalled: `{"jsonrpc":"1.0","method":"reconsiderblock","params":["123"],"id":1}`,
			unmarshalled: &acmjson.ReconsiderBlockCmd{
				BlockHash: "123",
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("searchrawtransactions", "1Address")
			},
			staticCmd: func() interface{} {
				return acmjson.NewSearchRawTransactionsCmd("1Address", nil, nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address"],"id":1}`,
			unmarshalled: &acmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     acmjson.Int(1),
				Skip:        acmjson.Int(0),
				Count:       acmjson.Int(100),
				VinExtra:    acmjson.Int(0),
				Reverse:     acmjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("searchrawtransactions", "1Address", 0)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSearchRawTransactionsCmd("1Address",
					acmjson.Int(0), nil, nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0],"id":1}`,
			unmarshalled: &acmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     acmjson.Int(0),
				Skip:        acmjson.Int(0),
				Count:       acmjson.Int(100),
				VinExtra:    acmjson.Int(0),
				Reverse:     acmjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("searchrawtransactions", "1Address", 0, 5)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSearchRawTransactionsCmd("1Address",
					acmjson.Int(0), acmjson.Int(5), nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5],"id":1}`,
			unmarshalled: &acmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     acmjson.Int(0),
				Skip:        acmjson.Int(5),
				Count:       acmjson.Int(100),
				VinExtra:    acmjson.Int(0),
				Reverse:     acmjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSearchRawTransactionsCmd("1Address",
					acmjson.Int(0), acmjson.Int(5), acmjson.Int(10), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10],"id":1}`,
			unmarshalled: &acmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     acmjson.Int(0),
				Skip:        acmjson.Int(5),
				Count:       acmjson.Int(10),
				VinExtra:    acmjson.Int(0),
				Reverse:     acmjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSearchRawTransactionsCmd("1Address",
					acmjson.Int(0), acmjson.Int(5), acmjson.Int(10), acmjson.Int(1), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1],"id":1}`,
			unmarshalled: &acmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     acmjson.Int(0),
				Skip:        acmjson.Int(5),
				Count:       acmjson.Int(10),
				VinExtra:    acmjson.Int(1),
				Reverse:     acmjson.Bool(false),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSearchRawTransactionsCmd("1Address",
					acmjson.Int(0), acmjson.Int(5), acmjson.Int(10), acmjson.Int(1), acmjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true],"id":1}`,
			unmarshalled: &acmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     acmjson.Int(0),
				Skip:        acmjson.Int(5),
				Count:       acmjson.Int(10),
				VinExtra:    acmjson.Int(1),
				Reverse:     acmjson.Bool(true),
				FilterAddrs: nil,
			},
		},
		{
			name: "searchrawtransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("searchrawtransactions", "1Address", 0, 5, 10, 1, true, []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return acmjson.NewSearchRawTransactionsCmd("1Address",
					acmjson.Int(0), acmjson.Int(5), acmjson.Int(10), acmjson.Int(1), acmjson.Bool(true), &[]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"searchrawtransactions","params":["1Address",0,5,10,1,true,["1Address"]],"id":1}`,
			unmarshalled: &acmjson.SearchRawTransactionsCmd{
				Address:     "1Address",
				Verbose:     acmjson.Int(0),
				Skip:        acmjson.Int(5),
				Count:       acmjson.Int(10),
				VinExtra:    acmjson.Int(1),
				Reverse:     acmjson.Bool(true),
				FilterAddrs: &[]string{"1Address"},
			},
		},
		{
			name: "sendrawtransaction",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("sendrawtransaction", "1122")
			},
			staticCmd: func() interface{} {
				return acmjson.NewSendRawTransactionCmd("1122", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122"],"id":1}`,
			unmarshalled: &acmjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: acmjson.Bool(false),
			},
		},
		{
			name: "sendrawtransaction optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("sendrawtransaction", "1122", false)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSendRawTransactionCmd("1122", acmjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendrawtransaction","params":["1122",false],"id":1}`,
			unmarshalled: &acmjson.SendRawTransactionCmd{
				HexTx:         "1122",
				AllowHighFees: acmjson.Bool(false),
			},
		},
		{
			name: "setgenerate",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("setgenerate", true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSetGenerateCmd(true, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true],"id":1}`,
			unmarshalled: &acmjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: acmjson.Int(-1),
			},
		},
		{
			name: "setgenerate optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("setgenerate", true, 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSetGenerateCmd(true, acmjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"setgenerate","params":[true,6],"id":1}`,
			unmarshalled: &acmjson.SetGenerateCmd{
				Generate:     true,
				GenProcLimit: acmjson.Int(6),
			},
		},
		{
			name: "stop",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("stop")
			},
			staticCmd: func() interface{} {
				return acmjson.NewStopCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stop","params":[],"id":1}`,
			unmarshalled: &acmjson.StopCmd{},
		},
		{
			name: "submitblock",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("submitblock", "112233")
			},
			staticCmd: func() interface{} {
				return acmjson.NewSubmitBlockCmd("112233", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233"],"id":1}`,
			unmarshalled: &acmjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options:  nil,
			},
		},
		{
			name: "submitblock optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("submitblock", "112233", `{"workid":"12345"}`)
			},
			staticCmd: func() interface{} {
				options := acmjson.SubmitBlockOptions{
					WorkID: "12345",
				}
				return acmjson.NewSubmitBlockCmd("112233", &options)
			},
			marshalled: `{"jsonrpc":"1.0","method":"submitblock","params":["112233",{"workid":"12345"}],"id":1}`,
			unmarshalled: &acmjson.SubmitBlockCmd{
				HexBlock: "112233",
				Options: &acmjson.SubmitBlockOptions{
					WorkID: "12345",
				},
			},
		},
		{
			name: "uptime",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("uptime")
			},
			staticCmd: func() interface{} {
				return acmjson.NewUptimeCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"uptime","params":[],"id":1}`,
			unmarshalled: &acmjson.UptimeCmd{},
		},
		{
			name: "validateaddress",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("validateaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return acmjson.NewValidateAddressCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"validateaddress","params":["1Address"],"id":1}`,
			unmarshalled: &acmjson.ValidateAddressCmd{
				Address: "1Address",
			},
		},
		{
			name: "verifychain",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("verifychain")
			},
			staticCmd: func() interface{} {
				return acmjson.NewVerifyChainCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[],"id":1}`,
			unmarshalled: &acmjson.VerifyChainCmd{
				CheckLevel: acmjson.Int32(3),
				CheckDepth: acmjson.Int32(288),
			},
		},
		{
			name: "verifychain optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("verifychain", 2)
			},
			staticCmd: func() interface{} {
				return acmjson.NewVerifyChainCmd(acmjson.Int32(2), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2],"id":1}`,
			unmarshalled: &acmjson.VerifyChainCmd{
				CheckLevel: acmjson.Int32(2),
				CheckDepth: acmjson.Int32(288),
			},
		},
		{
			name: "verifychain optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("verifychain", 2, 500)
			},
			staticCmd: func() interface{} {
				return acmjson.NewVerifyChainCmd(acmjson.Int32(2), acmjson.Int32(500))
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifychain","params":[2,500],"id":1}`,
			unmarshalled: &acmjson.VerifyChainCmd{
				CheckLevel: acmjson.Int32(2),
				CheckDepth: acmjson.Int32(500),
			},
		},
		{
			name: "verifymessage",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("verifymessage", "1Address", "301234", "test")
			},
			staticCmd: func() interface{} {
				return acmjson.NewVerifyMessageCmd("1Address", "301234", "test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifymessage","params":["1Address","301234","test"],"id":1}`,
			unmarshalled: &acmjson.VerifyMessageCmd{
				Address:   "1Address",
				Signature: "301234",
				Message:   "test",
			},
		},
		{
			name: "verifytxoutproof",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("verifytxoutproof", "test")
			},
			staticCmd: func() interface{} {
				return acmjson.NewVerifyTxOutProofCmd("test")
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifytxoutproof","params":["test"],"id":1}`,
			unmarshalled: &acmjson.VerifyTxOutProofCmd{
				Proof: "test",
			},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command
		// creation function.
		marshalled, err := acmjson.MarshalCmd(testID, test.staticCmd())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			t.Errorf("\n%s\n%s", marshalled, test.marshalled)
			continue
		}

		// Ensure the command is created without error via the generic
		// new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the command as created by the generic new command
		// creation function.
		marshalled, err = acmjson.MarshalCmd(testID, cmd)
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		var request acmjson.Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}

		cmd, err = acmjson.UnmarshalCmd(&request)
		if err != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}

// TestChainSvrCmdErrors ensures any errors that occur in the command during
// custom mashal and unmarshal are as expected.
func TestChainSvrCmdErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		result     interface{}
		marshalled string
		err        error
	}{
		{
			name:       "template request with invalid type",
			result:     &acmjson.TemplateRequest{},
			marshalled: `{"mode":1}`,
			err:        &json.UnmarshalTypeError{},
		},
		{
			name:       "invalid template request sigoplimit field",
			result:     &acmjson.TemplateRequest{},
			marshalled: `{"sigoplimit":"invalid"}`,
			err:        acmjson.Error{ErrorCode: acmjson.ErrInvalidType},
		},
		{
			name:       "invalid template request sizelimit field",
			result:     &acmjson.TemplateRequest{},
			marshalled: `{"sizelimit":"invalid"}`,
			err:        acmjson.Error{ErrorCode: acmjson.ErrInvalidType},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := json.Unmarshal([]byte(test.marshalled), &test.result)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Test #%d (%s) wrong error - got %T (%v), "+
				"want %T", i, test.name, err, err, test.err)
			continue
		}

		if terr, ok := test.err.(acmjson.Error); ok {
			gotErrorCode := err.(acmjson.Error).ErrorCode
			if gotErrorCode != terr.ErrorCode {
				t.Errorf("Test #%d (%s) mismatched error code "+
					"- got %v (%v), want %v", i, test.name,
					gotErrorCode, terr, terr.ErrorCode)
				continue
			}
		}
	}
}
