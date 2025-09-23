package db

import (
	"golang.org/x/sync/errgroup"
)

func Add[T any, T1 any](data *T, data1 *T1) error {
	g := new(errgroup.Group)

	g.Go(func() error {
		return AddToMysql(data) // 返回 error
	})

	g.Go(func() error {
		return AddToNeo4j(&data) // 返回 error
	})

	if err := g.Wait(); err != nil {
		// 任意一个失败，你可以在这里做补偿或回滚
		return err
	}
	return nil
}
