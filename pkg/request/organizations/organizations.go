package organizations

type Organizations struct {
	OrgUscc     string     `json:"orgUscc" binding:"required"`
	DueTime     string     `json:"dueTime" binding:"required"`
	RestartTime string     `json:"restartTime" binding:"required"`
	NodeList    []NodeList `json:"nodeList" binding:"required"`
}
type NodeList struct {
	NodeType      string `json:"nodeType" binding:"required"`
	NodeNumber    string `json:"nodeNumber" binding:"required"`
	NodeCore      string `json:"nodeCore" binding:"required"`
	NodeMemory    string `json:"nodeMemory" binding:"required"`
	NodeBandwidth string `json:"nodeBandwidth" binding:"required"`
	NodeDisk      string `json:"nodeDisk" binding:"required"`
}
