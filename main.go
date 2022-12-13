package main

import (
	"fmt"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/vmihailenco/msgpack/v5"
)

const broker = "175.27.192.58"

const port = 1883

func NewMQTTClient() mqtt.Client {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("admin")
	opts.SetPassword("public")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}

type GPS struct {
	Lat float32
	Lon float32
}

func Sub(client mqtt.Client) {
	topic := "event/gps/+/up"
	token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		var gps GPS
		msgpack.Unmarshal(msg.Payload(), &gps)
		fmt.Printf("%d - Received message: %+v from topic: %s\n", time.Now().Unix(), gps, msg.Topic())
	})
	token.Wait()
}

func main() {
	client := NewMQTTClient()
	Sub(client)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to my website!")
	})

	http.ListenAndServe(":8080", nil)

}
