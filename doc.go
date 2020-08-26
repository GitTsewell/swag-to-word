package main

import (
	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
)

func requestParam(doc *document.Document, param []ReqParam) {
	// request parameter list
	para := doc.AddParagraph()
	run := para.AddRun()
	run.AddText("Request parameter list")

	table := doc.AddTable()
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)

	hdrRow := table.AddRow()

	cell := hdrRow.AddCell()
	cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
	cellPara := cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcLeft)
	cellPara.AddRun().AddText("Parameter Name")

	cell = hdrRow.AddCell()
	cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("Parameter Type")

	cell = hdrRow.AddCell()
	cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("Required")

	cell = hdrRow.AddCell()
	cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcRight)
	cellPara.AddRun().AddText("Description")

	for _, v := range param {
		row := table.AddRow()
		row.AddCell().AddParagraph().AddRun().AddText(v.ParameterName)
		row.AddCell().AddParagraph().AddRun().AddText(v.ParameterType)
		row.AddCell().AddParagraph().AddRun().AddText(v.Required)
		row.AddCell().AddParagraph().AddRun().AddText(v.Description)
	}
}

func responseParam(doc *document.Document, param []RspParam) {
	// request parameter list
	para := doc.AddParagraph()
	run := para.AddRun()
	run.AddText("Response parameter list")

	table := doc.AddTable()
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)

	hdrRow := table.AddRow()

	cell := hdrRow.AddCell()
	cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
	cellPara := cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcLeft)
	cellPara.AddRun().AddText("Parameter Name")

	cell = hdrRow.AddCell()
	cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcCenter)
	cellPara.AddRun().AddText("Parameter Type")

	cell = hdrRow.AddCell()
	cell.Properties().SetShading(wml.ST_ShdSolid, color.LightGray, color.Auto)
	cellPara = cell.AddParagraph()
	cellPara.Properties().SetAlignment(wml.ST_JcRight)
	cellPara.AddRun().AddText("Description")

	for _, v := range param {
		row := table.AddRow()
		row.AddCell().AddParagraph().AddRun().AddText(v.ParameterName)
		row.AddCell().AddParagraph().AddRun().AddText(v.ParameterType)
		row.AddCell().AddParagraph().AddRun().AddText(v.Description)
	}
}

func searchReqData(dataMap map[string]interface{}, key string) (rsp []ReqParam) {
	param := dataMap["definitions"].(map[string]interface{})[key].(map[string]interface{})
	for key, value := range param["properties"].(map[string]interface{}) {
		var data ReqParam
		data.ParameterName = key
		data.ParameterType = value.(map[string]interface{})["type"].(string)
		if description, ok := value.(map[string]interface{})["description"]; ok {
			data.Description = description.(string)
		}

		if _, ok := param["required"]; !ok {
			data.Required = "O"
		} else {
			if find(param["required"].([]interface{}), key) {
				data.Required = "M"
			} else {
				data.Required = "O"
			}
		}

		rsp = append(rsp, data)
	}
	return
}

func searchRspData(dataMap map[string]interface{}, key string, exStr string) (rsp []RspParam) {
	param := dataMap["definitions"].(map[string]interface{})[key].(map[string]interface{})
	for key, value := range param["properties"].(map[string]interface{}) {
		var data RspParam
		data.ParameterName = exStr + key
		data.ParameterType = value.(map[string]interface{})["type"].(string)
		if description, ok := value.(map[string]interface{})["description"]; ok {
			data.Description = description.(string)
		}

		rsp = append(rsp, data)

		// type eq array, then find object
		if data.ParameterType == "array" {
			arrayKey := value.(map[string]interface{})["items"].(map[string]interface{})["$ref"].(string)[14:]
			arrayRsp := searchRspData(dataMap, arrayKey, "+")
			rsp = append(rsp, arrayRsp...)
		}
	}
	return
}

func find(slice []interface{}, val string) bool {
	for _, item := range slice {
		if item.(string) == val {
			return true
		}
	}
	return false
}
