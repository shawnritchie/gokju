package message

import (
	"github.com/shawnritchie/gokju/structs"
)

type Messenger interface {
	identifier() string
	payload() structs.Payload
	withMetaData(metaData map[string]structs.Payload) (Messenger, error)
	andMetaData(metaData map[string]structs.Payload) (Messenger, error)
}
