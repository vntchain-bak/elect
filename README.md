## Elect

Elect是VNT Chain公链Hubble NetWork见证人选举的命令行工具。

Elect并非一个区块链节点，Elect通过RPC连接到配置文件指定的Hubble network节点，查询信息、创建选举合约的底层交易，对交易签名后发送给网络节点，交易上链后，如果执行成功则完成了选举操作，操作的结果可以通过`elect query`命令查询。

Elect支持以下功能：
- 抵押和取消抵押VNT
- 注册和注销为见证人节点
- 为见证人节点投票和取消投票
- 账号开启和关闭投票代理功能，即成为和退出代理人
- 账号设置和取消设置其他账号为代理人
- 查询功能
  - 本用户抵押、投票信息
  - 见证人列表
  - 剩余VNT激励总量



## 安装

下载：

    mkdir -p $GOPATH/src/github.com/vntchain
    cd $GOPATH/src/github.com/vntchain
    git clone https://github.com/vntchain/elect
    cd elect

编译&安装：

    make install

## 运行

可以使用帮助命令查看elect所支持的功能：
    
    elect help
    
所支持功能的命令下：

    cancelProxy 取消投票代理
    cancelVote  取消对见证人的投票
    query       查询命令支持：抵押、投票、见证人列表
    register    注册成为见证人
    setProxy    设置某账户为代理自己投票
    stake       抵押代币
    startProxy  成为投票代理人
    stopProxy   退出投票代理人，不再代理其他人投票
    unregister  注销见证人
    unstake     取回抵押代币
    vote        为见证人投票，最多投30个见证人

运行命令前需要做3件事：

1. 创建工具运行目录

    ```
    cd ~
    mkdir vnt_elect && cd vnt_elect
    ```

2. 创建keystore目录，并把你的账户keystore文件放到keystore目录

    ```
    mkdir keystore
    cp path/to/your/keystore/file keystore
    ```

3. 设置配置文件[`config.json`](./config.json)

    ```json
    {
        "sender":"0x3dcf0b3787c31b2bdf62d5bc9128a79c2bb18829",
        "password":"",
        "keystoreDir":"./keystore",
        "rpcUrl":"http://localhost:8880",
        "chainID":0
    }
    ```
    
    需要在config.json替换为你当前的配置：
    - sender：要参与投票的账户地址，要与keytore文件的账户地址一样
    - password：账户的密码
    - keystoreDir：keystore文件所在的目录，即`./keystore`，你可以省略第2步，把你的keystore目录填写在此即可
    - rpcUrl：VNT网络上的任何开启RPC服务的节点的RPC URL（IP+端口），如果你本地运行了go-vnt节点，则填写`http://localhost:8880`
    - chainID：默认为0，即VNT Chain公链网络Hubble，如果你搭建了测试网，请填写你搭建网络chainID

## 文档

elect不仅是一个命令行工具，还可以作为package使用，接口文档请查看[这里](https://godoc.org/github.com/vntchain/elect)。

下面是查询抵押代币的样例：

```go
package main

import (
	"fmt"

	"github.com/vntchain/elect"
)

func main() {
	var (
		err error
		e   *elect.Election
	)

	if e, err = elect.NewElection("./config.json"); err != nil {
		panic(err)
	}

	if ret, err := e.QueryStake(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(ret))
	}
}

```


## 许可证

所有`elect`仓库生成的二进制程序都采用GNU General Public License v3.0许可证, 具体请查看[COPYING](./COPYING)。