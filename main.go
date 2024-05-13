package main

import (
    "fmt"
    "strings"
    "os"
    "os/signal"
    "os/exec"
    "time"

    "github.com/eclipse/paho.mqtt.golang"
)

const (
    MQTT_URL = "tcp://localhost:1883"
    MQTT_TOPIC_UPLINK = "topic/uplink"
    MQTT_TOPIC_DOWNLINK = "topic/downlink"
)

const (
    EXEC_CMD = "uname"
    EXEC_CMD_ARGS = "-a"
)

func on_connect(client mqtt.Client) {
    fmt.Println("Conneted MQTT Service.")
    fmt.Printf("\033[1;32mMessages waiting...\033[0m\n")
}

// Define message handler
func message_handler(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Upload payload: %s\n", msg.Payload())

    // Reply download message
    out, err := exec.Command(EXEC_CMD, EXEC_CMD_ARGS).Output()
    if err != nil {
        token := client.Publish(MQTT_TOPIC_DOWNLINK, 0, true, err)
        token.Wait()
    } else {
	    token := client.Publish(MQTT_TOPIC_DOWNLINK, 0, true, strings.TrimSpace(string(out)))
        token.Wait()
    }
    // TODO: It can also realize the control on hardware.
}

func listen(client mqtt.Client) {
    if token := client.Subscribe(MQTT_TOPIC_UPLINK, 1, message_handler); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    time.Sleep(100 * time.Millisecond)
}

func main() {
    // Create Mqtt client options
    opts := mqtt.NewClientOptions().AddBroker(MQTT_URL)
    opts.OnConnect = on_connect

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    // Create signal recieiver for safe closing connection
    sig_chan := make(chan os.Signal, 1)
    signal.Notify(sig_chan, os.Interrupt)

    go listen(client)

    // Waiting interrupt signal
    <-sig_chan

    token := client.Publish(MQTT_TOPIC_DOWNLINK, 0, true, "")
    token.Wait()

    fmt.Printf("\n\033[1;33mReceived interrupt signal and disconnecting...\033[0m\n")
    client.Disconnect(250)
}
