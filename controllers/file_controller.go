package controllers

// import (
// 	"context"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	middlewares "sodality/handlers"
// 	"sodality/ipfs"
// 	"sodality/models"

// 	iface "github.com/ipfs/interface-go-ipfs-core"
// )

// var Session iface.CoreAPI

// func init() {
// 	Session = ipfs.StartingIPFS(context.Background())
// 	ipfs.ConnectPeers(context.Background(), Session)
// }

// var UploadFile = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

// 	file, header, err := r.FormFile("file")
// 	if err != nil {
// 		middlewares.ErrorResponse("error uploading file:"+err.Error(), rw)
// 		return
// 	}
// 	defer file.Close()

// 	data, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		middlewares.ServerErrResponse(err.Error(), rw)
// 		return
// 	}

// 	path, err := ipfs.AddFilePath(data, r.Context(), Session)
// 	if err != nil {
// 		middlewares.ServerErrResponse(err.Error(), rw)
// 		return
// 	}

// 	var resp models.FileResp

// 	resp.IpfsURL = fmt.Sprintf("https://ipfs.io/ipfs/%v", path.Cid())
// 	resp.Filename = header.Filename
// 	resp.FileSize = header.Size

// 	middlewares.SuccessRespond(resp, rw)
// })

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	middlewares "sodality/handlers"
	"sodality/models"
)

var UploadFile = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("file")
	if err != nil {
		middlewares.ErrorResponse("error uploading file:"+err.Error(), rw)
		return
	}

	url := "http://18.117.141.216:2000/api/v1/file/upload"

	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err = header.Open()
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	defer file.Close()

	part1, err := writer.CreateFormFile("file", filepath.Base(header.Filename))
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	_, err = io.Copy(part1, file)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var resp models.FileResp
	err = json.Unmarshal(body, &resp)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	middlewares.SuccessArrRespond(resp.IpfsURL, rw)
})
