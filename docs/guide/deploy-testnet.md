# Deploy a testnet

This document describes 3 ways to setup a network of `terrad` nodes, each serving a different usecase:

1. Single-node, local, manual testnet
2. Multi-node, local, automated testnet
3. Multi-node, remote, automated testnet

Supporting code can be found in the [networks directory](https://github.com/terra-project/core/tree/develop/networks) and additionally the `local` or `remote` sub-directories.

> NOTE: The `remote` network bootstrapping may be out of sync with the latest releases and is not to be relied upon.

## Single-node, local, manual testnet

This guide helps you create a single validator node that runs a network locally for testing and other development related uses.

### Requirements

* [Install Terra](installation.md)
* [Install `jq`](https://stedolan.github.io/jq/download/) \(optional\)

### Create genesis file and start the network

```bash
# You can run all of these commands from your home directory
cd $HOME

# Initialize the genesis.json file that will help you to bootstrap the network
terrad init --chain-id=testing testing

# Create a key to hold your validator account
terracli keys add validator

# Add that key into the genesis.app_state.accounts array in the genesis file
# NOTE: this command lets you set the number of coins. Make sure this account has some coins
# with the genesis.app_state.staking.params.bond_denom denom, the default is staking
terrad add-genesis-account $(terracli keys show validator -a) 1000uluna,1000usdr

# Generate the transaction that creates your validator
terrad gentx --name validator --amount 100uluna

# Add the generated bonding transaction to the genesis file
terrad collect-gentxs

# Now its safe to start `terrad`
terrad start
```

This setup puts all the data for `terrad` in `~/.terrad`. You can examine the genesis file you created at `~/.terrad/config/genesis.json`. With this configuration `terracli` is also ready to use and has an account with tokens \(both staking and custom\).

## Multi-node, local, automated testnet

From the [networks/local directory](https://github.com/terra-project/core/tree/develop/networks/local):

### Requirements

* [Install terra](installation.md)
* [Install docker](https://docs.docker.com/engine/installation/)
* [Install docker-compose](https://docs.docker.com/compose/install/)

### Build

Build the `terrad` binary \(linux\) and the `tendermint/terradnode` docker image required for running the `localnet` commands. This binary will be mounted into the container and can be updated rebuilding the image, so you only need to build the image once. Go to the main project directory for Terra Core, and run:

```bash
# Build the linux binary in ./build
make build-linux

# Build tendermint/terradnode image
make build-docker-terradnode
```

### Run your testnet

To start a 4 node testnet run:

```text
make localnet-start
```

This command creates a 4-node network using the terradnode image. The ports for each node are found in this table:

| Node ID | P2P Port | RPC Port |
| :--- | :--- | :--- |
| `terranode0` | `26656` | `26657` |
| `terranode1` | `26659` | `26660` |
| `terranode2` | `26661` | `26662` |
| `terranode3` | `26663` | `26664` |

To update the binary, just rebuild it and restart the nodes:

```text
make build-linux build-docker-terradnode localnet-start
```

### Configuration

The `make localnet-start` creates files for a 4-node testnet in `./build` by calling the `terrad testnet` command. This outputs a handful of files in the `./build` directory:

```bash
$ tree -L 2 build/
build/
├── terracli
├── terrad
├── gentxs
│   ├── node0.json
│   ├── node1.json
│   ├── node2.json
│   └── node3.json
├── node0
│   ├── terracli
│   │   ├── key_seed.json
│   │   └── keys
│   └── terrad
│       ├── ${LOG:-terrad.log}
│       ├── config
│       └── data
├── node1
│   ├── terracli
│   │   └── key_seed.json
│   └── terrad
│       ├── ${LOG:-terrad.log}
│       ├── config
│       └── data
├── node2
│   ├── terracli
│   │   └── key_seed.json
│   └── terrad
│       ├── ${LOG:-terrad.log}
│       ├── config
│       └── data
└── node3
    ├── terracli
    │   └── key_seed.json
    └── terrad
        ├── ${LOG:-terrad.log}
        ├── config
        └── data
```

Each `./build/nodeN` directory is mounted to the `/terrad` directory in each container.

### Logging

Logs are saved under each `./build/nodeN/terrad/terra.log`. You can also watch logs directly via Docker, for example:

```text
docker logs -f terradnode0
```

### Keys & Accounts

To interact with `terracli` and start querying state or creating txs, you use the `terracli` directory of any given node as your `home`, for example:

```text
terracli keys list --home ./build/node0/terracli
```

Now that accounts exists, you may create new accounts and send those accounts funds!

{% hint style="info" %}
    Each node's seed is located at `./build/nodeN/terracli/key_seed.json` and can be restored to the CLI using the `terracli keys add --restore` command.
{% endhint %}

### Special binaries

If you have multiple binaries with different names, you can specify which one to run with the BINARY environment variable. The path of the binary is relative to the attached volume. For example:

```text
# Run with custom binary
BINARY=terrafoo make localnet-start
```

## Multi-node, remote, automated testnet

The following should be run from the [networks directory](https://github.com/terra-project/core/tree/develop/networks).

### Terraform & Ansible

Automated deployments are done using [Terraform](https://www.terraform.io/) to create servers on AWS then [Ansible](http://www.ansible.com/) to create and manage testnets on those servers.

### Prerequisites

* Install [Terraform](https://www.terraform.io/downloads.html) and [Ansible](http://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html) on a Linux machine.
* Create an [AWS API token](https://docs.aws.amazon.com/general/latest/gr/managing-aws-access-keys.html) with EC2 create capability.
* Create SSH keys

```text
export AWS_ACCESS_KEY_ID="2345234jk2lh4234"
export AWS_SECRET_ACCESS_KEY="234jhkg234h52kh4g5khg34"
export TESTNET_NAME="remotenet"
export CLUSTER_NAME= "remotenetvalidators"
export SSH_PRIVATE_FILE="$HOME/.ssh/id_rsa"
export SSH_PUBLIC_FILE="$HOME/.ssh/id_rsa.pub"
```

These will be used by both `terraform` and `ansible`.

### Create a remote network

```text
SERVERS=1 REGION_LIMIT=1 make validators-start
```

The testnet name is what's going to be used in --chain-id, while the cluster name is the administrative tag in AWS for the servers. The code will create SERVERS amount of servers in each availability zone up to the number of REGION\_LIMITs, starting at us-east-2. \(us-east-1 is excluded.\) The below BaSH script does the same, but sometimes it's more comfortable for input.

```text
./new-testnet.sh "$TESTNET_NAME" "$CLUSTER_NAME" 1 1
```

### Quickly see the /status endpoint

```text
make validators-status
```

### Delete servers

```text
make validators-stop
```

### Logging

You can ship logs to Logz.io, an Elastic stack \(Elastic search, Logstash and Kibana\) service provider. You can set up your nodes to log there automatically. Create an account and get your API key from the notes on [this page](https://app.logz.io/#/dashboard/data-sources/Filebeat), then:

```text
yum install systemd-devel || echo "This will only work on RHEL-based systems."
apt-get install libsystemd-dev || echo "This will only work on Debian-based systems."

go get github.com/mheese/journalbeat
ansible-playbook -i inventory/digital_ocean.py -l remotenet logzio.yml -e LOGZIO_TOKEN=ABCDEFGHIJKLMNOPQRSTUVWXYZ012345
```

### Monitoring

You can install the DataDog agent with:

```text
make datadog-install
```

