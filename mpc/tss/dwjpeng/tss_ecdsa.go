package dwjpeng

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/dwjpeng/mpc-tss-lib/common"
	"github.com/dwjpeng/mpc-tss-lib/ecdsa/keygen"
	"github.com/dwjpeng/mpc-tss-lib/ecdsa/signing"

	"github.com/dwjpeng/mpc-tss-lib/tss"
)

type PartyKeyData struct {
	ID   string
	Data keygen.LocalPartySaveData
}
type PartySignatureData struct {
	ID   string
	Data common.SignatureData
}

func getParticipantPartyIDs(num int) []*tss.PartyID {
	parties := []*tss.PartyID{}

	for i := 0; i < num; i++ {
		fmt.Printf("Create partyID %d \n ", (i + 1))
		id1 := fmt.Sprintf("%d", (i + 1))
		moniker1 := fmt.Sprintf("%d", (i + 1))
		uniqueKey1 := big.NewInt(int64(i + 1))
		partyID := tss.NewPartyID(id1, moniker1, uniqueKey1)
		parties = append(parties, partyID)
	}
	return parties
}

func routingMessage(outCh chan tss.Message, mapChs map[string]chan tss.ParsedMessage) {
	for msg := range outCh {
		bz, _, err := msg.WireBytes()
		if err != nil {
			fmt.Printf("From Party %s: Msg parse failed %s \n", msg.GetFrom().Id, err.Error())
			continue
		}
		pMsg, err := tss.ParseWireMessage(bz, msg.GetFrom(), msg.IsBroadcast())
		if err != nil {
			fmt.Printf("From Party %s: Msg ParseWireMessage failed %s \n", msg.GetFrom().Id, err.Error())
			continue
		}

		if msg.IsBroadcast() {
			fmt.Printf("From Party %s: broadcast msg....\n", msg.GetFrom().Id)
			for partyID, inCh := range mapChs {
				// fmt.Printf("From %s => Check Party %s \n", msg.GetFrom().Id, partyID)
				if partyID == msg.GetFrom().Id {
					// fmt.Printf("From %s => Bypass Party %s \n", msg.GetFrom().Id, partyID)
					continue
				}
				// fmt.Printf("From %s => send message to Party %s \n", msg.GetFrom().Id, partyID)
				inCh <- pMsg
				// fmt.Printf("From %s => send message to Party %s  done. \n", msg.GetFrom().Id, partyID)

			}
		} else {
			receivers := msg.GetTo()
			if len(receivers) == 0 {
				fmt.Printf("Party %s send msg no-one \n", msg.GetFrom().Id)
				continue
			}
			for _, receiver := range receivers {
				if receiver == msg.GetFrom() {
					continue
				}
				if inCh, ok := mapChs[receiver.GetId()]; ok {
					inCh <- pMsg
				} else {
					fmt.Printf("Party %s send msg to %s but not found \n", msg.GetFrom().Id, receiver.GetId())
					continue
				}
			}
		}
	}
}

func hashData(text string) (*big.Int, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(text))
	if err != nil {
		log.Fatalf("Cannot create hash")
		return nil, err
	}
	md := hasher.Sum(nil)
	msg := new(big.Int).SetBytes(md)
	return msg, nil
}

func GenKey(curve elliptic.Curve, threshold int, peers []*tss.PartyID) (map[string]keygen.LocalPartySaveData, error) {
	parties := tss.SortPartyIDs(peers)
	ctx := tss.NewPeerContext(parties)

	outCh := make(chan tss.Message)
	partyCh := make(chan PartyKeyData)

	mapParties := make(map[string]tss.Party, 0)

	mapInChs := make(map[string]chan tss.ParsedMessage, 0)
	mapEndChs := make(map[string]chan keygen.LocalPartySaveData, 0)

	//1. Create routing goroutine to handle messages between parties
	go func() {
		routingMessage(outCh, mapInChs)
	}()

	//2. Run Parties
	for i := 0; i < len(peers); i++ {
		// 2.1 Create party
		thisParty := parties[i]

		params := tss.NewParameters(curve, ctx, thisParty, len(parties), threshold)

		inCh := make(chan tss.ParsedMessage)
		endCh := make(chan keygen.LocalPartySaveData)

		preParams, err := keygen.GeneratePreParams(1 * time.Minute)
		if err != nil {
			log.Fatalf("GeneratePreParams Failed: %s ", err.Error())
			return nil, err
		}

		party := keygen.NewLocalParty(params, outCh, endCh, *preParams)

		//Save to share message
		mapParties[thisParty.Id] = party
		mapInChs[thisParty.Id] = inCh
		mapEndChs[thisParty.Id] = endCh

		//2.2 Start party
		go func() {
			time.Sleep(5 * time.Second)
			err := party.Start()
			if err != nil {
				fmt.Printf("Party %s finished err: %s \n", party.PartyID().GetId(), err.Error())
				return
			}
			fmt.Printf("Party %s finished \n", party.PartyID().GetId())
		}()

		// 2.3 Party process Message from other Peers
		go func() {
			for msg := range inCh {
				// fmt.Printf("Party %s receive from: %s => call party.Update \n", thisParty.Id, msg.GetFrom().GetId())
				go func() {
					ok, err := party.Update(msg)
					if err != nil {
						fmt.Printf("Party %s => update message failed: %s \n", thisParty.Id, err.Error())
					}
					if ok {
						fmt.Printf("Party %s => update success msg from: %s \n", thisParty.Id, msg.GetFrom().GetId())
					} else {
						fmt.Printf("Party %s => update msg failed at  %s \n", thisParty.Id, msg.GetFrom().GetId())
					}
				}()
			}
		}()

		// 2.4 Party process end message
		go func() {
			for msg := range endCh {
				partyCh <- PartyKeyData{
					ID:   thisParty.Id,
					Data: msg,
				}
			}
		}()

	}

	mapKeys := make(map[string]keygen.LocalPartySaveData, 0)

	for party := range partyCh {
		// fmt.Printf("Party %s LocalPartySaveData msg: %+v \n", party.ID, party.Data)
		// time.Sleep(1 * time.Second)
		mapKeys[party.ID] = party.Data
		if len(mapKeys) == len(peers) {
			break
		}
	}
	defer func() {
		// Close all channels
		for _, ch := range mapInChs {
			close(ch)
		}
		for _, ch := range mapEndChs {
			close(ch)
		}
		close(outCh)
		close(partyCh)
	}()
	return mapKeys, nil
}

func SignData(msg *big.Int, curve elliptic.Curve, threshold int, peers []*tss.PartyID, mapKeys map[string]keygen.LocalPartySaveData) (map[string]common.SignatureData, error) {

	outCh := make(chan tss.Message)
	mapInChs := make(map[string]chan tss.ParsedMessage, 0)

	parties := tss.SortPartyIDs(peers)
	ctx := tss.NewPeerContext(parties)

	partyCh := make(chan PartySignatureData)

	mapParties := make(map[string]tss.Party, 0)
	mapEndChs := make(map[string]chan common.SignatureData, 0)

	//1. Create routing goroutine to handle messages between parties
	go func() {
		routingMessage(outCh, mapInChs)
	}()

	for i := 0; i < len(peers); i++ {
		thisParty := parties[i]
		ourKeyData := mapKeys[thisParty.Id]
		params := tss.NewParameters(curve, ctx, thisParty, len(parties), threshold)

		inCh := make(chan tss.ParsedMessage)
		endCh := make(chan common.SignatureData)

		party := signing.NewLocalParty(msg, params, ourKeyData, outCh, endCh)

		//Save to share message
		mapParties[thisParty.Id] = party
		mapInChs[thisParty.Id] = inCh
		mapEndChs[thisParty.Id] = endCh

		//2.2 Start party
		go func() {
			time.Sleep(5 * time.Second)
			err := party.Start()
			if err != nil {
				fmt.Printf("Party %s finished err: %s \n", party.PartyID().GetId(), err.Error())
				return
			}
			fmt.Printf("Party %s finished \n", party.PartyID().GetId())
		}()

		// 2.3 Party process Message from other Peers
		go func() {
			for msg := range inCh {
				// fmt.Printf("Party %s receive from: %s => call party.Update \n", thisParty.Id, msg.GetFrom().GetId())
				go func() {
					ok, err := party.Update(msg)
					if err != nil {
						fmt.Printf("Party %s => update message failed: %s \n", thisParty.Id, err.Error())
					}
					if ok {
						fmt.Printf("Party %s => update success msg from: %s \n", thisParty.Id, msg.GetFrom().GetId())
					} else {
						fmt.Printf("Party %s => update msg failed at  %s \n", thisParty.Id, msg.GetFrom().GetId())
					}
				}()
			}
		}()

		// 2.4 Party process end message
		go func() {
			for msg := range endCh {
				partyCh <- PartySignatureData{
					ID:   thisParty.Id,
					Data: msg,
				}
			}
		}()
	}

	mapSignatures := make(map[string]common.SignatureData, 0)

	for party := range partyCh {
		// fmt.Printf("Party %s SignatureData msg: %+v \n", party.ID, party.Data)
		// time.Sleep(1 * time.Second)
		mapSignatures[party.ID] = party.Data
		if len(mapSignatures) == len(peers) {
			break
		}
	}
	defer func() {
		// Close all channels
		for _, ch := range mapInChs {
			close(ch)
		}
		for _, ch := range mapEndChs {
			close(ch)
		}
		close(outCh)
		close(partyCh)
	}()

	return mapSignatures, nil
}
func RUN_TSS_ECDSA() error {
	fmt.Printf("1. GenKey start.... \n")

	num := 3
	// use ECDSA
	curve := tss.S256()
	// or use EdDSA
	// curve := tss.Edwards()
	threshold := int(2)

	peers := getParticipantPartyIDs(num)

	//1. Start generate keys for peers
	keyData, err := GenKey(curve, threshold, peers)
	if err != nil {
		fmt.Printf("GenKey failed: %s \n", err.Error())
		return err
	}
	fmt.Printf("GenKey done. \n")

	time.Sleep(10 * time.Second)

	//2. Sign data
	fmt.Printf("2. Signing start.......... \n")
	msg := "Thu nghiem"
	hash, err := hashData(msg)
	if err != nil {
		fmt.Printf("hash data failed %s \n", err.Error())
		return err
	}
	signatureData, err := SignData(hash, curve, threshold, peers, keyData)
	if err != nil {
		fmt.Printf("GenSignDataKey failed: %s \n", err.Error())
		return err
	}
	fmt.Printf("2. Signing end. \n")

	// for partyId, data := range signatureData {
	// 	fmt.Printf("Party %s => signature: %+v \n", partyId, data)
	// }

	//3. Verify Signature
	pub := keyData["1"].ECDSAPub
	signData := signatureData["1"]

	pk := ecdsa.PublicKey{
		Curve: curve,
		X:     pub.X(),
		Y:     pub.Y(),
	}
	r := new(big.Int).SetBytes(signData.R)
	s := new(big.Int).SetBytes(signData.S)

	ok := ecdsa.Verify(&pk, hash.Bytes(), r, s)
	if !ok {
		fmt.Printf("Verify failed \n")
		return fmt.Errorf("signature verification failed")
	}
	fmt.Printf("Verify done... \n")

	return nil
}
