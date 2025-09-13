package pbc

import (
	"fmt"
	"time"

	"github.com/Nik-U/pbc"
)

// t是测试次数

func pbc_FTest(t int, pairing *pbc.Pairing) {
	a, b, c := pairing.NewZr(), pairing.NewZr(), pairing.NewZr()

	add_time := 0 * time.Millisecond
	sub_time := 0 * time.Millisecond
	neg_time := 0 * time.Millisecond
	mul_time := 0 * time.Millisecond
	div_time := 0 * time.Millisecond
	inv_time := 0 * time.Millisecond

	var start time.Time
	fmt.Println("Testing Field operations for", t, "times...")

	for i := 0; i < t; i++ {
		a = pairing.NewZr().Rand()
		b = pairing.NewZr().Rand()
		// fmt.Println("a:", a)
		// fmt.Println("b:", b)

		start = time.Now()
		c.Add(a, b)
		add_time += time.Since(start)
		// fmt.Println("c = a + b:", c)

		start = time.Now()
		c.Sub(a, b)
		sub_time += time.Since(start)
		// fmt.Println("c = a - b:", c)

		start = time.Now()
		c.Neg(a)
		neg_time += time.Since(start)
		// fmt.Println("c = -a:", c)

		start = time.Now()
		c.Mul(a, b)
		mul_time += time.Since(start)
		// fmt.Println("c = a * b:", c)

		start = time.Now()
		c.Div(a, b)
		div_time += time.Since(start)
		// fmt.Println("c = a / b:", c)

		start = time.Now()
		c.Invert(a)
		inv_time += time.Since(start)
		// fmt.Println("c = a^-1:", c)
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

func pbc_GTest(t int, pairing *pbc.Pairing) {
	P1 := pairing.NewG1().Rand()
	Q1 := pairing.NewG1().Rand()
	R1 := pairing.NewG1()
	P2 := pairing.NewG2().Rand()
	Q2 := pairing.NewG2().Rand()
	R2 := pairing.NewG2()
	Pt := pairing.NewGT().Rand()
	Qt := pairing.NewGT().Rand()
	Rt := pairing.NewGT()

	a := pairing.NewZr()
	b := pairing.NewZr()

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
		a = pairing.NewZr().Rand()
		b = pairing.NewZr().Rand()

		/* G1 */
		// G1 Addition
		start = time.Now()
		R1.Add(P1, Q1)
		add_time_G1 += time.Since(start)
		// G1 Subtraction
		start = time.Now()
		Q1.Sub(P1, Q1)
		sub_time_G1 += time.Since(start)
		// G1 Negation
		start = time.Now()
		R1.Neg(P1)
		neg_time_G1 += time.Since(start)
		// G1 Scalar Multiplication
		start = time.Now()
		P1.MulZn(P1, a)
		dot_time_G1 += time.Since(start)

		/* G2 */
		// G2 Addition
		start = time.Now()
		R2.Add(P2, Q2)
		add_time_G2 += time.Since(start)
		// G2 Subtraction
		start = time.Now()
		Q2.Sub(P2, Q2)
		sub_time_G2 += time.Since(start)
		// G2 Negation
		start = time.Now()
		R2.Neg(P2)
		neg_time_G2 += time.Since(start)
		// G2 Scalar Multiplication
		start = time.Now()
		P2.MulZn(P2, b)
		dot_time_G2 += time.Since(start)

		/* Gt */
		// Gt Multiplication
		start = time.Now()
		Rt.Mul(Pt, Qt)
		mul_time_Gt += time.Since(start)
		// Gt Division
		start = time.Now()
		Qt.Div(Pt, Qt)
		div_time_Gt += time.Since(start)
		// Gt Inversion
		start = time.Now()
		Rt.Invert(Pt)
		inv_time_Gt += time.Since(start)
		// Gt Exponentiation
		start = time.Now()
		Pt.PowZn(Pt, a)
		exp_time_Gt += time.Since(start)
	}

	fmt.Println("G1,G2,Gt operations completed.")
	fmt.Println("Average time for G1 operations:")
	fmt.Println("Add:", add_time_G1/time.Duration(t))
	fmt.Println("Sub:", sub_time_G1/time.Duration(t))
	fmt.Println("Neg:", neg_time_G1/time.Duration(t))
	fmt.Println("Dot:", dot_time_G1/time.Duration(t))

	fmt.Println("Average time for G2 operations:")
	fmt.Println("Add:", add_time_G2/time.Duration(t))
	fmt.Println("Sub:", sub_time_G2/time.Duration(t))
	fmt.Println("Neg:", neg_time_G2/time.Duration(t))
	fmt.Println("Dot:", dot_time_G2/time.Duration(t))

	fmt.Println("Average time for Gt operations:")
	fmt.Println("Mul:", mul_time_Gt/time.Duration(t))
	fmt.Println("Div:", div_time_Gt/time.Duration(t))
	fmt.Println("Inv:", inv_time_Gt/time.Duration(t))
	fmt.Println("Exp:", exp_time_Gt/time.Duration(t))
	return
}

func pbc_PairingTest(t int, pairing *pbc.Pairing) {
	g1 := pairing.NewG1().Rand()
	g2 := pairing.NewG2().Rand()
	gt := pairing.NewGT()

	a, b := pairing.NewZr(), pairing.NewZr()

	pairing_time := 0 * time.Millisecond

	var start time.Time
	fmt.Println("Testing Pairing operation for", t, "times...")

	for i := 0; i < t; i++ {
		a, b = pairing.NewZr().Rand(), pairing.NewZr().Rand()
		g1.PowZn(g1, a)
		g2.PowZn(g2, b)

		start = time.Now()
		gt.Pair(g1, g2)
		pairing_time += time.Since(start)
	}

	fmt.Println("Pairing operations completed.")
	fmt.Println("Average time for Pairing operation:")
	fmt.Println("Pairing:", pairing_time/time.Duration(t))
	return

}
