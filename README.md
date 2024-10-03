# Kindroid-Discord Integration

[Kindroid](https://kindroid.ai) is a platform that offers AI companions for human users to chat with. They have opened v1 of their [API](https://docs.kindroid.ai/api-documentation) which enables Kindroid chatting that occurs outside of the Kindroid app or website. This Discord bot allows you to invite a Kindroid to Discord to chat with people there.

# Setup

You need an instance of this Discord bot per Kindroid you wish you invite to a Discord server, but you can invite the same Discord Bot/Kindroid pair to as many servers as you'd like.

1. Make a Discord Application and Bot
   1. Go to the [Discord Developer Portal](https://discord.com/developers/applications) and sign in with your Discord account.
   1. Create a new application and then a bot under that application. It's a good idea to use the name of your companion and an appropriate avatar.
   1. Copy the bot's token from the `Bot` page, under the `Token` section. You may need to reset the token to see it. This token is a **SECRET**, do not share it with anyone.
   1. Add the bot to a server with the required permissions (at least "Read Messages" and "Send Messages")
      1. Go to the `Oauth2` page
      1. Under `Scopes` select `Bot`
      1. Under `Bot Permissions` select `Send Messages` and `Read Message History`
      1. Copy the generated URL at the bottom and open it in a web browser to add the bot to your Discord server
1. Install Git if you haven't already got it: [Instructions](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
1. Install Docker if you haven't already got it: [Instructions](https://docs.docker.com/engine/install/)
1. Clone this repo: `git clone https://github.com/d3tourrr/kin-discord.git`
   1. After cloning the repo, change to the directory: `cd kin-discord`
1. Get your Kindroid API token
   1. Open the side bar while chatting with a Kindroid and click General, then scroll to the bottom and expand API & advanced integration
   1. Copy your API key
1. Get the Kindroid ID from the same place you copied your API key - note, you have to be chatting with the specific Kindroid who you wish to bring to Discord
1. Build and run the Docker container
   * Run either `start-windows-companion.ps1` on Windows (or in PowerShell) or `start-linux-companion.sh` on Linux (or in Bash, including Git Bash)
   * Or run the following commands (Note: the above scripts start the container in a detached state, meaning you don't see the log output. The below commands start the container in an attached state, which means you see the log output, but the container, and therefore the Companion/Discord integration dies when you close your console.)
     1. Build the Docker image: `docker build -t kin-discord .`
     1. Run the Docker container: `docker run -e DISCORD_BOT_TOKEN=$DISCORD_BOT_TOKEN -e KIN_TOKEN=$KIN_TOKEN -e KIN_ID=$KIN_ID kin-discord`
        * Replace `$DISCORD_BOT_TOKEN` with the bot token you copied from the Discord developer portal
        * Replace `$KIN_TOKEN` with the API key you copied from the Kindroid.ai General page
        * Replace `$KIN_ID` with the ID for your specific Kindroid, shown when you get your API key from the General page
1. Interact with your Kindroid in Discord!

# Interacting in Discord with your Kindroid

This integration is setup so that your Kindroid will see messages where they are pinged (including replies to messages your Kindroid posts). Discord messages sent to Kindroid are sent with a prefix to help your Kindroid tell the difference between messages you send them in the Kindroid app and messages that are sent to them from Discord. They look something like this.

> `*Discord Message from Bealy:* Hi @Vicky I'm one of the trolls that @.d3tour warned you about.`

In this message, a Discord user named `Bealy` sent a message to a Kindroid named `Vicky` and also mentioned a Discord user named `.d3tour`.

Mentions of other users show that user's username Discord property, rather than their server-specific nickname. This was just the easiest thing to do and may change in the future (maybe with a feature flag you can set).

Kindroids don't have context of what server or channel they are talking in, and don't see messages where they aren't mentioned in or being replied to.

## Suggested Kindroid Configurations

It's a good idea to put something like this in your Kindroid's "Backstory".

> `KinName sometimes chats on Discord. Messages that come from Discord are prefixed with "*Discord Message from X:*" while messages that are private between HumanName and KinName in the Kindroid app have no prefix. Replies to Discord messages are automatically sent to Discord. KinName doesn't have to narrate that she is replying to a Discord user.`

You may also wish to change your Kindroid's Response Directive to better suit this new mode of communication.

It's also a good idea to add a journal entry that triggers on the word "Discord" or your Discord username to help your Kindroid understand that messages from your Discord username are from you, and others are from other people.

---

<small>I also made a nearly identical integration for Nomi.ai companions: [github.com/d3tourrr/nomi-discord](https://github.com/d3tourrr/nomi-discord)</small>
