package main

import (
	"context"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// RegistrationDataのメンバDateの日付フォーマット
const dateFormat = "2006-01-02"

type Group struct {
	GrainDishes       int // 主菜
	VegetableDishes   int // 副菜
	FishAndMealDishes int // 主菜
	Milk              int // 乳製品
	Fruit             int // 果物
}

// RegistrationData 登録する内容
type RegistrationData struct {
	UserID       string
	Date         string
	Name         string
	TimeZone     int
	BalanceGroup Group
}

// Registration ユーザーへのアンケート結果を登録する
type Registration struct {
	projectID  string
	entityType string // datastore entity type
}

// NewRegistationData 登録するデータ
func NewRegistationData(userid string, date string, name string, timezone int, group Group) *RegistrationData {
	return &RegistrationData{
		UserID:       userid,
		Date:         date,
		Name:         name,
		TimeZone:     timezone,
		BalanceGroup: group,
	}
}

// NewRegistration 登録を実行するクラス
func NewRegistration(projectID string, entityType string) *Registration {
	return &Registration{
		projectID:  projectID,
		entityType: entityType,
	}
}

// Put エンティティに登録する
func (r *Registration) Put(ctx context.Context, key *datastore.Key, entity *RegistrationData) error {
	client, err := datastore.NewClient(ctx, r.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Put(ctx, key, entity)
	return err
}

// Get クエリが一致するエンティティを取り出す
func (r *Registration) Get(ctx context.Context, query *datastore.Query, entity *RegistrationData) (*datastore.Key, error) {
	client, err := datastore.NewClient(ctx, r.projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	it := client.Run(ctx, query)
	key, err := it.Next(entity)
	if err == iterator.Done { // datastoreにデータがない場合はキーを発行する
		key = datastore.NameKey("RegistrationData", "", nil)
		err = nil
	}
	return key, err
}

// GetAll 全てのエンティティを取得する
func (r *Registration) GetAll(ctx context.Context, query *datastore.Query, entity *[]RegistrationData) error {
	client, err := datastore.NewClient(ctx, r.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	if _, err := client.GetAll(ctx, query, entity); err != nil {
		return err
	}
	return nil
}
