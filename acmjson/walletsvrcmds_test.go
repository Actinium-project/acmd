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
)

// TestWalletSvrCmds tests all of the wallet server commands marshal and
// unmarshal into valid results include handling of optional fields being
// omitted in the marshalled command, while optional fields with defaults have
// the default assigned on unmarshalled commands.
func TestWalletSvrCmds(t *testing.T) {
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
			name: "addmultisigaddress",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return acmjson.NewAddMultisigAddressCmd(2, keys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &acmjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   nil,
			},
		},
		{
			name: "addmultisigaddress optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"}, "test")
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return acmjson.NewAddMultisigAddressCmd(2, keys, acmjson.String("test"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"],"test"],"id":1}`,
			unmarshalled: &acmjson.AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   acmjson.String("test"),
			},
		},
		{
			name: "addwitnessaddress",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("addwitnessaddress", "1address")
			},
			staticCmd: func() interface{} {
				return acmjson.NewAddWitnessAddressCmd("1address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"addwitnessaddress","params":["1address"],"id":1}`,
			unmarshalled: &acmjson.AddWitnessAddressCmd{
				Address: "1address",
			},
		},
		{
			name: "createmultisig",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("createmultisig", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return acmjson.NewCreateMultisigCmd(2, keys)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createmultisig","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &acmjson.CreateMultisigCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
			},
		},
		{
			name: "dumpprivkey",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("dumpprivkey", "1Address")
			},
			staticCmd: func() interface{} {
				return acmjson.NewDumpPrivKeyCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpprivkey","params":["1Address"],"id":1}`,
			unmarshalled: &acmjson.DumpPrivKeyCmd{
				Address: "1Address",
			},
		},
		{
			name: "encryptwallet",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("encryptwallet", "pass")
			},
			staticCmd: func() interface{} {
				return acmjson.NewEncryptWalletCmd("pass")
			},
			marshalled: `{"jsonrpc":"1.0","method":"encryptwallet","params":["pass"],"id":1}`,
			unmarshalled: &acmjson.EncryptWalletCmd{
				Passphrase: "pass",
			},
		},
		{
			name: "estimatefee",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("estimatefee", 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewEstimateFeeCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatefee","params":[6],"id":1}`,
			unmarshalled: &acmjson.EstimateFeeCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "estimatepriority",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("estimatepriority", 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewEstimatePriorityCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatepriority","params":[6],"id":1}`,
			unmarshalled: &acmjson.EstimatePriorityCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "getaccount",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getaccount", "1Address")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetAccountCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccount","params":["1Address"],"id":1}`,
			unmarshalled: &acmjson.GetAccountCmd{
				Address: "1Address",
			},
		},
		{
			name: "getaccountaddress",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getaccountaddress", "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetAccountAddressCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccountaddress","params":["acct"],"id":1}`,
			unmarshalled: &acmjson.GetAccountAddressCmd{
				Account: "acct",
			},
		},
		{
			name: "getaddressesbyaccount",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getaddressesbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetAddressesByAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddressesbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &acmjson.GetAddressesByAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "getbalance",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getbalance")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBalanceCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":[],"id":1}`,
			unmarshalled: &acmjson.GetBalanceCmd{
				Account: nil,
				MinConf: acmjson.Int(1),
			},
		},
		{
			name: "getbalance optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getbalance", "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBalanceCmd(acmjson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct"],"id":1}`,
			unmarshalled: &acmjson.GetBalanceCmd{
				Account: acmjson.String("acct"),
				MinConf: acmjson.Int(1),
			},
		},
		{
			name: "getbalance optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getbalance", "acct", 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetBalanceCmd(acmjson.String("acct"), acmjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct",6],"id":1}`,
			unmarshalled: &acmjson.GetBalanceCmd{
				Account: acmjson.String("acct"),
				MinConf: acmjson.Int(6),
			},
		},
		{
			name: "getnewaddress",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getnewaddress")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetNewAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":[],"id":1}`,
			unmarshalled: &acmjson.GetNewAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getnewaddress optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getnewaddress", "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetNewAddressCmd(acmjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":["acct"],"id":1}`,
			unmarshalled: &acmjson.GetNewAddressCmd{
				Account: acmjson.String("acct"),
			},
		},
		{
			name: "getrawchangeaddress",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getrawchangeaddress")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetRawChangeAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":[],"id":1}`,
			unmarshalled: &acmjson.GetRawChangeAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getrawchangeaddress optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getrawchangeaddress", "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetRawChangeAddressCmd(acmjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":["acct"],"id":1}`,
			unmarshalled: &acmjson.GetRawChangeAddressCmd{
				Account: acmjson.String("acct"),
			},
		},
		{
			name: "getreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getreceivedbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetReceivedByAccountCmd("acct", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &acmjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: acmjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaccount optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getreceivedbyaccount", "acct", 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetReceivedByAccountCmd("acct", acmjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct",6],"id":1}`,
			unmarshalled: &acmjson.GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: acmjson.Int(6),
			},
		},
		{
			name: "getreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getreceivedbyaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetReceivedByAddressCmd("1Address", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address"],"id":1}`,
			unmarshalled: &acmjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: acmjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaddress optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getreceivedbyaddress", "1Address", 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetReceivedByAddressCmd("1Address", acmjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address",6],"id":1}`,
			unmarshalled: &acmjson.GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: acmjson.Int(6),
			},
		},
		{
			name: "gettransaction",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("gettransaction", "123")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123"],"id":1}`,
			unmarshalled: &acmjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "gettransaction optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("gettransaction", "123", true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetTransactionCmd("123", acmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123",true],"id":1}`,
			unmarshalled: &acmjson.GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: acmjson.Bool(true),
			},
		},
		{
			name: "getwalletinfo",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getwalletinfo")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetWalletInfoCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"getwalletinfo","params":[],"id":1}`,
			unmarshalled: &acmjson.GetWalletInfoCmd{},
		},
		{
			name: "importprivkey",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("importprivkey", "abc")
			},
			staticCmd: func() interface{} {
				return acmjson.NewImportPrivKeyCmd("abc", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc"],"id":1}`,
			unmarshalled: &acmjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   nil,
				Rescan:  acmjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("importprivkey", "abc", "label")
			},
			staticCmd: func() interface{} {
				return acmjson.NewImportPrivKeyCmd("abc", acmjson.String("label"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label"],"id":1}`,
			unmarshalled: &acmjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   acmjson.String("label"),
				Rescan:  acmjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("importprivkey", "abc", "label", false)
			},
			staticCmd: func() interface{} {
				return acmjson.NewImportPrivKeyCmd("abc", acmjson.String("label"), acmjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false],"id":1}`,
			unmarshalled: &acmjson.ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   acmjson.String("label"),
				Rescan:  acmjson.Bool(false),
			},
		},
		{
			name: "keypoolrefill",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("keypoolrefill")
			},
			staticCmd: func() interface{} {
				return acmjson.NewKeyPoolRefillCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[],"id":1}`,
			unmarshalled: &acmjson.KeyPoolRefillCmd{
				NewSize: acmjson.Uint(100),
			},
		},
		{
			name: "keypoolrefill optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("keypoolrefill", 200)
			},
			staticCmd: func() interface{} {
				return acmjson.NewKeyPoolRefillCmd(acmjson.Uint(200))
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[200],"id":1}`,
			unmarshalled: &acmjson.KeyPoolRefillCmd{
				NewSize: acmjson.Uint(200),
			},
		},
		{
			name: "listaccounts",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listaccounts")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListAccountsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[],"id":1}`,
			unmarshalled: &acmjson.ListAccountsCmd{
				MinConf: acmjson.Int(1),
			},
		},
		{
			name: "listaccounts optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listaccounts", 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListAccountsCmd(acmjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[6],"id":1}`,
			unmarshalled: &acmjson.ListAccountsCmd{
				MinConf: acmjson.Int(6),
			},
		},
		{
			name: "listaddressgroupings",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listaddressgroupings")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListAddressGroupingsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listaddressgroupings","params":[],"id":1}`,
			unmarshalled: &acmjson.ListAddressGroupingsCmd{},
		},
		{
			name: "listlockunspent",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listlockunspent")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListLockUnspentCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listlockunspent","params":[],"id":1}`,
			unmarshalled: &acmjson.ListLockUnspentCmd{},
		},
		{
			name: "listreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listreceivedbyaccount")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListReceivedByAccountCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[],"id":1}`,
			unmarshalled: &acmjson.ListReceivedByAccountCmd{
				MinConf:          acmjson.Int(1),
				IncludeEmpty:     acmjson.Bool(false),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listreceivedbyaccount", 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListReceivedByAccountCmd(acmjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6],"id":1}`,
			unmarshalled: &acmjson.ListReceivedByAccountCmd{
				MinConf:          acmjson.Int(6),
				IncludeEmpty:     acmjson.Bool(false),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listreceivedbyaccount", 6, true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListReceivedByAccountCmd(acmjson.Int(6), acmjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true],"id":1}`,
			unmarshalled: &acmjson.ListReceivedByAccountCmd{
				MinConf:          acmjson.Int(6),
				IncludeEmpty:     acmjson.Bool(true),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional3",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listreceivedbyaccount", 6, true, false)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListReceivedByAccountCmd(acmjson.Int(6), acmjson.Bool(true), acmjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true,false],"id":1}`,
			unmarshalled: &acmjson.ListReceivedByAccountCmd{
				MinConf:          acmjson.Int(6),
				IncludeEmpty:     acmjson.Bool(true),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listreceivedbyaddress")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListReceivedByAddressCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[],"id":1}`,
			unmarshalled: &acmjson.ListReceivedByAddressCmd{
				MinConf:          acmjson.Int(1),
				IncludeEmpty:     acmjson.Bool(false),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listreceivedbyaddress", 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListReceivedByAddressCmd(acmjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6],"id":1}`,
			unmarshalled: &acmjson.ListReceivedByAddressCmd{
				MinConf:          acmjson.Int(6),
				IncludeEmpty:     acmjson.Bool(false),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listreceivedbyaddress", 6, true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListReceivedByAddressCmd(acmjson.Int(6), acmjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true],"id":1}`,
			unmarshalled: &acmjson.ListReceivedByAddressCmd{
				MinConf:          acmjson.Int(6),
				IncludeEmpty:     acmjson.Bool(true),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional3",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listreceivedbyaddress", 6, true, false)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListReceivedByAddressCmd(acmjson.Int(6), acmjson.Bool(true), acmjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true,false],"id":1}`,
			unmarshalled: &acmjson.ListReceivedByAddressCmd{
				MinConf:          acmjson.Int(6),
				IncludeEmpty:     acmjson.Bool(true),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listsinceblock",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listsinceblock")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListSinceBlockCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":[],"id":1}`,
			unmarshalled: &acmjson.ListSinceBlockCmd{
				BlockHash:           nil,
				TargetConfirmations: acmjson.Int(1),
				IncludeWatchOnly:    acmjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listsinceblock", "123")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListSinceBlockCmd(acmjson.String("123"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123"],"id":1}`,
			unmarshalled: &acmjson.ListSinceBlockCmd{
				BlockHash:           acmjson.String("123"),
				TargetConfirmations: acmjson.Int(1),
				IncludeWatchOnly:    acmjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listsinceblock", "123", 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListSinceBlockCmd(acmjson.String("123"), acmjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6],"id":1}`,
			unmarshalled: &acmjson.ListSinceBlockCmd{
				BlockHash:           acmjson.String("123"),
				TargetConfirmations: acmjson.Int(6),
				IncludeWatchOnly:    acmjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional3",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listsinceblock", "123", 6, true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListSinceBlockCmd(acmjson.String("123"), acmjson.Int(6), acmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6,true],"id":1}`,
			unmarshalled: &acmjson.ListSinceBlockCmd{
				BlockHash:           acmjson.String("123"),
				TargetConfirmations: acmjson.Int(6),
				IncludeWatchOnly:    acmjson.Bool(true),
			},
		},
		{
			name: "listtransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listtransactions")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListTransactionsCmd(nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":[],"id":1}`,
			unmarshalled: &acmjson.ListTransactionsCmd{
				Account:          nil,
				Count:            acmjson.Int(10),
				From:             acmjson.Int(0),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listtransactions", "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListTransactionsCmd(acmjson.String("acct"), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct"],"id":1}`,
			unmarshalled: &acmjson.ListTransactionsCmd{
				Account:          acmjson.String("acct"),
				Count:            acmjson.Int(10),
				From:             acmjson.Int(0),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listtransactions", "acct", 20)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListTransactionsCmd(acmjson.String("acct"), acmjson.Int(20), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20],"id":1}`,
			unmarshalled: &acmjson.ListTransactionsCmd{
				Account:          acmjson.String("acct"),
				Count:            acmjson.Int(20),
				From:             acmjson.Int(0),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional3",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listtransactions", "acct", 20, 1)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListTransactionsCmd(acmjson.String("acct"), acmjson.Int(20),
					acmjson.Int(1), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1],"id":1}`,
			unmarshalled: &acmjson.ListTransactionsCmd{
				Account:          acmjson.String("acct"),
				Count:            acmjson.Int(20),
				From:             acmjson.Int(1),
				IncludeWatchOnly: acmjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional4",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listtransactions", "acct", 20, 1, true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListTransactionsCmd(acmjson.String("acct"), acmjson.Int(20),
					acmjson.Int(1), acmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1,true],"id":1}`,
			unmarshalled: &acmjson.ListTransactionsCmd{
				Account:          acmjson.String("acct"),
				Count:            acmjson.Int(20),
				From:             acmjson.Int(1),
				IncludeWatchOnly: acmjson.Bool(true),
			},
		},
		{
			name: "listunspent",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listunspent")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListUnspentCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[],"id":1}`,
			unmarshalled: &acmjson.ListUnspentCmd{
				MinConf:   acmjson.Int(1),
				MaxConf:   acmjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listunspent", 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListUnspentCmd(acmjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6],"id":1}`,
			unmarshalled: &acmjson.ListUnspentCmd{
				MinConf:   acmjson.Int(6),
				MaxConf:   acmjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listunspent", 6, 100)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListUnspentCmd(acmjson.Int(6), acmjson.Int(100), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100],"id":1}`,
			unmarshalled: &acmjson.ListUnspentCmd{
				MinConf:   acmjson.Int(6),
				MaxConf:   acmjson.Int(100),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional3",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listunspent", 6, 100, []string{"1Address", "1Address2"})
			},
			staticCmd: func() interface{} {
				return acmjson.NewListUnspentCmd(acmjson.Int(6), acmjson.Int(100),
					&[]string{"1Address", "1Address2"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100,["1Address","1Address2"]],"id":1}`,
			unmarshalled: &acmjson.ListUnspentCmd{
				MinConf:   acmjson.Int(6),
				MaxConf:   acmjson.Int(100),
				Addresses: &[]string{"1Address", "1Address2"},
			},
		},
		{
			name: "lockunspent",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("lockunspent", true, `[{"txid":"123","vout":1}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []acmjson.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				return acmjson.NewLockUnspentCmd(true, txInputs)
			},
			marshalled: `{"jsonrpc":"1.0","method":"lockunspent","params":[true,[{"txid":"123","vout":1}]],"id":1}`,
			unmarshalled: &acmjson.LockUnspentCmd{
				Unlock: true,
				Transactions: []acmjson.TransactionInput{
					{Txid: "123", Vout: 1},
				},
			},
		},
		{
			name: "move",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("move", "from", "to", 0.5)
			},
			staticCmd: func() interface{} {
				return acmjson.NewMoveCmd("from", "to", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5],"id":1}`,
			unmarshalled: &acmjson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     acmjson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "move optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("move", "from", "to", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewMoveCmd("from", "to", 0.5, acmjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5,6],"id":1}`,
			unmarshalled: &acmjson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     acmjson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "move optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("move", "from", "to", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return acmjson.NewMoveCmd("from", "to", 0.5, acmjson.Int(6), acmjson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"move","params":["from","to",0.5,6,"comment"],"id":1}`,
			unmarshalled: &acmjson.MoveCmd{
				FromAccount: "from",
				ToAccount:   "to",
				Amount:      0.5,
				MinConf:     acmjson.Int(6),
				Comment:     acmjson.String("comment"),
			},
		},
		{
			name: "sendfrom",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("sendfrom", "from", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSendFromCmd("from", "1Address", 0.5, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5],"id":1}`,
			unmarshalled: &acmjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     acmjson.Int(1),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSendFromCmd("from", "1Address", 0.5, acmjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6],"id":1}`,
			unmarshalled: &acmjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     acmjson.Int(6),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return acmjson.NewSendFromCmd("from", "1Address", 0.5, acmjson.Int(6),
					acmjson.String("comment"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment"],"id":1}`,
			unmarshalled: &acmjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     acmjson.Int(6),
				Comment:     acmjson.String("comment"),
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional3",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return acmjson.NewSendFromCmd("from", "1Address", 0.5, acmjson.Int(6),
					acmjson.String("comment"), acmjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment","commentto"],"id":1}`,
			unmarshalled: &acmjson.SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     acmjson.Int(6),
				Comment:     acmjson.String("comment"),
				CommentTo:   acmjson.String("commentto"),
			},
		},
		{
			name: "sendmany",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("sendmany", "from", `{"1Address":0.5}`)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return acmjson.NewSendManyCmd("from", amounts, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5}],"id":1}`,
			unmarshalled: &acmjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     acmjson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return acmjson.NewSendManyCmd("from", amounts, acmjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6],"id":1}`,
			unmarshalled: &acmjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     acmjson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6, "comment")
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return acmjson.NewSendManyCmd("from", amounts, acmjson.Int(6), acmjson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6,"comment"],"id":1}`,
			unmarshalled: &acmjson.SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     acmjson.Int(6),
				Comment:     acmjson.String("comment"),
			},
		},
		{
			name: "sendtoaddress",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("sendtoaddress", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSendToAddressCmd("1Address", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5],"id":1}`,
			unmarshalled: &acmjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   nil,
				CommentTo: nil,
			},
		},
		{
			name: "sendtoaddress optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("sendtoaddress", "1Address", 0.5, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return acmjson.NewSendToAddressCmd("1Address", 0.5, acmjson.String("comment"),
					acmjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5,"comment","commentto"],"id":1}`,
			unmarshalled: &acmjson.SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   acmjson.String("comment"),
				CommentTo: acmjson.String("commentto"),
			},
		},
		{
			name: "setaccount",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("setaccount", "1Address", "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewSetAccountCmd("1Address", "acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"setaccount","params":["1Address","acct"],"id":1}`,
			unmarshalled: &acmjson.SetAccountCmd{
				Address: "1Address",
				Account: "acct",
			},
		},
		{
			name: "settxfee",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("settxfee", 0.0001)
			},
			staticCmd: func() interface{} {
				return acmjson.NewSetTxFeeCmd(0.0001)
			},
			marshalled: `{"jsonrpc":"1.0","method":"settxfee","params":[0.0001],"id":1}`,
			unmarshalled: &acmjson.SetTxFeeCmd{
				Amount: 0.0001,
			},
		},
		{
			name: "signmessage",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("signmessage", "1Address", "message")
			},
			staticCmd: func() interface{} {
				return acmjson.NewSignMessageCmd("1Address", "message")
			},
			marshalled: `{"jsonrpc":"1.0","method":"signmessage","params":["1Address","message"],"id":1}`,
			unmarshalled: &acmjson.SignMessageCmd{
				Address: "1Address",
				Message: "message",
			},
		},
		{
			name: "signrawtransaction",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("signrawtransaction", "001122")
			},
			staticCmd: func() interface{} {
				return acmjson.NewSignRawTransactionCmd("001122", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122"],"id":1}`,
			unmarshalled: &acmjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   nil,
				PrivKeys: nil,
				Flags:    acmjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("signrawtransaction", "001122", `[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []acmjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				}

				return acmjson.NewSignRawTransactionCmd("001122", &txInputs, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[{"txid":"123","vout":1,"scriptPubKey":"00","redeemScript":"01"}]],"id":1}`,
			unmarshalled: &acmjson.SignRawTransactionCmd{
				RawTx: "001122",
				Inputs: &[]acmjson.RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				},
				PrivKeys: nil,
				Flags:    acmjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("signrawtransaction", "001122", `[]`, `["abc"]`)
			},
			staticCmd: func() interface{} {
				txInputs := []acmjson.RawTxInput{}
				privKeys := []string{"abc"}
				return acmjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],["abc"]],"id":1}`,
			unmarshalled: &acmjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]acmjson.RawTxInput{},
				PrivKeys: &[]string{"abc"},
				Flags:    acmjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional3",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("signrawtransaction", "001122", `[]`, `[]`, "ALL")
			},
			staticCmd: func() interface{} {
				txInputs := []acmjson.RawTxInput{}
				privKeys := []string{}
				return acmjson.NewSignRawTransactionCmd("001122", &txInputs, &privKeys,
					acmjson.String("ALL"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],[],"ALL"],"id":1}`,
			unmarshalled: &acmjson.SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]acmjson.RawTxInput{},
				PrivKeys: &[]string{},
				Flags:    acmjson.String("ALL"),
			},
		},
		{
			name: "walletlock",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("walletlock")
			},
			staticCmd: func() interface{} {
				return acmjson.NewWalletLockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletlock","params":[],"id":1}`,
			unmarshalled: &acmjson.WalletLockCmd{},
		},
		{
			name: "walletpassphrase",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("walletpassphrase", "pass", 60)
			},
			staticCmd: func() interface{} {
				return acmjson.NewWalletPassphraseCmd("pass", 60)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrase","params":["pass",60],"id":1}`,
			unmarshalled: &acmjson.WalletPassphraseCmd{
				Passphrase: "pass",
				Timeout:    60,
			},
		},
		{
			name: "walletpassphrasechange",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("walletpassphrasechange", "old", "new")
			},
			staticCmd: func() interface{} {
				return acmjson.NewWalletPassphraseChangeCmd("old", "new")
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrasechange","params":["old","new"],"id":1}`,
			unmarshalled: &acmjson.WalletPassphraseChangeCmd{
				OldPassphrase: "old",
				NewPassphrase: "new",
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
