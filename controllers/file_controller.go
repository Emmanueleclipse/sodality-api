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
