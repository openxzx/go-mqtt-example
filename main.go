package main

import (
    "fmt"
    "os"
    "os/signal"
    "time"

    "github.com/eclipse/paho.mqtt.golang"
)

const (
    MQTT_URL = "tcp://localhost:1883"
    MQTT_TOPIC_UPLOAD = "topic/upload"
    MQTT_TOPIC_DOWNLOAD = "topic/download"
)

// Define message handler
func message_handler(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Upload Topic: %s, Payload: %s\n", msg.Topic(), msg.Payload())

    // Reply download message
    token := client.Publish(MQTT_TOPIC_DOWNLOAD, 0, true, "Reply download message.")
    token.Wait()
    // TODO: It can also realize the control on hardware.
}

func listen(client mqtt.Client) {
    if token := client.Subscribe(MQTT_TOPIC_UPLOAD, 1, message_handler); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    time.Sleep(100 * time.Millisecond)
}

func main() {
    // Create Mqtt client options
    opts := mqtt.NewClientOptions().AddBroker(MQTT_URL)
    client := mqtt.NewClient(opts)

    // Connect service
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    // Create signal recieiver for safe closing connection
    sig_chan := make(chan os.Signal, 1)
    signal.Notify(sig_chan, os.Interrupt)

    go listen(client)
    fmt.Printf("\033[1;32mMessages waiting...\033[0m\n")

    // Waiting interrupt signal
    <-sig_chan

    token := client.Publish(MQTT_TOPIC_DOWNLOAD, 0, true, "")
    token.Wait()

    fmt.Printf("\n\033[0;33mReceived interrupt signal and disconnecting...\033[0m\n")
    client.Disconnect(250)
}
