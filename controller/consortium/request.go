package consortium

type NewReq struct {
	ManagerPolicy string         `json:"manager_policy"`
	LedgerQuery   bool           `json:"ledger_query"`
	Name          string         `json:"Name"`
	Scenes        string         `json:"scenes"`
	Desc          string         `json:"desc"`
	Group         map[string]Org `json:"group"`
	MemberName    string         `json:"member_name"`
	StoreType     uint           `json:"store_type"`
	Token         string         `json:"token"`
	Keystore      string         `json:"keystore"`  // 私钥
	Initiator     string         `json:"initiator"` // 发起组织uscc
}

type Org struct {
	Role        string   `json:"role"` // 组织联盟权限
	CommitTx    bool     `json:"commitTx"`
	TxSignature bool     `json:"txSignature"`
	Orderer     []string `json:"orderer"`
	Peer        []struct {
		Domain    string `json:"domain"`
		SyncBlock int    `json:"syncBlock"`
	} `json:"peer"`
}
