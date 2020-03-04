// Copyright (c) 2014-2017 The btcsuite developers
// Copyright (c) 2015-2017 The Decred developers
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

// TestChainSvrWsCmds tests all of the chain server websocket-specific commands
// marshal and unmarshal into valid results include handling of optional fields
// being omitted in the marshalled command, while optional fields with defaults
// have the default assigned on unmarshalled commands.
func TestChainSvrWsCmds(t *testing.T) {
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
			name: "authenticate",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("authenticate", "user", "pass")
			},
			staticCmd: func() interface{} {
				return acmjson.NewAuthenticateCmd("user", "pass")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"authenticate","params":["user","pass"],"id":1}`,
			unmarshalled: &acmjson.AuthenticateCmd{Username: "user", Passphrase: "pass"},
		},
		{
			name: "notifyblocks",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("notifyblocks")
			},
			staticCmd: func() interface{} {
				return acmjson.NewNotifyBlocksCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"notifyblocks","params":[],"id":1}`,
			unmarshalled: &acmjson.NotifyBlocksCmd{},
		},
		{
			name: "stopnotifyblocks",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("stopnotifyblocks")
			},
			staticCmd: func() interface{} {
				return acmjson.NewStopNotifyBlocksCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stopnotifyblocks","params":[],"id":1}`,
			unmarshalled: &acmjson.StopNotifyBlocksCmd{},
		},
		{
			name: "notifynewtransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("notifynewtransactions")
			},
			staticCmd: func() interface{} {
				return acmjson.NewNotifyNewTransactionsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"notifynewtransactions","params":[],"id":1}`,
			unmarshalled: &acmjson.NotifyNewTransactionsCmd{
				Verbose: acmjson.Bool(false),
			},
		},
		{
			name: "notifynewtransactions optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("notifynewtransactions", true)
			},
			staticCmd: func() interface{} {
				return acmjson.NewNotifyNewTransactionsCmd(acmjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"notifynewtransactions","params":[true],"id":1}`,
			unmarshalled: &acmjson.NotifyNewTransactionsCmd{
				Verbose: acmjson.Bool(true),
			},
		},
		{
			name: "stopnotifynewtransactions",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("stopnotifynewtransactions")
			},
			staticCmd: func() interface{} {
				return acmjson.NewStopNotifyNewTransactionsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stopnotifynewtransactions","params":[],"id":1}`,
			unmarshalled: &acmjson.StopNotifyNewTransactionsCmd{},
		},
		{
			name: "notifyreceived",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("notifyreceived", []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return acmjson.NewNotifyReceivedCmd([]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"notifyreceived","params":[["1Address"]],"id":1}`,
			unmarshalled: &acmjson.NotifyReceivedCmd{
				Addresses: []string{"1Address"},
			},
		},
		{
			name: "stopnotifyreceived",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("stopnotifyreceived", []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return acmjson.NewStopNotifyReceivedCmd([]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"stopnotifyreceived","params":[["1Address"]],"id":1}`,
			unmarshalled: &acmjson.StopNotifyReceivedCmd{
				Addresses: []string{"1Address"},
			},
		},
		{
			name: "notifyspent",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("notifyspent", `[{"hash":"123","index":0}]`)
			},
			staticCmd: func() interface{} {
				ops := []acmjson.OutPoint{{Hash: "123", Index: 0}}
				return acmjson.NewNotifySpentCmd(ops)
			},
			marshalled: `{"jsonrpc":"1.0","method":"notifyspent","params":[[{"hash":"123","index":0}]],"id":1}`,
			unmarshalled: &acmjson.NotifySpentCmd{
				OutPoints: []acmjson.OutPoint{{Hash: "123", Index: 0}},
			},
		},
		{
			name: "stopnotifyspent",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("stopnotifyspent", `[{"hash":"123","index":0}]`)
			},
			staticCmd: func() interface{} {
				ops := []acmjson.OutPoint{{Hash: "123", Index: 0}}
				return acmjson.NewStopNotifySpentCmd(ops)
			},
			marshalled: `{"jsonrpc":"1.0","method":"stopnotifyspent","params":[[{"hash":"123","index":0}]],"id":1}`,
			unmarshalled: &acmjson.StopNotifySpentCmd{
				OutPoints: []acmjson.OutPoint{{Hash: "123", Index: 0}},
			},
		},
		{
			name: "rescan",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("rescan", "123", `["1Address"]`, `[{"hash":"0000000000000000000000000000000000000000000000000000000000000123","index":0}]`)
			},
			staticCmd: func() interface{} {
				addrs := []string{"1Address"}
				ops := []acmjson.OutPoint{{
					Hash:  "0000000000000000000000000000000000000000000000000000000000000123",
					Index: 0,
				}}
				return acmjson.NewRescanCmd("123", addrs, ops, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"rescan","params":["123",["1Address"],[{"hash":"0000000000000000000000000000000000000000000000000000000000000123","index":0}]],"id":1}`,
			unmarshalled: &acmjson.RescanCmd{
				BeginBlock: "123",
				Addresses:  []string{"1Address"},
				OutPoints:  []acmjson.OutPoint{{Hash: "0000000000000000000000000000000000000000000000000000000000000123", Index: 0}},
				EndBlock:   nil,
			},
		},
		{
			name: "rescan optional",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("rescan", "123", `["1Address"]`, `[{"hash":"123","index":0}]`, "456")
			},
			staticCmd: func() interface{} {
				addrs := []string{"1Address"}
				ops := []acmjson.OutPoint{{Hash: "123", Index: 0}}
				return acmjson.NewRescanCmd("123", addrs, ops, acmjson.String("456"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"rescan","params":["123",["1Address"],[{"hash":"123","index":0}],"456"],"id":1}`,
			unmarshalled: &acmjson.RescanCmd{
				BeginBlock: "123",
				Addresses:  []string{"1Address"},
				OutPoints:  []acmjson.OutPoint{{Hash: "123", Index: 0}},
				EndBlock:   acmjson.String("456"),
			},
		},
		{
			name: "loadtxfilter",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("loadtxfilter", false, `["1Address"]`, `[{"hash":"0000000000000000000000000000000000000000000000000000000000000123","index":0}]`)
			},
			staticCmd: func() interface{} {
				addrs := []string{"1Address"}
				ops := []acmjson.OutPoint{{
					Hash:  "0000000000000000000000000000000000000000000000000000000000000123",
					Index: 0,
				}}
				return acmjson.NewLoadTxFilterCmd(false, addrs, ops)
			},
			marshalled: `{"jsonrpc":"1.0","method":"loadtxfilter","params":[false,["1Address"],[{"hash":"0000000000000000000000000000000000000000000000000000000000000123","index":0}]],"id":1}`,
			unmarshalled: &acmjson.LoadTxFilterCmd{
				Reload:    false,
				Addresses: []string{"1Address"},
				OutPoints: []acmjson.OutPoint{{Hash: "0000000000000000000000000000000000000000000000000000000000000123", Index: 0}},
			},
		},
		{
			name: "rescanblocks",
			newCmd: func() (interface{}, error) {
				return acmjson.NewCmd("rescanblocks", `["0000000000000000000000000000000000000000000000000000000000000123"]`)
			},
			staticCmd: func() interface{} {
				blockhashes := []string{"0000000000000000000000000000000000000000000000000000000000000123"}
				return acmjson.NewRescanBlocksCmd(blockhashes)
			},
			marshalled: `{"jsonrpc":"1.0","method":"rescanblocks","params":[["0000000000000000000000000000000000000000000000000000000000000123"]],"id":1}`,
			unmarshalled: &acmjson.RescanBlocksCmd{
				BlockHashes: []string{"0000000000000000000000000000000000000000000000000000000000000123"},
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
