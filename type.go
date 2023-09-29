package wdb_drive

import (
	"fmt"
)

type CreateDataReq struct {
	Key        string   `json:"key"`
	Categories []string `json:"categories"`
	Content    string   `json:"content"`
	Time       uint64   `json:"time"`
	Sg         string   `json:"sg"`
}

type UpdateDataReq struct {
	Key     string `json:"key"`
	Content string `json:"content"`
	Time    uint64 `json:"time"`
	Sg      string `json:"sg"`
}

type UploadReq struct {
	Path       string   `json:"path"`
	Key        string   `json:"key"`
	Categories []string `json:"categories"`
	Time       uint64   `json:"time"`
	Sg         string   `json:"sg"`
}

type DownReq struct {
	Path string `json:"path"`
	Key  string `json:"key"`
	Time uint64 `json:"time"`
	Sg   string `json:"sg"`
}

type TransBeginReq struct {
	Keys []string `json:"keys"`
	Time uint64   `json:"time"`
	Sg   string   `json:"sg"`
}

type TransIdReq struct {
	Tsid string `json:"tsid"`
	Time uint64 `json:"time"`
	Sg   string `json:"sg"`
}

type TransCreateDataReq struct {
	Tsid       string   `json:"tsid"`
	Key        string   `json:"key"`
	Categories []string `json:"categories"`
	Content    string   `json:"content"`
	Time       uint64   `json:"time"`
	Sg         string   `json:"sg"`
}

type CreateIndexReq struct {
	Indexkey []string `json:"indexkey"`
	Key      string   `json:"key"`
	Indexraw []string `json:"indexraw"`
	Time     uint64   `json:"time"`
	Sg       string   `json:"sg"`
}

type UpdateIndexReq struct {
	Oindexkey []string `json:"oindexkey"`
	Cindexkey []string `json:"cindexkey"`
	Key       string   `json:"key"`
	Indexraw  []string `json:"indexraw"`
	Time      uint64   `json:"time"`
	Sg        string   `json:"sg"`
}

type DelIndexReq struct {
	Indexkey []string `json:"indexkey"`
	Key      string   `json:"key"`
	Time     uint64   `json:"time"`
	Sg       string   `json:"sg"`
}

type ListIndexReq struct {
	Indexkey  string `json:"indexkey"`
	Condition string `json:"condition"`
	Offset    uint64 `json:"offset"`
	Limit     uint64 `json:"limit"`
	Order     string `json:"order"`
	Time      uint64 `json:"time"`
	Sg        string `json:"sg"`
}

type TransUpdateDataReq struct {
	Tsid    string `json:"tsid"`
	Key     string `json:"key"`
	Content string `json:"content"`
	Time    uint64 `json:"time"`
	Sg      string `json:"sg"`
}

type ApiRsp struct {
	Code uint64 `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type RawRsp struct {
	Code uint64 `json:"code"`
	Msg  string `json:"msg"`
	Size uint64 `json:"size"`
	Raw  []byte `json:"raw"`
}

type ListRsp struct {
	Code  uint64   `json:"code"`
	Msg   string   `json:"msg"`
	Total uint64   `json:"total"`
	List  []string `json:"list"`
}

func rsp_err(err error) ApiRsp {
	return ApiRsp{
		400,
		fmt.Sprintf("%s", err),
		"",
	}
}

func rsp_ok(code uint64, data string) ApiRsp {
	return ApiRsp{
		200,
		"",
		data,
	}
}

func list_err(err string) ListRsp {
	return ListRsp{
		400,
		err,
		0,
		[]string{},
	}
}
