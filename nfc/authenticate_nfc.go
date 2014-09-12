// Package for handling NFC related stuff
package nfc

import (
	"fmt"
	"github.com/fuzxxl/freefare/0.2/freefare"
	"github.com/fuzxxl/nfc/1.0/nfc"
	"github.com/fuzxxl/openkey/0.1/openkey"
	"time"
)

// As the function name says - authenticate NFC
func authenticate_nfc(basedir string, nfc_c chan string) {

	var err error
	var openkeyContext openkey.Context
	var device nfc.Device
	var tags = []freefare.Tag{}
	var cardId, lastCardId string

	INFO.Printf("%s \n", basedir)

	openkeyContext = openkey.New()
	if err != nil {
		ERROR.Println(err)
	}
	err = openkeyContext.AddRole(openkey.CardAuthenticator, basedir)
	if err != nil {
		ERROR.Println(err)
	}

	if openkeyContext.PrepareAuthenticator() != true {
		ERROR.Println("Couldn't prepare the card authenticator\n")
	}

	device, err = nfc.Open("")
	if err != nil {
		ERROR.Println(err)
	}

	for true {
		tags, err = freefare.GetTags(device)
		if err != nil {
			ERROR.Println(err)
		} else if len(tags) == 0 {
			lastCardId = ""
			TRACE.Println("no tags found")
		} else {
			cardId, err = openkeyContext.AuthenticateCard(tags[0].(freefare.DESFireTag), nil)
			if err != nil {
				ERROR.Println(err)
			} else if cardId != lastCardId {
				nfc_c <- fmt.Sprintf("%s", cardId)
				lastCardId = cardId
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
}
