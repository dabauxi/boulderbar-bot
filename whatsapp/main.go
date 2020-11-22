package whatsapp

import (
	"encoding/gob"
	"os"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
)

type WhatsAppClient struct {
	Connection *whatsapp.Conn
}

func InitClient() *WhatsAppClient {
	wac, err := whatsapp.NewConn(20 * time.Second)

	if err != nil {
		panic(err)
	}

	var client = WhatsAppClient{Connection: wac}
	return &client
}

func Login(client *WhatsAppClient) error {
	session, err := loadSession()
	if err != nil {
		session, err = login(client)
		err = saveSession(session)
	}
	if err != nil {
		return err
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
