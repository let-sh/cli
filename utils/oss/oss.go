package oss

import (
	"bufio"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/denormal/go-gitignore"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var bar *mpb.Bar

func UploadFileToCodeSource(filepath, filename, projectName string) {
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
		mpb.BarRemoveOnComplete(),
	)

	stsToken, err := requests.GetStsToken("buildBundle", projectName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 创建OSSClient实例
	endpoint := strings.Join(strings.Split(stsToken.Host, ".")[1:], ".")
	client, err := oss.New(endpoint, stsToken.AccessKeyID, stsToken.AccessKeySecret, oss.SecurityToken(stsToken.SecurityToken))
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

	err = bucket.PutObject(filename, proxyReader, oss.Progress(&OssProgressListener{}))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	bar.Completed()
}

func UploadDirToStaticSource(dirPath, projectName, bundleID string) error {
	log.BPause()
	stsToken, err := requests.GetStsToken("static", projectName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	// 创建OSSClient实例
	endpoint := strings.Join(strings.Split(stsToken.Host, ".")[1:], ".")
	client, err := oss.New(endpoint, stsToken.AccessKeyID, stsToken.AccessKeySecret, oss.SecurityToken(stsToken.SecurityToken))

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

	// Read directory and close.
	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	dir.Close()

	// respect .gitignore and .letignore
	if _, err := os.Stat(dirPath + ".gitignore"); err == nil {
		// match a file against a particular .gitignore
		ignore, _ := gitignore.NewFromFile(dirPath + ".gitignore")

		tmp := []string{}
		for _, v := range names {
			match := ignore.Match(v)
			if match != nil {
				if !match.Ignore() {
					tmp = append(tmp, v)
				}
			}
		}

		names = tmp
	}

	// .letignore
	if _, err := os.Stat(dirPath + ".letignore"); err == nil {
		// match a file against a particular .gitignore
		ignore, _ := gitignore.NewFromFile(dirPath + ".gitignore")

		tmp := []string{}
		for _, v := range names {
			match := ignore.Match(v)
			if match != nil {
				if !match.Ignore() {
					tmp = append(tmp, v)
				}
			}
		}

		names = tmp
	}

	// Copy names to a channel for workers to consume. Close the
	// channel so that workers stop when all work is complete.

	namesChan := make(chan string, len(names))
	for _, name := range names {
		namesChan <- name
	}
	close(namesChan)

	// Create a maximum of 8 workers

	workers := 8
	if len(names) < workers {
		workers = len(names)
	}

	errChan := make(chan error, 1)
	resChan := make(chan *error, len(names))

	// Run workers

	for i := 0; i < workers; i++ {
		go func() {
			// Consume work from namesChan. Loop will end when no more work.
			for name := range namesChan {
				if err != nil {
					select {
					case errChan <- err:
						// will break parent goroutine out of loop
					default:
						// don't care, first error wins
					}
					return
				}
				objKey := bundleID + "/" + name
				filePath := filepath.Join(dirPath, name)
				err = bucket.PutObjectFromFile(objKey, filePath)

				if err != nil {
					select {
					case errChan <- err:
						// will break parent goroutine out of loop
					default:
						// don't care, first error wins
					}
					return
				}
				resChan <- &err
			}
		}()
	}

	// Collect results from workers
	for i := 0; i < len(names); i++ {
		select {
		case res := <-resChan:
			// collect result
			_ = res
		case err := <-errChan:
			return err
		}
	}
	log.S.Suffix(" deploying ")
	log.BUnpause()
	return nil
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
