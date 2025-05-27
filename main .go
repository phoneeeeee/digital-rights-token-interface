package main

import (
	"encoding/json"
	"log"
	"math"
	"strconv"
	"chainmaker/pb/protogo"
	"chainmaker/shim"
)

// TokenContract 合约对象
type TokenContract struct{}

//Fact 存证对象,存证合约的数据内容
type Fact struct {
	FileHash string
	FileName string
	Time     int
}

//NewFact 新建存证对象
func NewFact(fileHash, fileName string, time int) *Fact {
	return &Fact{
		FileHash: fileHash,
		FileName: fileName,
		Time:     time,
	}
}

// TokenIssue 通证初始化结构体（5.1 通证初始化）
type TokenIssue struct {
	Account       string      `json:"account"`         // 动态发币账户
	Publisher     string      `json:"publisher"`       // ERC-721通证的发行账户地址
	Token         string      `json:"token"`           // ERC-721通证的发行名称
	Number        int         `json:"number"`          // ERC-721通证发行数量
	Flag          int         `json:"flag"`            // 标志位，0表示可以流通，1表示不可以流通
	Version       string      `json:"version"`         // 版本号，固定值："v1"或"v2"
	Roles         []TokenRole `json:"roles,omitempty"` // 控制token权限列表
	ReferenceFlag int         `json:"reference_flag"`  // 许可/通证标识: v1(0/1等), v2(1/2/3)
}

// TokenRole 描述 TokenIssue 中的 roles 数组的单项
type TokenRole struct {
	Role string `json:"role"` // 权限账户地址
	Type int    `json:"type"` // 权限标识(含义 v1/v2 在文档中不同，这里仅做示例)
}

// tokenObject 对应文档中一般通证发行时的参数结构
type TokenObject struct {
	Flag                int                   `json:"flag"`                // 0=可流通，1=不可流通
	TokenId             string                `json:"tokenId"`             // 版权通证的唯一标识符 (hash256)
	AuthenticationInfos []AuthenticationInfo  `json:"authenticationInfos"` // 确权信息列表(可选)
	CopyrightType       int                   `json:"copyrightType"`
	CopyrightGetType    int                   `json:"copyrightGetType"`
	CopyrightUnits      []CopyrightUnit       `json:"copyrightUnits"` // 权利主体组
	ConstraintExplain   string                `json:"constraintExplain"`
	ConstraintExpand    int                   `json:"constraintExpand"`
	CopyrightConstraint []CopyrightConstraint `json:"copyrightConstraint"`
	ApprConstraint      []ApprConstraint      `json:"apprConstraint"`
	LicenseConstraint   []LicenseConstraint   `json:"licenseConstraint"`
	WorkId              string                `json:"workId"`
	CopyrightStatus     []CopyrightStatus     `json:"copyrightStatus"`
}

// AuthenticationInfo 单个确权信息
type AuthenticationInfo struct {
	AuthenticationInstitudeName string `json:"authenticationInstitudeName"`
	AuthenticationId            string `json:"authenticationId"`
	AuthenticatedDate           string `json:"authenticatedDate"`
}

// CopyrightUnit 单个版权单元
type CopyrightUnit struct {
	Address          string `json:"address"`
	Proportion       string `json:"proportion"`
	CopyrightExplain string `json:"copyrightExplain"`
}

// CopyrightConstraint 版权限制
type CopyrightConstraint struct {
	CopyrightLimit int `json:"copyrightLimit"`
	// 如果还需要更多字段，可在此扩展
}

// ApprConstraint 授权约束
type ApprConstraint struct {
	Channel       string `json:"channel"`
	Area          string `json:"area"`
	Time          string `json:"time"`
	TransferType  int    `json:"transferType"`
	ReapproveType int    `json:"reapproveType"`
	// 其他可选字段...
}

// LicenseConstraint 许可约束
type LicenseConstraint struct {
	Type string `json:"type"`
	Area string `json:"area"`
	Time string `json:"time"`
}

// CopyrightStatus 发行状态信息
type CopyrightStatus struct {
	PublishStatus  int    `json:"publishStatus"`
	PublishCity    string `json:"publishCity"`
	PublishCountry string `json:"publishCountry"`
	PublishDate    string `json:"publishDate"`
	ComeoutStatus  int    `json:"comeoutStatus"`
	ComeoutCity    string `json:"comeoutCity"`
	ComeoutCountry string `json:"comeoutCountry"`
	ComeoutDate    string `json:"comeoutDate"`
	IssueStatus    int    `json:"issueStatus"`
	IssueCity      string `json:"issueCity"`
	IssueCountry   string `json:"issueCountry"`
	IssueDate      string `json:"issueDate"`
}

// ApproveToken 授权通证发行结构
type ApproveToken struct {
	Publisher          string              `json:"publisher"`          // 发行者账户地址
	Receiver           string              `json:"receiver"`           // 授权通证的接收者
	Token              string              `json:"token"`              // ERC-721通证名称
	TokenId            string              `json:"tokenId"`            // 授权通证id，hash256格式
	ReferenceID        string              `json:"referenceID"`        // 关联的版权通证tokenId, hash256格式
	ApproveType        int                 `json:"approveType"`        // 授权类型
	ApproveConstraints []ApproveConstraint `json:"approveConstraints"` // 约束信息（Array）
	Duty               []DutyInfo          `json:"duty"`               // 计酬信息（Array）
}

// ApproveConstraint 单个授权约束内容
type ApproveConstraint struct {
	ApproveChannel int `json:"approveChannel"` // 授权渠道
	ApproveArea    int `json:"approveArea"`    // 授权范围
	ApproveTime    int `json:"approveTime"`    // 授权时间
	ApproveStatus  int `json:"approveStatus"`  // 授权状态
	ReapproveType  int `json:"reapproveType"`  // 再授权类型
}

// DutyInfo 计酬信息
type DutyInfo struct {
	DistributionMethod int    `json:"distributionMethod"` // 计酬方式
	DistributionDesc   string `json:"distributionDesc"`   // 计酬描述
	ReceivablePayment  string `json:"receivablePayment"`  // 应收酬金
	ReceivedPayment    string `json:"receivedPayment"`    // 已收酬金
	ToReceivePayment   string `json:"toReceivePayment"`   // 未收酬金
	BalanceDate        string `json:"balanceDate"`        // 结算日期
}

// PubTokenTx 表示构建好的通证许可交易数据
type PubTokenTx struct {
	Publisher   string      `json:"publisher"`
	Receiver    string      `json:"receiver"`
	Token       string      `json:"token"`
	TokenId     string      `json:"tokenId"`
	ReferenceId string      `json:"referenceId"`
	TokenInfos  []TokenInfo `json:"tokenInfos,omitempty"`

	// 为了演示 ownerSign 的结果，我们额外加一些字段
	OwnerSigned    bool   `json:"ownerSigned"`    // 是否完成owner签名
	OwnerAccount   string `json:"ownerAccount"`   // 版权通证拥有者账户
	OwnerSignature string `json:"ownerSignature"` // owner签名(演示用,可存放签名后的字符串)
}

// TokenInfo 通证的属性
type TokenInfo struct {
	Type string `json:"type"` // 属性名称
	Data string `json:"data"` // 属性内容
}

// TokenDetail 示例：存储单个通证的完整信息
type TokenDetail struct {
	TokenId             string               `json:"tokenId"`                       // 通证ID
	Version             string               `json:"version"`                       // v1或v2
	Flag                int                  `json:"flag"`                          // v2版本的字段 (1=版权,2=授权,3=操作许可)；v1可忽略
	OwnerAccount        string               `json:"ownerAccount"`                  // 当前持有者(若是NFT，一般只有一个owner)
	Frozen              bool                 `json:"frozen"`                        // 是否冻结
	AuthenticationInfos []AuthenticationInfo `json:"authenticationInfos,omitempty"` // 确权信息数组

	// 这里可扩展更多字段
	// ...

	TokenInfos []TokenInfo `json:"tokenInfos,omitempty"` // 可选的属性信息列表
	// 版权单元数组
	CopyrightUnits []CopyrightUnit `json:"copyrightUnits,omitempty"`

	// 版权约束
	CopyrightConstraint []CopyrightConstraint `json:"copyrightConstraint,omitempty"`
	ApprConstraint      []ApprConstraint      `json:"apprConstraint,omitempty"`
	LicenseConstraint   []LicenseConstraint   `json:"licenseConstraint,omitempty"`
}

// ConstraintUpdate 更新约束时需要的一些数据
type ConstraintUpdate struct {
	TokenId    string `json:"tokenId"`
	Constraint struct {
		// 这里对应文档中需要修改的字段
		// "copyrightLimit", "apprConstraint", "licenseConstraint"等
		CopyrightLimit    int               `json:"copyrightLimit"`
		ApprConstraint    ApprConstraint    `json:"apprConstraint"`
		LicenseConstraint LicenseConstraint `json:"licenseConstraint"`
		ConstraintExplain string            `json:"constraintExplain"`
		ConstraintExpand  int               `json:"constraintExpand"`
	} `json:"constraint"`

	// 多签时，需要版权单元全部签名
	Signers []string `json:"signers"` // 本次提交的签名者账户列表
	// ...
}

// InitContract 合约初始化方法
func (tc *TokenContract) InitContract(stub shim.CMStubInterface) protogo.Response {
	return shim.Success([]byte("TokenContract Init Success"))
}

// UpgradeContract 合约升级方法
func (tc *TokenContract) UpgradeContract(stub shim.CMStubInterface) protogo.Response {
	return shim.Success([]byte("TokenContract Upgrade Success"))
}

//InvokeContract 调用合约
func (tc *TokenContract) InvokeContract(stub shim.CMStubInterface) protogo.Response {

	//获取调用合约哪个方法
	method := string(stub.GetArgs()["method"])

	// 这里必须写成 switch {case "a": ... [case "b": ...[...]] default:...} 形式
	// 而且case后面的内容必须是字符串,不能是常量
	// 这里必须写成 switch {case "a": ... [case "b": ...[...]] default:...} 形式
	// 而且case后面的内容必须是字符串,不能是常量
	// 这里必须写成 switch {case "a": ... [case "b": ...[...]] default:...} 形式
	// 而且case后面的内容必须是字符串,不能是常量

	// 如果 method == "save", 执行FactContract的save方法
	// 如果 method == "findByFileHash", 执行FactContract的findByFileHash方法
	// 如果没有对应的 case 语句，返回错误
	switch method {
	// 5.1 通证初始化
	case "buildTokenIssueTx":
		return tc.BuildTokenIssueTx(stub)

	// 5.2 (1) 一般通证发行
	case "buildPublishTokenTx":
		return tc.BuildPublishTokenTx(stub)
		// 5.2 (2) 授权通证发行
	case "buildPublishApproveTokenTx":
		return tc.BuildPublishApproveTokenTx(stub)

	// 5.3 通证许可（先组织交易，再 ownerSign）
	case "buildPubTokenTx":
		return tc.BuildPubTokenTx(stub)
	case "ownerSign":
		return tc.OwnerSign(stub)

	// 5.4 通证信息查询 (1) 查询账户所持通证 (2) 查询单个通证
	case "requestAccountToken":
		return tc.RequestAccountToken(stub)
	case "requestTokenInfo":
		return tc.RequestTokenInfo(stub)

	// 5.5 通证信息修改
	// (1) 修改通证标识位（冻结/解冻）
	case "BuildModifyCopyrightTokenFlagTx":
		// 如果你方法里写的名字是 buildModifyCopyrightTokenFlagTx，
		// 请保持一致:
		return tc.BuildModifyCopyrightTokenFlagTx(stub)

	// (2) 修改通证确权信息
	case "buildModifyAuthenticationInfoTx":
		return tc.BuildModifyAuthenticationInfoTx(stub)

	// (3) 修改版权通证的权利主体组
	case "buildModifyCopyrightUnitTx":
		return tc.BuildModifyCopyrightUnitTx(stub)

	// (4) 修改通证约束（多签）
	case "buildModifyConstraintTx":
		return tc.BuildModifyConstraintTx(stub)

	// (5) 通证变更方法
	case "buildTokenChangeTx":
		return tc.BuildTokenChangeTx(stub)

	// (6) 版权份额转让方法
	case "buildTransferProportionTx":
		return tc.BuildTransferProportionTx(stub)

	default:
		// 未匹配到任何已知方法，返回错误
		return shim.Error("[TokenContract] invalid method: " + method)
	}
}

// BuildTokenIssueTx 通证初始化（5.1 通证初始化）
func (tc *TokenContract) BuildTokenIssueTx(stub shim.CMStubInterface) protogo.Response {
	args := stub.GetArgs()

	// 读取参数
	account := string(args["account"])
	publisher := string(args["publisher"])
	tokenName := string(args["token"])
	numberStr := string(args["number"])
	flagStr := string(args["flag"])
	version := string(args["version"])
	rolesStr := string(args["roles"])                  // 输入需要是JSON数组，例如[{     "role": "0x1234567890abcdef",     "type": 1   }]
	referenceFlagStr := string(args["reference_flag"]) // 新增: referenceFlag

	// 基础校验
	if account == "" || publisher == "" || tokenName == "" ||
		numberStr == "" || flagStr == "" || version == "" || referenceFlagStr == "" {
		return shim.Error("some required params are empty, please check 'account','publisher','token','number','flag','version','reference_flag'")
	}

	// 转换 number
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		msg := "[buildTokenIssueTx] number must be integer, got: " + numberStr
		stub.Log(msg)
		return shim.Error(msg)
	}
	// 转换 flag
	flag, err := strconv.Atoi(flagStr)
	if err != nil {
		msg := "[buildTokenIssueTx] flag must be integer, got: " + flagStr
		stub.Log(msg)
		return shim.Error(msg)
	}
	// 转换 referenceFlag
	referenceFlag, err := strconv.Atoi(referenceFlagStr)
	if err != nil {
		msg := "[buildTokenIssueTx] reference_flag must be integer, got: " + referenceFlagStr
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 解析 roles
	var roles []TokenRole
	if rolesStr != "" {
		if err := json.Unmarshal([]byte(rolesStr), &roles); err != nil {
			msg := "[buildTokenIssueTx] fail to parse 'roles' json array: " + err.Error()
			stub.Log(msg)
			return shim.Error(msg)
		}
	}
	// ---------- 【关键】ERC‑721 注册占位 ----------
	// 说明：
	// ‑ 在大多数链/网关实现里，首次发行前需要把 tokenName
	//   注册到链下元数据，典型接口如:
	//      RegisterTokenERC721(tokenName, publisher)
	// ‑ 如果你已有 SDK 调用，请把下行替换成实际函数；
	//   若 ChainMaker 底层已自动允许“先发行后查询”，可注释掉。

	// 构造 TokenIssue 对象
	tokenIssue := &TokenIssue{
		Account:       account,
		Publisher:     publisher,
		Token:         tokenName,
		Number:        number,
		Flag:          flag,
		Version:       version,
		Roles:         roles,
		ReferenceFlag: referenceFlag,
	}

	// 序列化存储
	issueBytes, err := json.Marshal(tokenIssue)
	if err != nil {
		msg := "[buildTokenIssueTx] fail to marshal TokenIssue: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 定义存储键
	key := "token_issue_" + tokenName
	err = stub.PutStateFromKeyByte(key, issueBytes)
	if err != nil {
		msg := "[buildTokenIssueTx] fail to PutState: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 输出日志
	stub.Log("[buildTokenIssueTx] success, key=" + key)

	// 返回执行成功
	return shim.Success([]byte("buildTokenIssueTx success: " + tokenName))
}

// BuildPublishTokenTx 一般通证发行
// 文档：buildPublishTokenTx({})
func (tc *TokenContract) BuildPublishTokenTx(stub shim.CMStubInterface) protogo.Response {
	// 1. 获取参数
	args := stub.GetArgs()
	publisher := string(args["publisher"])            // 发行账户地址
	receiver := string(args["receiver"])              // 接收账户地址
	tokenName := string(args["token"])                // 通证名称
	referenceFlagStr := string(args["referenceFlag"]) // 通证标识(1=版权通证,2=授权通证,3=操作许可通证)
	tokenObjectStr := string(args["tokenObject"])     // JSON格式字符串

	// 2. 基础校验
	if publisher == "" || receiver == "" || tokenName == "" || referenceFlagStr == "" || tokenObjectStr == "" {
		return shim.Error("[buildPublishTokenTx] missing required params: 'publisher','receiver','token','referenceFlag','tokenObject'")
	}

	// 3. 转换 referenceFlag -> int
	referenceFlag, err := strconv.Atoi(referenceFlagStr)
	if err != nil {
		msg := "[buildPublishTokenTx] referenceFlag must be integer, got: " + referenceFlagStr
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 4. 解析 tokenObject
	var tokenObj TokenObject
	err = json.Unmarshal([]byte(tokenObjectStr), &tokenObj)
	if err != nil {
		msg := "[buildPublishTokenTx] fail to parse tokenObject JSON, err: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 5. 这里可根据 referenceFlag 判断通证类型, 做一些业务逻辑分支(可选)
	//    例如：1=版权通证 -> 要求必须有copyrightType...
	//          2=授权通证 -> ...
	//    若暂时不需要特别区分，可略过
	stub.Log("[buildPublishTokenTx] referenceFlag: " + referenceFlagStr)

	// 6. 组装要写入状态的结构
	//    你可以直接把前面解析好的数据打包到一个临时结构，然后序列化写入
	issueData := map[string]interface{}{
		"publisher":     publisher,
		"receiver":      receiver,
		"tokenName":     tokenName,
		"referenceFlag": referenceFlag,
		"tokenObject":   tokenObj,
	}

	dataBytes, err := json.Marshal(issueData)
	if err != nil {
		msg := "[buildPublishTokenTx] marshal data fail: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 7. 存储到区块链状态, key可自定义
	//    例如: "publish_token_" + tokenObj.TokenId
	//    如果 tokenObj.TokenId 为空，可以用别的拼装
	storeKey := "publish_token_" + tokenObj.TokenId
	if err := stub.PutStateFromKeyByte(storeKey, dataBytes); err != nil {
		msg := "[buildPublishTokenTx] PutState failed: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 8. (可选) 记录日志，或发事件
	stub.Log("[buildPublishTokenTx] success with key: " + storeKey)
	//    发事件： stub.EmitEvent("event_publish_token", []string{tokenName, tokenObj.TokenId})

	// 9. 返回成功
	return shim.Success([]byte("[buildPublishTokenTx] success for token: " + tokenObj.TokenId))
}

// BuildPublishApproveTokenTx 授权通证发行
// 对应文档：buildPublishApproveTokenTx({})
func (tc *TokenContract) BuildPublishApproveTokenTx(stub shim.CMStubInterface) protogo.Response {
	// 1. 从stub获取调用参数
	args := stub.GetArgs()
	publisher := string(args["publisher"])        // 发行者账户地址
	receiver := string(args["receiver"])          // 授权通证接收者
	tokenName := string(args["token"])            // ERC-721通证名称
	tokenId := string(args["tokenId"])            // 授权通证ID (hash256)
	referenceID := string(args["referenceID"])    // 关联的版权通证ID (hash256)
	approveTypeStr := string(args["approveType"]) // 授权类型 (Number)

	// 这两个是数组结构，用 JSON 解析
	approveConstraintsStr := string(args["approveConstraints"]) // JSON数组
	dutyStr := string(args["duty"])                             // JSON数组

	// 2. 基础校验
	if publisher == "" || receiver == "" || tokenName == "" ||
		tokenId == "" || referenceID == "" || approveTypeStr == "" {
		return shim.Error("[buildPublishApproveTokenTx] missing required params: 'publisher','receiver','token','tokenId','referenceID','approveType'")
	}

	// 3. 转换approveType -> int
	approveType, err := strconv.Atoi(approveTypeStr)
	if err != nil {
		msg := "[buildPublishApproveTokenTx] approveType must be integer, got: " + approveTypeStr
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 4. 解析 approveConstraints 数组(JSON)
	var approveConstraints []ApproveConstraint
	if approveConstraintsStr != "" {
		if err := json.Unmarshal([]byte(approveConstraintsStr), &approveConstraints); err != nil {
			msg := "[buildPublishApproveTokenTx] fail to parse 'approveConstraints': " + err.Error()
			stub.Log(msg)
			return shim.Error(msg)
		}
	}

	// 5. 解析 duty 数组(JSON)
	var dutyList []DutyInfo
	if dutyStr != "" {
		if err := json.Unmarshal([]byte(dutyStr), &dutyList); err != nil {
			msg := "[buildPublishApproveTokenTx] fail to parse 'duty': " + err.Error()
			stub.Log(msg)
			return shim.Error(msg)
		}
	}

	// 6. 构造 ApproveToken 对象
	approveToken := &ApproveToken{
		Publisher:          publisher,
		Receiver:           receiver,
		Token:              tokenName,
		TokenId:            tokenId,
		ReferenceID:        referenceID,
		ApproveType:        approveType,
		ApproveConstraints: approveConstraints,
		Duty:               dutyList,
	}

	// 7. 序列化
	tokenBytes, err := json.Marshal(approveToken)
	if err != nil {
		msg := "[buildPublishApproveTokenTx] fail to marshal ApproveToken: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 8. 构造存储Key，并写入区块链状态
	//    你可以自行决定如何拼接, 比如 "approve_token_" + tokenId
	storeKey := "approve_token_" + tokenId
	err = stub.PutStateFromKeyByte(storeKey, tokenBytes)
	if err != nil {
		msg := "[buildPublishApproveTokenTx] PutState failed: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 9. (可选) 发送合约事件
	// stub.EmitEvent("event_approve_token", []string{tokenId, referenceID})

	// 10. 返回执行成功
	stub.Log("[buildPublishApproveTokenTx] success with key: " + storeKey)
	return shim.Success([]byte("[buildPublishApproveTokenTx] success for token: " + tokenName))
}

// BuildPubTokenTx 通证许可 - 第一步
// 对应文档: buildPubTokenTx({publisher,receiver,token,tokenId,tokenInfos,referenceId})
func (tc *TokenContract) BuildPubTokenTx(stub shim.CMStubInterface) protogo.Response {
	args := stub.GetArgs()

	publisher := string(args["publisher"])
	receiver := string(args["receiver"])
	tokenName := string(args["token"])
	tokenId := string(args["tokenId"])
	referenceId := string(args["referenceId"])  // 关联的版权通证ID
	tokenInfosStr := string(args["tokenInfos"]) // JSON数组

	// 校验必填项
	if publisher == "" || receiver == "" || tokenName == "" || tokenId == "" || referenceId == "" {
		return shim.Error("[buildPubTokenTx] missing required params: 'publisher','receiver','token','tokenId','referenceId'")
	}

	// 解析 tokenInfos (可选)
	var tokenInfos []TokenInfo
	if tokenInfosStr != "" {
		if err := json.Unmarshal([]byte(tokenInfosStr), &tokenInfos); err != nil {
			msg := "[buildPubTokenTx] fail to parse tokenInfos: " + err.Error()
			stub.Log(msg)
			return shim.Error(msg)
		}
	}

	// 构造 PubTokenTx 对象
	pubTx := &PubTokenTx{
		Publisher:   publisher,
		Receiver:    receiver,
		Token:       tokenName,
		TokenId:     tokenId,
		ReferenceId: referenceId,
		TokenInfos:  tokenInfos,

		// 初始状态下，还没有owner签名
		OwnerSigned:    false,
		OwnerAccount:   "",
		OwnerSignature: "",
	}

	// 序列化
	txBytes, err := json.Marshal(pubTx)
	if err != nil {
		msg := "[buildPubTokenTx] fail to marshal pubTx: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 存储到区块链状态中
	// key 可以自定义, 比如 "pub_token_tx_" + tokenId
	storeKey := "pub_token_tx_" + tokenId
	if err := stub.PutStateFromKeyByte(storeKey, txBytes); err != nil {
		msg := "[buildPubTokenTx] PutState failed: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	stub.Log("[buildPubTokenTx] success with key: " + storeKey)

	// 返回成功
	return shim.Success([]byte("[buildPubTokenTx] success for tokenId: " + tokenId))
}

// OwnerSign 通证许可 - 第二步(Owner签名)
// 对应文档: ownerSign({account, secret})
func (tc *TokenContract) OwnerSign(stub shim.CMStubInterface) protogo.Response {
	args := stub.GetArgs()

	account := string(args["account"])
	secret := string(args["secret"])
	if account == "" || secret == "" {
		return shim.Error("[ownerSign] missing required params: 'account','secret'")
	}

	// 这里你可能需要知道要签的是哪个许可交易
	// 文档中只提到 ownerSign 需要 {account, secret}，但没有说如何区分哪条交易
	// 通常我们还需要加一个 tokenId(或交易ID)参数来知道签的是哪条记录
	// 假如我们通过 "tokenId" 来区分：
	tokenId := string(args["tokenId"])
	if tokenId == "" {
		return shim.Error("[ownerSign] missing required param: 'tokenId'")
	}

	// 从区块链状态中把 pub_token_tx_ + tokenId 取出来
	storeKey := "pub_token_tx_" + tokenId
	oldData, err := stub.GetStateFromKeyByte(storeKey)
	if err != nil {
		msg := "[ownerSign] GetState failed: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if len(oldData) == 0 {
		return shim.Error("[ownerSign] no such pubTokenTx for tokenId: " + tokenId)
	}

	// 反序列化为 PubTokenTx
	var pubTx PubTokenTx
	if err := json.Unmarshal(oldData, &pubTx); err != nil {
		msg := "[ownerSign] Unmarshal pubTx failed: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 这里应该校验一下 account 是否真的是版权通证 owner
	// 或者检查 referenceId(也就是关联的版权通证)对应的 owner 是否匹配本 account
	// 由于文档没给更详细的校验逻辑，我们先省略，只演示简单记录

	// 假装把 secret 做成一个"签名"字符串
	signature := "SignatureOf:" + account + ":" + secret

	// 更新 pubTx
	pubTx.OwnerSigned = true
	pubTx.OwnerAccount = account
	pubTx.OwnerSignature = signature

	// 重新序列化并写回状态
	newData, err := json.Marshal(pubTx)
	if err != nil {
		msg := "[ownerSign] marshal updated pubTx failed: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if err := stub.PutStateFromKeyByte(storeKey, newData); err != nil {
		msg := "[ownerSign] PutState failed: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	stub.Log("[ownerSign] success for tokenId: " + tokenId)
	return shim.Success([]byte("[ownerSign] success, tokenId=" + tokenId))
}

// RequestAccountToken 查询账户所持有的通证
// 文档: requestAccountToken({account, version, flag?})
func (tc *TokenContract) RequestAccountToken(stub shim.CMStubInterface) protogo.Response {
	args := stub.GetArgs()

	account := string(args["account"])
	version := string(args["version"])
	flagStr := string(args["flag"]) // 只在 v2 有效

	if account == "" || version == "" {
		return shim.Error("[RequestAccountToken] missing required params: 'account','version'")
	}

	// 默认值(仅当 version=="v2" 时才使用), 若为空则默认=1
	var flag int
	if version == "v2" {
		if flagStr == "" {
			flag = 1 // 默认=1(版权通证)
		} else {
			f, err := strconv.Atoi(flagStr)
			if err != nil {
				return shim.Error("[RequestAccountToken] invalid 'flag' param: " + flagStr)
			}
			flag = f
		}
	}

	// 1. 读取 account_tokens_<account>
	indexKey := "account_tokens_" + account
	indexBytes, err := stub.GetStateFromKeyByte(indexKey)
	if err != nil {
		msg := "[RequestAccountToken] fail to GetState for " + indexKey + ": " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	if len(indexBytes) == 0 {
		// 说明该账户没有任何通证
		// 根据文档，可以返回一个空数组
		return shim.Success([]byte("[]"))
	}

	// 2. indexBytes 中存的是JSON数组，比如 ["tokenId1","tokenId2"...]
	var tokenIds []string
	if err := json.Unmarshal(indexBytes, &tokenIds); err != nil {
		msg := "[RequestAccountToken] fail to Unmarshal index array: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 3. 遍历tokenIds, 读取每个token详情, 根据version/flag过滤
	var resultTokens []TokenDetail
	for _, tid := range tokenIds {
		detailKey := "publish_token_" + tid
		detailBytes, err := stub.GetStateFromKeyByte(detailKey)
		if err != nil {
			stub.Log("[RequestAccountToken] skip tokenId=" + tid + ", error: " + err.Error())
			continue
		}
		if len(detailBytes) == 0 {
			// 数据不存在,跳过
			continue
		}

		var detail TokenDetail
		if err := json.Unmarshal(detailBytes, &detail); err != nil {
			stub.Log("[RequestAccountToken] skip tokenId=" + tid + ", unmarshal error: " + err.Error())
			continue
		}

		// 按 version & flag 进行过滤 (若version=="v1"不比较flag)
		if detail.Version != version {
			continue
		}
		if version == "v2" {
			if detail.Flag != flag {
				continue
			}
		}

		// 符合筛选条件, 放进结果数组
		resultTokens = append(resultTokens, detail)
	}

	// 4. 序列化 resultTokens 返回
	retBytes, err := json.Marshal(resultTokens)
	if err != nil {
		msg := "[RequestAccountToken] fail to marshal result: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	return shim.Success(retBytes)
}

// RequestTokenInfo 查询单个通证详情
// 文档: requestTokenInfo({ tokenId, version })
func (tc *TokenContract) RequestTokenInfo(stub shim.CMStubInterface) protogo.Response {
	args := stub.GetArgs()

	tokenId := string(args["tokenId"])
	version := string(args["version"])
	if tokenId == "" || version == "" {
		return shim.Error("[RequestTokenInfo] missing required param: 'tokenId','version'")
	}

	// 1. 读取通证详情
	detailKey := "publish_token_" + tokenId
	detailBytes, err := stub.GetStateFromKeyByte(detailKey)
	if err != nil {
		msg := "[RequestTokenInfo] fail to GetState for " + detailKey + ": " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if len(detailBytes) == 0 {
		// 不存在
		return shim.Error("[RequestTokenInfo] no token found for tokenId: " + tokenId)
	}

	// 2. 反序列化
	var detail TokenDetail
	if err := json.Unmarshal(detailBytes, &detail); err != nil {
		msg := "[RequestTokenInfo] unmarshal error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 3. 判断版本是否匹配 (如果需要强制校验)
	if detail.Version != version {
		// 也可以返回错误,或直接返回该通证(看你业务要求)
		return shim.Error("[RequestTokenInfo] version mismatch, the token stored version=" + detail.Version + ", but input=" + version)
	}

	// 4. 返回详情
	return shim.Success(detailBytes)
}

// BuildModifyCopyrightTokenFlagTx
// (1) 修改通证标识位(冻结/解冻)
// 文档: buildModifyCopyrightTokenFlagTx({account, tokenId, flag})
func (tc *TokenContract) BuildModifyCopyrightTokenFlagTx(stub shim.CMStubInterface) protogo.Response {
	args := stub.GetArgs()

	account := string(args["account"]) // 有修改权限的账户(监管机构白名单)
	tokenId := string(args["tokenId"]) // 通证ID
	flagStr := string(args["flag"])    // 0=解冻,1=冻结

	// 基础校验
	if account == "" || tokenId == "" || flagStr == "" {
		return shim.Error("[BuildModifyCopyrightTokenFlagTx] missing params: account, tokenId, flag")
	}
	// 转成 int
	freezeFlag, err := strconv.Atoi(flagStr)
	if err != nil {
		return shim.Error("[BuildModifyCopyrightTokenFlagTx] flag must be int (0 or 1)")
	}

	// 1. 读取通证详情
	storeKey := "publish_token_" + tokenId
	detailBytes, err := stub.GetStateFromKeyByte(storeKey)
	if err != nil {
		msg := "[ModifyTokenFlag] fail to GetState, err: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if len(detailBytes) == 0 {
		return shim.Error("[ModifyTokenFlag] token not found for tokenId=" + tokenId)
	}

	var detail TokenDetail
	if err := json.Unmarshal(detailBytes, &detail); err != nil {
		msg := "[ModifyTokenFlag] unmarshal error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 2. 这里可以校验是否 account 在"监管机构白名单"内, 具体逻辑看你业务
	//    for example:
	//    if !checkRegulatorWhitelist(account) {
	//        return shim.Error("[ModifyTokenFlag] account not in regulator whitelist: " + account)
	//    }

	// 3. 根据flag=0/1进行冻结或解冻
	if freezeFlag == 0 {
		// 解冻
		detail.Frozen = false
	} else {
		// 冻结
		detail.Frozen = true
	}

	// 4. 序列化并写回
	updatedBytes, err := json.Marshal(detail)
	if err != nil {
		msg := "[ModifyTokenFlag] marshal error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if err := stub.PutStateFromKeyByte(storeKey, updatedBytes); err != nil {
		msg := "[ModifyTokenFlag] PutState failed: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	stub.Log("[ModifyTokenFlag] success for tokenId=" + tokenId)
	return shim.Success([]byte("[ModifyTokenFlag] success"))
}

func (tc *TokenContract) BuildModifyAuthenticationInfoTx(stub shim.CMStubInterface) protogo.Response {
	args := stub.GetArgs()

	account := string(args["account"])                // 有修改通证权限的账户(确权白名单)
	tokenId := string(args["tokenId"])                // 通证ID
	authInfoStr := string(args["authenticationInfo"]) // JSON对象

	// 校验
	if account == "" || tokenId == "" || authInfoStr == "" {
		return shim.Error("[BuildModifyAuthenticationInfoTx] missing params: 'account','tokenId','authenticationInfo'")
	}

	// 1. 解析 authenticationInfo
	var newAuthInfo AuthenticationInfo
	if err := json.Unmarshal([]byte(authInfoStr), &newAuthInfo); err != nil {
		msg := "[ModifyAuthInfo] fail to parse authenticationInfo: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 2. 读取通证详情
	storeKey := "publish_token_" + tokenId
	detailBytes, err := stub.GetStateFromKeyByte(storeKey)
	if err != nil {
		msg := "[ModifyAuthInfo] fail to GetState, err: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if len(detailBytes) == 0 {
		return shim.Error("[ModifyAuthInfo] token not found for tokenId=" + tokenId)
	}

	var detail TokenDetail
	if err := json.Unmarshal(detailBytes, &detail); err != nil {
		msg := "[ModifyAuthInfo] unmarshal error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 3. 校验 'account' 是否在确权白名单
	//    if !checkAuthWhitelist(account) {
	//        return shim.Error("[ModifyAuthInfo] account not in auth-whitelist: " + account)
	//    }

	// 4. 更新通证的 "AuthenticationInfos"
	//    (文档说“一次只能上传一个”，可以把它append到列表中，或者做替换等)
	detail.AuthenticationInfos = append(detail.AuthenticationInfos, newAuthInfo)

	// 5. 序列化并写回
	updatedBytes, err := json.Marshal(detail)
	if err != nil {
		msg := "[ModifyAuthInfo] marshal error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if err := stub.PutStateFromKeyByte(storeKey, updatedBytes); err != nil {
		msg := "[ModifyAuthInfo] PutState failed: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	stub.Log("[ModifyAuthInfo] success for tokenId=" + tokenId)
	return shim.Success([]byte("[ModifyAuthInfo] success"))
}

// BuildModifyCopyrightUnitTx
// (3) 修改版权通证的权利主体组/修改版权单元
// 文档: buildModifyCopyrightUnitTx({account, tokenId, address})
func (tc *TokenContract) BuildModifyCopyrightUnitTx(stub shim.CMStubInterface) protogo.Response {
	args := stub.GetArgs()

	account := string(args["account"]) // 有修改通证权限的账户, owner
	tokenId := string(args["tokenId"])
	newAddress := string(args["address"])

	if account == "" || tokenId == "" || newAddress == "" {
		return shim.Error("[BuildModifyCopyrightUnitTx] missing params: 'account','tokenId','address'")
	}

	// 1. 读取通证详情
	storeKey := "publish_token_" + tokenId
	detailBytes, err := stub.GetStateFromKeyByte(storeKey)
	if err != nil {
		msg := "[ModifyCopyrightUnit] fail to GetState: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if len(detailBytes) == 0 {
		return shim.Error("[ModifyCopyrightUnit] no token found for tokenId=" + tokenId)
	}

	var detail TokenDetail
	if err := json.Unmarshal(detailBytes, &detail); err != nil {
		msg := "[ModifyCopyrightUnit] unmarshal error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 2. 校验权限：文档说要由owner来修改
	//    if detail.OwnerAccount != account {
	//       return shim.Error("[ModifyCopyrightUnit] account is not the owner")
	//    }

	// 3. 这里按照文档“address: 新的版权单元地址，替换旧的”。
	//    但并未说明替换哪个？ 需要我们自己设计匹配逻辑。
	//    例如：假设只替换第一个或所有?
	//    这里假设：找第一个匹配account的版权单元地址，然后改成 newAddress

	replaced := false
	for i, cu := range detail.CopyrightUnits {
		if cu.Address == account {
			// 找到匹配的旧地址,替换为 newAddress
			detail.CopyrightUnits[i].Address = newAddress
			replaced = true
		}
	}

	if !replaced {
		// 如果没找到，也许返回错误，也许直接append(看你业务需求)
		// 这里演示返回错误
		return shim.Error("[ModifyCopyrightUnit] no matched old address with account=" + account)
	}

	// 4. 重新序列化写回
	updatedBytes, err := json.Marshal(detail)
	if err != nil {
		msg := "[ModifyCopyrightUnit] marshal error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if err := stub.PutStateFromKeyByte(storeKey, updatedBytes); err != nil {
		msg := "[ModifyCopyrightUnit] PutState error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	stub.Log("[ModifyCopyrightUnit] success for tokenId=" + tokenId)
	return shim.Success([]byte("[ModifyCopyrightUnit] success"))
}

// BuildModifyConstraintTx
// (4) 修改通证约束（需多签提交）
// 文档: buildModifyConstraintTx({tokenId, constraint})
// 这里演示一次性提交所有签名, 若签名不足, 交易失败
func (tc *TokenContract) BuildModifyConstraintTx(stub shim.CMStubInterface) protogo.Response {
	args := stub.GetArgs()

	tokenId := string(args["tokenId"])
	constraintStr := string(args["constraint"])

	if tokenId == "" || constraintStr == "" {
		return shim.Error("[BuildModifyConstraintTx] missing params: 'tokenId','constraint'")
	}

	// 解析 constraint
	var constraintUpdate ConstraintUpdate
	if err := json.Unmarshal([]byte(constraintStr), &constraintUpdate); err != nil {
		msg := "[BuildModifyConstraintTx] fail to parse constraint: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 1. 读取 tokenDetail
	storeKey := "publish_token_" + tokenId
	detailBytes, err := stub.GetStateFromKeyByte(storeKey)
	if err != nil {
		msg := "[BuildModifyConstraintTx] GetState failed: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if len(detailBytes) == 0 {
		return shim.Error("[BuildModifyConstraintTx] no token found for tokenId=" + tokenId)
	}

	var detail TokenDetail
	if err := json.Unmarshal(detailBytes, &detail); err != nil {
		msg := "[BuildModifyConstraintTx] unmarshal detail error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 2. 多签检查: 版权单元记录的所有地址
	//    简化逻辑: 逐个检查 detail.CopyrightUnits
	//    真实业务中，可能需要2/3签名之类的规则
	signerMap := make(map[string]bool)
	for _, s := range constraintUpdate.Signers {
		signerMap[s] = true
	}

	for _, cu := range detail.CopyrightUnits {
		if !signerMap[cu.Address] {
			// 说明 cu.Address 未签名, 多签不通过
			return shim.Error("[BuildModifyConstraintTx] multi-sign failed, missing address: " + cu.Address)
		}
	}

	// 3. 多签通过，执行修改
	//    把 constraintUpdate.Constraint 的内容写到 detail 里
	detail.CopyrightConstraint = []CopyrightConstraint{
		{
			CopyrightLimit: constraintUpdate.Constraint.CopyrightLimit,
			// 如果你需要更多字段, 在这里赋值
		},
	}
	detail.ApprConstraint = []ApprConstraint{
		constraintUpdate.Constraint.ApprConstraint,
	}
	detail.LicenseConstraint = []LicenseConstraint{
		constraintUpdate.Constraint.LicenseConstraint,
	}
	// 你也可以把 constraintUpdate.Constraint.ConstraintExplain / Expand
	// 存到 detail 里, 视业务而定

	// 4. 写回
	updatedBytes, err := json.Marshal(detail)
	if err != nil {
		msg := "[BuildModifyConstraintTx] marshal detail error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	if err := stub.PutStateFromKeyByte(storeKey, updatedBytes); err != nil {
		msg := "[BuildModifyConstraintTx] PutState error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	stub.Log("[BuildModifyConstraintTx] success, tokenId=" + tokenId)
	return shim.Success([]byte("[BuildModifyConstraintTx] success"))
}

// BuildTokenChangeTx 通证变更方法
// 用于变更版权通证持有者, 修改冻结标志, 以及可选的 tokenInfos
// 文档: buildTokenChangeTx({ account, tokenId, flags, tokenInfos? })
func (tc *TokenContract) BuildTokenChangeTx(stub shim.CMStubInterface) protogo.Response {
	args := stub.GetArgs()

	account := string(args["account"])          // 有修改权限的账户（这里视为新的持有者）
	tokenId := string(args["tokenId"])          // 通证ID
	flagsStr := string(args["flags"])           // 0=解冻,1=冻结
	tokenInfosStr := string(args["tokenInfos"]) // JSON数组,可选

	if account == "" || tokenId == "" {
		return shim.Error("[BuildTokenChangeTx] missing required params: 'account','tokenId'")
	}

	// 1. 读取通证详情
	storeKey := "publish_token_" + tokenId
	detailBytes, err := stub.GetStateFromKeyByte(storeKey)
	if err != nil {
		msg := "[BuildTokenChangeTx] fail to GetState: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if len(detailBytes) == 0 {
		return shim.Error("[BuildTokenChangeTx] no token found for tokenId=" + tokenId)
	}

	var detail TokenDetail
	if err := json.Unmarshal(detailBytes, &detail); err != nil {
		msg := "[BuildTokenChangeTx] unmarshal error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 2. 校验权限: 提取合约调用者的address，判断“account”是否真有修改该token的权限？
	//    业务可自定义. 这里示例直接允许
	//    if detail.OwnerAccount != ??? { ... }

	// 3. 若 flagsStr 不空, 解析并进行冻结/解冻操作
	if flagsStr != "" {
		flags, err := strconv.Atoi(flagsStr)
		if err != nil {
			msg := "[BuildTokenChangeTx] invalid flags param: " + flagsStr
			stub.Log(msg)
			return shim.Error(msg)
		}
		if flags == 0 {
			detail.Frozen = false
		} else {
			detail.Frozen = true
		}
	}

	// 4. 更新 ownerAccount => 以 “account” 作为新的持有者
	detail.OwnerAccount = account

	// 5. 如果有 tokenInfos, 解析并更新
	if tokenInfosStr != "" {
		var tInfos []TokenInfo
		if err := json.Unmarshal([]byte(tokenInfosStr), &tInfos); err != nil {
			msg := "[BuildTokenChangeTx] fail to parse tokenInfos: " + err.Error()
			stub.Log(msg)
			return shim.Error(msg)
		}
		// 可以覆盖或追加, 看你业务需求
		detail.TokenInfos = tInfos
	}

	// 6. 写回
	updatedBytes, err := json.Marshal(detail)
	if err != nil {
		msg := "[BuildTokenChangeTx] marshal error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if err := stub.PutStateFromKeyByte(storeKey, updatedBytes); err != nil {
		msg := "[BuildTokenChangeTx] PutState error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	stub.Log("[BuildTokenChangeTx] success, tokenId=" + tokenId)
	return shim.Success([]byte("[BuildTokenChangeTx] success"))
}

// BuildTransferProportionTx
// (6) 版权份额转让
// 文档: buildTransferProportionTx({ account, tokenId, copyrightUnits })
func (tc *TokenContract) BuildTransferProportionTx(stub shim.CMStubInterface) protogo.Response {
	args := stub.GetArgs()

	account := string(args["account"]) // 需转移的 address 账户
	tokenId := string(args["tokenId"])
	cuStr := string(args["copyrightUnits"]) // JSON数组

	if account == "" || tokenId == "" || cuStr == "" {
		return shim.Error("[BuildTransferProportionTx] missing params: 'account','tokenId','copyrightUnits'")
	}

	// 1. 解析copyrightUnits
	var newUnits []CopyrightUnit
	if err := json.Unmarshal([]byte(cuStr), &newUnits); err != nil {
		msg := "[BuildTransferProportionTx] fail to parse 'copyrightUnits': " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 2. 读取 tokenDetail
	storeKey := "publish_token_" + tokenId
	detailBytes, err := stub.GetStateFromKeyByte(storeKey)
	if err != nil {
		msg := "[BuildTransferProportionTx] GetState failed: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}
	if len(detailBytes) == 0 {
		return shim.Error("[BuildTransferProportionTx] no token found for tokenId=" + tokenId)
	}

	var detail TokenDetail
	if err := json.Unmarshal(detailBytes, &detail); err != nil {
		msg := "[BuildTransferProportionTx] unmarshal detail error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	// 3. 找到 'account' 对应的版权单元, 获取其 proportion, 并删除/清零
	var oldProportionStr string
	oldUnits := detail.CopyrightUnits
	foundIndex := -1
	for i, cu := range oldUnits {
		if cu.Address == account {
			oldProportionStr = cu.Proportion
			foundIndex = i
			break
		}
	}
	_ = oldProportionStr
	if foundIndex < 0 {
		return shim.Error("[BuildTransferProportionTx] the account does not hold any proportion: " + account)
	}

	// 移除 oldUnits[foundIndex]
	oldProportion := oldUnits[foundIndex].Proportion
	detail.CopyrightUnits = append(
		oldUnits[:foundIndex],
		oldUnits[foundIndex+1:]...,
	)

	// 4. 校验“必须全部转让”
	//    oldProportion 可能是 "0.3" 或 "NA"
	//    你可以转换为浮点数(需注意精度)或做字符串对比
	//    这里示例中假设 oldProportion 是一个可解析的浮点数
	oldP, err := strconv.ParseFloat(oldProportion, 64)
	if err != nil {
		return shim.Error("[BuildTransferProportionTx] old proportion is not numeric: " + oldProportion)
	}

	// 5. 把 newUnits 追加到 detail, 并验证 newUnits 合计
	var sumNew float64
	for _, nu := range newUnits {
		p, err := strconv.ParseFloat(nu.Proportion, 64)
		if err != nil {
			return shim.Error("[BuildTransferProportionTx] new proportion invalid: " + nu.Proportion)
		}
		sumNew += p
	}

	// 如果你想严格要求 sumNew == oldP，可以检查:
	if math.Abs(sumNew-oldP) > 1e-9 {
		// 不相等
		return shim.Error("[BuildTransferProportionTx] sum of new proportions != old proportion to be transferred")
	}

	// 把 newUnits 追加到 detail
	detail.CopyrightUnits = append(
		detail.CopyrightUnits,
		newUnits...,
	)

	// 6. 写回
	updatedBytes, err := json.Marshal(detail)
	if err != nil {
		msg := "[BuildTransferProportionTx] marshal detail error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	if err := stub.PutStateFromKeyByte(storeKey, updatedBytes); err != nil {
		msg := "[BuildTransferProportionTx] PutState error: " + err.Error()
		stub.Log(msg)
		return shim.Error(msg)
	}

	stub.Log("[BuildTransferProportionTx] success for tokenId=" + tokenId)
	return shim.Success([]byte("[BuildTransferProportionTx] success"))
}

func main() {

	//运行合约
	err := shim.Start(new(TokenContract))
	if err != nil {
		log.Fatal(err)
	}
}
