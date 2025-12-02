# Tagger Slack Bot v2

This is a replacement for the old [bots-tagger](https://github.com/srv1054/bots-tagger) bot that was based on the original slack RTM.   v2 now leverages Slack API via websockets.  The new bot also includes a number of new features allowing for the ability to add key words and spray cans (emoji tags) to the configuration file via the slack bot from within slack.    Currently you can control the ability of who can delete items based on locking deletes to a specific channel GUID specified in the `config.json`

For example if there is a Spray Can called "business-cat" (which will match the exact name of a slack emoji on the server (without colons)) and inside that spray can (in the tags.json) one of the trigger words is "synergy".   Any time the word synergy is used in a post in a slack channel tagger is a member of, it will stick the :business-cat: emoji as a reaction to that post.   Hilarity can insue when properly setup.

Be very careful about keywords.  Its easy to over use and make things annoying.   Use some subtelty that will make a good time for the users of the server.

To get this running you can:
1. Pull the binary for your OS from the latest Release down to your machine
2. Pull the config.json and tags.json example files down to the same directory
3. Edit them accordingly and fire it up.
   
   *OR*
   
1. git clone -v \<repo url\>
2. go mod tidy
3. go build
4. Put the tagger binary and the `config.json` and `tags.json` files in the same location and edit both JSON accordingly
5. Fire it up

If you are running this bot in *NIX operating systems, its best to run it inside `screen` or once you are happy its working built it into a systemctl service to run at startup

### Running as a service
 - sudo mkdir -p /opt/tagger  
 - sudo cp /path/to/tagger /opt/tagger/  
 - sudo cp /path/to/config.json /opt/tagger  

Create `/etc/systemd/system/tagger.service`

```[Unit]
Description=Tagger Slack Bot
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/opt/tagger
ExecStart=/opt/tagger/tagger -c /opt/tagger/config.json -j /opt/tagger/tags.json
Restart=on-failure
RestartSec=5s
# Optional: limit resources / file handles
LimitNOFILE=65536

# The app already writes daily logs (tagger-YYYY-MM-DD.log) to WorkingDirectory.
# Stdout/stderr also go to journald:
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target```

- sudo systemctl daemon-reload
- sudo systemctl enable --now tagger
- sudo systemctl status tagger --no-pager
- # Follow logs
- journalctl -u tagger -f

## The inner workings!

### Tagger Commands (from within slack)
You can get a list of these commands from within slack by asking `tagger` for help: `@tagger` help\
`@tagger` list spray cans - List all availabe tags\
`@tagger` add spray can - Add a tag\
`@tagger` delete spray can - Delete a tag\
`@tagger` reload spray cans - Reload tags.json\
`@tagger` add spray can word - Add keyword to a spray can (tag)\
&nbsp;&nbsp; &nbsp; You must specify an existing Spray Can\
&nbsp;&nbsp; &nbsp; Word must be in (" ") quotation marks to allow for spaces.\
`@tagger` delete word - Delete a spray can (tag)

Specifics for Adding Words to Spray Cans:\
&nbsp; &nbsp; `@tagger` add word \<spray can> "\<new word>"\
&nbsp;&nbsp; &nbsp; &nbsp;&nbsp; &nbsp; e.g.: `@taggerbot add word smile "happyness"`\
&nbsp;&nbsp; &nbsp; &nbsp;&nbsp; &nbsp; The \<spray can> must exist as a real slack emoji.

Specifics for Adding new Spray Cans:\
&nbsp;&nbsp; `@tagger` add spray can <emoji name (no colons)>\
&nbsp;&nbsp; &nbsp; &nbsp;&nbsp; &nbsp; e.g.: @taggerbot add spray can catwave\
&nbsp;&nbsp; &nbsp; &nbsp;&nbsp; &nbsp; The <emoji name> must exist as a real slack emoji.
   
### Tagger CLI Variables
`-h` - Command line help\
`-v` - Show current version and exit\
`-cp` - Path to configuration file\
`-jp` - Path to SprayCan JSON file

Paths should be in quote if they contain spaces.\
Path specifications should include filename   

### configs.json
	"slackhook": "",  - Slack App Webhook
	"slackapptoken": "", - Slack App APP TOKEN
	"slackbottoken": "", - Slack App BOT TOKEN
	"botid": "", - Currently un-used
	"botname": "", - The name your bot will display as when it posts to slack
	"teamid": "", - Currently un-used
	"teamname": "", - Currently un-used
	"logchannel": "", - Slack channel for tagger to send log messages to
	"sprayjsonpath": "", - Currently un-used
	"debug": false, - Turn on excessive logging information
	"allowdeletefrom": "" - List of Slack UIDs that can delete JSON entries from within Slack app
   

### tags.json
See [/configs/tas.json.example](https://github.com/srv1054/tagger-v2/tree/main/configs) for a starter file and how to use the format.   You can edit this file manually or via the slack app interface.  Large edits are easier to do manually.   Don't forget to validate your JSON formatting, and after editing you must either restart `tagger` or from slack do an `@tagger reload spray cans`

### /emoji directory
There are two graphics in the source [/emoji](https://github.com/srv1054/tagger-v2/tree/main/emoji) directory that should be used when configuring slack.   The 512x512 is setup for the Slack App Bot configuration which requires a graphic of 512x512 or 256x256 square.   The smaller emoji graphic should be added into your slack server as a standard emoji called :tagger:

### Work TODO
- [ ] - BUILD TESTS!  Ya, we got none and that's dumb
- [ ] - Add ability to lock down adds to a specific channel (same as deletes)
- [ ] - Add ability to open up adds/deletes to anyone (scarey)
- [ ] - Add ability to lock things down by user GUID
- [ ] - It would be awesome if you could have a tags.json file *per channel*  so tagger could do different reactions based on the channel its in.   Hawt.

To support tagger, please branch out and put in pull requests.   If you intend to do major amounts of work its possible to add you as a collaborator if that is of interest.

@tagger-v2 &copy; 2023-2024 srv1054 
