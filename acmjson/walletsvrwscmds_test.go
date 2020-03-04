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

// TestWalletSvrWsCmds tests all of the wallet server websocket-specific
// commands marshal and unmarshal into valid results include handling of
// optional fields being omitted in the marshalled command, while optional
// fields with defaults have the default assigned on unmarshalled commands.
func TestWalletSvrWsCmds(t *testing.T) {
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
			name: "createencryptedwallet",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("createencryptedwallet", "pass")
			},
			staticCmd: func() interface{} {
				return acmjson.NewCreateEncryptedWalletCmd("pass")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"createencryptedwallet","params":["pass"],"id":1}`,
			unmarshalled: &acmjson.CreateEncryptedWalletCmd{Passphrase: "pass"},
		},
		{
			name: "exportwatchingwallet",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("exportwatchingwallet")
			},
			staticCmd: func() interface{} {
				return acmjson.NewExportWatchingWalletCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"exportwatchingwallet","params":[],"id":1}`,
			unmarshalled: &acmjson.ExportWatchingWalletCmd{
				Account:  nil,
				Download: acmjson.Bool(false),
			},
		},
		{
			name: "exportwatchingwallet optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("exportwatchingwallet", "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewExportWatchingWalletCmd(acmjson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"exportwatchingwallet","params":["acct"],"id":1}`,
			unmarshalled: &acmjson.ExportWatchingWalletCmd{
				Account:  acmjson.String("acct"),
				Download: acmjson.Bool(false),
			},
		},
		{
			name: "exportwatchingwallet optional2",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("exportwatchingwallet", "acct", true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewExportWatchingWalletCmd(acmjson.String("acct"),
					acmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"exportwatchingwallet","params":["acct",true],"id":1}`,
			unmarshalled: &acmjson.ExportWatchingWalletCmd{
				Account:  acmjson.String("acct"),
				Download: acmjson.Bool(true),
			},
		},
		{
			name: "getunconfirmedbalance",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getunconfirmedbalance")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetUnconfirmedBalanceCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getunconfirmedbalance","params":[],"id":1}`,
			unmarshalled: &acmjson.GetUnconfirmedBalanceCmd{
				Account: nil,
			},
		},
		{
			name: "getunconfirmedbalance optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("getunconfirmedbalance", "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewGetUnconfirmedBalanceCmd(acmjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getunconfirmedbalance","params":["acct"],"id":1}`,
			unmarshalled: &acmjson.GetUnconfirmedBalanceCmd{
				Account: acmjson.String("acct"),
			},
		},
		{
			name: "listaddresstransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listaddresstransactions", `["1Address"]`)
			},
			staticCmd: func() interface{} {
				return acmjson.NewListAddressTransactionsCmd([]string{"1Address"}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaddresstransactions","params":[["1Address"]],"id":1}`,
			unmarshalled: &acmjson.ListAddressTransactionsCmd{
				Addresses: []string{"1Address"},
				Account:   nil,
			},
		},
		{
			name: "listaddresstransactions optional1",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listaddresstransactions", `["1Address"]`, "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListAddressTransactionsCmd([]string{"1Address"},
					acmjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaddresstransactions","params":[["1Address"],"acct"],"id":1}`,
			unmarshalled: &acmjson.ListAddressTransactionsCmd{
				Addresses: []string{"1Address"},
				Account:   acmjson.String("acct"),
			},
		},
		{
			name: "listalltransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listalltransactions")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListAllTransactionsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listalltransactions","params":[],"id":1}`,
			unmarshalled: &acmjson.ListAllTransactionsCmd{
				Account: nil,
			},
		},
		{
			name: "listalltransactions optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("listalltransactions", "acct")
			},
			staticCmd: func() interface{} {
				return acmjson.NewListAllTransactionsCmd(acmjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listalltransactions","params":["acct"],"id":1}`,
			unmarshalled: &acmjson.ListAllTransactionsCmd{
				Account: acmjson.String("acct"),
			},
		},
		{
			name: "recoveraddresses",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("recoveraddresses", "acct", 10)
			},
			staticCmd: func() interface{} {
				return acmjson.NewRecoverAddressesCmd("acct", 10)
			},
			marshalled: `{"jsonrpc":"1.0","method":"recoveraddresses","params":["acct",10],"id":1}`,
			unmarshalled: &acmjson.RecoverAddressesCmd{
				Account: "acct",
				N:       10,
			},
		},
		{
			name: "walletislocked",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("walletislocked")
			},
			staticCmd: func() interface{} {
				return acmjson.NewWalletIsLockedCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletislocked","params":[],"id":1}`,
			unmarshalled: &acmjson.WalletIsLockedCmd{},
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
