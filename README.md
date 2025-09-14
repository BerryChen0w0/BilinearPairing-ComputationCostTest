对双线性映射 $e:\mathbb{G}_1 \times \mathbb{G}_2 \to \mathbb{G}_T$ 中涉及的各种运算，测试其计算开销。运算包括：

* 域上运算：加减乘除和求逆元
* 群上运算：群加法、标量乘法（对乘法群就是群上乘法和群指数运算）
* 配对运算：pairing运算

使用常用的go语言代码库，对常用的椭圆曲线进行测试。



## 符号约定

代数结构 $[S,+,\times]$ ，对应的加法求逆元为 $neg$ ，逆运算为 $-$ ；乘法求逆元为 $inv$ ，逆运算为 $\div$ 。用 $\cdot$ 表示加法群的标量乘法， $\exp$ 表示乘法群的指数运算（标量乘法）。

在密码学论文中，一般把 $e:\mathbb{G}_1 \times \mathbb{G}_2 \to \mathbb{G}_T$ 中的 $\mathbb{G}_1$ 和 $\mathbb{G}_2$ 写作加法群， $\mathbb{G}_T$ 群写作乘法群。 $e(a\cdot P,b\cdot Q)=e(P,Q)^{ab}$ 。本测试使用的代码库中，也都是这样规定的。



## 测试环境

**硬件**

- **CPU**: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
- **内存**: 32 GB DDR4
- **存储**: 500 GB NVMe SSD

**软件**

- **操作系统**: Ubuntu 24.04.3 LTS
- **内核**: 6.14.0-29-generic
- **Go 版本**: go version go1.23.11 linux/amd64

**依赖**

- `github.com/consensys/gnark-crypto`: v0.18.0
- `github.com/Nik-U/pbc`: v0.0.0-20181205041846-3e516ca0c5d6

**测试运行环境**

- **物理机**: 笔记本电脑，未在虚拟机中运行。
- **系统负载**: 测试运行期间无其他高负载应用。



## 工具库-gnark

> 基于gnark-crypto库实现。https://github.com/Consensys/gnark-crypto
>
> 对bn254和bls12381两个曲线做测试。

> 注：这里都是用`G1Affine`和`G2Affine`来做运算。实际使用时，对存在大量标量乘法、点加法时，可以先转化为Jacobian点（`G1Jac`, `G2Jac`）进行运算，最后再转化为Affine点。
>
> | 维度     | Affine (x,y)                      | Jacobian (X:Y:Z)                           |
> | -------- | --------------------------------- | ------------------------------------------ |
> | 坐标数   | 2                                 | 3                                          |
> | 运算成本 | 每次加/倍点都要做求逆             | 大部分情况只需乘法/加法，最后再做 1 次求逆 |
> | 序列化   | 最紧凑，标准格式（压缩/非压缩点） | 不能直接传输，需要先转回仿射               |
> | 典型用途 | 存储、公钥/签名交换、验证输入输出 | 内部计算、批量标量乘                       |
>
> 常见的运算模式为：
>
> * **生成/读入公钥** → 仿射；
> * **做标量乘/多点运算** → 转成 Jacobian，内部快算；
> * **结果输出/传输** → 转回仿射。

> 在gnark-crypto库里面， $\mathbb{G}_1$ 和 $\mathbb{G}_2$ 分别在`g1.go`和`g2.go`里实现，实现为加法群。 $\mathbb{G}_T$ 在`e12.go`中实现，这个e12本身是一个域，实现了加和乘两套计算的函数。我们在pairing里面用乘法的那一套【加法：Add, Sub，没有提供标量乘法和求加法逆元的函数。乘法：Mul, Div, Inverse, Exp】



测试结果如下：

| 运算类型                                         | BN254      | BLS12381   |
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
| $\mathbb{G}_2^-$, `Sub()`                        | 2.338 µs   | 4.464 µs   |
| $\mathbb{G}_2^{neg}$, `Neg()`                    | 26 ns      | 27 ns      |
| $\mathbb{G}_2^{\cdot}$, `ScalarMultiplication()` | 103.125 µs | 191.075 µs |
|                                                  |            |            |
| $\mathbb{G}_T^{\times}$, `Mul()`                 | 1.365 µs   | 2.301 µs   |
| $\mathbb{G}_T^{\div}$, `Div()`                   | 5.177 µs   | 8.846 µs   |
| $\mathbb{G}_T^{inv}$, `Inverse()`                | 3.502 µs   | 6.167 µs   |
| $\mathbb{G}_T^{\exp}$, `Exp()`                   | 380.91 µs  | 636.181 µs |
|                                                  |            |            |
| $e(\cdot,\cdot)$, `Pair()`                       | 381.444 µs | 583.345 µs |



## 工具库-pbc

> 基于https://github.com/Nik-U/pbc库实现。
>
> 对TypeA，TypeD和TypeF曲线做测试，参照官方文档配置参数：https://pkg.go.dev/github.com/Nik-U/pbc

> 注：对 $\mathbb{G}_1$ 和 $\mathbb{G}_2$ 上的标量乘法，因为这两个群是加法群，因此应使用`MulZn()`函数。不过在代码中，因为`PowZn()`函数也被解释为“重复群运算”，因此，在 $\mathbb{G}_1, \mathbb{G}_2$ 上使用`PowZn()`函数，会获得与`MulZn()`一样的结果。但是，我们不建议在 $\mathbb{G}_1, \mathbb{G}_2$ 上使用`PowZn()`，因为这样的写法可能与 $\mathbb{G}_T$ 上的群指数运算混淆，不利于阅读和维护代码。建议是，在 $\mathbb{G}_1,\mathbb{G}_2$ 上使用`Add(), Sub(), Neg(), MulZn()`，在 $\mathbb{G}_T$ 上使用`Mul(), Div(), Invert(), PowZn()`。



测试结果如下：

| 运算类型                          | TypeA(160,512) | TypeD(9563, 160, 171, 500) | TypeF(160)   |
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
| $\mathbb{G}_1^{neg}$, `Neg()`     | 179 ns         | 268 ns                     | 174 ns       |
| $\mathbb{G}_1^{\cdot}$, `MulZn()` | 836.449 µs     | 463.295 µs                 | 325.381 µs   |
|                                   |                |                            |              |
| $\mathbb{G}_2^+$, `Add()`         | 4.271 µs       | 22.044 µs                  | 3.542 µs     |
| $\mathbb{G}_2^-$, `Sub()`         | 4.259 µs       | 19.792 µs                  | 3.031 µs     |
| $\mathbb{G}_2^{neg}$, `Neg()`     | 155 ns         | 272 ns                     | 172 ns       |
| $\mathbb{G}_2^{\cdot}$, `MulZn()` | 824.753 µs     | 3.789244 ms                | 665.947 µs   |
|                                   |                |                            |              |
| $\mathbb{G}_T^{\times}$, `Mul()`  | 786 ns         | 6.078 µs                   | 16.795 µs    |
| $\mathbb{G}_T^{\div}$, `Div()`    | 4.509 µs       | 17.531 µs                  | 83.905 µs    |
| $\mathbb{G}_T^{inv}$, `Invert()`  | 3.278 µs       | 10.803 µs                  | 63.739 µs    |
| $\mathbb{G}_T^{\exp}$, `PowZn()`  | 74.586 µs      | 875.686 µs                 | 2.73738 ms   |
|                                   |                |                            |              |
| $e(\cdot,\cdot)$, `Pair()`        | 523.865 µs     | 2.786758 ms                | 11.742492 ms |

