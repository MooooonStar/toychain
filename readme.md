## 程序员眼中的区块链

### (一) 区块、交易和挖矿

#### 前言
最近区块链热度又起来了, 但多数人的理解也仅限于“区块链”这个词，或者一些高大上的概念，细究起来仍是不明所以。为搞懂这一概念，我查阅了一些文章和代码，试图从一些代码片段和底层数据结构上，来深入理解下区块链。

#### 区块链(blockchain)
区块链，英文为blockchain，其实是两个词，block和chain，分别是方块和链条的意思, 于是中文就叫做了区块链。我个人对于区块链的理解是： **区块链是一个个区块构成的链表，而区块是一些交易的集合**。

#### 区块(block)
先来看下什么是区块，以 ```0000000000000bae09a7a393a8acded75aa67e46cb81f7acaa5ad94f9eacd103``` 这个块为例，其基本结构如下。
```json5 
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
 - prev_block, 上一个块的哈希。学过C的应该知道，这表明区块链的数据结构是链表，并且是个前向列表。从任意一个节点的位置可以一直往前遍历直到表头，这也是区块链号称可以追根溯源的原因。
 - mrkl_root, 默克尔树根，为交易列表的哈希。
 - time, 区块生成的时间，unix时间戳
 - bits, 难度
 - nonce, 随机数，和难度一起定义了区块链链表增长的规则，即设置了一个规则，只有满足条件的块才能被添加到链上去。
 - tx, 交易列表

所以，我们可以这么理解一个区块： **区块就是一些交易数据的集合，并且区块头必须满足一定的规则**。

(注: mrkl_root, bits, nonce 字段的介绍见后文)

#### 交易

再来看下什么是交易，以`b6f6991d03df0e2e04dafffcd6bc418aac66049e2cd74b80f14ac86db1e3f0da`这笔交易为例

```json5
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

这里的数据结构，和常规的交易很不一样。常规的中Alice向Bob转账了200，表现形式应该是这样的

```json5
{
  "from": "Alice",
  "to": "Bob",
  "amount": 200
}

```
而这里完全没有类似的字段。

现实生活中的交易采用的是账户模型，比如Alice、Bob账户上起初余额分别是500,100, Alice向Bob转账200，即在数据库里起个事务,把Alice账户扣除200,同时把Bob账户增加200，一笔交易就算完成了。而BTC采用了完全不同的模型，即UTXO(unspent transaction outputs)模型，称为"未花费的交易输出"。要理解这一概念，首先得先理解交易输出和交易输入。

##### 交易输出

交易输出包含两部分:

- value, 数量，单位聪(satoshis), 1BTC = 10^8 satoshis
- scriptPubKey, 锁定脚本。以`76a91429d6a3540acfa0a950bef2bfdc75cd51c24390fd88ac`为例，它实际上是
  
``` OP_DUP OP_HASH160 29d6a3540acfa0a950bef2bfdc75cd51c24390fd OP_EQUALVERIFY OP_CHECKSIG```

其中十六进制部分再处理下就是地址 `14pDqB95GWLWCjFxM4t96H2kXH7QMKSsgG`(当然这里涉及到公钥和BTC地址的转换关系，后文会介绍)，剩余的部分是定义的一些运算，简单得理解为加减乘除就好了。

简而言之，**交易输出就是在地址上放一定数量的BTC和锁定脚本， 需要解锁脚本才能使用**，形象地可以理解为在房间里放了东西，并且上了把锁。

##### 交易输入

交易输入包含三部分:
    
- txid: 引用的交易id.
- vout: 引用交易的输出索引。如一笔交易产生了5个交易输出，这里vout可以是0,1,2,3,4.
- scriptSig: 解锁脚本，包含一个签名和公钥，和锁定脚本拼起来会是如下形式
  
``` <Signature> <PublicKey> | DUP HASH160 <PubKeyHash> EQUALVERIFY CHECKSIG```

BTC里面的支付称为对公钥哈希的付款(P2PKH, Pay-to-Public-Key-Hash),原因就在于此。

交易A的输出里会放一个公钥的哈希，掌握私钥的人可以引用交易A, 并在交易B的输入里用私钥对交易B签名，同时提供要转出的地址。而上面的脚本拼起来其实做了两件事，一是先对输入的公钥进行哈希，看是否和输出的一致，相当于先验证你要转账的人有没有填错，二是用公钥验证签名对不对，即再验证使用者是否拥有使用权。

总而言之，**交易输入即引用交易输出，并签名使用**。而所谓的交易，就是这样一个打开一个房间，取走东西，放入另一个房间的过程。

##### 创币交易

交易输入和交易输出会存在一个鸡生蛋还是蛋生鸡的问题，没有输入显然就没有输出，而输出又需要引用输入。

事实上是，每个区块的第一个是个创币交易(coinbase)，只有输出没有输入，也就是不需要解锁脚本，该字段会被coinbase数据代替。如创世块中，中本聪就填入了"The Times 03/Jan/ 2009 Chancellor on brink of second bailout for banks".创币交易的数量是被给定的规则限定的，经常提到的BTC产量每4年减半，就是指创币交易的数量，而输出会被矿工光明正大且合法地填上自己的地址。由创币交易就可以衍生其他交易，保证链条不断增长下去。

##### 手续费

通常一笔普通交易会有两个去处，一个是待转账的地址，一个是自己的地址，因为一笔交易的输入之和不会恰好等于输出，于是会存在一个指向自己的找零。输入之和是大于输出之和的，中间的差值会被系统奖励给矿工，也就是交易手续费。手续费高的区块会优先被打包成区块得到确认，并且连接到链上去。

#### 挖矿

新创建的交易会进入内存池，等待矿工打包成区块。而所谓的矿工, 就是打包交易的节点。挖矿主要和区块头的哈希有关。

##### 区块头
 
区块头包含六个字段,Version、Previous Block Hash、 Merkle Root、 Timestamp、 Target、 Nonce。其中前五个字段可以看做是相对固定, 而最后的Nonce值则是矿工要完成的工作。

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

默克尔树可以快速验证区块中是否存在某笔交易，并且可以快速校验区块交易数据的一致性。平时在下载文件时，会有一个哈希校验值，原理和作用与此类似。

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

我们每次改变了字符串最后的数字，每次结果都千差万别，看起来像是完全随机。这就是哈希函数的作用，输出均匀且单向不可逆。现在我们设定一个目标，找到一个数字，使得哈希后的十六进制数以0开头。在多增加几次循环后，我们找到了目标，即

```
I am Satoshi Nakamoto13 => 0ebc56d59a34f5082aaef3d66b37a661696c2b618e62432727216ba9531041a5
```

从概率上讲，输出是均匀分布的，而十六进制数为0-F, 所以平均十六次试验可以找到1个以0开头的哈希值。即小于0x1000000000000000000000000000000000000000000000000000000000000000的值。如果我们不断的减小这个阈值，那么平均试验次数将以16的指数爆炸性增长。

因而，可以从实现目标的难度估算出所需的工作量。当算法是基于诸如SHA256的确定性函数的时候，输入本身就成为证据，必须要以一定的工作量才能产生低于目标的结果。因此，称之为工作量证明。

区块中难度用系数/指数的形式表示，在高度为277,316的区块中，难度以0x1903a30c为例，其中0x19为幂，0x03a30为系数，代表的难度目标为
```
target = 22,829,202,948,393,929,850,749,706,076,701,368,331,072,452,018,388,575,715,328
```
转为为十六进制后为```0x0000000000000003A30C00000000000000000000000000000000000000000000``` ,其二进制数中前60位都是0, 平均需计算16^60=2^240次。

一旦某个矿工找到了满足条件的Nonce值，便可以将区块头和交易列表打包成区块，添加到区块链里去。而打包这个区块的矿工,就获得了创币的特权，并且得到了交易中包含的手续费。

所以我们看到挖矿实际上是通过工作量证明将交易打包成区块，并链接到链上使得交易得到确认的过程，是BTC里不可缺少的一部分。没有了挖矿，所有的交易都只会留着节点的内存里。

#### 私钥和地址

BTC里的私钥其实就是一个256bit长的随机数字，然后通过椭囿曲线乘法可以由私钥得到公钥。然后对公钥进行SHA256哈希运算，再进行ripemd160哈希运算，结果前面加上版本，后面加上校验和，最后使用base58编码就得到了形如`17XLQvEM5uKPvuFPqfN8op2GQ6zs51Rqjv`的地址。

这里的公私钥和常用的RSA加密性质上是一样的，地址呢就是公钥先哈希为256bit，再哈希为160bit，加上版本和校验和，最后编码成不容易认错的base58的形式。


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
	// 创世块, 奖励100
	addr := "17XLQvEM5uKPvuFPqfN8op2GQ6zs51Rqjv"
	message := "Long Live The Bitcoin"
	bc := core.NewBlockchain(addr, message)

	// A挖到了第一个块, 奖励100
	keyA := core.NewKeyPair()
	tx0 := core.NewCoinbaseTx(keyA.Address(), "A got it")
	block1 := core.NewBlock([]*core.Transaction{tx0}, bc.Current.Hash())
	nonce1 := core.ProofOfWork(*block1)
	block1.Nonce = nonce1
	bc.AddBlock(block1)

	// B挖到了第二个块，并且A向C矿工转账了6。 A: 94, B: 100, C: 6
	keyB, keyC := core.NewKeyPair(), core.NewKeyPair()
	tx1 := core.NewCoinbaseTx(keyB.Address(), "B got it")
	tx2, err := core.NewTransaction(keyC.Address(), 6, keyA, bc)
	if err != nil {
		log.Fatal(err)
	}
	block2 := core.NewBlock([]*core.Transaction{tx1, tx2}, bc.Current.Hash())
	nonce2 := core.ProofOfWork(*block2)
	block1.Nonce = nonce2
	bc.AddBlock(block2)

	// D挖到了第三个块, C->D 2, A->B 4,  A: 90, B: 104, C: 4, D: 102
	keyD := core.NewKeyPair()
	tx3 := core.NewCoinbaseTx(keyD.Address(), "D got it")
	tx4, err := core.NewTransaction(keyD.Address(), 2, keyC, bc)
	if err != nil {
		log.Fatal(err)
	}
	tx5, err := core.NewTransaction(keyB.Address(), 4, keyA, bc)
	if err != nil {
		log.Fatal(err)
	}
	block3 := core.NewBlock([]*core.Transaction{tx3, tx4, tx5}, bc.Current.Hash())
	nonce3 := core.ProofOfWork(*block3)
	block3.Nonce = nonce3
	bc.AddBlock(block3)

	// 显示区块的结构
	bt, _ := prettyjson.Marshal(bc.GetBlocks())
	log.Println(string(bt))

	// 查看余额是否为 A: 90, B: 104, C: 4, D: 102
	for _, key := range []*core.KeyPair{keyA, keyB, keyC, keyD} {
		_, err := core.NewTransaction(addr, 1000, key, bc)
		if err != nil {
			log.Println("show balance:", err)
		}
	}
}
```

输出为:

```json5
[
  // 创世块
  {
    "MerkleRoot": "6efadffa729e5290a2420b46aabf5d712f2f52750894facc8827d3da3be77322",
    "Nonce": 0,
    "PrevBlock": "0000000000000000000000000000000000000000000000000000000000000000",
    "Target": 8,
    "Timestamp": 1573910049,
    "Transactions": [
      {
        "ID": "e10de29ba6cccf2155c8e1535d8613280e3e660509c9e692980d01ffbf01b222",
        "Vin": [
          {
            "PubKey": "TG9uZyBMaXZlIFRoZSBCaXRjb2lu",  // base64解码即为 Long Live The Bitcoin
            "Signature": null,
            "TxID": "0000000000000000000000000000000000000000000000000000000000000000",
            "Vout": -1
          }
        ],
        "Vout": [
          {
            "PubKeyHash": "R41wLUp70Wr7+Lfm7rlQ8aNEG4A=",
            "Value": 100
          }
        ]
      }
    ],
    "Version": 0
  },
  {
    "MerkleRoot": "161500bed1f416de91231b4dc622f81d7bc8fee3aae504bf12b7f7e0e14db506",
    "Nonce": 249,
    "PrevBlock": "554283088c0f8a9cea43bfdcaa064cf4dd832d243463a92eb37b1cdeee7d676e",
    "Target": 8,
    "Timestamp": 1573910049,
    "Transactions": [
      // 创币交易
      {
        "ID": "4255ec21fb59ff25f2d87ac07a727b68f547c2ca6bad900f2699e620e2ea5089",
        "Vin": [
          {
            "PubKey": "QSBnb3QgaXQ=",
            "Signature": null,
            "TxID": "0000000000000000000000000000000000000000000000000000000000000000",
            "Vout": -1
          }
        ],
        "Vout": [
          {
            "PubKeyHash": "nCf+aqgHXJa2MmkOiUlhnmwo8XE=",
            "Value": 100
          }
        ]
      }
    ],
    "Version": 0
  },
  {
    "MerkleRoot": "9ad6460a770dafd9cfe077cd67f6a0c3715c056c82a26fbb10d7a8bb034ade1a",
    "Nonce": 249,
    "PrevBlock": "63f6a51ad522adf384494c9ea6a85421c02c316048c43f73c11073a9faf822e7",
    "Target": 8,
    "Timestamp": 1573910049,
    "Transactions": [
      {
        "ID": "4399d728a35e72ee6674cef6eb9170a852d92d0835ed27b07a10773e4ddb17b7",
        "Vin": [
          {
            "PubKey": "QiBnb3QgaXQ=",
            "Signature": null,
            "TxID": "0000000000000000000000000000000000000000000000000000000000000000",
            "Vout": -1
          }
        ],
        "Vout": [
          {
            "PubKeyHash": "62TkdhVTHvx89l+js5noLZyQTBE=",
            "Value": 100
          }
        ]
      },
      {
        "ID": "ad1102af9d1a614aaf55842ae4e089c120dd7ca8317fbac94591ec0cc92a1a65",
        "Vin": [
          {
            "PubKey": "ErPS5vv/exUazYe9TlQZYIBQ09vATLMBKXFff+0hJG0oJ90gJR/vLmDEYEUCNMyRNg/ebJsxoWfPgJB3k/UVGw==",
            "Signature": "8FAGtjYIJ4CE8kATnTI29hh++uhM02eopa2TERNlltu9bJlEeAMkLiQEUoXnWZvM1Xe+diniDo48ojL4bz6qnA==",
            "TxID": "4255ec21fb59ff25f2d87ac07a727b68f547c2ca6bad900f2699e620e2ea5089",  //引用了第一个块中的第0笔交易
            "Vout": 0                                                                    //引用了交易的第0个输出
          }
        ],
        "Vout": [
          {
            "PubKeyHash": "oWN5N5JDcaLmzwqed6QXYxyAk+c=",
            "Value": 6
          },
          {
            "PubKeyHash": "nCf+aqgHXJa2MmkOiUlhnmwo8XE=",                               // 找零
            "Value": 94
          }
        ]
      }
    ],
    "Version": 0
  },
  {
    "MerkleRoot": "995a716c388bf697ed5bde1ce647f3f5e14babc6e0b0465c94e9411ad9bc0dfa",
    "Nonce": 398,
    "PrevBlock": "506f63e923797f8b1b84806588a188511e803402ff8df7682b0398653f10cc8d",
    "Target": 8,
    "Timestamp": 1573910049,
    "Transactions": [
      {
        "ID": "6ecae96d41670bc41fa4da632febb7ee3ba1a6eee67838a0d7ec97131fdbaeb2",
        "Vin": [
          {
            "PubKey": "RCBnb3QgaXQ=",
            "Signature": null,
            "TxID": "0000000000000000000000000000000000000000000000000000000000000000",
            "Vout": -1
          }
        ],
        "Vout": [
          {
            "PubKeyHash": "6foao7bLzoZS1mhzI/YiLc0899M=",
            "Value": 100
          }
        ]
      },
      {
        "ID": "a4da32136945e400cde1adfc69847918180a655941d53093fb84586f5a58756a",
        "Vin": [
          {
            "PubKey": "R5TWIEdq/fF+fsiIhwppJu1QSBMfC9RiyFZHzzHp0NpItYyagAg61EnPK76L01Q7fb32g0e71dBHJAeo4fXVbA==",
            "Signature": "/zC5vKvI7MLVhJAbFvDtJ8YrvyJcRMzpGruWg7xBbCXc873MyR/kAO6nAdjNzNGJ/BXNbog8tKdokIpXcFKYOg==",
            "TxID": "ad1102af9d1a614aaf55842ae4e089c120dd7ca8317fbac94591ec0cc92a1a65",         // 引用了第2个块中的第1笔交易
            "Vout": 0                                                                           // 引用了交易的第0个输出
          }
        ],
        "Vout": [
          {
            "PubKeyHash": "6foao7bLzoZS1mhzI/YiLc0899M=",
            "Value": 2
          },
          {
            "PubKeyHash": "oWN5N5JDcaLmzwqed6QXYxyAk+c=",
            "Value": 4
          }
        ]
      },
      {
        "ID": "ba5caa697b4b05e3799c3321e58c8fce58370b450d356052c9fb0eb8290b2760",
        "Vin": [
          {
            "PubKey": "ErPS5vv/exUazYe9TlQZYIBQ09vATLMBKXFff+0hJG0oJ90gJR/vLmDEYEUCNMyRNg/ebJsxoWfPgJB3k/UVGw==",
            "Signature": "2BEe63BWBUfiyvsWUmNylulEGxDdmlvHkwt6hrKs0oYQq4qlqJhC6Mny+JyAY6v93souYTWGniJQoEop2eLUWg==",
            "TxID": "ad1102af9d1a614aaf55842ae4e089c120dd7ca8317fbac94591ec0cc92a1a65",       // 引用了第2个块中的第1笔交易
            "Vout": 1                                                                         // 引用了交易的第1个输出
          }
        ],
        "Vout": [
          {
            "PubKeyHash": "62TkdhVTHvx89l+js5noLZyQTBE=",
            "Value": 4
          },
          {
            "PubKeyHash": "nCf+aqgHXJa2MmkOiUlhnmwo8XE=",
            "Value": 90
          }
        ]
      }
    ],
    "Version": 0
  }
]
```

#### 最后
1. 代码参考了 ```https://github.com/Jeiwan/blockchain_go```， 如有雷同，算我抄他
2. 书籍参考了《精通比特币（第二版)》
3. 打赏: ```17XLQvEM5uKPvuFPqfN8op2GQ6zs51Rqjv```
4. TODO: 加入P2P网络

