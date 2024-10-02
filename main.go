package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    "os"
    "time"
    "sync"

    "github.com/bwmarrin/discordgo"
)

type QueuedMessage struct {
    Message   *discordgo.MessageCreate
    Session   *discordgo.Session
}

type MessageQueue struct {
    messages []QueuedMessage
    mu       sync.Mutex
}

func (q *MessageQueue) Enqueue(message QueuedMessage) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.messages = append(q.messages, message)
}

func (q *MessageQueue) Dequeue() (QueuedMessage, bool) {
    q.mu.Lock()
    defer q.mu.Unlock()

    if len(q.messages) == 0 {
        return QueuedMessage{}, false
    }

    message := q.messages[0]
    q.messages = q.messages[1:]
    return message, true
}

func (q *MessageQueue) ProcessMessages() {
    for {
        queuedMessage, ok := q.Dequeue()
        if !ok {
            time.Sleep(1 * time.Second) // No messages in queue, sleep for a while
            continue
        }

        err := sendMessageToAPI(queuedMessage.Session, queuedMessage.Message)
        if err != nil {
            log.Printf("Failed to send message to Kindroid API: %v", err)
            q.Enqueue(queuedMessage) // Requeue the message if failed
        }

        time.Sleep(5 * time.Second) // Try to keep from sending messages toooo quickly
    }
}

func sendMessageToAPI(s *discordgo.Session, m *discordgo.MessageCreate) error {
    // Ignore messages from the bot itself - this should be filtered out already but you never know
    if m.Author.ID == s.State.User.ID {
        return nil
    }

    // Check if the message mentions the bot
    for _, user := range m.Mentions {
        if user.ID == s.State.User.ID {
            kinToken := os.Getenv("KIN_TOKEN")
            if kinToken == "" {
                fmt.Println("No Kindroid token provided. Set KIN_TOKEN environment variable.")
                return nil
            }

            kinId := os.Getenv("KIN_ID")
            if kinId == "" {
                fmt.Println("No Kindroid AI ID provided. Set KIN_ID environment variable.")
                return nil
            }

            url := "https://api.kindroid.ai/v1/send-message"

            // Replacing mentions makes it so the Kin sees the usernames instead of <@userID> syntax
            updatedMessage, err := m.ContentWithMoreMentionsReplaced(s)
            if err != nil {
                log.Printf("Error replacing Discord mentions with usernames: %v", err)
            }

            // Prefix messages sent to the Kin so they know who they're from and that it's Discord
            // and not a normal Kin app message
            updatedMessage = "*Discord Message from " + m.Author.Username + ":* " + updatedMessage

            headers := map[string]string{
                "Authorization": "Bearer " + kinToken,
                "Content-Type": "application/json",
            }

            bodyMap := map[string]string{
                "message": updatedMessage,
                "ai_id": kinId,
            }
            jsonBody, err := json.Marshal(bodyMap)
            jsonString := string(jsonBody)
            fmt.Printf("Sending message to Kin API: %v", jsonString)

            req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
            if err != nil {
                log.Fatalf("Error reading HTTP request: %v", err)
            }

            req.Header.Set("Authorization", headers["Authorization"])
            req.Header.Set("Content-Type", headers["Content-Type"])

            client := &http.Client{}
            resp, err := client.Do(req)
            if err != nil {
                log.Fatalf("Error sending HTTP request: %v", err)
            }

            defer resp.Body.Close()

            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                log.Fatalf("Error reading HTTP response: %v", err)
            }


            if resp.StatusCode != http.StatusOK {
                log.Printf("Error response from Kin API: %v", string(body))
            }

            kinReply := string(body)
            log.Printf("Received reply message from Kin API, sending to Discord: %v", kinReply)
            // Send as a reply to the message that triggered the response, helps keep things orderly
            _, sendErr := s.ChannelMessageSendReply(m.ChannelID, kinReply, m.Reference())
            if sendErr != nil {
                fmt.Println("Error sending message: ", err)
            }
            return nil
        }
    }
    return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    // Ignore messages from the bot itself
    if m.Author.ID == s.State.User.ID {
        return
    }

    message := QueuedMessage{
        Message: m,
        Session: s,
    }

    queue.Enqueue(message)
}

var queue MessageQueue

func main() {
    botToken := os.Getenv("DISCORD_BOT_TOKEN")
    if botToken == "" {
        fmt.Println("No bot token provided. Set DISCORD_BOT_TOKEN environment variable.")
        return
    }

    dg, err := discordgo.New("Bot " + botToken)
    if err != nil {
        log.Fatalf("Error creating Discord session: %v", err)
    }

    dg.AddHandler(messageCreate)

    err = dg.Open()
    if err != nil {
        log.Fatalf("Error opening Discord connection: %v", err)
    }

    go queue.ProcessMessages()

    fmt.Println("Bot is now running. Press CTRL+C to exit.")
    select {}
}
