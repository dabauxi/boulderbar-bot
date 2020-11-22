package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"os/signal"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
)

type WhatsAppClient struct {
	Connection *whatsapp.Conn
}

func main() {
	var whatsappClient = NewClient()

	Login(whatsappClient)
	whatsappClient.Connection.AddHandler(MessageHandler{Client: whatsappClient})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	fmt.Println("Press ctrl+c to exit.")

	<-sigs
	fmt.Println("Shutdown.")
	os.Exit(0)

}

func NewClient() *WhatsAppClient {
	wac, err := whatsapp.NewConn(20 * time.Second)
	wac.SetClientName("Boulderbar Whatsapp Bot", "Boulderbar Whatsapp Bot", "0.1.0")
	if err != nil {
		os.Exit(1)
		return nil
	}

	var client = WhatsAppClient{Connection: wac}
	return &client
}

func Login(client *WhatsAppClient) error {
	session, err := loadSession()
	if err == nil {
		session, err = client.Connection.RestoreWithSession(session)
		if err != nil {
			return err
		}
	} else {
		session, err = login(client)
		if err != nil {
			return err
		}
		err = saveSession(session)
		if err != nil {
			return err
		}
	}

	return nil
}

func login(client *WhatsAppClient) (whatsapp.Session, error) {
	qr := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New()
		terminal.Get(<-qr).Print()
	}()
	return client.Connection.Login(qr)
}

func loadSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "session.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func saveSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "session.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}
