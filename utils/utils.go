package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/spf13/viper"
)

func NameToID(client godo.Client, droplet_name string) (int, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 25,
	}
	droplets, _, err := client.Droplets.List(opt)
	if err != nil {
		return 0, err
	}

	for _, element := range droplets {
		if element.Name == droplet_name {
			return element.ID, nil
		}
	}
	return 0, errors.New("No droplet matches for given name")
}

//https://gist.github.com/m4ng0squ4sh/3dcbb0c8f6cfe9c66ab8008f55f8f28b
func AskForConfirmation(s string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/N]: ", s)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func ViperGetSSHKeys() ([]godo.DropletCreateSSHKey, error) {
	if !viper.IsSet("ssh_keys") {
		return nil, errors.New("No SSH Keys array specified in config file under 'ssh_keys'")
	}

	ssh_keys := viper.GetStringSlice("ssh_keys")
	droplet_keys := make([]godo.DropletCreateSSHKey, len(ssh_keys))
	for i, ssh_key := range ssh_keys {
		droplet_keys[i].Fingerprint = ssh_key
	}
	return droplet_keys, nil
}
