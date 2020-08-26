package main

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
)

type ReqParam struct {
	ParameterName string
	ParameterType string
	Required      string
	Description   string
}

type RspParam struct {
	ParameterName string
	ParameterType string
	Description   string
}

var licenseKey = ``
var company = ""
var swaggerFile = ""
var docxFile = ""

func init() {
	err := license.SetLicenseKey(licenseKey, company)
	if err != nil {
		panic(err)
	}
}

func main() {
	// read swagger json
	dataByte, err := ioutil.ReadFile(swaggerFile)
	if err != nil {
		return
	}

	// str to map
	dataMap, err := jsonToMap(dataByte)
	if err != nil {
		return
	}

	doc := document.New()
	defer doc.Close()

	//  Title
	para := doc.AddParagraph()
	run := para.AddRun()
	para.SetStyle("Title")
	run.AddText(dataMap["info"].(map[string]interface{})["title"].(string))

	doc.AddParagraph().AddRun().AddText("")
	doc.AddParagraph().AddRun().AddText("")
	doc.AddParagraph().AddRun().AddText("")

	// group api by tag
	tag := make(map[string][]string)
	for k, v := range dataMap["paths"].(map[string]interface{}) {
		key := reflect.ValueOf(v.(map[string]interface{})).MapKeys()[0].String()
		value := v.(map[string]interface{})[key].(map[string]interface{})

		tagKey := value["tags"].([]interface{})[0].(string)
		if _, ok := tag[tagKey]; !ok {
			s := []string{k}
			tag[tagKey] = s
		} else {
			tag[tagKey] = append(tag[tagKey], k)
		}
	}

	num := 1
	for k, v := range tag {
		// tag
		para = doc.AddParagraph()
		para.SetStyle("Heading1")
		run = para.AddRun()
		run.AddText(fmt.Sprintf("%d . %s", num, k))

		for i, value := range v {
			api := dataMap["paths"].(map[string]interface{})[value].(map[string]interface{})
			key := reflect.ValueOf(api).MapKeys()[0].String()

			// summary
			para := doc.AddParagraph()
			para.SetStyle("Heading2")
			run := para.AddRun()
			run.AddText(fmt.Sprintf("%d . %s", i+1, api[key].(map[string]interface{})["summary"].(string)))

			// description
			para = doc.AddParagraph()
			run = para.AddRun()
			run.AddText(fmt.Sprintf("Description : %s", api[key].(map[string]interface{})["description"].(string)))

			// path
			para = doc.AddParagraph()
			run = para.AddRun()
			run.AddText(fmt.Sprintf("Path : %s", value))

			// content-type
			para = doc.AddParagraph()
			run = para.AddRun()
			run.AddText(fmt.Sprintf("Content-Type : %s", api[key].(map[string]interface{})["consumes"].([]interface{})[0].(string)))

			doc.AddParagraph().AddRun().AddText("")

			// request params
			var reqData []ReqParam
			params := api[key].(map[string]interface{})["parameters"].([]interface{})[0].(map[string]interface{})
			if _, ok := params["schema"]; ok {
				reqStr := params["schema"].(map[string]interface{})["$ref"].(string)[14:]
				reqData = searchReqData(dataMap, reqStr)
			}
			requestParam(doc, reqData)

			doc.AddParagraph().AddRun().AddText("")

			// response params
			var rspData []RspParam
			response := api[key].(map[string]interface{})["responses"].(map[string]interface{})["200"].(map[string]interface{})
			if _, ok := response["schema"].(map[string]interface{})["$ref"]; ok {
				rspStr := response["schema"].(map[string]interface{})["$ref"].(string)[14:]
				rspData = searchRspData(dataMap, rspStr, "")
			}
			responseParam(doc, rspData)

			doc.AddParagraph().AddRun().AddText("")
		}

		doc.AddParagraph().AddRun().AddText("")
		num++
	}

	if err := doc.SaveToFile(docxFile); err != nil {
		fmt.Printf("swag to docx failed err : %s \n", err.Error())
	}
}
