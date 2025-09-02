package main

import (
	ns "awesomeProject/notification_service"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func notification_service() {
	rand.Seed(time.Now().UnixNano())
	app := ns.NewApp()

	noti := &ns.Notification{
		ID:       "notif_" + strconv.Itoa(rand.Intn(1000000)),
		Title:    "Welcome to our service",
		Body:     "Hello {{name}}, thank you for joining us!",
		Channels: []ns.Channel{ns.ChannelEmail, ns.ChannelSMS},
		Data:     map[string]interface{}{"name": "John"},
		Recipients: []ns.Recipient{
			{ID: "user_1"},
			{ID: "user_2"},
		},
		Status: ns.StatusDraft,
	}

	id, status, err := app.CreateNotification(noti)
	if err != nil {
		fmt.Println("failed to create notification: %v", err)
	}
	fmt.Println("Notification created:%s status: %s", id, status)
	nots := app.ListNotifications()
	for _, n := range nots {
		fmt.Println("Notification created:%s status: %s", n.ID, n.Status)
	}

	// shutdown on SIGINT
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	log.Println("shutting down")
}
