# Twitch-Raidban
Program to assist with quickly banning lists of bots on Twitch

## Building

`go install`
This command will build the program and output the binary as `twitch-raidban` to your `$GOBIN` directory. By default this will be `~/go/bin`.

## Usage

To run this program you will need to obtain an oauth token for your user. I typically make use of https://twitchapps.com/tmi/ which is documented as a recommended solution at https://dev.twitch.tv/docs/irc

```bash
twitch-raidban -channel mychannel -username mytwitchuser -token myouthtoken -file mybotslist.txt
```

| Argument Name  | Description |
| ------------- | ------------- |
| channel  | The twitch channel to apply bans to                             |
| username | The username for the twitch account to apply bans as            |
| token    | The oauth token used for authenticating with twitch as the user |
| file     | The path to a txt file with a bot name on each line             |

This command loads the bots file and will parse each line as a name, and then will connect to the given channel and issue `/ban` commands. It will show output similar to:

```bash
Size of list: 200
Processing entry 1/200
Processing entry 2/200
Processing entry 3/200
Processing entry 4/200
Processing entry 5/200
...
Processing entry 200/200
Disconnecting from channel
```

You can view progress by watching the output counter as well as the list in the channel's twitch chat. The program is not smart enough to check if users are already banned before running but will output any messages from Twitch listing already banned users.