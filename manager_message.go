package mtproto

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
)

type Message struct {
	Flags         MessageFlags          `json:"flags"`
	Type          string                `json:"type"`
	ID            int32                 `json:"id"`
	From          int32                 `json:"from"`
	To            *Peer                 `json:"to"`
	Date          int32                 `json:"date"`
	Body          string                `json:"body"`
	MediaType     string                `json:"media_type"`
	Action        *MessageAction        `json:"action"`
	ForwardHeader *MessageForwardHeader `json:"forward_header"`
	Entities      []MessageEntity       `json:"entities"`
	Views         int32                 `json:"views"`
	Media         interface{}           `json:"media"`
}
type MessageFlags struct {
	Out         bool `json:"out"`           // flags_1?true
	Mentioned   bool `json:"mentioned"`     // flags_4?true
	MediaUnread bool `json:"media_unheard"` // flags_5?true
	Silent      bool `json:"silent"`        // flags_13?true
	Post        bool `json:"post"`          // flags_14?true
}

func (f *MessageFlags) loadFlags(flags int32) {
	if flags&1<<1 != 0 {
		f.Out = true
	}
	if flags&1<<4 != 0 {
		f.Mentioned = true
	}
	if flags&1<<5 != 0 {
		f.MediaUnread = true
	}
	if flags&1<<13 != 0 {
		f.Silent = true
	}
	if flags&1<<14 != 0 {
		f.Post = true
	}
}

type MessageAction struct {
	Type      string  `json:"type"`
	Title     string  `json:"title"`
	ChatID    int32   `json:"chat_id"`
	ChannelID int32   `json:"channel_id"`
	GameID    int64   `json:"game_id"`
	GameScore int32   `json:"game_score"`
	UserID    int32   `json:"user_id"`
	UserIDs   []int32 `json:"user_ids"`
	Photo     *Photo  `json:"photo"`
}
type MessageEntity struct {
	Type     string `json:"type"`
	Offset   int32  `json:"offset"`
	Length   int32  `json:"length"`
	Url      string `json:"url"`
	language string `json:"language"`
	UserID   int32  `json:"user_id"`
}
type MessageForwardHeader struct {
	From        int32  `json:"from"`
	Date        int32  `json:"date"`
	ChannelID   int32  `json:"channel_id"`
	ChannelPost int32  `json:"channel_post"`
	Author      string `json:"author"`
}
type MessageMedia interface{}
type MessageMediaPhoto struct {
	Caption string `json:"caption"`
	Photo   Photo  `json:"photo"`
}
type MessageMediaContact struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	UserID    int32  `json:"user_id"`
	Phone     string `json:"phone"`
}
type MessageMediaDocument struct {
	Caption  string   `json:"caption"`
	Document Document `json:"document"`
}
type MessageReplyMarkup struct {
}

// NewMessage
// input
//	1. TL_message
//	2. TL_messageService
func NewMessage(input TL) (m *Message) {
	m = new(Message)
	switch x := input.(type) {
	case TL_messageEmpty:
		return nil
	case TL_message:
		m.Flags.loadFlags(x.Flags)
		m.Type = MESSAGE_TYPE_NORMAL
		m.ID = x.Id
		m.Date = x.Date
		m.From = x.From_id
		m.Body = x.Message
		m.To = NewPeer(x.To_id)
		m.Views = x.Views
		if x.Media != nil {
			m.Media = NewMessageMedia(x.Media)
		}
		if x.Fwd_from != nil {
			m.ForwardHeader = NewMessageForwardHeader(x.Fwd_from)
		}
		m.Entities = make([]MessageEntity, 0, 0)
		for _, e := range x.Entities {
			m.Entities = append(m.Entities, *NewMessageEntity(e))
		}
	case TL_messageService:
		m.Flags.loadFlags(x.Flags)
		m.Type = MESSAGE_TYPE_SERVICE
		m.ID = x.Id
		m.Date = x.Date
		m.From = x.From_id
		m.To = NewPeer(x.To_id)
		m.Action = NewMessageAction(x.Action)
		m.ForwardHeader = new(MessageForwardHeader)
	default:
		fmt.Println("NewMessage::UnSupported Input Format", reflect.TypeOf(x).String())
		return nil
	}
	return
}

// NewMessageAction
// input:
//	1. TL_messageActionEmpty
//	2. TL_messageActionChannelCreate
//	3. TL_messageActionChannelMigrateFrom
//	4. TL_messageActionChatCreate
//	5. TL_messageActionChatAddUser
//	6. TL_messageActionChatDeleteUser
//	7. TL_messageActionChatDeleteUser
//	8. TL_messageActionChatEditPhoto
//	9. TL_messageActionChatEditTitle
//	10. TL_messageActionChatJoinedByLink
//	11.	TL_messageActionChatMigrateTo
//	12.	TL_messageActionGameScore
//	13. TL_messageActionHistoryClear
//	14. TL_messageActionPinMessage
//	15. TL_messageActionPhoneCall
func NewMessageAction(input TL) (m *MessageAction) {
	m = new(MessageAction)
	switch x := input.(type) {
	case TL_messageActionEmpty:
	case TL_messageActionChannelCreate:
		m.Type = MESSAGE_ACTION_CHANNEL_CREATED
		m.Title = x.Title
	case TL_messageActionChannelMigrateFrom:
		m.Type = MESSAGE_ACTION_CHANNEL_MIGRATE_FROM
		m.Title = x.Title
		m.ChatID = x.Chat_id
	case TL_messageActionChatCreate:
		m.Type = MESSAGE_ACTION_CHAT_CREATED
		m.Title = x.Title
		m.UserIDs = x.Users
	case TL_messageActionChatAddUser:
		m.Type = MESSAGE_ACTION_CHAT_ADD_USER
		m.UserIDs = x.Users
	case TL_messageActionChatDeleteUser:
		m.Type = MESSAGE_ACTION_CHAT_DELETE_USER
		m.UserID = x.User_id
	case TL_messageActionChatDeletePhoto:
		m.Type = MESSAGE_ACTION_CHAT_DELETE_PHOTO
	case TL_messageActionChatEditPhoto:
		m.Type = MESSAGE_ACTION_CHAT_EDIT_PHOTO
		m.Photo = NewPhoto(x.Photo)
	case TL_messageActionChatEditTitle:
		m.Type = MESSAGE_ACTION_CHAT_EDIT_TITLE
		m.Title = x.Title
	case TL_messageActionChatJoinedByLink:
		m.Type = MESSAGE_ACTION_CHAT_JOINED_BY_LINK
		m.UserID = x.Inviter_id
	case TL_messageActionChatMigrateTo:
		m.Type = MESSAGE_ACTION_CHAT_MIGRATE_TO
		m.ChannelID = x.Channel_id
	case TL_messageActionGameScore:
		m.Type = MESSAGE_ACTION_GAME_SCORE
		m.GameID = x.Game_id
		m.GameScore = x.Score
	case TL_messageActionHistoryClear:
		m.Type = MESSAGE_ACTION_HISTORY_CLEAN
	case TL_messageActionPinMessage:
	case TL_messageActionPhoneCall:
		m.Type = MESSAGE_ACTION_PHONE_CALL
	default:
		fmt.Println("NewMessageAction::UnSupported Input Format", reflect.TypeOf(x).String())
	}
	return
}

func NewMessageEntity(input TL) (e *MessageEntity) {
	e = new(MessageEntity)
	switch x := input.(type) {
	case TL_messageEntityBold:
		e.Type = MESSAGE_ENTITY_BOLD
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityEmail:
		e.Type = MESSAGE_ENTITY_EMAIL
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityBotCommand:
		e.Type = MESSAGE_ENTITY_BOT_COMMAND
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityHashtag:
		e.Type = MESSAGE_ENTITY_HASHTAG
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityCode:
		e.Type = MESSAGE_ENTITY_CODE
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityItalic:
		e.Type = MESSAGE_ENTITY_ITALIC
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityMention:
		e.Type = MESSAGE_ENTITY_MENTION
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityUrl:
		e.Type = MESSAGE_ENTITY_URL
		e.Offset, e.Length = x.Offset, x.Length
	case TL_messageEntityTextUrl:
		e.Type = MESSAGE_ENTITY_TEXT_URL
		e.Offset, e.Length = x.Offset, x.Length
		e.Url = x.Url
	case TL_messageEntityPre:
		e.Type = MESSAGE_ENTITY_PRE
		e.Offset, e.Length = x.Offset, x.Length
		e.language = x.Language
	case TL_messageEntityMentionName:
		e.Type = MESSAGE_ENTITY_MENTION_NAME
		e.Offset, e.Length = x.Offset, x.Length
		e.UserID = x.User_id
	default:
		fmt.Println("NewMessageEntity::UnSupported Input Format", reflect.TypeOf(x).String())
	}
	return e
}

func NewMessageForwardHeader(input TL) (fwd *MessageForwardHeader) {
	fwd = new(MessageForwardHeader)
	fwdHeader := input.(TL_messageFwdHeader)
	fwd.Date = fwdHeader.Date
	fwd.From = fwdHeader.From_id
	fwd.ChannelID = fwdHeader.Channel_id
	fwd.ChannelPost = fwdHeader.Channel_post
	fwd.Author = fwdHeader.Post_author
	return
}

// NewMessageMedia
// input:
//	1. TL_messageMediaPhoto
//	2. TL_messageMediaContact
//	3. TL_messageMediaDocument
//
func NewMessageMedia(input TL) interface{} {
	switch x := input.(type) {
	case TL_messageMediaPhoto:
		mm := new(MessageMediaPhoto)
		mm.Caption = x.Caption
		mm.Photo = *NewPhoto(x.Photo)
		return mm
	case TL_messageMediaContact:
		mm := new(MessageMediaContact)
		mm.UserID = x.User_id
		mm.Firstname = x.First_name
		mm.Lastname = x.Last_name
		mm.Phone = x.Phone_number
		return mm
	case TL_messageMediaDocument:
		mm := new(MessageMediaDocument)
		mm.Caption = x.Caption
		mm.Document = *NewDocument(x.Document)
		return mm
	case TL_messageMediaWebPage:
		// TODO:: implement it
	default:
		fmt.Println("NewMessageMedia::UnSupported Input Format", reflect.TypeOf(x).String())
	}
	return nil
}

func (m *MTProto) Messages_SendMessage(text string, peer TL, reply_to int32) (interface{}, error) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_sendMessage{
			0,
			peer,
			reply_to,
			text,
			rand.Int63(),
			TL_null{},
			nil,
		},
		resp,
	}
	x := <-resp
	switch r := x.(type) {
	default:
		log.Println(reflect.TypeOf(r))
		return r, nil

	}

}

func (m *MTProto) Messages_ImportChatInvite(hash string) *Chat {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_importChatInvite{
			hash,
		},
		resp,
	}
	x := <-resp
	switch r := x.(type) {
	case TL_updates:
		chat := NewChat(r.Chats[0])
		return chat
	case TL_rpc_error:
		log.Println(r.error_code, r.error_message)
	default:
		log.Println(reflect.TypeOf(r))
	}
	return nil
}

func (m *MTProto) Messages_GetHistory(inputPeer TL, limit, min_id, max_id int32) ([]Message, int32) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_getHistory{
			Peer:   inputPeer,
			Limit:  limit,
			Min_id: min_id,
			Max_id: max_id,
		},
		resp,
	}
	x := <-resp
	messages := make([]Message, 0, 20)
	switch input := x.(type) {
	case TL_messages_messages:
		for _, m := range input.Messages {
			msg := NewMessage(m)
			if msg != nil {
				messages = append(messages, *msg)
			}
		}
		return messages, int32(len(messages))
	case TL_messages_messagesSlice:
		for _, m := range input.Messages {
			msg := NewMessage(m)
			if msg != nil {
				messages = append(messages, *msg)
			}
		}
		return messages, input.Count
	case TL_messages_channelMessages:
		for _, m := range input.Messages {
			msg := NewMessage(m)
			if msg != nil {
				messages = append(messages, *msg)
			}
		}
		return messages, input.Count
	case TL_rpc_error:
		fmt.Println(input.error_message, input.error_code)
		return messages, 0
	default:
		fmt.Println(reflect.TypeOf(input).String())
		return messages, 0
	}

}

func (m *MTProto) Messages_GetChats(chatIDs []int32) ([]Chat, error) {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_getChats{
			Id: chatIDs,
		},
		resp,
	}
	x := <-resp
	chats := make([]Chat, 0, len(chatIDs))
	switch input := x.(type) {
	case TL_messages_chats:
		for _, ch := range input.Chats {
			chats = append(chats, *NewChat(ch))
		}
		return chats, nil
	case TL_rpc_error:
		fmt.Println(input.error_code, input.error_message)
		return chats, fmt.Errorf("TL_rpc_error: %d - %s", input.error_code, input.error_message)
	default:
		fmt.Println(reflect.TypeOf(input).String())
		return chats, fmt.Errorf("Don't know how to handle response: %s - %v", reflect.TypeOf(input).String(), input)
	}
}

func (m *MTProto) Messages_GetFullChat(chatID int32) *Chat {
	resp := make(chan TL, 1)
	m.queueSend <- packetToSend{
		TL_messages_getFullChat{
			Chat_id: chatID,
		},
		resp,
	}
	x := <-resp
	chat := new(Chat)
	switch input := x.(type) {
	case TL_messages_chatFull:
		chat = NewChat(input)
	default:

	}
	return chat
}
