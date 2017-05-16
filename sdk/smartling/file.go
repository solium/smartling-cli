
package smartling


import (
	"net/url"
	"encoding/json"
	"fmt"
	"log"
)

// API endpoints
const (
	fileApiList    = "/files-api/v2/projects/%v/files/list"
)

type FileType string

const (
	Android        FileType = "android"
	Ios            FileType = "ios"
	Gettext        FileType = "gettext"
	Html           FileType = "html"
	JavaProperties FileType = "javaProperties"
	Yaml           FileType = "yaml"
	Xliff          FileType = "xliff"
	Xml            FileType = "xml"
	Json           FileType = "json"
	Docx           FileType = "docx"
	Pptx           FileType = "pptx"
	Xlsx           FileType = "xlsx"
	Idml           FileType = "idml"
	Qt             FileType = "qt"
	Resx           FileType = "resx"
	Plaintext      FileType = "plaintext"
	Csv            FileType = "csv"
	Stringsdict    FileType = "stringsdict"
)

type FileStatus struct {
	FileUri              string
	LastUploaded         *Iso8601Time
	FileType             FileType
	HasInstructions      bool
}

// list files api call response data
type FilesList struct {
	TotalCount int64
	Items      []FileStatus
}

type FileListRequest struct {
	UriMask            string
	FileTypes          []FileType
	LastUploadedAfter  Iso8601Time
	LastUploadedBefore Iso8601Time
	Offset             int64
	Limit              int64
}

// list files for a project
func (c *Client) ListFiles(projectId string, listRequest FileListRequest) (list FilesList, err error) {

	header, err := c.auth.AccessHeader(c)
	if err != nil {
		return
	}

	// prepare the url
	urlObject, err := url.Parse(c.baseUrl + fmt.Sprintf(fileApiList, projectId))
	if err != nil {
		return
	}
	urlObject.RawQuery = listRequest.RawQuery()

	log.Printf(urlObject.String())
	bytes, statusCode, err := c.doGetRequest(urlObject.String(), header)
	if err != nil {
		return
	}

	if statusCode != 200 {
		err = fmt.Errorf("ListFiles call returned unexpected status code: %v", statusCode)
		return
	}

	// unmarshal transport header
	apiResponse, err := unmarshalTransportHeader(bytes)
	if err != nil {
		return
	}

	// unmarshal projects array
	err = json.Unmarshal(apiResponse.Response.Data, &list)
	if err != nil {
		// TODO: special error here
		return
	}

	log.Printf("List files - received %v status code", statusCode)

	return
}
