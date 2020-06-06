package controllers

import (
	"encoding/json"
	"fmt"
	"helloweb/logger"
	"io/ioutil"
	"net/http"
)

// ReplyProto 后端响应数据通信协议
type ReplyProto struct {
	Status   int         `json:"status"` //状态 0正常，小于0出错，大于0可能有问题
	Msg      string      `json:"msg"`    //状态信息
	Data     interface{} `json:"data"`
	API      string      `json:"API"`    //api接口
	Method   string      `json:"method"` //post,put,get,delete
	SN       int         `json:"SN"`
	RowCount int         `json:"rowCount"` //Data若是数组，算其长度
}

// ReqProto 前端请求数据通讯协议
type ReqProto struct {
	Action   string              `json:"action"` //请求类型GET/POST/PUT/DELETE
	Data     interface{}         `json:"data"`   //请求数据
	Sets     []string            `json:"sets"`
	OrderBy  []map[string]string `json:"orderBy"`  //排序要求
	Filter   interface{}         `json:"filter"`   //筛选条件
	Page     int                 `json:"page"`     //分页
	PageSize int                 `json:"pageSize"` //分页大小
}


// ErrorResp 错误响应
func ErrorResp(w http.ResponseWriter, r *http.Request, errorMsg string) error {
	if w == nil {
		err := fmt.Errorf(" ErrorResp w is nil")
		logger.Z.Error(err.Error())
		return err
	}
	if r == nil {
		err := fmt.Errorf(" ErrorResp r is nil")
		logger.Z.Error(err.Error())
		return err
	}
	if errorMsg == "" {
		err := fmt.Errorf(" ErrorResp errorMsg is null")
		logger.Z.Error(err.Error())
		return err
	}

	// 发送响应数据
	resp := &ReplyProto{}
	resp.Status = -1
	resp.Msg = errorMsg
	resp.Data = nil
	resp.API = r.URL.String()
	resp.Method = r.Method
	resp.RowCount = 0
	response, err := json.Marshal(resp)
	if err != nil {
		err = fmt.Errorf("ErrorResp json.Marshal error:%v", err)
		logger.Z.Error(err.Error())
		return err
	}
	_, err = w.Write(response)
	if err != nil {
		err = fmt.Errorf("ErrorResp w.Write error:%v", err)
		logger.Z.Error(err.Error())
		return err
	}
	logger.Z.Info("错误发送成功")
	return nil
}

// SuccessResp 成功响应
func SuccessResp(w http.ResponseWriter, r *http.Request, result interface{}) error {
	if w == nil {
		err := fmt.Errorf(" SuccessResp w is nil")
		logger.Z.Error(err.Error())
		return err
	}
	if r == nil {
		err := fmt.Errorf(" SuccessRespr is nil")
		logger.Z.Error(err.Error())
		return err
	}
	// 发送响应数据
	resp := ReplyProto{}
	resp.Status = 0
	resp.Msg = ""
	resp.Data = result
	resp.API = r.URL.String()
	resp.Method = r.Method
	resp.RowCount = 0

	response, err := json.Marshal(resp)
	if err != nil {
		err = fmt.Errorf(" SuccessResp json.Marshal error:%v", err)
		logger.Z.Error(err.Error())
		return err
	}
	_, err = w.Write(response)
	if err != nil {
		err = fmt.Errorf(" SuccessResp w.Write error:%v", err)
		logger.Z.Error(err.Error())
		return err
	}
	return nil
}

// GetBodyData post请求数据处理//读取请求body 这里课配置
func GetBodyData(r *http.Request) (respData map[string]interface{}, err error) {
	if r == nil {
		err := fmt.Errorf("GetBodyData of argument r should not be nil")
		logger.Z.Error(err.Error())
		return nil, err
	}
	defer r.Body.Close()
	// 用来接收前端发送过来的数据的变量
	data := make(map[string]interface{})
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err := fmt.Errorf("GetBodyData ioutil.ReadAll error:%v", err)
		logger.Z.Error(err.Error())
		return nil, err
	}
	if len(body) == 0 {
		err := fmt.Errorf("没有传来数据，你想干啥子啊！")
		logger.Z.Error(err.Error())
		return nil, err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		err := fmt.Errorf("GetBodyData json.Unmarshal error:%v", err)
		logger.Z.Error(err.Error())
		return nil, err
	}
	respData, ok := data["data"].(map[string]interface{})
	if !ok {
		errMsg := fmt.Errorf("表单格式不正确,有没有好好用通信协议啊！！")
		logger.Z.Error(errMsg.Error())
		return nil, err
	}
	return respData, nil
}

// analyFilter  解析Filter
func analyFilter(reqFilter interface{}) ([]string, error) {
	var err error
	if reqFilter == nil {
		err = fmt.Errorf("reqFilter is nll")
		logger.Z.Error(err.Error())
		return nil, err
	}
	filter, ok := reqFilter.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("Filter must be {string:object")
		logger.Z.Error(err.Error())
		return nil, err
	}
	fmt.Println(filter["Status"])
	filterStatus, ok := filter["Status"].(map[string]interface{})
	if !ok {
		err = fmt.Errorf("filter['Status'] must be {string:object}")
		logger.Z.Error(err.Error())
		return nil, err
	}
	inStatus := filterStatus["IN"].([]interface{})
	if !ok {
		err = fmt.Errorf("filterStatus['IN'] must be []string")
		logger.Z.Error(err.Error())
		return nil, err
	}
	StatusSlice:= make([]string,0)
	for i:=0;i<len(inStatus);i++{
		data,ok := inStatus[i].(string)
		if !ok{
			err = fmt.Errorf("analyFilter inStatus[i] must string")
			logger.Z.Info(err.Error())
			return nil,err
		}
		StatusSlice = append(StatusSlice,data)
	}
	return StatusSlice,nil
}



