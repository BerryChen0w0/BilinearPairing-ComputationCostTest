// 测试gnark里面BN254库的各种运算性能

package gnark

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	bn254 "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// t是测试次数

func FTest_bn254(t int) {
	var a, b fr.Element
	var c fr.Element

	add_time := 0 * time.Millisecond
	sub_time := 0 * time.Millisecond
	neg_time := 0 * time.Millisecond
	mul_time := 0 * time.Millisecond
	div_time := 0 * time.Millisecond
	inv_time := 0 * time.Millisecond

	var start time.Time
	fmt.Println("Testing Field operations for", t, "times...")

	for i := 0; i < t; i++ {
		a.SetRandom()
		b.SetRandom()
		// fmt.Println("a:", a)
		// fmt.Println("b:", b)

		start = time.Now()
		c.Add(&a, &b)
		add_time += time.Since(start)
		// fmt.Println("c = a + b:", c)

		start = time.Now()
		c.Sub(&a, &b)
		sub_time += time.Since(start)
		// fmt.Println("c = a - b:", c)

		start = time.Now()
		c.Neg(&a)
		neg_time += time.Since(start)
		// fmt.Println("c = -a:", c)

		start = time.Now()
		c.Mul(&a, &b)
		mul_time += time.Since(start)
		// fmt.Println("c = a * b:", c)

		start = time.Now()
		c.Div(&a, &b)
		div_time += time.Since(start)
		// fmt.Println("c = a / b:", c)

		start = time.Now()
		c.Inverse(&a)
		inv_time += time.Since(start)
		// fmt.Println("c = 1/a:", c)
	}

	fmt.Println("Field operations completed.")
	fmt.Println("Average time for Field operations:")
	fmt.Println("Add:", add_time/time.Duration(t))
	fmt.Println("Sub:", sub_time/time.Duration(t))
	fmt.Println("Neg:", neg_time/time.Duration(t))
	fmt.Println("Mul:", mul_time/time.Duration(t))
	fmt.Println("Div:", div_time/time.Duration(t))
	fmt.Println("Inv:", inv_time/time.Duration(t))
	return
}

func GTest_bn254(t int) {
	var err error
	_, _, g1, g2 := bn254.Generators()

	// prepare g1, g2, gt
	var P1, Q1, R1 bn254.G1Affine
	var P2, Q2, R2 bn254.G2Affine
	var Pt, Qt, Rt bn254.GT
	var a, b *big.Int

	bitLength := 256

	a, err = rand.Prime(rand.Reader, bitLength)
	if err != nil {
		panic(err)
	}
	b, err = rand.Prime(rand.Reader, bitLength)
	if err != nil {
		panic(err)
	}
	g1.ScalarMultiplication(&g1, a)
	g2.ScalarMultiplication(&g2, b)
	gt, err := bn254.Pair([]bn254.G1Affine{g1}, []bn254.G2Affine{g2})
	if err != nil {
		panic(err)
	}

	// test variables to accumulate time

	add_time_G1 := 0 * time.Millisecond
	sub_time_G1 := 0 * time.Millisecond
	neg_time_G1 := 0 * time.Millisecond
	dot_time_G1 := 0 * time.Millisecond

	add_time_G2 := 0 * time.Millisecond
	sub_time_G2 := 0 * time.Millisecond
	neg_time_G2 := 0 * time.Millisecond
	dot_time_G2 := 0 * time.Millisecond

	mul_time_Gt := 0 * time.Millisecond
	div_time_Gt := 0 * time.Millisecond
	inv_time_Gt := 0 * time.Millisecond
	exp_time_Gt := 0 * time.Millisecond

	var start time.Time
	fmt.Println("Testing G1,G2,Gt operations for", t, "times...")

	for i := 0; i < t; i++ {
		a, err = rand.Prime(rand.Reader, bitLength)
		if err != nil {
			panic(err)
		}
		b, err = rand.Prime(rand.Reader, bitLength)
		if err != nil {
			panic(err)
		}

		/* G1 */

		// G1 Scalar Multiplication
		start = time.Now()
		P1.ScalarMultiplication(&g1, a)
		dot_time_G1 += time.Since(start)
		Q1.ScalarMultiplication(&g1, b)
		// G1 Addition
		start = time.Now()
		R1.Add(&P1, &Q1)
		add_time_G1 += time.Since(start)

		// G1 Subtraction
		start = time.Now()
		R1.Sub(&P1, &Q1)
		sub_time_G1 += time.Since(start)
		// G1 Negation (additive inverse)
		start = time.Now()
		R1.Neg(&R1)
		neg_time_G1 += time.Since(start)

		/* G2 */

		// G2 Scalar Multiplication
		start = time.Now()
		P2.ScalarMultiplication(&g2, a)
		dot_time_G2 += time.Since(start)
		Q2.ScalarMultiplication(&g2, b)
		// G2 Addition
		start = time.Now()
		R2.Add(&P2, &Q2)
		add_time_G2 += time.Since(start)

		// G2 Subtraction
		start = time.Now()
		R2.Sub(&P2, &Q2)
		sub_time_G2 += time.Since(start)
		// G2 Negation (additive inverse)
		start = time.Now()
		R2.Neg(&R2)
		neg_time_G2 += time.Since(start)

		/* Gt. Note that Gt is a multiplicative group */

		// Gt Scalar Multiplication
		start = time.Now()
		Pt.Exp(gt, a)
		exp_time_Gt += time.Since(start)
		Qt.Exp(gt, b)
		// Gt Addition
		start = time.Now()
		Rt.Mul(&Pt, &Qt)
		mul_time_Gt += time.Since(start)

		// Gt Subtraction
		start = time.Now()
		Rt.Div(&Pt, &Qt)
		div_time_Gt += time.Since(start)
		// Gt Negation (additive inverse)
		start = time.Now()
		Rt.Inverse(&Rt)
		inv_time_Gt += time.Since(start)
	}

	fmt.Println("G1,G2,Gt operations completed.")
	fmt.Println("Average time for G1,G2,Gt operations:")

	fmt.Println("G1 Add:", add_time_G1/time.Duration(t))
	fmt.Println("G1 Sub:", sub_time_G1/time.Duration(t))
	fmt.Println("G1 Neg:", neg_time_G1/time.Duration(t))
	fmt.Println("G1 Mul:", dot_time_G1/time.Duration(t))
	fmt.Println()
	fmt.Println("G2 Add:", add_time_G2/time.Duration(t))
	fmt.Println("G2 Sub:", sub_time_G2/time.Duration(t))
	fmt.Println("G2 Neg:", neg_time_G2/time.Duration(t))
	fmt.Println("G2 Mul:", dot_time_G2/time.Duration(t))
	fmt.Println()
	fmt.Println("Gt Mul:", mul_time_Gt/time.Duration(t))
	fmt.Println("Gt Div:", div_time_Gt/time.Duration(t))
	fmt.Println("Gt Inv:", inv_time_Gt/time.Duration(t))
	fmt.Println("Gt Exp:", exp_time_Gt/time.Duration(t))

	return
}

func PairingTest_bn254(t int) {
	var err error
	_, _, g1, g2 := bn254.Generators()

	var a, b *big.Int
	bitLength := 256

	a, err = rand.Prime(rand.Reader, bitLength)
	if err != nil {
		panic(err)
	}
	b, err = rand.Prime(rand.Reader, bitLength)
	if err != nil {
		panic(err)
	}
	g1.ScalarMultiplication(&g1, a)
	g2.ScalarMultiplication(&g2, b)

	pairing_time := 0 * time.Millisecond

	var start time.Time
	fmt.Println("Testing Pairing operations for", t, "times...")

	for i := 0; i < t; i++ {
		a, err = rand.Prime(rand.Reader, bitLength)
		if err != nil {
			panic(err)
		}
		b, err = rand.Prime(rand.Reader, bitLength)
		if err != nil {
			panic(err)
		}

		g1.ScalarMultiplication(&g1, a)
		g2.ScalarMultiplication(&g2, b)

		// Pairing
		start = time.Now()
		_, err = bn254.Pair([]bn254.G1Affine{g1}, []bn254.G2Affine{g2})
		if err != nil {
			panic(err)
		}
		pairing_time += time.Since(start)
	}

	fmt.Println("Pairing operations completed.")
	fmt.Println("Average time for Pairing operations:")
	fmt.Println("Pairing:", pairing_time/time.Duration(t))
	return
}

func BN254_Test() {
	const t = 1e4
	fmt.Println("BN254 Computation Cost Test")
	fmt.Println("")
	FTest_bn254(t)
	fmt.Println("")
	GTest_bn254(t)
	fmt.Println("")
	PairingTest_bn254(t)
	fmt.Println("")
	fmt.Println("BN254 Computation Cost Test Completed")
}
