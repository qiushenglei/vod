package entity

type BindGroupReq struct {
}

type DoPlayReq struct {
}

type SendCommandReq struct {
	Pos  int64 `json:"pos"`
	Sign int   `json:"sign"`
}
