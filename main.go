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

type Event struct {
	Rate int
}

func Sub(client mqtt.Client) {
	topic := "event/gps/+/up"
	token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		var gps GPS
		msgpack.Unmarshal(msg.Payload(), &gps)
		fmt.Printf(
			"%d - Received message: %+v\n",
			time.Now().Unix(),
			gps,
		)
	})
	token.Wait()
}

func main() {
	client := NewMQTTClient()
	Sub(client)

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		sn := r.FormValue("sn")
		rate := r.FormValue("rate")

		client.Publish("event/gps/"+sn+"/down", 2, true, rate)
		fmt.Println("down event to " + sn + ", rate is " + rate)
	})

	http.ListenAndServe(":8080", nil)
}
