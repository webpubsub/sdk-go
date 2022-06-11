package webpubsub

// WPSPublishMessage is the part of the message struct used in Publish File
type WPSPublishMessage struct {
	Text string `json:"text"`
}

// WPSFileInfoForPublish is the part of the message struct used in Publish File
type WPSFileInfoForPublish struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// WPSPublishFileMessage is the message struct used in Publish File
type WPSPublishFileMessage struct {
	WPSMessage *WPSPublishMessage     `json:"message"`
	WPSFile    *WPSFileInfoForPublish `json:"file"`
}

// WPSFileInfo is the File Upload API struct returned on for each file.
type WPSFileInfo struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	Size    int    `json:"size"`
	Created string `json:"created"`
}

// WPSFileData is used in the responses to show File ID
type WPSFileData struct {
	ID string `json:"id"`
}

// WPSFileUploadRequest is used to store the info related to file upload to S3
type WPSFileUploadRequest struct {
	URL        string         `json:"url"`
	Method     string         `json:"method"`
	FormFields []WPSFormField `json:"form_fields"`
}

// WPSFormField is part of the struct used in file upload to S3
type WPSFormField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// WPSFileDetails is used in the responses to show File Info
type WPSFileDetails struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	URL  string
}

// WPSFileMessageAndDetails is used to store the file message and file info
type WPSFileMessageAndDetails struct {
	WPSMessage WPSPublishMessage `json:"message"`
	WPSFile    WPSFileDetails    `json:"file"`
}

// ParseFileInfo is a function extract file info and add to the struct WPSFileMessageAndDetails
func ParseFileInfo(filesPayload map[string]interface{}) (WPSFileDetails, WPSPublishMessage) {
	var data map[string]interface{}
	resp := &WPSFileMessageAndDetails{}
	resp.WPSMessage = WPSPublishMessage{}
	resp.WPSFile = WPSFileDetails{}

	//"message":{"text":"test file"},"file":{"name":"test_file_upload_name_32899","id":"9076246e-5036-42af-b3a3-767b514c93c8"}}
	if o, ok := filesPayload["file"]; ok {
		if o != nil {
			data = o.(map[string]interface{})
			if d, ok := data["id"]; ok {
				resp.WPSFile.ID = d.(string)
			}
			if d, ok := data["name"]; ok {
				resp.WPSFile.Name = d.(string)
			}
		}
	}
	if m, ok := filesPayload["message"]; ok {
		if m != nil {
			data = m.(map[string]interface{})
			if d, ok := data["text"]; ok {
				resp.WPSMessage.Text = d.(string)
			}
		}
	}
	return resp.WPSFile, resp.WPSMessage
}
