依赖项目：

```
git clone https://github.com/unichainplatform/unichain.git
git clone https://github.com/spf13/jwalterweatherman.git
```

根据rpc功能划分了不同的域：
- accountAPI: 包括账户、资产操作等功能
- dposAPI: 包括注册候选者，投票等功能
- blockAPI: 包括区块查询等功能
- transactionAPI: 包括交易查询功能
- p2pAPI: 包括增加、删除节点等功能
- utils: 包括获取nonce，发送交易等功能
> rpc文档：https://github.com/unichainplatform/unichain/wiki/JSON-RPC
