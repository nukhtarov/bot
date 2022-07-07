package main

import (
	"github.com/SevereCloud/vksdk/api"
	"os"
)

func kickUser(id string) {
	token := os.Getenv("TOKEN_VK")
	vk := api.NewVK(token)
	vk.MessagesRemoveChatUser(api.Params{
		"chat_id": 1,
		"user_id": id,
	})
	return
}
