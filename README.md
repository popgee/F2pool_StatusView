# F2pool_StatusView 简介
简单的鱼池矿工状态监控。通过Get方法请求F2pool的api，将获取到的数据保存到变量中，网页访问时将变量的内容替换到网页中。  
> *第一次接触go，语法水平低下，见谅。*

## 使用方法：
1. 将f2status.go文件中第72行的`***username_or_WalletAddress***`替换为F2Pool的用户名或钱包地址。
    * 例如：将`https://api.f2pool.com/ethereum/***username_or_WalletAddress***`修改为`https://api.f2pool.com/ethereum/abc`或`https://api.f2pool.com/ethereum/0x123456abcdefg`
2. 以Linux服务器为例：
    * 首先请自行安装Go环境
3. 将文件上传到某个目录下，路径：
```
└─F2Pool_StatusView
    │─f2status.go
    └─web
        └─index.html
```
4. 进入到`F2Pool_StatusView`目录下，执行`go build`进行编译
5. 执行`./f2status`来运行生成的文件
6. 浏览器输入`http://服务器ip:8000`来访问网页
