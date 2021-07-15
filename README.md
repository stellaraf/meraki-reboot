<div align="center">
  <br/>
  <img src="https://res.cloudinary.com/stellaraf/image/upload/v1604277355/stellar-logo-gradient.svg" width=300 />
  <br/>
  <h3>meraki-reboot</h3>
  <br/>
  <a href="https://github.com/stellaraf/meraki-reboot/actions?query=workflow%3Agoreleaser">
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/stellaraf/meraki-reboot/goreleaser?color=9100fa&style=for-the-badge">
  </a>
  <br/>
  
  Reboot multiple Cisco Meraki devices

</div>

## Usage

### Download the latest [release](https://github.com/stellaraf/meraki-reboot/releases/latest)

There are multiple builds of the release, for different CPU architectures/platforms:

There are multiple builds of the release, for different CPU architectures/platforms. Download and unpack the release for your platform:

```shell
wget <release url>
tar xvfz <release file> meraki-reboot
```

### Run the binary

```console
$ ./meraki-reboot --help

meraki-reboot 0.1.0
  Reboot shit tons of Meraki devices because Meraki is terrible

Options:

  -h, --help         display help information
  -o, --org          Meraki Organization Name
  -n, --network      Meraki Network Name
  -e, --exclusions   Comma-separated list of tags to exclude from the results

Commands:

  help         Display Help Information
  devices      List all matched devices
  reboot       Reboot one device
  reboot-all   Reboot all devices
  slack        Send a test Slack message

```

### Environment Variables

All of the below environment variables are required for meraki-reboot to run.

| Name                      | Default                      | Description                                                                                     |
| :------------------------ | :--------------------------- | :---------------------------------------------------------------------------------------------- |
| `MERAKI_API_KEY`          |                              | Meraki Dashboard API Key                                                                        |
| `MERAKI_REBOOT_LOG_FILE`  | `/var/log/meraki-reboot.log` | Path to a log file to which logs will be written.                                               |
| `MERAKI_REBOOT_SLACK_URL` |                              | Slack incoming webhook URL. Success or failure messages will be posted to slack after each run. |

## Creating a New Release

This project uses [GoReleaser](https://goreleaser.com/) to manage releases. After completing code changes and committing them via Git, be sure to tag the release before pushing:

```
git tag <release>
```

Once a new tag is pushed, GoReleaser will automagically create a new build & release.
