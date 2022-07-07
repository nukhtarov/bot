package main

import (
	"context"
	"github.com/SevereCloud/vksdk/object"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"os"
	"strconv"
	"strings"
)

var admins = []int{237286647, 162667568}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

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
		switch command[0] {
		case "изменить":
			if contains(admins, obj.Message.FromID) {
				if len(command) != 1 {
					name := command[1]
					b.PeerID(obj.Message.PeerID)
					b.RandomID(0)
					vk.MessagesEditChat(api.Params{
						"chat_id": 1,
						"title":   name,
					})
					b.Message("изменил")
					_, err = vk.MessagesSend(b.Params)
					if err != nil {
						log.Fatal(err)
					}

				}
			}
		case "кик":
			if contains(admins, obj.Message.FromID) {
				if len(command) != 1 {
					id := command[1]
					b.PeerID(obj.Message.PeerID)
					b.RandomID(0)
					kickUser(id)
					b.Message("kicked")
					_, err = vk.MessagesSend(b.Params)
					if err != nil {
						log.Fatal(err)
					}

				}
			}

		case "создать":
			if contains(admins, obj.Message.FromID) {
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
					log.Println(link)
					b.RandomID(0)
					b.Message(link.Link)

					_, err = vk.MessagesSend(b.Params)
					if err != nil {
						log.Fatal(err)
					}

				}
			}
		case "бан":
			if contains(admins, obj.Message.FromID) {
				if len(command) != 1 {
					if _, err := strconv.Atoi(command[1]); err == nil {
						id := command[1]
						users, err := vk.UsersGet(api.Params{
							"user_ids": id,
						})
						log.Println(users)
						if err != nil {
							log.Fatal(err)

						}
						b.Message(`@id` + id + ` (Пользователь) с id: ` + id + ` заблокирован`)

						b.RandomID(0)
						b.PeerID(obj.Message.PeerID)

						_, err = vk.MessagesSend(b.Params)
					} else {
						b.Message(`ID пользователя может быть только числовым`)
						b.RandomID(0)
						b.PeerID(obj.Message.PeerID)

						_, err = vk.MessagesSend(b.Params)
						if err != nil {
							log.Fatal(err)

						}
					}

				} else {
					b.Message(`Не хватает ID пользователя`)
					b.RandomID(0)
					b.PeerID(obj.Message.PeerID)

					_, err = vk.MessagesSend(b.Params)
					if err != nil {
						log.Fatal(err)

					}
				}

			} else {
				b.Message(`Вы не админ! Напишите админу @ri`)
				b.RandomID(0)
				b.PeerID(obj.Message.PeerID)

				_, err = vk.MessagesSend(b.Params)
				if err != nil {
					log.Fatal(err)
				}
			}

		default:

			b.RandomID(0)
			b.Message("hello")
			b.PeerID(obj.Message.PeerID)
			k := object.NewMessagesKeyboardInline()
			k.AddRow()
			k.AddOpenLinkButton(`https://vk.com`, `привет`, ``)
			k.AddTextButton(`label`, ``, `primary`)
			k.AddRow()
			k.AddTextButton(`label`, ``, `primary`)
			k.AddTextButton(`label`, ``, `primary`)
			k.AddTextButton(`label`, ``, `primary`)
			k.AddTextButton(`label`, ``, `primary`)
			k.AddTextButton(`label`, ``, `primary`)
			b.Keyboard(k)
			_, err = vk.MessagesSend(b.Params)

			if err != nil {
				log.Fatal(err)
			}
		}

	})

	// Run Bots Long Poll
	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
