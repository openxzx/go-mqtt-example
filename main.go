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
    MQTT_TOPIC_UPLINK = "topic/uplink"
    MQTT_TOPIC_DOWNLINK = "topic/downlink"
)

func onConnect(client mqtt.Client) {
    fmt.Println("Conneted MQTT Service.")
    fmt.Printf("\033[1;32mMessages waiting...\033[0m\n")
}

// Define message handler
func messageHandler(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Upload payload: %s\n", msg.Payload())

    // Ctrl Led and reply download message
    ret := CtrlRGB(string(msg.Payload()))
    if ret == 0 {
        token := client.Publish(MQTT_TOPIC_DOWNLINK, 0, true, "Yes")
        token.Wait()
    } else {
        token := client.Publish(MQTT_TOPIC_DOWNLINK, 0, true, "No")
        token.Wait()
    }
}

func listen(client mqtt.Client) {
    if token := client.Subscribe(MQTT_TOPIC_UPLINK, 1, messageHandler); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    time.Sleep(100 * time.Millisecond)
}

func main() {
    // Create Mqtt client options
    opts := mqtt.NewClientOptions().AddBroker(MQTT_URL)
    opts.OnConnect = onConnect

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    // Create pipe type signal recieiver for safe closing connection
    sig_chan := make(chan os.Signal, 1)
    signal.Notify(sig_chan, os.Interrupt)

    go listen(client)

    // Waiting interrupt signal
    <-sig_chan

    fmt.Printf("\n\033[1;33mReceived interrupt signal and disconnecting...\033[0m\n")
    client.Disconnect(300)
}
