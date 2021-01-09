package oss

import (
	"bufio"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
	"os"
	"strings"
	"time"
)

var bar *mpb.Bar

func UploadFileToCodeSource(filepath, filename, projectName string) {
	log.BPause()

	// create and start new bar
	fi, _ := os.Stat(filepath)

	file, _ := os.Open(filepath)
	r := bufio.NewReader(file)

	p := mpb.New(
		mpb.WithWidth(64),
		mpb.WithRefreshRate(200*time.Millisecond),
	)

	bar = p.AddBar(fi.Size(), mpb.BarStyle("[=>-|"),
		mpb.PrependDecorators(
			decor.CountersKiloByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKB, "% .2f", 1024),
		),
	)

	stsToken, err := requests.GetStsToken("buildBundle", projectName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 创建OSSClient实例
	endpoint := strings.Join(strings.Split(stsToken.Host, ".")[1:], ".")
	client, err := oss.New(endpoint, stsToken.AccessKeyID, stsToken.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	bucketName := strings.Replace(strings.Split(stsToken.Host, ".")[0], "https://", "", 1)

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// create proxy reader
	proxyReader := bar.ProxyReader(r)
	defer proxyReader.Close()

	bucket.PutObject(filename, proxyReader, oss.Progress(&OssProgressListener{}))

	bar.Abort(true)

	log.S.Suffix(" deploying ")
	log.BUnpause()
}

type OssProgressListener struct {
}

func (listener *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	//bar.SetTotal(event.TotalBytes, false)
	//bar.SetCurrent(event.ConsumedBytes)
	switch event.EventType {
	case oss.TransferStartedEvent:
		//fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
		//	event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
		//fmt.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
		//	event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)

	case oss.TransferCompletedEvent:
		//fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
		//	event.ConsumedBytes, event.TotalBytes)
	case oss.TransferFailedEvent:
		//fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
		//	event.ConsumedBytes, event.TotalBytes)
	default:
	}
}
