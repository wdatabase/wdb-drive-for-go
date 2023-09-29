package wdb_drive

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Wdb struct {
	Host string
	Key  string
}

var logger = log.New(os.Stdout, "<wdb drive>", log.Lshortfile|log.Ldate|log.Ltime)

func (self *Wdb) CreateObj(key string, data string, categories []string) ApiRsp {
	var req CreateDataReq
	req.Key = key
	req.Categories = categories
	req.Content = data
	tm := uint64(time.Now().Unix())
	req.Time = tm
	req.Sg = self.sign(fmt.Sprintf("%s%d%s%s", self.Key, tm, key, data))

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/create", str_req)
}

func (self *Wdb) UpdateObj(key string, data string) ApiRsp {
	var req UpdateDataReq
	req.Key = key
	req.Content = data
	tm := uint64(time.Now().Unix())
	req.Time = tm
	req.Sg = self.sign(fmt.Sprintf("%s%d%s%s", self.Key, tm, key, data))

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/update", str_req)
}

func (self *Wdb) GetObj(key string) ApiRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%d%s", self.Key, tm, key))
	return self.getData(fmt.Sprintf("/wdb/api/get?key=%s&time=%d&sg=%s", key, tm, sg))
}

func (self *Wdb) DelObj(key string) ApiRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%d%s", self.Key, tm, key))
	return self.getData(fmt.Sprintf("/wdb/api/del?key=%s&time=%d&sg=%s", key, tm, sg))
}

func (self *Wdb) ListObj(category string, offset uint64, limit uint64, order string) ListRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%s%d", self.Key, category, tm))
	apiRsp := self.getData(fmt.Sprintf("/wdb/api/list?category=%s&offset=%d&limit=%d&order=%s&time=%d&sg=%s", category, offset, limit, order, tm, sg))
	if apiRsp.Code != 200 {
		logger.Println(apiRsp.Msg)
		return list_err(apiRsp.Msg)
	}

	var clist struct {
		Total uint64   `json:"total,omitempty"`
		List  []string `json:"list,omitempty"`
	}
	if err := json.Unmarshal([]byte(apiRsp.Data), &clist); err != nil {
		logger.Println(err)
		return list_err(fmt.Sprintf("%s", err))
	}

	return ListRsp{
		200,
		"",
		clist.Total,
		clist.List,
	}
}

func (self *Wdb) UploadByPath(path string, key string, categories []string) ApiRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%d%s", self.Key, tm, key))

	var req UploadReq
	req.Path = path
	req.Key = key
	req.Categories = categories
	req.Time = uint64(tm)
	req.Sg = sg

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/upload", str_req)
}

func (self *Wdb) DownToPath(path string, key string) ApiRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%d%s", self.Key, tm, key))

	var req DownReq
	req.Path = path
	req.Key = key
	req.Time = uint64(tm)
	req.Sg = sg

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/down", str_req)
}

func (self *Wdb) TransBegin(keys []string) ApiRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%d", self.Key, tm))

	var req TransBeginReq
	req.Keys = keys
	req.Time = uint64(tm)
	req.Sg = sg

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/trans/begin", str_req)
}

func (self *Wdb) TransCreateObj(tsid string, key string, data string, categories []string) ApiRsp {
	var req TransCreateDataReq
	req.Tsid = tsid
	req.Key = key
	req.Categories = categories
	req.Content = data
	tm := uint64(time.Now().Unix())
	req.Time = tm
	req.Sg = self.sign(fmt.Sprintf("%s%d%s%s%s", self.Key, tm, key, data, tsid))

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/trans/create", str_req)
}

func (self *Wdb) TransUpdateObj(tsid string, key string, data string) ApiRsp {
	var req TransUpdateDataReq
	req.Tsid = tsid
	req.Key = key
	req.Content = data
	tm := uint64(time.Now().Unix())
	req.Time = tm
	req.Sg = self.sign(fmt.Sprintf("%s%d%s%s%s", self.Key, tm, key, data, tsid))

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/trans/update", str_req)
}

func (self *Wdb) TransGet(tsid string, key string) ApiRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%d%s%s", self.Key, tm, key, tsid))
	return self.getData(fmt.Sprintf("/wdb/api/trans/get?tsid=%s&key=%s&time=%d&sg=%s", tsid, key, tm, sg))
}

func (self *Wdb) TransDelObj(tsid string, key string) ApiRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%d%s%s", self.Key, tm, key, tsid))
	return self.getData(fmt.Sprintf("/wdb/api/trans/del?tsid=%s&key=%s&time=%d&sg=%s", tsid, key, tm, sg))
}

func (self *Wdb) TransCommit(tsid string) ApiRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%d%s", self.Key, tm, tsid))

	var req TransIdReq
	req.Tsid = tsid
	req.Time = uint64(tm)
	req.Sg = sg

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/trans/commit", str_req)
}

func (self *Wdb) TransRollBack(tsid string) ApiRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%d%s", self.Key, tm, tsid))

	var req TransIdReq
	req.Tsid = tsid
	req.Time = uint64(tm)
	req.Sg = sg

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/trans/roll_back", str_req)
}

func (self *Wdb) CreateIndex(indexkeys []string, key string, indexraw []string) ApiRsp {
	var req CreateIndexReq
	req.Indexkey = indexkeys
	req.Key = key
	req.Indexraw = indexraw
	tm := uint64(time.Now().Unix())
	req.Time = tm
	req.Sg = self.sign(fmt.Sprintf("%s%d%s", self.Key, tm, key))

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/index/create", str_req)
}

func (self *Wdb) UpdateIndex(oindexkeys []string, cindexkeys []string, key string, indexraw []string) ApiRsp {
	var req UpdateIndexReq
	req.Oindexkey = oindexkeys
	req.Cindexkey = cindexkeys
	req.Key = key
	req.Indexraw = indexraw
	tm := uint64(time.Now().Unix())
	req.Time = tm
	req.Sg = self.sign(fmt.Sprintf("%s%d%s", self.Key, tm, key))

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/index/update", str_req)
}

func (self *Wdb) DelIndex(indexkeys []string, key string) ApiRsp {
	var req DelIndexReq
	req.Indexkey = indexkeys
	req.Key = key
	tm := uint64(time.Now().Unix())
	req.Time = tm
	req.Sg = self.sign(fmt.Sprintf("%s%d%s", self.Key, tm, key))

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/index/del", str_req)
}

func (self *Wdb) ListIndex(indexkey string, condition string, offset uint64, limit uint64, order string) ListRsp {
	var req ListIndexReq
	req.Indexkey = indexkey
	req.Condition = condition
	req.Offset = offset
	req.Limit = limit
	req.Order = order
	tm := uint64(time.Now().Unix())
	req.Time = tm
	req.Sg = self.sign(fmt.Sprintf("%s%s%d", self.Key, indexkey, tm))

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return list_err(fmt.Sprintf("%s", err))
	}

	apiRsp := self.postData("/wdb/api/index/list", str_req)
	if apiRsp.Code != 200 {
		logger.Println(apiRsp.Msg)
		return list_err(apiRsp.Msg)
	}

	var clist struct {
		Total uint64   `json:"total,omitempty"`
		List  []string `json:"list,omitempty"`
	}
	fmt.Println(apiRsp.Data)
	if err := json.Unmarshal([]byte(apiRsp.Data), &clist); err != nil {
		logger.Println(err)
		return list_err(fmt.Sprintf("%s", err))
	}

	return ListRsp{
		200,
		"",
		clist.Total,
		clist.List,
	}
}

func (self *Wdb) CreateRawData(key string, data []byte, categories []string) ApiRsp {
	content := base64.StdEncoding.EncodeToString(data)
	var req CreateDataReq
	req.Key = key
	req.Categories = categories
	req.Content = content
	tm := uint64(time.Now().Unix())
	req.Time = tm
	req.Sg = self.sign(fmt.Sprintf("%s%d%s%s", self.Key, tm, key, content))

	str_req, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return self.postData("/wdb/api/create_raw", str_req)
}

func (self *Wdb) GetRawData(key string) RawRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%d%s", self.Key, tm, key))
	path := fmt.Sprintf("/wdb/api/get_raw?key=%s&time=%d&sg=%s", key, tm, sg)
	apiRsp := self.getData(path)
	var rawRsp RawRsp
	if apiRsp.Code != 200 {
		logger.Println(apiRsp.Msg)
		rawRsp.Code = 500
		rawRsp.Msg = apiRsp.Msg
		return rawRsp
	}

	raw, err := base64.StdEncoding.DecodeString(apiRsp.Data)
	if err != nil {
		logger.Println(err)
		rawRsp.Code = 500
		rawRsp.Msg = fmt.Sprintf("%s", err)
		return rawRsp
	}

	rawRsp.Code = 200
	rawRsp.Size = uint64(bytes.Count(raw, nil) - 1)
	rawRsp.Raw = raw

	return rawRsp
}

func (self *Wdb) GetRangeData(key string, offset uint64, limit uint64) RawRsp {
	tm := time.Now().Unix()
	sg := self.sign(fmt.Sprintf("%s%d%s", self.Key, tm, key))
	path := fmt.Sprintf("/wdb/api/get_range?key=%s&offset=%d&limit=%d&time=%d&sg=%s", key, offset, limit, tm, sg)
	apiRsp := self.getData(path)
	var rawRsp RawRsp
	if apiRsp.Code != 200 {
		logger.Println(apiRsp.Msg)
		rawRsp.Code = 500
		rawRsp.Msg = apiRsp.Msg
		return rawRsp
	}

	var rgdata struct {
		AllSize uint64 `json:"all_size,omitempty"`
		Data    string `json:"data,omitempty"`
	}
	if err := json.Unmarshal([]byte(apiRsp.Data), &rgdata); err != nil {
		logger.Println(err)
		rawRsp.Code = 500
		rawRsp.Msg = fmt.Sprintf("%s", err)
		return rawRsp
	}

	raw, err := base64.StdEncoding.DecodeString(rgdata.Data)
	if err != nil {
		logger.Println(err)
		rawRsp.Code = 500
		rawRsp.Msg = fmt.Sprintf("%s", err)
		return rawRsp
	}

	rawRsp.Code = 200
	rawRsp.Msg = ""
	rawRsp.Size = rgdata.AllSize
	rawRsp.Raw = raw

	return rawRsp
}

func (self *Wdb) sign(text string) string {
	sum := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", sum)
}

func (self *Wdb) postData(path string, data []byte) ApiRsp {
	url := fmt.Sprintf("%s%s", self.Host, path)
	resp, err := http.Post(url,
		"application/json",
		bytes.NewBuffer(data))
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if len(body) == 0 {
		msg := "body is empty"
		logger.Println(msg)
		return rsp_err(errors.New(msg))
	}
	logger.Println(url, "------->", string(data), " === ", string(body))
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	crsp := ApiRsp{}
	if err := json.Unmarshal(body, &crsp); err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return crsp
}

func (self *Wdb) getData(path string) ApiRsp {
	url := fmt.Sprintf("%s%s", self.Host, path)
	resp, err := http.Get(url)
	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	blen := len(body)
	logger.Println(url, "---------->bac#", blen)

	if blen == 0 {
		msg := "body is empty"
		logger.Println(msg)
		return rsp_err(errors.New(msg))
	}

	if err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	rsp := ApiRsp{}
	if err := json.Unmarshal(body, &rsp); err != nil {
		logger.Println(err)
		return rsp_err(err)
	}

	return rsp
}
