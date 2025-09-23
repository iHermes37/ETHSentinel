//package main
//
//import (
//	"context"
//	"fmt"
//	"github.com/machinebox/graphql"
//	"log"
//)
//
//func main() {
//	client := graphql.NewClient("https://gateway.thegraph.com/api/subgraphs/id/5zvR82QoaXYFyDEKLZ9t6v9adgnptxYpKpSbxtgVENFV")
//
//	req := graphql.NewRequest(`
//		query {
//			factories(first: 5) {
//				id
//				poolCount
//				txCount
//				totalVolumeUSD
//			}
//			bundles(first: 5) {
//				id
//				ethPriceUSD
//			}
//		}
//	`)
//
//	req.Header.Set("Authorization", "Bearer df5d393ba8219b65e3eea66df2242e6b")
//
//	ctx := context.Background()
//
//	var resp struct {
//		Factories []struct {
//			ID             string
//			PoolCount      string
//			TxCount        string
//			TotalVolumeUSD string
//		}
//		Bundles []struct {
//			ID          string
//			EthPriceUSD string
//		}
//	}
//
//	if err := client.Run(ctx, req, &resp); err != nil {
//		log.Fatal(err)
//	}
//
//	for _, factory := range resp.Factories {
//		fmt.Printf("Factory ID: %s, Pools: %s, TxCount: %s, VolumeUSD: %s\n",
//			factory.ID, factory.PoolCount, factory.TxCount, factory.TotalVolumeUSD)
//	}
//
//	for _, bundle := range resp.Bundles {
//		fmt.Printf("Bundle ID: %s, ETH Price USD: %s\n", bundle.ID, bundle.EthPriceUSD)
//	}
//}
