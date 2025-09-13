package pbc

import (
	"fmt"

	"github.com/Nik-U/pbc"
)

func pbcTestWorker(t int, params *pbc.Params) {
	pairing := params.NewPairing()
	// fmt.Println("Pairing parameters:", params.String())
	fmt.Println()
	pbc_FTest(t, pairing)
	fmt.Println()
	pbc_GTest(t, pairing)
	fmt.Println()
	pbc_PairingTest(t, pairing)
	fmt.Println()
}

func TypeA_Test() {
	const t = 1e4
	params := pbc.GenerateA(160, 512)
	fmt.Println("Starting Type A Computation Cost Test...")
	pbcTestWorker(t, params)
	fmt.Println("Type A Computation Cost Test Completed")
}

func TypeD_Test() {
	const t = 1e4
	params, err := pbc.GenerateD(9563, 160, 171, 500)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting Type D Computation Cost Test...")
	pbcTestWorker(t, params)
	fmt.Println("Type D Computation Cost Test Completed")
}

func TypeF_Test() {
	const t = 1e4
	params := pbc.GenerateF(160)
	fmt.Println("Starting Type F Computation Cost Test...")
	pbcTestWorker(t, params)
	fmt.Println("Type F Computation Cost Test Completed")
}
