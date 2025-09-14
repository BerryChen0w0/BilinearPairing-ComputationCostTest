
# Benchmarking Computational Cost of Bilinear Pairing Operations

This document benchmarks the computational cost of various operations involved in the bilinear pairing $e:\mathbb{G}_1 \times \mathbb{G}_2 \to \mathbb{G}_T$. The operations include:

*   Field operations: addition, subtraction, multiplication, division, and inversion.
*   Group operations: group addition, scalar multiplication (which corresponds to group multiplication and exponentiation for multiplicative groups).
*   Pairing operation: the pairing function itself.

The tests are conducted on popular elliptic curves using common Go language libraries.

## Symbol Convention

For an algebraic structure $[S, +, \times]$, the additive inverse is denoted by $neg$, and its inverse operation is $-$. The multiplicative inverse is denoted by $inv$, and its inverse operation is $\div$. We use $\cdot$ for scalar multiplication in an additive group and $\exp$ for exponentiation (scalar multiplication) in a multiplicative group.

In cryptography literature, for a pairing $e:\mathbb{G}_1 \times \mathbb{G}_2 \to \mathbb{G}_T$, $\mathbb{G}_1$ and $\mathbb{G}_2$ are typically written as additive groups, while $\mathbb{G}_T$ is written as a multiplicative group, satisfying $e(a\cdot P, b\cdot Q) = e(P, Q)^{ab}$. The libraries used in this benchmark follow this convention.

## Test Environment

**Hardware**

- **CPU**: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
- **Memory**: 32 GB DDR4
- **Storage**: 500 GB NVMe SSD

**Software**

- **Operating System**: Ubuntu 24.04.3 LTS
- **Kernel**: 6.14.0-29-generic
- **Go Version**: go version go1.23.11 linux/amd64

**Dependencies**

- `github.com/consensys/gnark-crypto`: v0.18.0
- `github.com/Nik-U/pbc`: v0.0.0-20181205041846-3e516ca0c5d6

**Test Runtime Environment**

- **Machine**: Physical laptop, not running in a virtual machine.
- **System Load**: No other high-load applications were running during the tests.

## Library - gnark

> Implemented based on the `gnark-crypto` library: https://github.com/Consensys/gnark-crypto
>
> Benchmarks are performed on the BN254 and BLS12381 curves.

> Note: Operations here are performed using `G1Affine` and `G2Affine`. In practice, for a large number of scalar multiplications or point additions, it is more efficient to convert points to Jacobian coordinates (`G1Jac`, `G2Jac`) for computation and then convert them back to Affine coordinates.
>
> | Dimension        | Affine (x,y)                               | Jacobian (X:Y:Z)                                     |
> | ---------------- | ------------------------------------------ | ---------------------------------------------------- |
> | Coordinates      | 2                                          | 3                                                    |
> | Operation Cost   | Inversion required for each add/double     | Most operations only need mul/add; 1 inversion at the end |
> | Serialization    | Most compact, standard format (compressed/uncompressed) | Cannot be directly transmitted; must convert back to Affine |
> | Typical Use Case | Storage, public key/signature exchange, I/O validation | Internal computation, batch scalar multiplication    |
>
> A common operational pattern is:
>
> *   **Generate/Read Public Key** → Affine;
> *   **Perform Scalar Multiplication/Multi-point Operations** → Convert to Jacobian for fast internal computation;
> *   **Output/Transmit Result** → Convert back to Affine.

> In the `gnark-crypto` library, $\mathbb{G}_1$ and $\mathbb{G}_2$ are implemented as additive groups in `g1.go` and `g2.go`, respectively. $\mathbb{G}_T$ is implemented in `e12.go`. This `e12` itself is a field, providing functions for both addition and multiplication. For pairings, we use the multiplicative set of functions. [Addition: `Add`, `Sub`; no functions for scalar multiplication or additive inverse are provided. Multiplication: `Mul`, `Div`, `Inverse`, `Exp`].

Test results are as follows:

| Operation Type                                   | BN254      | BLS12381   |
| ------------------------------------------------ | ---------- | ---------- |
| $\mathbb{F}^+$, `Add()`                          | 25 ns      | 25 ns      |
| $\mathbb{F}^-$, `Sub()`                          | 25 ns      | 25 ns      |
| $\mathbb{F}^{neg}$, `Neg()`                      | 23 ns      | 20 ns      |
| $\mathbb{F}^{\times}$, `Mul()`                   | 36 ns      | 36 ns      |
| $\mathbb{F}^{\div}$, `Div()`                     | 2.59 µs    | 2.57 µs    |
| $\mathbb{F}^{inv}$, `Inverse()`                  | 2.534 µs   | 2.524 µs   |
|                                                  |            |            |
| $\mathbb{G}_1^+$, `Add()`                        | 2.712 µs   | 4.606 µs   |
| $\mathbb{G}_1^-$, `Sub()`                        | 2.057 µs   | 3.819 µs   |
| $\mathbb{G}_1^{neg}$, `Neg()`                    | 25 ns      | 27 ns      |
| $\mathbb{G}_1^{\cdot}$, `ScalarMultiplication()` | 54.534 µs  | 87.697 µs  |
|                                                  |            |            |
| $\mathbb{G}_2^+$, `Add()`                        | 3.009 µs   | 5.301 µs   |
| $\mathbb{G}_2^-`, `Sub()`                        | 2.338 µs   | 4.464 µs   |
| $\mathbb{G}_2^{neg}$, `Neg()`                    | 26 ns      | 27 ns      |
| $\mathbb{G}_2^{\cdot}$, `ScalarMultiplication()` | 103.125 µs | 191.075 µs |
|                                                  |            |            |
| $\mathbb{G}_T^{\times}$, `Mul()`                 | 1.365 µs   | 2.301 µs   |
| $\mathbb{G}_T^{\div}$, `Div()`                   | 5.177 µs   | 8.846 µs   |
| $\mathbb{G}_T^{inv}$, `Inverse()`                | 3.502 µs   | 6.167 µs   |
| $\mathbb{G}_T^{\exp}$, `Exp()`                   | 380.91 µs  | 636.181 µs |
|                                                  |            |            |
| $e(\cdot,\cdot)$, `Pair()`                       | 381.444 µs | 583.345 µs |

## Library - pbc

> Implemented based on the `pbc` library: https://github.com/Nik-U/pbc
>
> Benchmarks are performed on Type A, Type D, and Type F curves, with parameters configured according to the official documentation: https://pkg.go.dev/github.com/Nik-U/pbc

> Note: For scalar multiplication on $\mathbb{G}_1$ and $\mathbb{G}_2$, the `MulZn()` function should be used since these are additive groups. However, in the code, the `PowZn()` function is also interpreted as "repeated group operation," so using `PowZn()` on $\mathbb{G}_1$ and $\mathbb{G}_2$ yields the same result as `MulZn()`. Nevertheless, we do not recommend using `PowZn()` on $\mathbb{G}_1$ and $\mathbb{G}_2$, as this notation can be confused with group exponentiation in $\mathbb{G}_T$, harming code readability and maintenance. The recommendation is to use `Add(), Sub(), Neg(), MulZn()` for $\mathbb{G}_1, \mathbb{G}_2$ and `Mul(), Div(), Invert(), PowZn()` for $\mathbb{G}_T$.

Test results are as follows:

| Operation Type                    | TypeA(160,512) | TypeD(9563, 160, 171, 500) | TypeF(160)   |
| --------------------------------- | -------------- | -------------------------- | ------------ |
| $\mathbb{F}^+$, `Add()`           | 216 ns         | 223 ns                     | 196 ns       |
| $\mathbb{F}^-$, `Sub()`           | 143 ns         | 145 ns                     | 153 ns       |
| $\mathbb{F}^{neg}$, `Neg()`       | 128 ns         | 134 ns                     | 132 ns       |
| $\mathbb{F}^{\times}$, `Mul()`    | 200 ns         | 195 ns                     | 186 ns       |
| $\mathbb{F}^{\div}$, `Div()`      | 1.462 µs       | 1.512 µs                   | 1.427 µs     |
| $\mathbb{F}^{inv}$, `Invert()`    | 1.259 µs       | 1.347 µs                   | 1.243 µs     |
|                                   |                |                            |              |
| $\mathbb{G}_1^+$, `Add()`         | 4.865 µs       | 3.179 µs                   | 2.576 µs     |
| $\mathbb{G}_1^-$, `Sub()`         | 4.304 µs       | 2.648 µs                   | 1.85 µs      |
| $\mathbb{G}_1^{neg}`, `Neg()`     | 179 ns         | 268 ns                     | 174 ns       |
| $\mathbb{G}_1^{\cdot}$, `MulZn()` | 836.449 µs     | 463.295 µs                 | 325.381 µs   |
|                                   |                |                            |              |
| $\mathbb{G}_2^+$, `Add()`         | 4.271 µs       | 22.044 µs                  | 3.542 µs     |
| $\mathbb{G}_2^-`, `Sub()`         | 4.259 µs       | 19.792 µs                  | 3.031 µs     |
| $\mathbb{G}_2^{neg}`, `Neg()`     | 155 ns         | 272 ns                     | 172 ns       |
| $\mathbb{G}_2^{\cdot}$, `MulZn()` | 824.753 µs     | 3.789244 ms                | 665.947 µs   |
|                                   |                |                            |              |
| $\mathbb{G}_T^{\times}$, `Mul()`  | 786 ns         | 6.078 µs                   | 16.795 µs    |
| $\mathbb{G}_T^{\div}`, `Div()`    | 4.509 µs       | 17.531 µs                  | 83.905 µs    |
| $\mathbb{G}_T^{inv}`, `Invert()`  | 3.278 µs       | 10.803 µs                  | 63.739 µs    |
| $\mathbb{G}_T^{\exp}$, `PowZn()`  | 74.586 µs      | 875.686 µs                 | 2.73738 ms   |
|                                   |                |                            |              |
| $e(\cdot,\cdot)$, `Pair()`        | 523.865 µs     | 2.786758 ms                | 11.742492 ms |
