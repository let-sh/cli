package firebase

//
//import (
//	"context"
//	"fmt"
//	"google.golang.org/api/firebaseremoteconfig/v1"
//)
//
//var service *firebaseremoteconfig.Service
//
//func init() {
//
//}
//
//func GetConfig() {
//	var err error
//	ctx := context.Background()
//	service, err = firebaseremoteconfig.NewService(ctx)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	conf, err := service.Projects.GetRemoteConfig("letcli").Do()
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	fmt.Println(conf)
//}
