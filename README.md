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

## The inner workings!

### Tagger Commands (from within slack)
 - one
 - two
   
### Tagger CLI Variables
 - one
 - two
   

### configs.json
 - one
 - two
   

### tags.json
 - one
 - two
   

### Work TODO
- [ ] - BUILD TESTS!  Ya, we got none and that's dumb
- [ ] - Add ability to lock down adds to a specific channel (same as deletes)
- [ ] - Add ability to open up adds/deletes to anyone (scarey)
- [ ] - Add ability to lock things down by user GUID
- [ ] - It would be awesome if you could have a tags.json file *per channel*  so tagger could do different reactions based on the channel its in.   Hawt.

To support tagger, please branch out and put in pull requests.   If you intend to do major amounts of work its possible to add you as a collaborator if that is of interest.

@tagger-v2 &copy; 2023-2024 srv1054 
