<!-- toc -->

# Getting Started With Asimov

This document will introduce how to develop simple contracts and
deploying to Asimov blockchain in few minutes.

# Try Asimov

Your can start a solo instance of asimov node for the usage of
development. The following programs is currently available for macOS and
Linux.

- Asimov Full Node ```asimovd```
- Asimov Toolchain ```asimov-tool```

## Asimov Full Node

### Clone Asimov

```sh
go get github.com/AsimovNetwork/asimov
```

## Install Dependencies:
```sh
govendor sync
```

## Run the CLI:
```sh
./asimovd
```

## Run Your First Transaction:
```sh
```

## Config ```asimovd.conf``` in your working path.

Copy sample config file
```sh
cp asimovd.sample.conf asimovd.conf
```

Edit config file
```sh
; devnet,testnet,regtest,simnet
testnet=1
; solo,poa,region
consensustype=solo
privatekey=<your_privatekey>
; rpc port
rpclisten=:18334
; rpc username & password
rpcuser=<whatever_username_you_want>
rpcpass=<whatever_password_you_want>
```


## Config ```genesis.json``` in your working path.

Copy sample json file
```sh
cp genesis.sample.json genesis.json
```

Edit config file in testnet or devnet, use a same value in all nodes.
```sh
  "collectHeight":204,
  "chainStartTime":1576221900
```

## Generate your private key.

Build genkeys and generate your private key

```sh
cd cmd
go build
./wallet -c genKey -n test
  New key pair (priv, pubkey) (format:hex)
      { 0xf57dcd236c78368d64eb455b35d98c6d676394a1bb246f88fddd4d9103132729 , 0x020e470848be43cc2a9927af3ea9919e36a1e1b3319bc3845a31e1dcb5bcde79b8 }
      Compressed pubkey hash address: 0x664ceee539d3d5314824750b5f164cd8f8f87109b4
  Bitcoin base58check encode secret:  cVouSDdm3zVvTg5DuBqFaatghYQ1cApyEYBSdoSAyscmTCB5CK55
  Bitcoin legacy type address:  mnXjqNr44AbA5HbxbySBqf7b8Ho37Wqpv6
  Bitcoin p2sh-segwit type address:  2N1unojExgURPZ56sgwHetWVMsWqFPZD2wS
  Bitcoin bech32 type address:  tb1qfnhw2wwn65c5sfr4pd03vnxclru8zzd50s8fcy
```

## Toolchain

Clone and build

```sh
git clone https://github.com/AsimovNetwork/developer.git
npm install
npm run serve
```

Or visit online [Dev Center](https://developer.asimov.network) to use
the web IDE for asimov blockchain contract development.

## Cleanup

You can cleanup the state data when you want to reset all the states of
the test blockchain.

```sh
# stop asimovd
cd asimov
rm -rf data
rm -rf state
rm -rt logs
# restart asimov
```

## Asimov Testnet

Asimov Testnet is now available for developers.

Every one can connect the testnet with AsiLink and Web IDE.

# Wallet and Account

Using [AsiLink](https://www.asimov.network/wallet) Chrome Extension to
manager your wallet and account.

# Contract Template

Before deploy your contract,it is required to upload your contract
source code as template.

(image)

# Deploy Contract

Click Deploy in web IDE.

(image)

# Call Contract Functions

(image)

