package organizations

import (
	"github.com/hxx258456/pyramidel-chain-baas/pkg/utils/localtime"
)

// Organizations 新增组织
type Organizations struct {
	OrgUscc      string               `json:"orgUscc" binding:"required"`
	DueTime      *localtime.LocalTime `json:"dueTime" binding:"required"`
	RestartTime  *localtime.LocalTime `json:"-"`
	NodeList     []NodeList           `json:"nodeList" binding:"required"`
	OrgPackageId uint64               `json:"orgPackageId" binding:"required"`
}
type NodeList struct {
	NodeType      int  `json:"nodeType" binding:"required"`
	NodeNumber    uint `json:"nodeNumber" binding:"required"`
	NodeCore      uint `json:"nodeCore" binding:"required"`
	NodeMemory    uint `json:"nodeMemory" binding:"required"`
	NodeBandwidth uint `json:"nodeBandwidth" binding:"required"`
	NodeDisk      uint `json:"nodeDisk" binding:"required"`
}
