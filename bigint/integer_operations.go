package bigint

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// test integer operations
// with diffferent bit length of big.Int

// Test interger operations with integers of bigLength bits length. Repeat t times
func Test_integer_operations_with_bitlength(bitLength int, t int) {
	var err error
	var a, b, p *big.Int

	add_time := 0 * time.Millisecond
	sub_time := 0 * time.Millisecond
	mul_time := 0 * time.Millisecond
	div_time := 0 * time.Millisecond
	mod_time := 0 * time.Millisecond
	modexp_time := 0 * time.Millisecond

	p, err = rand.Prime(rand.Reader, bitLength)
	if err != nil {
		panic(err)
	}

	var start time.Time
	fmt.Println("Testing integer operations with", bitLength, "bits length for", t, "times...")

	for i := 0; i < t; i++ {
		a, err = rand.Int(rand.Reader, p)
		if err != nil {
			panic(err)
		}
		b, err = rand.Int(rand.Reader, p)
		if err != nil {
			panic(err)
		}

		// Addition
		start = time.Now()
		_ = new(big.Int).Add(a, b)
		add_time += time.Since(start)

		// Subtraction
		start = time.Now()
		_ = new(big.Int).Sub(a, b)
		sub_time += time.Since(start)

		// Multiplication
		start = time.Now()
		_ = new(big.Int).Mul(a, b)
		mul_time += time.Since(start)

		// Division
		start = time.Now()
		_ = new(big.Int).Div(a, b)
		div_time += time.Since(start)

		// Modular reduction
		start = time.Now()
		_ = new(big.Int).Mod(a, p)
		mod_time += time.Since(start)

		// Modular exponentiation
		start = time.Now()
		_ = new(big.Int).Exp(a, b, p)
		modexp_time += time.Since(start)
	}

	fmt.Println("Addition:", add_time/time.Duration(t))
	fmt.Println("Subtraction:", sub_time/time.Duration(t))
	fmt.Println("Multiplication:", mul_time/time.Duration(t))
	fmt.Println("Division:", div_time/time.Duration(t))
	fmt.Println("Modular reduction:", mod_time/time.Duration(t))
	fmt.Println("Modular exponentiation:", modexp_time/time.Duration(t))
	fmt.Println()
	return

}

func Test_integer_operations() {
	// To achieve 128-bit security using a group with the Discrete Logarithm Problem (DLP), you would need to use either:
	// * Elliptic curve groups such as bn254 and bls12381.
	// * An integer multiplication group using a large prime of about 3072 bits.

	const t = 1e5

	bitLengths := []int{64, 128, 256, 512, 1024, 2048, 3072, 4096}
	for _, bl := range bitLengths {
		Test_integer_operations_with_bitlength(bl, t)
	}
}
