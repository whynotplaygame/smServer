package model

type OpenCollectionReq struct {
}

type OpenCollectionRsp struct {
	Limit    int8  `json:"limit"`
	CurTimes int8  `json:"cur_times"`
	NextTime int64 `json:"next_time"`
}

type CollectionReq struct {
}

type CollectionRsp struct {
	Gold     int   `json:"gold"`
	Limit    int8  `json:"limit"`
	CurTimes int8  `json:"cur_times"`
	NextTime int64 `json:"next_time"`
}
