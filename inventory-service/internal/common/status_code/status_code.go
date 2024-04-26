package status_code

import "inventory-system/common/pkg/dto"

type StatusCode string

const (
	IMS101 dto.StatusCode = "IMS101:Error occurred while saving configuration"
	IMS102 dto.StatusCode = "IMS102:Error while creating collection"
	IMS103 dto.StatusCode = "IMS103:Error while creating index"
	IMS104 dto.StatusCode = "IMS104:Inventory Configuration already exists with given name"
	IMS105 dto.StatusCode = "IMS105:Error while fetching inventory configuration"
	IMS106 dto.StatusCode = "IMS106:Error while fetching all inventory configuration"
	IMS306 dto.StatusCode = "IMS306:Error while fetching inventory item"
	IMS107 dto.StatusCode = "IMS107:Error while deleting inventory configuration"
	IMS108 dto.StatusCode = "IMS108:Inventory already exists with given identifier"
	IMS109 dto.StatusCode = "IMS109:Document Validation failed"
	IMS110 dto.StatusCode = "IMS110:Error while fetching configuration"
	IMS112 dto.StatusCode = "IMS112:Multiple filters not allowed"
	IMS113 dto.StatusCode = "IMS113:Filter attribute not present in inventory identifiers"
	IMS115 dto.StatusCode = "IMS115:Filter attribute not present in unique inventory identifiers"
	IMS116 dto.StatusCode = "IMS116:Invalid Subject Id Provided"
	IMS117 dto.StatusCode = "IMS117:Invalid Topic Id Provided"
	IMS118 dto.StatusCode = "IMS118:Invalid Lesson Name or Subject Id"
	IMS119 dto.StatusCode = "IMS119:Invalid Resource Id Provided"
	IMS114 dto.StatusCode = "IMS114:Inventory Configuration not found"

	IMS200 dto.StatusCode = "IMS200:success"
	IMS204 dto.StatusCode = "IMS204:Inventory Configuration deleted"
	IMS205 dto.StatusCode = "IMS205:Inventory Configuration already deleted"
	IMS400 dto.StatusCode = "IMS400:Bad request"
	IMS404 dto.StatusCode = "IMS404:Not found"
	IMS500 dto.StatusCode = "IMS500:Internal server error"
	IMS600 dto.StatusCode = "IMS600:Unable to find Resource details for the provided ID"
)
