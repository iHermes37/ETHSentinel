package RealtimeMonitor

//import (
//	"context"
//	"github.com/machinebox/graphql"
//)
//
//query:=""
//
//type NewPair_ResponseData struct{
//
//}
//
//
//// 泛型函数
//func FetchGraphQL[T any](client *graphql.Client, query string) (*T, error) {
//	req := graphql.NewRequest(query)
//	var respData T
//	if err := client.Run(context.Background(), req, &respData); err != nil {
//		return nil, err
//	}
//	return &respData, nil
//}
