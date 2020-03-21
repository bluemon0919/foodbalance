package main

import (
	"context"

	"cloud.google.com/go/datastore"
)

// MealのメンバDateの日付フォーマット
const dateFormat = "2006-01-02"

type Meal struct {
	Date         string
	Name         string
	TimeZone     int
	BalanceGroup Group
}
type Group struct {
	GrainDishes       int // 主菜
	VegetableDishes   int // 副菜
	FishAndMealDishes int // 主菜
	Milk              int // 乳製品
	Fruit             int // 果物
}

// Sum adds the entered Group to your own Group
func (g *Group) Sum(in Group) {
	g.GrainDishes += in.GrainDishes
	g.VegetableDishes += in.VegetableDishes
	g.FishAndMealDishes += in.FishAndMealDishes
	g.Milk += in.Milk
	g.Fruit += in.Fruit
}

// SumGroup returns sum of groups
func SumGroup(ms []Meal) Group {
	var group Group
	for _, m := range ms {
		group.Sum(m.BalanceGroup)
	}
	return group
}

func Create(date string, name string, timeZone int, grain int, vegetable int, fish int, milk int, fruit int) Meal {
	return Meal{
		Date:     date,
		Name:     name,
		TimeZone: timeZone,
		BalanceGroup: Group{
			GrainDishes:       grain,
			VegetableDishes:   vegetable,
			FishAndMealDishes: fish,
			Milk:              milk,
			Fruit:             fruit,
		},
	}
}

func Put(m Meal) error {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, myProjectID)
	if err != nil {
		return err
	}
	defer client.Close()

	key := datastore.NameKey("Meal", "", nil)
	_, err = client.Put(ctx, key, &m)
	return err
}

func GetAll() ([]Meal, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, myProjectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	var ms []Meal
	q := datastore.NewQuery("Meal")
	if _, err := client.GetAll(ctx, q, &ms); err != nil {
		return nil, err
	}
	return ms, nil
}
