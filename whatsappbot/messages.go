package main

import (
	"fmt"
	"os"

	whatsapp "github.com/Rhymen/go-whatsapp"
)

type MessageHandler struct {
	Client *WhatsAppClient
}

func (MessageHandler) HandleError(err error) {
	fmt.Fprintf(os.Stderr, "%v", err)
}

func (m MessageHandler) HandleTextMessage(message whatsapp.TextMessage) {
	if message.Text == "/status" {
		msgID, err := m.Client.Connection.Send(m.createWhatsappTextMsg("schisschas", message.Info.RemoteJid))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v", err)
		} else {
			fmt.Println("Message Sent -> ID : " + msgID)
		}
	}
	m.Client.Connection.Send(m.createWhatsappTextMsg(message.Text, message.Info.RemoteJid))

}

func (MessageHandler) createWhatsappTextMsg(text, remoteJid string) whatsapp.TextMessage {
	return whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: remoteJid,
		},
		Text: text,
	}
}
