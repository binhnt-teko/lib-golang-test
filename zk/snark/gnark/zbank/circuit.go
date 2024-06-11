package zbank

import "github.com/consensys/gnark/frontend"

//
// Circuit to ensure that a transfer is legit
//

type Circuit struct {
	// note: Alice only has 500 tokens in her account
	AliceBalance frontend.Variable `gnark:",public"`
	// note: Bob has 0 tokens in his account
	BobBalance      frontend.Variable `gnark:",public"`
	NewBobBalance   frontend.Variable `gnark:",public"`
	NewAliceBalance frontend.Variable
	Transfer        frontend.Variable
}

func (circuit *Circuit) Define(api frontend.API) error {
	// init
	gkrBalance := NewBalanceGKR(api, 1)

	// transfer is legit?
	api.AssertIsLessOrEqual(circuit.Transfer, circuit.AliceBalance)

	// new balance for Alice
	negated := api.Neg(circuit.Transfer)
	newAliceBalance := gkrBalance.AddCircuit(circuit.AliceBalance, negated)
	api.AssertIsEqual(newAliceBalance, circuit.NewAliceBalance)

	// new balance for Bob
	newBobBalance := gkrBalance.AddCircuit(circuit.BobBalance, circuit.Transfer)
	api.AssertIsEqual(newBobBalance, circuit.NewBobBalance)

	// GKR verifier
	err := gkrBalance.VerifyGKR(circuit.AliceBalance, circuit.BobBalance)
	if err != nil {
		panic(err)
	}
	return nil
}
