package s3

import (
	"bufio"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/beyondstorage/go-service-cos/v2"
	"github.com/beyondstorage/go-storage/v4/pairs"
	"github.com/beyondstorage/go-storage/v4/types"
	"github.com/let-sh/cli/log"
	"github.com/let-sh/cli/requests"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/sirupsen/logrus"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func NewCos() (types.Storager, error) {
	pwd, _ := os.Getwd()
	return cos.NewStorager(
		// work_dir: https://beyondstorage.io/docs/go-storage/pairs/work_dir
		//
		// Relative operations will be based on this WorkDir.
		pairs.WithWorkDir(pwd),
		// credential: https://beyondstorage.io/docs/go-storage/pairs/credential
		//
		// Credential could be fetched from service's console: https://console.cloud.tencent.com/cam/capi
		//
		// Example Value: hmac:access_key_id:secret_access_key
		pairs.WithCredential(os.Getenv("STORAGE_COS_CREDENTIAL")),
		// location: https://beyondstorage.io/docs/go-storage/pairs/location
		//
		// Available location: https://cloud.tencent.com/document/product/436/6224
		pairs.WithLocation(os.Getenv("STORAGE_COS_LOCATION")),
		// name: https://beyondstorage.io/docs/go-storage/pairs/name
		//
		// name is the bucket name.
		pairs.WithName(os.Getenv("STORAGE_COS_NAME")),
	)
}

func UploadFileToS3CodeSource(filedir, filename, projectName string, cn bool) {
	// create and start new bar
	fi, _ := os.Stat(filedir)

	file, _ := os.Open(filedir)
	r := bufio.NewReader(file)

	p := mpb.New(
		mpb.WithWidth(64),
		mpb.WithRefreshRate(200*time.Millisecond),
	)

	bar = p.AddBar(fi.Size(),
		mpb.PrependDecorators(
			decor.Name("uploading file: "),
			//decor.Counters(decor.UnitKiB, "% .1f / % .1f"),
		),

		//mpb.NewBarFiller(mpb.BarStyle("[=>-|")),
		mpb.PrependDecorators(
			decor.CountersKiloByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
			decor.Name(" ] "),
			//decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.EwmaSpeed(decor.UnitKB, "% .2f", 1024),
		),
		mpb.BarRemoveOnComplete(),
	)
	bar.SetTotal(fi.Size(), false)

	stsToken, err := requests.GetStsToken("buildBundle", projectName, cn)
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

	logrus.WithFields(logrus.Fields{
		"objKey": filename,
	}).Debug("put object from file")

	err = bucket.PutObject(filename, proxyReader, oss.Progress(&OssProgressListener{}))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	bar.Completed()
}

func UploadDirToS3(dirPath, projectName, bundleID string, cn bool) error {
	log.BPause()
	uploadStatus = make(map[string]fileUplaodStatus)
	stsToken, err := requests.GetStsToken("static", projectName, cn)
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

	// Read directory files
	var names []string
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			names = append(names, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// respect .gitignore and .letignore
	if _, err := os.Stat(filepath.Join(dirPath, ".gitignore")); err == nil {
		// match a file against a particular .gitignore
		i, _ := ignore.CompileIgnoreFile(filepath.Join(dirPath, ".gitignore"))

		tmp := []string{}
		for _, v := range names {
			if !i.MatchesPath(v) {
				tmp = append(tmp, v)
			}
		}

		names = tmp
	}

	// .letignore
	if _, err := os.Stat(filepath.Join(dirPath + ".letignore")); err == nil {
		// match a file against a particular .gitignore
		i, _ := ignore.CompileIgnoreFile(filepath.Join(dirPath + ".letignore"))

		tmp := []string{}
		for _, v := range names {
			if !i.MatchesPath(v) {
				fi, _ := os.Stat(v)
				mutex.Lock()
				// register upload status
				uploadStatus[v] = struct {
					FilePath     string
					ConsumedSize int64
					TotalSize    int64
				}{FilePath: fi.Name(), ConsumedSize: 0, TotalSize: fi.Size()}
				mutex.Unlock()

				tmp = append(tmp, v)
			}
		}

		names = tmp
	}

	// fill in files info
	var totalFilesSize int64
	for _, v := range names {
		fi, _ := os.Stat(v)
		mutex.Lock()
		// register upload status
		uploadStatus[v] = struct {
			FilePath     string
			ConsumedSize int64
			TotalSize    int64
		}{FilePath: fi.Name(), ConsumedSize: 0, TotalSize: fi.Size()}

		totalFilesSize += fi.Size()
		mutex.Unlock()
	}
	status := uploadStatus
	logrus.Debug(status)

	p := mpb.New(
		mpb.WithWidth(64),
		mpb.WithRefreshRate(200*time.Millisecond),
	)

	// init progress bar
	{
		bar = p.AddBar(totalFilesSize,
			mpb.PrependDecorators(
				decor.Name("uploading files: "),
				//decor.Counters(decor.UnitKiB, "% .1f / % .1f"),
			),

			//mpb.NewBarFiller(mpb.BarStyle("[=>-|")),
			mpb.PrependDecorators(
				decor.CountersKiloByte("% .2f / % .2f"),
			),
			mpb.AppendDecorators(
				decor.Percentage(),
				decor.Name(" ] "),
				//decor.EwmaETA(decor.ET_STYLE_GO, 90),
				//decor.EwmaSpeed(decor.UnitKB, "% .2f", 1024),
			),
			mpb.BarRemoveOnComplete(),
		)
		bar.SetTotal(totalFilesSize, false)
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

				objKey := filepath.Join(bundleID, strings.Replace(name, dirPath, "", 1))
				filePath := name

				// skip dir
				fi, err := os.Stat(filePath)

				if err != nil {
					fmt.Println(err)
					resChan <- &err
					return
				}
				if fi.IsDir() {
					resChan <- &err
					return
				}

				logrus.WithFields(logrus.Fields{
					"objKey":   objKey,
					"filePath": filePath,
				}).Debug("put object from file")

				// TODO:
				// * check file exists in previous deployment
				// * if matched etag, copy file
				// * else upload file
				err = bucket.PutObjectFromFile(func() string {
					if runtime.GOOS == "windows" {
						return filepath.ToSlash(objKey)
					}
					return objKey
				}(), filePath, oss.Progress(&OssProgressListener{filepath: filePath, totalFilesSize: totalFilesSize, currentTime: time.Now()}))
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

	bar.Completed()
	bar.Abort(true)
	log.BUnpause()
	log.S.Suffix(" deploying ")
	return nil
}
