# docore
Golang interface to CoreOS cluster hosted on Digital Ocean

## Config
Place config file at $HOME/.docore.yaml
NOTE: ssh_keys array will currently only accpet fingerprints of ssh keys. If you place a full public key in the array, it will fail to create a droplet

## Usage
```
 $ docore --help
Docore is a CLI to manage a CoreOS cluster hosted on DigitalOcean.
This application provides cluster-level management interfaces as well as fine-
grained droplet-level control.

Usage:
  docore [command]

Available Commands:
  droplet     Interface for DigitalOcean Droplets

Flags:
      --config string   config file (default is $HOME/.docore.yaml)
  -t, --toggle          Help message for toggle

Use "docore [command] --help" for more information about a command.
```

### Create a droplet
```
 $ docore droplet create --name joco
Using config file: /Users/jconard/.docore.yaml
godo.DropletCreateRequest{Name:"joco", Region:"nyc1", Size:"512mb", Image:godo.DropletCreateImage{ID:0, Slug:"coreos-stable"}, SSHKeys:[godo.DropletCreateSSHKey{ID:0, Fingerprint:"7a:e9:13:23:95:71:b2:21:dd:7c:e4:db:33:bc:b2:df"}], Backups:false, IPv6:false, PrivateNetworking:false, Monitoring:false, UserData:""}
Are you sure you want to create this droplet [y/n]: y
Droplet  joco  created. Currently provisioning...
```

### Get all information on a droplet
```
 $ docore droplet get joco
Using config file: /Users/jconard/.docore.yaml
godo.Droplet{ID:37649649, Name:"joco", Memory:512, Vcpus:1, Disk:20, Region:godo.Region{Slug:"nyc1", Name:"New York 1", Sizes:["512mb" "1gb" "2gb" "4gb" "8gb" "16gb" "m-16gb" "32gb" "m-32gb" "48gb" "m-64gb" "64gb" "m-128gb" "m-224gb"], Available:true, Features:["private_networking" "backups" "ipv6" "metadata" "install_agent" "storage"]}, Image:godo.Image{ID:22099398, Name:"1235.6.0 (stable)", Type:"snapshot", Distribution:"CoreOS", Slug:"coreos-stable", Public:true, Regions:["nyc1" "sfo1" "nyc2" "ams2" "sgp1" "lon1" "nyc3" "ams3" "fra1" "tor1" "sfo2" "blr1"], MinDiskSize:20, Created:"2017-01-11T03:06:26Z"}, Size:godo.Size{Slug:"512mb", Memory:512, Vcpus:1, Disk:20, PriceMonthly:5, PriceHourly:0.00744, Regions:["ams1" "ams2" "ams3" "blr1" "fra1" "lon1" "nyc1" "nyc2" "nyc3" "sfo1" "sfo2" "sgp1" "tor1"], Available:true, Transfer:1}, SizeSlug:"512mb", BackupIDs:[], SnapshotIDs:[], Locked:false, Status:"active", Networks:godo.Networks{V4:[godo.NetworkV4{IPAddress:"67.205.160.47", Netmask:"255.255.240.0", Gateway:"67.205.160.1", Type:"public"}], V6:[]}, Created:"2017-01-18T21:55:37Z", Tags:[], VolumeIDs:[]}
```
