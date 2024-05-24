# Go MQTT Example

Writing MQTT communication program via go, in order to be better familiar with go.

## Install dependencies

```
./install_dep.sh
```

## Run

```
./run.sh
```

## Send/Receive MQTT Information

### Send Command

If controling led for red, the following command:

```
mosquitto_pub -h localhost -t topic/uplink -m "red"
```

### Receive Command

Yes: Controlling led is normal; No: Controlling led is not normal.

```
mosquitto_sub -h localhost -t topic/downlink
```
