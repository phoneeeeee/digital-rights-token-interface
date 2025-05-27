# 🚀  一种媒体版权及溯源管理区块链技术架构及数字权利通证接口规范（参考代码）


## ✅ 功能介绍

- :white_check_mark: **支持多用户操作**：允许多个用户在链上进行通证的注册、查询和管理。
- :lock: **数据加密存储**：所有通证相关数据在区块链上加密存储，确保隐私和安全。
- :bar_chart: **实时数据展示**：用户可随时查询通证的详细信息，包括发行状态、持有者、冻结状态等。
- 🚀 **灵活的通证管理**：支持通证的初始化、发行、授权、冻结/解冻、信息修改等全生命周期管理。
- :heavy_plus_sign: **多签约束功能**：对于关键操作（如修改约束），支持多签名验证，增强安全性。
- :memo: **详细信息记录**：完整记录通证的确权信息、授权信息、权利主体组等详细数据。

## 🛠️ 技术栈

- :rocket: **开发平台**：ChainMaker SmartIDE
- :scroll: **开发语言**：GO 1.2.1
- 🚀 **区块链**：长安链 v2.1.0+
- :file_folder: **核心库**：chainmaker/shim、protogo

## 📖 代码结构介绍

### 📁 合约对象
```go
type TokenContract struct{}
```
- :white_check_mark: 定义智能合约对象，是所有通证操作的入口。

### 📁 存证对象
```go
type Fact struct {
    FileHash string
    FileName string
    Time     int
}
```
- :memo: 用于存储文件的哈希值、文件名和时间戳，确保文件的真实性和不可篡改。

### 📁 通证初始化结构
```go
type TokenIssue struct {
    Account       string      `json:"account"`
    Publisher     string      `json:"publisher"`
    Token         string      `json:"token"`
    Number        int         `json:"number"`
    Flag          int         `json:"flag"`
    Version       string      `json:"version"`
    Roles         []TokenRole `json:"roles,omitempty"`
    ReferenceFlag int         `json:"reference_flag"`
}
```
- :heavy_plus_sign: 包含通证的基本信息，如账户、发行者、通证名称、数量、标志位、版本等。

### 📁 通证许可交易数据
```go
type PubTokenTx struct {
    Publisher   string      `json:"publisher"`
    Receiver    string      `json:"receiver"`
    Token       string      `json:"token"`
    TokenId     string      `json:"tokenId"`
    ReferenceId string      `json:"referenceId"`
    TokenInfos  []TokenInfo `json:"tokenInfos,omitempty"`
    OwnerSigned bool        `json:"ownerSigned"`
    OwnerAccount string     `json:"ownerAccount"`
    OwnerSignature string   `json:"ownerSignature"`
}
```
- :memo: 存储通证许可交易的详细信息，包括发布者、接收者、通证ID、引用ID和签名状态等。

### 📁 合约初始化方法
```go
func (tc *TokenContract) InitContract(stub shim.CMStubInterface) protogo.Response {
    return shim.Success([]byte("TokenContract Init Success"))
}
```
- :white_check_mark: 初始化智能合约，确保合约在区块链上正确部署。

### 📁 通证初始化实现
```go
func (tc *TokenContract) BuildTokenIssueTx(stub shim.CMStubInterface) protogo.Response {
    // 实现逻辑...
}
```
- :heavy_plus_sign: 实现通证的初始化逻辑，包括参数校验、数据序列化和存储。

### 📁 通证发行方法
```go
func (tc *TokenContract) BuildPublishTokenTx(stub shim.CMStubInterface) protogo.Response {
    // 实现逻辑...
}
```
- :scroll: 处理通证的发行逻辑，支持一般通证和授权通证的发行。

### 📁 通证信息查询
```go
func (tc *TokenContract) RequestTokenInfo(stub shim.CMStubInterface) protogo.Response {
    // 实现逻辑...
}
```
- :mag: 提供通证信息查询功能，用户可按通证ID查询详细信息。

### 📁 通证信息修改
```go
func (tc *TokenContract) BuildModifyCopyrightTokenFlagTx(stub shim.CMStubInterface) protogo.Response {
    // 实现逻辑...
}
```
- :pencil: 允许修改通证的状态（如冻结/解冻）和相关信息。

## 🛠️ 使用指南

1. **初始化合约**：部署合约后，调用`InitContract`方法进行初始化。
2. **通证初始化**：通过`BuildTokenIssueTx`方法初始化通证，设置通证的基本属性。
3. **通证发行**：使用`BuildPublishTokenTx`或`BuildPublishApproveTokenTx`方法发行通证。
4. **通证查询**：利用`RequestTokenInfo`方法查询特定通证的信息。
5. **通证管理**：使用各类`BuildModify...Tx`方法管理通证的状态和信息。

## 🤝 社区参与

- :star: 如果你喜欢这个项目，请给它一个星标！
- :memo: 提交你的改进建议或问题报告到Issues。
- :package: 欢迎提交Pull Request，共同完善项目。

## 📝 开发者文档

所有方法的详细文档和参数说明，请参考[API文档](https://docs.chainmaker.org.cn/)。
