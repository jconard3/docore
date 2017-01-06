package main

//import "os"
import "fmt"
import "errors"

//import "strings"

import "github.com/digitalocean/godo"

func ListDroplets(client godo.Client) []string {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 25,
	}
	var dropletNames []string
	droplets, _, _ := client.Droplets.List(opt)
	for _, element := range droplets {
		dropletNames = append(dropletNames, element.Name)
	}
	return dropletNames
}

func DeleteDroplet(client godo.Client, id int) error {
	_, err := client.Droplets.Delete(id)
	return err
}

func RetrieveDroplet(client godo.Client, id int) (*godo.Droplet, error) {
	droplet, _, err := client.Droplets.Get(id)
	return droplet, err
}

func IDLookupDroplet(client godo.Client, name string) (int, error) {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 25,
	}
	droplets, _, _ := client.Droplets.List(opt)
	for _, element := range droplets {
		if element.Name == name {
			return element.ID, nil
		}
	}
	return 0, errors.New("No droplet matched the given name.\n")
}

func NameLookupDroplet(client godo.Client, id int) (string, error) {
	droplet, err := RetrieveDroplet(client, id)
	if err == nil {
		return droplet.Name, nil
	}
	return "", errors.New("No droplet matched the given ID.\n")
}

func main() {
	client, _ := CreateClient()
	names := ListDroplets(client)
	for _, name := range names {
		//if strings.Contains(name, "coreos") {
		//	ID, _ := IDLookupDroplet(client, name)
		//	_ = DeleteDroplet(client, ID)
		//}
		fmt.Println(name)
	}
}
