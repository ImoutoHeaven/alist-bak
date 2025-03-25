package doubao

type BaseResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type NodeInfoResp struct {
	BaseResp
	Data struct {
		NodeInfo   NodeInfo   `json:"node_info"`
		Children   []NodeInfo `json:"children"`
		NextCursor string     `json:"next_cursor"`
		HasMore    bool       `json:"has_more"`
	} `json:"data"`
}

type NodeInfo struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	Key                 string `json:"key"`
	NodeType            int    `json:"node_type"` // 0: 文件, 1: 文件夹
	Size                int    `json:"size"`
	Source              int    `json:"source"`
	NameReviewStatus    int    `json:"name_review_status"`
	ContentReviewStatus int    `json:"content_review_status"`
	RiskReviewStatus    int    `json:"risk_review_status"`
	ConversationID      string `json:"conversation_id"`
	ParentID            string `json:"parent_id"`
	CreateTime          int    `json:"create_time"`
	UpdateTime          int    `json:"update_time"`
}
