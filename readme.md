## 程序员眼中的区块链

### (一) 区块、交易和挖矿

#### 前言
最近区块链热度又起来了, 但多数人的理解也仅限于“区块链”这个词，或者一些高大上的概念，细究起来

仍是不明所以。为搞懂这一概念，我查阅了一些文章和代码，试着从一些代码片段和数据结构上，来重新

认识一下区块链。

#### 区块链(blockchain)
区块链，英文为blockchain，其实是两个词，block和chain，分别是方块和链条的意思, 于是中文就叫做

了区块链。我个人对于区块链的定义是： **区块链是一个个区块构成的链表，而区块是一些交易的集合**。

#### 区块(block)
先来看下什么是区块，以 ```0000000000000bae09a7a393a8acded75aa67e46cb81f7acaa5ad94f9eacd103``` 这个块为例，

其基本结构如下。
```json 
{
    "ver":1,
    "prev_block":"00000000000007d0f98d9edca880a6c124e25095712df8952e0439ac7409738a",
    "mrkl_root":"935aa0ed2e29a4b81e0c995c39e06995ecce7ddbebb26ed32d550a72e8200bf5",
    "time":1322131230,
    "bits":437129626,
    "nonce":2964215930,
    "tx":[--Array of Transactions--]
}
```
 其中各字段的含义：
 - ver, 版本
 - prev_block, 上一个块的哈希。学过C的应该知道，这表明区块链的数据结构是链表，并且是个前向列表。从任意一个节点的位置

可以一直往前遍历直到表头，这也是区块链号称可以追根溯源的原因。
 - mrkl_root, 默克尔树根，为交易列表的哈希。
 - time, 区块生成的时间，unix时间戳
 - bits, 难度
 - nonce, 随机数，和难度一起定义了区块链链表增长的规则，即设置了一个规则，只有满足条件的块才能被添加到链上去。
 - tx, 交易列表

所以，我们可以这么理解一个区块： **区块里包含了一些交易数据，并定义了区块增长的规则**。

(注: mrkl_root, bits, nonce 字段的介绍见后文。)

#### 交易

再来看下什么是交易，以`b6f6991d03df0e2e04dafffcd6bc418aac66049e2cd74b80f14ac86db1e3f0da`这笔交易为例

```json 
{
  "inputs": [
    {
      "txid": ,
      "vout": 2,
      "scriptSig": "48304502210098a2851420e4daba656fd79cb60cb565bd7218b6b117fda9a512ffbf17f8f178022005c61f31fef3ce3f906eb672e05b65f506045a65a80431b5eaf28e0999266993014104f0f86fa57c424deb160d0fc7693f13fce5ed6542c29483c51953e4fa87ebf247487ed79b1ddcf3de66b182217fcaf3fcef3fcb44737eb93b1fcb8927ebecea26"
    }
  ],
  "out": [
    {
      "value": "98000000",
      "scriptPubKey": "76a91429d6a3540acfa0a950bef2bfdc75cd51c24390fd88ac"
    },
    {
      "value": "2000000",
      "scriptPubKey": "76a91417b5038a413f5c5ee288caa64cfab35a0c01914e88ac"
    }
  ]
}
```

这种结果，和常规的交易很不一样。常规的, 比如Alice向Bob转账了200，结构应该是这样的

```json
{
  "from": "Alice",
  "to": "Bob",
  "amount": 200
}

```
而实际的交易结构里完全没有类似的字段。现实生活中采用的是账户模型，比如Alice、Bob账户上

起初余额分别是500,100, Alice向Bob转账200，即在数据库里起个事务,把Alice账户扣除200,

同时把Bob账户增加200，一笔交易就算完成了。而BTC采用了完全不同的模型，即

UTXO(unspent transaction outputs)模型，称为"未花费的交易输出"。要理解这一概念，

首先得先理解交易输入和交易输出。

##### 交易输出

交易输出包含两部分:

- value, 数量，单位聪(satoshis), 1BTC = 10^8 satoshis
- scriptPubKey, 锁定脚本。以`76a91429d6a3540acfa0a950bef2bfdc75cd51c24390fd88ac`为例，它实际上是
  
``` OP_DUP OP_HASH160 29d6a3540acfa0a950bef2bfdc75cd51c24390fd OP_EQUALVERIFY OP_CHECKSIG```

其中十六进制部分就是个地址 `14pDqB95GWLWCjFxM4t96H2kXH7QMKSsgG`(当然这里涉及到公钥和BTC地址的转换关系，留到后文介绍)，

剩余的部分是定义的一些运算，简单得理解为加减乘除就好了。

简而言之，**交易输出就是在某个地址上放一定数量的BTC和加密脚本， 需要解锁脚本才能使用**，即在房间里放了东西，并且上了把锁。

##### 交易输入

交易输入包含三部分:
    
- txid: 引用的交易id.
- vout: 引用交易的输出索引。如一笔交易产生了5个交易输出，这里vout可以是0,1,2,3,4.
- scriptSig: 解锁脚本，包含一个签名和公钥，和锁定脚本拼起来会是如下形式
  ``` <Signature> <PublicKey> | DUP HASH160 <PubKeyHash> EQUALVERIFY CHECKSIG

BTC里面的支付称为对公钥哈希的付款(P2PKH, Pay-to-Public-Key-Hash),原因就在于此。交易A的输出里会放一个公钥

的哈希，使用者需要引用交易A, 并在交易B的输入里用私钥对这笔交易B签名，同时提供要转出的地址。

上面拼起来的脚本其实做了两件事，一是先对输入的公钥进行哈希，看是否和输出的一致，相当于先验证你要转账的人有没有填错，

二是用公钥验证签名对不对，即验证你是否拥有使用权。

总而言之，**交易输入即引用交易输出，并签名使用**。也即打开自己的房间，取走东西，并放入另一个房间里。而所谓的交易，就是这样

一个打开一个房间，取走东西，放入另一个房间的过程。

##### 创币交易

交易输入和交易输出会存在一个鸡生蛋还是蛋生鸡的问题，没有输入显然就没有输出，而输出又需要引用输入。实际上是，

每个区块的第一个是个创币交易(coinbase)，只有输出没有输入，也就是不需要解锁脚本，该字段会被coinbase数据代替。

如创世块中，中本聪就填入了"The Times 03/Jan/ 2009 Chancellor on brink of second bailout for banks".

创币交易的数量是被给定的规则限定的，经常提到的BTC产量每4年减半，就是指创币交易的数量，而输出会被矿工光明正大且合法地

填上自己的地址。由创币交易就可以衍生出其他交易，保证链条不断增长下去。

##### 手续费

通常一笔普通交易会有两个去处，一个是待转账的地址，一个是自己的地址，因为一笔交易的输入之和不会恰好等于输出，

于是会存在一个指向自己的找零。输入之和是大于输出之和的，中间的差值会被系统奖励给矿工，也就是交易手续费。

手续费高的区块会优先被打包成区块，连接到链上去。

#### 挖矿

新创建的交易会进入内存池，等待矿工打包成区块。所谓的矿工, 就是打包交易的节点。

##### 区块头
 
 区块头包含六个字段,Version、Previous Block Hash、 Merkle Root、 Timestamp、 Target、 Nonce。

 ##### 默克尔树

 Merkle 树是一种哈希二叉树，假设区块里存在4笔交易Ta、Tb、Tc、 Td, 则其默克尔树是这样子的。
 ```
    Ta         Tb       Tc        Td
    |          |        |         |
Hash(Ta)  Hash(Tb)  Hash(Tc)   Hash(Td)
    |        |          |          |
    | _______|          |__________|  
         |                    |
     Hash(Tab)           Hash(Tcd)
         |                    |
         |____________________|  
                    |
              Hash(Tabcd)
```
即两两哈希，直到树根，若节点数目为奇数，则复制一个相同节点。

默克尔树可以快速验证区块中是否存在某笔交易，并且可以快速校验区块交易数据的一致性。平时在下载文件时，会有一个

哈希校验值，原理和作用与此类似。

##### 工作量证明 (Proof-Of-Work, POW)

简单点来说，挖矿就是遍历，使得区块头的哈希匹配的过程。为理解挖矿，可做如下试验。

在python里运行如下脚本
```python2.7
import hashlib

s = "I am Satoshi Nakamoto"
for i in range(0,10):
   t = s + str(i)
   print(t,hashlib.sha256(t).hexdigest())
```
可以得到结果
```
I am Satoshi Nakamoto0 => a80a81401765c8eddee25df36728d732acb6d135bcdee6c2f87a3784279cfaed
I am Satoshi Nakamoto1 => f7bc9a6304a4647bb41241a677b5345fe3cd30db882c8281cf24fbb7645b6240
I am Satoshi Nakamoto2 => ea758a8134b115298a1583ffb80ae62939a2d086273ef5a7b14fbfe7fb8a799e
I am Satoshi Nakamoto3 => bfa9779618ff072c903d773de30c99bd6e2fd70bb8f2cbb929400e0976a5c6f4
I am Satoshi Nakamoto4 => bce8564de9a83c18c31944a66bde992ff1a77513f888e91c185bd08ab9c831d5
I am Satoshi Nakamoto5 => eb362c3cf3479be0a97a20163589038e4dbead49f915e96e8f983f99efa3ef0a
I am Satoshi Nakamoto6 => 4a2fd48e3be420d0d28e202360cfbaba410beddeebb8ec07a669cd8928a8ba0e
I am Satoshi Nakamoto7 => 790b5a1349a5f2b909bf74d0d166b17a333c7fd80c0f0eeabf29c4564ada8351
I am Satoshi Nakamoto8 => 702c45e5b15aa54b625d68dd947f1597b1fa571d00ac6c3dedfa499f425e7369
I am Satoshi Nakamoto9 => 7007cf7dd40f5e933cd89fff5b791ff0614d9c6017fbe831d63d392583564f74
```

我们每次改变了字符串最后的数字，每次结果都千差万别，看起来像是完全随机。这就是哈希函数的作用，输出均匀且单向

不可逆。现在我们设定一个目标，找到一个数字，使得哈希后的十六进制数以0开头。

在多增加几次循环后，我们找到了目标，即
```
I am Satoshi Nakamoto13 => 0ebc56d59a34f5082aaef3d66b37a661696c2b618e62432727216ba9531041a5
```

从概率上讲，输出是均匀分布的，而十六进制数为0-F, 所以平均十六次试验可以找到1个以0开头的哈希值。即小于0x1000000000000000000000000000000000000000000000000000000000000000的值。如果我们不断的减小
这个阈值，那么平均试验次数将以16的指数爆炸性增长。

因而，可以从实现目标的难度估算出所需的工作量。当算法是基于诸如SHA256的确定性函数的时候，输入本身就成为
证据，必须要以一定的工作量才能产生低于目标的结果。因此，称之为工作量证明。

区块中难度用系数/指数的形式表示，在高度为277,316的区块中，难度以0x1903a30c为例，其中0x19为幂，0x03a30为系数，
代表的难度目标为
```
target = 22,829,202,948,393,929,850,749,706,076,701,368,331,072,452,018,388,575,715,328
```
转为为十六进制后为：
```0x0000000000000003A30C00000000000000000000000000000000000000000000``` ,其二进制数中前60位都是0，
平均需计算16^60=2^240次。

一旦某个矿工找到了满足条件的Nonce值，便可以将区块头和交易列表打包成区块，添加到区块链里去。而打包这个区块的矿工,

就获得了创币的特权，并且得到了交易中包含的手续费。

所以我们看到挖矿实际上是通过工作量证明将交易打包成区块，并链接到链上使得交易得到确认的过程，是BTC里不可缺少的一部分。

没有了挖矿，所有的交易都只会留着节点的内存里。

#### 私钥和地址

BTC里的私钥其实就是一个256bit长的随机数字，然后通过椭囿曲线乘法可以由私钥得到公钥。然后对公钥进行SHA256哈希运算，再进行

ripemd160哈希运算，结果前面加上版本，后面加上校验和，最后使用base58编码就得到了形如`17XLQvEM5uKPvuFPqfN8op2GQ6zs51Rqjv`

的地址。

这里的公私钥和常用的RSA加密性质上是一样的，地址呢就是公钥先哈希为256bit，再哈希为160bit，加上版本和校验和，最后编码

成不容易认错的base58的形式。


#### show me the code

理解里上面这些概念，现在我们就可以做个单机版的玩具链出来啦！

```golang
package main

import (
	"log"

	"github.com/MooooonStar/toychain/core"
	"github.com/hokaccha/go-prettyjson"
)

func main() {
  // 创世块
	addr := "17XLQvEM5uKPvuFPqfN8op2GQ6zs51Rqjv"
	key := core.NewKeyPair()
	to, message := key.Address(), "Long Live The Bitcoin"
	bc := core.NewBlockchain(to, message)

  //交易
	key1 := core.NewKeyPair()
	tx1, err := core.NewTransaction(key1.Address(), 6, key, bc)
	if err != nil {
		log.Panic(err)
	}

 //挖矿
	block1 := core.NewBlock([]*core.Transaction{tx1}, bc.Current.Hash())
	nonce1 := core.ProofOfWork(*block1)
	block1.Nonce = nonce1
	bc.AddBlock(block1)

  //交易和挖矿
	tx2, err := core.NewTransaction(addr, 2, key1, bc)
	if err != nil {
		log.Panic(err)
	}
	tx3, err := core.NewTransaction(addr, 3, key, bc)
	if err != nil {
		log.Panic(err)
	}
	block2 := core.NewBlock([]*core.Transaction{tx2, tx3}, bc.Current.Hash())
	nonce2 := core.ProofOfWork(*block2)
	block2.Nonce = nonce2
	bc.AddBlock(block2)

	bt, _ := prettyjson.Marshal(bc.GetBlocks())
	log.Println(string(bt))
}
```

输出为:

```json
[
  {
    "MerkleRoot": "aae3caa1e04e66cfee2f5101e1d1af6c8064834bf489f504047b17deb6a42943",
    "Nonce": 0,
    "PrevBlock": "0000000000000000000000000000000000000000000000000000000000000000",
    "Target": 16,
    "Timestamp": 1573872817,
    "Transactions": [
      {
        "ID": "2e0e92f7db707b2557e692d696d6365a0175a31e35fe98a03ee9c4c3354454cd",
        "Vin": [
          {
            "PubKey": "TG9uZyBMaXZlIFRoZSBCaXRjb2lu",          //base64解码后为"Long Live The Bitcoin"，即创世块中留下的信息
            "Signature": null,
            "TxID": "0000000000000000000000000000000000000000000000000000000000000000",
            "Vout": -1                                         // TxID和Vout为0, 创币交易不需要输入
          }
        ],
        "Vout": [
          {
            "PubKeyHash": "NdyIkMdzWmRpeL67cvyl9ERKUK8=",      //地址17XLQvEM5uKPvuFPqfN8op2GQ6zs51Rqjv的base58解码，去除前面的版本和后面的校验和
            "Value": 10
          }
        ]
      }
    ],
    "Version": 0
  },
  {
    "MerkleRoot": "5def27d49f27ec653e5589ee49c0712eea9f30e38f446f5a3c47cb227cb7c81c",
    "Nonce": 6063,                                             // 工作量证明
    "PrevBlock": "3e29183afd578491d12a3886bed0237854665c801bc6439b10c628f82cd26a76",
    "Target": 16,                                              // 难度
    "Timestamp": 1573872817,
    "Transactions": [
      {
        "ID": "0ecbf1a50a78790c1fba6c94df4a372548e509496865bfb4025a71377c5fe92b",
        "Vin": [
          {
            "PubKey": "5vTK/ovpDDnaVx62ft1gGabJDIg9Jc7FY88+LO3N2c3080mVMng3+AkxBzRxNmb9P5NUIpudwpjRMK/XgVUrYw==",
            "Signature": "jQugfLYTDoCRMN8dKIDeXgaBQjuAfcmFb+1mIOnTT9NTQMhi6Xq21wIL8AajnMytE2fj7pRi1GOPxv8Yc9X/oA==",
            "TxID": "2e0e92f7db707b2557e692d696d6365a0175a31e35fe98a03ee9c4c3354454cd",      // 引用创世块的交易
            "Vout": 0                                                                        // 引用交易的第一个输出,10
          }
        ],
        "Vout": [
          {
            "PubKeyHash": "BNCnUT7QA81EAEj4ZA6lKpHHhOg=",                                    //向A地址转账6
            "Value": 6
          },
          {
            "PubKeyHash": "NdyIkMdzWmRpeL67cvyl9ERKUK8=",                                    //给自己找零4
            "Value": 4
          }
        ]
      }
    ],
    "Version": 0
  },
  {
    "MerkleRoot": "59f0a8aa7f027eb80cb664c31b41f496ed310f5b79aaea432d57ead78a72a403",
    "Nonce": 33056,                                                                           //工作量证明
    "PrevBlock": "0000e76e73e357e0e88bea1e9cff21804fdc992c5ec4995a6caa2f2ba68cb8dc",
    "Target": 16,
    "Timestamp": 1573872817,
    "Transactions": [
      {
        "ID": "541de4b819e748cf24efc7064b45d2feb0da30fa767ab1e5284c85c8c6873c8b",
        "Vin": [
          {
            "PubKey": "5tbRH49GgwyJ7utRSe+ftXgKS8R2je9d6SoKRkQ4rSKJO7wJr1/mt8G4qJzaxXA82MYbIiE68h+pXgyKy0TLnw==",
            "Signature": "Fg5Xgb+58+f7dduzqj4Zzg03QEDgu4Gtm0wevlDjmvktIIB/UYNX4uL5XrDRPHpK0FAdWPTxwilDCVYBXPROAQ==",
            "TxID": "0ecbf1a50a78790c1fba6c94df4a372548e509496865bfb4025a71377c5fe92b",
            "Vout": 0
          }
        ],
        "Vout": [
          {
            "PubKeyHash": "R41wLUp70Wr7+Lfm7rlQ8aNEG4A=",
            "Value": 2
          },
          {
            "PubKeyHash": "BNCnUT7QA81EAEj4ZA6lKpHHhOg=",
            "Value": 4
          }
        ]
      },
      {
        "ID": "ec9d8d25d0a9fd47de7cc31c8859923260f0266374980f42a1429b80ad19eae7",
        "Vin": [
          {
            "PubKey": "5vTK/ovpDDnaVx62ft1gGabJDIg9Jc7FY88+LO3N2c3080mVMng3+AkxBzRxNmb9P5NUIpudwpjRMK/XgVUrYw==",
            "Signature": "C81Tun//47F7MfDTjAil7odINSh5RV7Mnf1wARiOozxifceufXLCnn72k7uBY8rhOpKH9MeLI1O7JJy4MBQgSg==",
            "TxID": "0ecbf1a50a78790c1fba6c94df4a372548e509496865bfb4025a71377c5fe92b",
            "Vout": 0
          }
        ],
        "Vout": [
          {
            "PubKeyHash": "R41wLUp70Wr7+Lfm7rlQ8aNEG4A=",
            "Value": 3
          },
          {
            "PubKeyHash": "NdyIkMdzWmRpeL67cvyl9ERKUK8=",
            "Value": 1
          }
        ]
      }
    ],
    "Version": 0
  }
]
```

