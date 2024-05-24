package util

import (
	"fmt"

	"github.com/miekg/pkcs11"
)

type CancelFunc func(pkcs11.SessionHandle)

func HSM_Login(hsmLib, pin, token string) (*pkcs11.Ctx, CancelFunc, pkcs11.SessionHandle, error) {
	p := pkcs11.New(hsmLib)
	err := p.Initialize()
	close := func(pkcs11.SessionHandle) {}
	if err != nil {
		return nil, close, pkcs11.SessionHandle(0), err
	}
	close = func(session pkcs11.SessionHandle) {
		p.Finalize()
		p.Destroy()
	}
	slots, err := p.GetSlotList(true)
	if err != nil {
		return nil, close, pkcs11.SessionHandle(0), err
	}
	slotId := uint(0)
	for _, slot := range slots {
		tokenInfo, _ := p.GetTokenInfo(slot)
		if tokenInfo.Label == token {
			slotId = slot
			break
		}
	}

	session, err := p.OpenSession(slotId, pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return nil, close, pkcs11.SessionHandle(0), err
	}
	close = func(session pkcs11.SessionHandle) {
		p.CloseSession(session)
		p.Finalize()
		p.Destroy()
	}

	err = p.Login(session, pkcs11.CKU_USER, pin)
	if err != nil {
		return nil, close, session, err
	}
	close = func(session pkcs11.SessionHandle) {
		p.Logout(session)
		p.CloseSession(session)
		p.Finalize()
		p.Destroy()
	}

	return p, close, session, nil
}

func HSM_Find(p *pkcs11.Ctx, session pkcs11.SessionHandle, tpls []*pkcs11.Attribute, numObj int) ([]pkcs11.ObjectHandle, error) {
	if err := p.FindObjectsInit(session, tpls); err != nil {
		fmt.Printf("FindObjectsInit failed: %s \n", err)
		return nil, err
	}
	handles, moreAvailable, err := p.FindObjects(session, numObj)
	if err != nil {
		fmt.Printf("FindObjects failed: %s \n", err)
		return nil, err
	}
	if moreAvailable {
		fmt.Printf("Too many object return from token \n")
		return nil, fmt.Errorf("NUMBER OBJECT OVER %d", numObj)
	}
	if err = p.FindObjectsFinal(session); err != nil {
		p.Destroy()
		fmt.Printf("FindObjectsFinal failed: %s \n", err.Error())
		return nil, err
	}
	return handles, nil
}
func HSM_ListSlot(libPath string) error {
	p := pkcs11.New(libPath)
	err := p.Initialize()
	if err != nil {
		fmt.Printf("Initialize failed %s \n", err.Error())
		return err
	}
	defer p.Destroy()
	defer p.Finalize()
	slots, err := p.GetSlotList(true)
	if err != nil {
		fmt.Printf("GetSlotList failed %s \n", err.Error())
		return err
	}
	for _, slot := range slots {
		slotInfo, err := p.GetSlotInfo(slot)
		if err != nil {
			continue
		}
		fmt.Printf("Slot: %s\n", slotInfo.SlotDescription)
	}
	return nil
}

func getSlotId(p *pkcs11.Ctx, slots []uint, token string) uint {
	slotID := slots[0]
	for _, slot := range slots {
		tokenInfo, err := p.GetTokenInfo(slot)
		if err != nil {
			continue
		}
		if tokenInfo.Label == token {
			slotID = slot
			break
		}
	}
	return slotID
}
