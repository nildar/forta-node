package cmd

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/fatih/color"
	"github.com/forta-network/forta-node/config"
	"github.com/spf13/cobra"
)

func handleFortaInit(cmd *cobra.Command, args []string) error {
	if isInitialized() {
		greenBold("Already initialized - please ensure that your configuration at %s is correct!\n", cfg.ConfigPath)
		return nil
	}

	if !isDirInitialized() {
		if err := os.Mkdir(cfg.FortaDir, 0755); err != nil {
			return err
		}
	}

	if !isConfigFileInitialized() {
		tmpl, err := template.New("config-template").Parse(defaultConfig)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, config.GetEnvDefaults(cfg.Development)); err != nil {
			return err
		}
		if err := os.WriteFile(cfg.ConfigPath, buf.Bytes(), 0644); err != nil {
			return err
		}
	}

	if !isKeyDirInitialized() {
		if err := os.Mkdir(cfg.KeyDirPath, 0755); err != nil {
			return err
		}
	}

	if !isKeyInitialized() {
		ks := keystore.NewKeyStore(cfg.KeyDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
		acct, err := ks.NewAccount(cfg.Passphrase)
		if err != nil {
			return err
		}
		printScannerAddress(acct.Address.Hex())
	}

	color.Green("\nSuccessfully initialized at %s\n", cfg.FortaDir)

	return nil
}

func printScannerAddress(address string) {
	fmt.Printf("\nScanner address: %s\n", color.New(color.FgYellow).Sprintf(address))
}

const defaultConfig = `# Auto generated by 'forta init' - safe to modify

registry:
  ethereum:
    websocketUrl: <fill in this and remove other>
    jsonRpcUrl: <fill in this and remove other>
  contractAddress: {{ .DefaultRegistryContractAddress }} # Auto-set to default
  containerRegistry: https://{{ .DiscoSubdomain }}.forta.network # Auto-set
  username: discouser
  password: discopass

scanner:
  scannerImage: forta-network/forta-scanner
  chainId: <fill in correctly>
  # startBlock: 123 <if you do not fill this in, it will start from the latest block>
  # endBlock: 123 <fill this in and enable if you really need>
  ethereum:
    websocketUrl: <fill in this and remove other>
    jsonRpcUrl: <fill in this and remove other>

jsonRpcProxy:
  jsonRpcImage: forta-network/forta-json-rpc
  ethereum:
    websocketUrl: <fill in this and remove other>
    jsonRpcUrl: <fill in this and remove other>

trace:
  enabled: false # Set this to true and set the Ethereum fields
  # ethereum:
  #   websocketUrl: <fill in this and remove other>
  #   jsonRpcUrl: <fill in this and remove other>

query:
  queryImage: forta-network/forta-query
  port: 8778
  publishTo:
    skipPublish: true # Make this false when your config is ready for publishing
    batch:
      skipEmpty: false
      intervalSeconds: 15
      maxAlerts: 100
    contractAddress: {{ .DefaultAlertContractAddress }} # Auto-set to default
    # ipfs:
    #   gatewayUrl: <set this to get ready for publishing>
    #   username: <set if needed>
    #   password: <set if needed>
    ethereum:
      websocketUrl: <fill in this and remove other>
      jsonRpcUrl: <fill in this and remove other>
    # testAlerts: # Configuration for logging alerts from test agents
    #   disable: false
    #   webhookUrl: '' # Does POST <webhookUrl> for each test alert if this field is non-empty 

log:
  level: info
  maxLogSize: 50m
  maxLogFiles: 10
`

func isDirInitialized() bool {
	info, err := os.Stat(cfg.FortaDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isConfigFileInitialized() bool {
	info, err := os.Stat(cfg.ConfigPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func isKeyDirInitialized() bool {
	info, err := os.Stat(cfg.KeyDirPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isKeyInitialized() bool {
	if !isKeyDirInitialized() {
		return false
	}
	entries, err := os.ReadDir(cfg.KeyDirPath)
	if err != nil {
		return false
	}
	for i, entry := range entries {
		if i > 0 {
			return false // There must be one key file
		}
		return !entry.IsDir() // so it should be a geth key file
	}
	return false // No keys found in dir
}

func isInitialized() bool {
	return isDirInitialized() && isConfigFileInitialized() && isKeyInitialized()
}
