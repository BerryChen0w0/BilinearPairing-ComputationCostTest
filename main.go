package main

import (
	"main/bigint"
	"main/gnark"
	"main/pbc"
)

func main() {
	gnark.BN254_Test()
	gnark.BLS12381_Test()

	pbc.TypeA_Test()
	pbc.TypeD_Test()
	pbc.TypeF_Test()

	bigint.Test_integer_operations()
}
