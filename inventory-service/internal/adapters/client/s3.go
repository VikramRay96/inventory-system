package client

type IS3 interface {
	UploadFile(serviceName string, serviceType string, topicName string, topicId string, resource []byte, fileRequestId string, contentType string, fileExtension string, fType string) error
}
