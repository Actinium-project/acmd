// Copyright (c) 2014-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"time"

	"github.com/Actinium-project/acmd/chaincfg/chainhash"
	"github.com/Actinium-project/acmd/wire"
)

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the main network, regression test network, and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
	Version: 1,
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  chainhash.Hash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x04, 0xff, 0xff, 0x00, 0x1d, 0x01, 0x04, 0x45, // |.......E|
				0x4e, 0x59, 0x20, 0x54, 0x69, 0x6d, 0x65, 0x73, // |NY Times|
				0x20, 0x32, 0x34, 0x2f, 0x41, 0x70, 0x72, 0x2f, // | 24/Apr/|
				0x32, 0x30, 0x31, 0x38, 0x20, 0x54, 0x6f, 0x72, // |2018 Tor|
				0x6f, 0x6e, 0x74, 0x6f, 0x20, 0x56, 0x61, 0x6e, // |onto Van|
				0x20, 0x41, 0x74, 0x74, 0x61, 0x63, 0x6b, 0x20, // | Attack |
				0x53, 0x75, 0x73, 0x70, 0x65, 0x63, 0x74, 0x20, // |Suspect |
				0x45, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x65, // |Expresse|
				0x64, 0x20, 0x41, 0x6e, 0x67, 0x65, 0x72, 0x20, // |d Anger |
				0x61, 0x74, 0x20, 0x57, 0x6f, 0x6d, 0x65, 0x6e, // |at Women|

			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value: 0x12a05f200,
			PkScript: []byte{
				0x41, 0x04, 0x67, 0x8a, 0xfd, 0xb0, 0xfe, 0x55, /* |A.g....U| */
				0x48, 0x27, 0x19, 0x67, 0xf1, 0xa6, 0x71, 0x30, /* |H'.g..q0| */
				0xb7, 0x10, 0x5c, 0xd6, 0xa8, 0x28, 0xe0, 0x39, /* |..\..(.9| */
				0x09, 0xa6, 0x79, 0x62, 0xe0, 0xea, 0x1f, 0x61, /* |..yb...a| */
				0xde, 0xb6, 0x49, 0xf6, 0xbc, 0x3f, 0x4c, 0xef, /* |..I..?L.| */
				0x38, 0xc4, 0xf3, 0x55, 0x04, 0xe5, 0x1e, 0xc1, /* |8..U....| */
				0x12, 0xde, 0x5c, 0x38, 0x4d, 0xf7, 0xba, 0x0b, /* |..\8M...| */
				0x8d, 0x57, 0x8a, 0x4c, 0x70, 0x2b, 0x6b, 0xf1, /* |.W.Lp+k.| */
                0x1d, 0x5f, 0xac, /* |._.| */
			},
		},
	},
	LockTime: 0,
}

// genesisHash is the hash of the first block in the block chain for the main
// network (genesis block).
var genesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x24, 0xb8, 0x00, 0x8b, 0xe7, 0xfd, 0x5c, 0x03,
	0xb6, 0x21, 0xcf, 0xfd, 0xe3, 0xbb, 0xc1, 0x23,
	0x66, 0x27, 0x2c, 0x79, 0xbe, 0xa1, 0x49, 0x2f,
	0x56, 0x14, 0x37, 0xe2, 0x72, 0x78, 0xd7, 0x28,
 })

 // 28d77872e23714562f49a1be792c276623c1bbe3fdcf21b6035cfde78b00b824

// genesisMerkleRoot is the hash of the first transaction in the genesis block
// for the main network.
var genesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x91, 0x91, 0x2c, 0xef, 0xda, 0x3a, 0x88, 0x13,
	0x95, 0x28, 0xd6, 0xa4, 0xe6, 0x13, 0x7c, 0x68,
	0x12, 0xd9, 0xf4, 0x6d, 0xf8, 0x0d, 0xb4, 0x8c,
	0xf8, 0xad, 0x22, 0x2f, 0x0e, 0xb1, 0x55, 0xec,
 })

// genesisBlock defines the genesis block of the block chain which serves as the
// public transaction ledger for the main network.
var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,        // ec55b10e2f22adf88cb40df86df4d912687c13e6a4d6289513883adaef2c9191
		Timestamp:  time.Unix(1524649713, 0),
		Bits:       0x1e0ffff0,
		Nonce:      743522,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// regTestGenesisHash is the hash of the first block in the block chain for the
// regression test network (genesis block).
var regTestGenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xd2, 0x8a, 0x3f, 0x3b, 0x3c, 0x01, 0x48, 0xe3,
	0xa9, 0xc0, 0xea, 0x3b, 0xc7, 0x05, 0x13, 0xda,
	0xf0, 0xbd, 0xdd, 0x52, 0xc1, 0x3d, 0xfb, 0x1e,
	0x66, 0x8d, 0xc1, 0x31, 0x4f, 0xd1, 0x5c, 0xf0,
})

// regTestGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the regression test network.  It is the same as the merkle root for
// the main network.
var regTestGenesisMerkleRoot = genesisMerkleRoot

// regTestGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the regression test network.
var regTestGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: regTestGenesisMerkleRoot, // ec55b10e2f22adf88cb40df86df4d912687c13e6a4d6289513883adaef2c9191
		Timestamp:  time.Unix(1524650028, 0), // 2011-02-02 23:16:42 +0000 UTC
		Bits:       0x1e0ffff0,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      771587,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// testNet4GenesisHash is the hash of the first block in the block chain for the
// test network (version 4).
var testNet4GenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0x1b, 0x1f, 0xf8, 0xf9, 0x99, 0xf0, 0x9e, 0xd1,
	0xc2, 0xaf, 0xb1, 0x9a, 0x35, 0x71, 0xe8, 0xc4,
	0x0c, 0x85, 0x6b, 0x6d, 0x21, 0x36, 0x78, 0x49,
	0xf1, 0xbf, 0xc2, 0x4d, 0xc8, 0x7c, 0x61, 0x7c,
})

// 7c617cc84dc2bff1497836216d6b850cc4e871359ab1afc2d19ef099f9f81f1b

// testNet3GenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the test network (version 3).  It is the same as the merkle root
// for the main network.
var testNet3GenesisMerkleRoot = genesisMerkleRoot

// testNet4GenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the test network (version 4).  It is the same as the merkle root
// for the main network.
var testNet4GenesisMerkleRoot = genesisMerkleRoot

// testNet4GenesisBlock defines the genesis block of the block chain which
// serves as the public transaction ledger for the test network (version 4).
var testNet4GenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},          // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: testNet4GenesisMerkleRoot, // ec55b10e2f22adf88cb40df86df4d912687c13e6a4d6289513883adaef2c9191
		Timestamp:  time.Unix(1524650001, 0),
		Bits:       0x1e0ffff0,
		Nonce:      19905,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// simNetGenesisHash is the hash of the first block in the block chain for the
// simulation test network.
var simNetGenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
	0xbe, 0xa4, 0x3b, 0x3a, 0xb4, 0xbc, 0x9d, 0x86,
	0x2b, 0xfa, 0xbc, 0x0d, 0x45, 0xa3, 0xc5, 0xc9,
	0xd0, 0x9f, 0x66, 0x20, 0xce, 0x7d, 0xb6, 0x78,
	0xb5, 0xdc, 0x1d, 0xb8, 0x19, 0xc9, 0x35, 0x12,
})

// simNetGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the simulation test network.  It is the same as the merkle root for
// the main network.
var simNetGenesisMerkleRoot = genesisMerkleRoot

// simNetGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the simulation test network.
var simNetGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  chainhash.Hash{},         // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: simNetGenesisMerkleRoot,  // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1401292357, 0), // 2014-05-28 15:52:37 +0000 UTC
		Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      2,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}
