![build](https://travis-ci.org/alexshemesh/claptrap.svg?branch=master)
# claptrap
General fintech bot

# Command line interface
```
claptrap main.go --help
Fintech aytomation toolkit, can be used as CLI and as a server

Usage:
  claptrap [command]

Available Commands:
  help        Help about any command
  kuna        Set of commands ot manage Kuna market
  vault       Set of commands to manage vault

Flags:
      --config string   config file (default is $HOME/.claptrap.yaml)
  -h, --help            help for claptrap
  -t, --toggle          Help message for toggle
```

```
claptrap kuna --help
https://kuna.io is Ukrainian cryptocurency market. This is root command to manage it

Usage:
  claptrap kuna [flags]
  claptrap kuna [command]

Available Commands:
  ordersbook  Retruns book of orders frm Kuna market

Flags:
  -h, --help   help for kuna

Global Flags:
      --config string   config file (default is $HOME/.claptrap.yaml)
```

```
claptrap main.go vault --help

Hashicorp Vault is used for storing sensitive information
	account details, configurations etc.
	This set of commands allows to initialize new vault, change and monitor values in vault

Usage:
  claptrap vault [flags]
  claptrap vault [command]

Available Commands:
  getkey      A brief description of your command
  init        Init vault
  setkey      Saves secret to vault

Flags:
  -h, --help               help for vault
      --vaultAddr string   address of vault server (default "http://127.0.0.1")

Global Flags:
      --config string   config file (default is $HOME/.claptrap.yaml)
```
