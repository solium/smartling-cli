package smartling

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/joeshaw/iso8601"
)

type UploadRequest struct {
	FileUri          string            `url:"fileUri,"`
	Approved         bool              `url:"approved,omitempty"`
	LocalesToApprove []string          `url:"localesToApprove,omitempty"`
	FileType         FileType          `url:"fileType"`
	CallbackUrl      *url.URL          `url:"callbackUrl,omitempty"`
	ParserConfig     map[string]string `url:"-"`
}

func (u *UploadRequest) AddParserConfig(k, v string) {
	u.ParserConfig[k] = v
}

func (u *UploadRequest) urlValues() (url.Values, error) {
	uv, err := query.Values(u)
	if err != nil {
		return nil, err
	}
	for k, v := range u.ParserConfig {
		uv.Add("smartling."+k, v)
	}

	return uv, nil
}

type UploadResponse struct {
	OverWritten bool `json:"overWritten"`
	StringCount int  `json:"stringCount"`
	WordCount   int  `json:"wordCount"`
}

func newfileUploadRequest(url string, uv url.Values, formFileKey, formFilePath string) (*http.Request, error) {
	f, err := os.Open(formFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for k, uvv := range uv {
		for _, v := range uvv {
			err = writer.WriteField(k, v)
			if err != nil {
				return nil, err
			}
		}
	}

	part, err := writer.CreateFormFile("file", filepath.Base(formFilePath))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, f)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, err
}

func (c *Client) Upload(localFilePath string, req *UploadRequest) (*UploadResponse, error) {
	v, err := req.urlValues()
	if err != nil {
		return nil, err
	}
	c.addCredentials(&v)

	request, err := newfileUploadRequest(c.BaseUrl+"/file/upload", v, "file", localFilePath)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	ur := UploadResponse{}
	err = unmarshalResponseData(response, &ur)

	return &ur, err
}

type GetRequest struct {
	FileUri                string        `url:"fileUri"`
	Locale                 string        `url:"locale,omitempty"`
	RetrievalType          RetrievalType `url:"retrievalType,omitempty"`
	IncludeOriginalStrings bool          `url:"includeOriginalStrings,omitempty"`
}

func (c *Client) Get(req *GetRequest) ([]byte, error) {
	v, err := query.Values(req)
	if err != nil {
		return nil, err
	}
	c.addCredentials(&v)

	response, err := c.httpClient.PostForm(c.BaseUrl+"/file/get", v)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		err := unmarshalResponseData(response, nil)
		if err != nil {
			return nil, err
		}
	}

	return ioutil.ReadAll(response.Body)
}

type ListCondition string

const (
	HaveAtLeastOneUnapproved ListCondition = "haveAtLeastOneUnapproved"
	HaveAtLeastOneApproved   ListCondition = "haveAtLeastOneApproved"
	HaveAtLeastOneTranslated ListCondition = "haveAtLeastOneTranslated"
	HaveAllTranslated        ListCondition = "haveAllTranslated"
	HaveAllApproved          ListCondition = "haveAllApproved"
	HaveAllUnapproved        ListCondition = "haveAllUnapproved"
)

type ListRequest struct {
	Locale            string          `url:"locale,omitempty"`
	UriMask           string          `url:"uriMask,omitempty"`
	FileTypes         []FileType      `url:"fileTypes,omitempty"`
	LastUploadedAfter time.Time       `url:"lastUploadedAfter,omitempty"`
	Offset            int             `url:"offset,omitempty"`
	Limit             int             `url:"limit,omitempty"`
	Conditions        []ListCondition `url:"conditions,omitempty"`
	OrderBy           string          `url:"orderBy,omitempty"`
}

type File struct {
	FileUri              string       `json:"fileUri"`
	StringCount          int          `json:"stringCount"`
	WordCount            int          `json:"wordCount"`
	ApprovedStringCount  int          `json:"approvedStringCount"`
	CompletedStringCount int          `json:"completedStringCount"`
	LastUploaded         iso8601.Time `json:"lastUploaded"`
	FileType             FileType     `json:"fileType"`
}

type ListResponse struct {
	FileCount int    `json:"fileCount"`
	Files     []File `json:"fileList"`
}

func (c *Client) List(req ListRequest) ([]File, error) {
	r := ListResponse{}
	err := c.doRequestAndUnmarshalData("/file/list", req, &r)

	return r.Files, err
}

type statusRequest struct {
	FileUri string `url:"fileUri"`
	Locale  string `url:"locale"`
}

func (c *Client) Status(fileUri, locale string) (File, error) {
	req := statusRequest{fileUri, locale}
	r := File{}
	err := c.doRequestAndUnmarshalData("/file/status", req, &r)

	return r, err
}

type renameRequest struct {
	FileUri    string `url:"fileUri"`
	NewFileUri string `url:"newFileUri"`
}

func (c *Client) Rename(oldFileUri, newFileUri string) error {
	req := renameRequest{oldFileUri, newFileUri}
	return c.doRequestAndUnmarshalData("/file/rename", req, nil)
}

type deleteRequest struct {
	FileUri string `url:"fileUri"`
}

func (c *Client) Delete(fileUri string) error {
	req := deleteRequest{fileUri}
	return c.doRequestAndUnmarshalData("/file/delete", req, nil)
}

type LastModifiedRequest struct {
	FileUri           string `url:"fileUri"`
	LastModifiedAfter string `url:"lastModifiedAfter,omitempty"`
	Locale            string `url:"locale"`
}

type LastModifiedItem struct {
	Locale       string       `json:"locale"`
	LastModified iso8601.Time `json:"lastModified"`
}

type LastModifiedResponse struct {
	Items []LastModifiedItem `json:"items"`
}

func (c *Client) LastModified(req LastModifiedRequest) ([]LastModifiedItem, error) {
	r := LastModifiedResponse{}
	err := c.doRequestAndUnmarshalData("/file/last_modified", req, &r)

	return r.Items, err
}
