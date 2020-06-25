package download

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/cheggaaa/pb/v3"
)

func downloadFile(downloadURL string) ([]byte, error) {
	res, err := http.Get(downloadURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	downloadSize, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	bar := pb.Full.Start64(int64(downloadSize))
	barReader := bar.NewProxyReader(res.Body)
	io.Copy(&buf, barReader)
	bar.Finish()
	file := buf.Bytes()
	if len(file) == 0 {
		return nil, errors.New("soundcloud-cli: downloaded track corrupted")
	}
	return file, nil
}
