package main

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	token := os.Getenv("TOKEN_VK")
	vk := api.NewVK(token)

	// get information about the group
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initializing Long Poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	// New message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)
		b := params.NewMessagesSendBuilder()
		command := strings.Split(obj.Message.Text, " ")
		b.PeerID(obj.Message.PeerID)
		b.RandomID(0)
		chatId := obj.Message.PeerID - 2000000000

		switch command[0] {
		case "изменить":
			if isAdmin(obj.Message.FromID) {
				if len(command) != 1 {
					newName := command[1]
					vk.MessagesEditChat(api.Params{
						"chat_id": chatId,
						"title":   newName,
					})
					b.Message(fmt.Sprintf("Изменено название чата на: %v", newName))
					_, err = vk.MessagesSend(b.Params)
					if err != nil {
						log.Println(err)
					}

				}
			} else {
				b.Message(`Вы не админ! Напишите админу @ri`)
				b.RandomID(0)
				b.PeerID(obj.Message.PeerID)

				_, err = vk.MessagesSend(b.Params)
				if err != nil {
					log.Println(err)
				}
			}
		case "создать":
			if isAdmin(obj.Message.FromID) {
				if len(command) != 0 {
					//id := command[1]
					//intVar, err := strconv.Atoi(id)
					//b.UserID(intVar)
					b.PeerID(obj.Message.PeerID)
					vkChatID, err := vk.MessagesCreateChat(api.Params{
						"title": command[1],
					})
					link, _ := vk.MessagesGetInviteLink(api.Params{
						"peer_id": 2000000000 + vkChatID,
					})
					b.Message(fmt.Sprintf("Чат:`%v создан \n Ссылка: %v", command[1], link.Link))

					_, err = vk.MessagesSend(b.Params)
					if err != nil {
						log.Fatal(err)
					}

				}
			}
		case "инфо":
			if isAdmin(obj.Message.FromID) {
				if len(command) != 1 {
					if _, err := strconv.Atoi(command[1]); err == nil {
						id := command[1]
						users, err := vk.UsersGet(api.Params{
							"user_ids": id,
							"fields": []string{
								"screen_name"}})

						log.Println(users[0].Status)
						if err != nil {
							log.Println(err)

						}

						b.Message(fmt.Sprintf("Пользователь %v %v \n ID: %v \n Короткое имя: %v", users[0].FirstName, users[0].LastName, users[0].ID, users[0].ScreenName))
						b.DisableMentions(true)
						_, err = vk.MessagesSend(b.Params)
						if err != nil {
							log.Println(err)

						}
					} else {
						b.Message(`ID пользователя может быть только числовым`)
						b.RandomID(0)
						b.PeerID(obj.Message.PeerID)

						_, err = vk.MessagesSend(b.Params)
						if err != nil {
							log.Println(err)

						}
					}

				} else {
					b.Message(`Не хватает ID пользователя`)
					b.RandomID(0)
					b.PeerID(obj.Message.PeerID)

					_, err = vk.MessagesSend(b.Params)
					if err != nil {
						log.Println(err)

					}
				}

			} else {
				b.Message(`Вы не админ! Напишите админу @ri`)
				b.RandomID(0)
				b.PeerID(obj.Message.PeerID)

				_, err = vk.MessagesSend(b.Params)
				if err != nil {
					log.Println(err)
				}
			}

		default:

			b.RandomID(0)
			b.Message("hello")
			b.PeerID(obj.Message.PeerID)
			_, err = vk.MessagesSend(b.Params)

			if err != nil {
				log.Println(err)
			}
		}

	})

	// Run Bots Long Poll
	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
