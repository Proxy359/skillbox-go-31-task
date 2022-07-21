package storage

import (
	"cde/internal/entities"
	"context"
	"errors"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	Store *mongo.Database
}

type Counter struct {
	Value int `bson:"value"`
}

func (s *MongoStorage) Save(u *entities.MongoUser) (int, error) {
	if u == nil {
		return 0, errors.New("недостаточно данных для создания пользователя")
	}

	var counter Counter
	err := s.Store.Collection("counter").FindOneAndUpdate(
		context.TODO(),
		bson.M{},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
		options.FindOneAndUpdate().SetUpsert(true),
	).Decode(&counter)
	if err != nil {
		log.Println(err)
	}

	u.ID = counter.Value

	s.Store.Collection("users").InsertOne(context.TODO(), u)
	return u.ID, nil
}

func (s *MongoStorage) Add(sourceId, targetId int) (string, error) {

	firstRest := s.Store.Collection("users").FindOne(context.TODO(), bson.M{"id": sourceId})
	if firstRest.Err() != nil {
		return "", firstRest.Err()
	}

	firstUser := entities.MongoUser{}
	if err := firstRest.Decode(&firstUser); err != nil {
		return "", err
	}

	secondRest := s.Store.Collection("users").FindOne(context.TODO(), bson.M{"id": targetId})
	if secondRest.Err() != nil {
		return "", secondRest.Err()
	}

	secondUser := entities.MongoUser{}
	if err := secondRest.Decode(&secondUser); err != nil {
		return "", err
	}

	firstUser.Friends = append(firstUser.Friends, secondUser.ID)
	secondUser.Friends = append(secondUser.Friends, firstUser.ID)

	s.Store.Collection("users").UpdateOne(context.TODO(), bson.M{"id": firstUser.ID}, bson.M{"$set": bson.M{"friends": firstUser.Friends}})
	s.Store.Collection("users").UpdateOne(context.TODO(), bson.M{"id": secondUser.ID}, bson.M{"$set": bson.M{"friends": secondUser.Friends}})

	ans := firstUser.Name + " и " + secondUser.Name + " теперь друзья"
	return ans, nil
}

func (s *MongoStorage) DelieteU(targetId int) (string, error) {

	res := s.Store.Collection("users").FindOne(context.TODO(), bson.M{"id": targetId})
	if res.Err() != nil {
		return "", res.Err()
	}

	user := entities.MongoUser{}
	if err := res.Decode(&user); err != nil {
		return "", err
	}

	friends := user.Friends

	if len(friends) != 0 {
		for i := 0; i < len(friends); i++ {
			cycleRes := s.Store.Collection("users").FindOne(context.TODO(), bson.M{"id": friends[i]})
			cycleUser := entities.MongoUser{}
			if err := cycleRes.Decode(&cycleUser); err != nil {
				return "", err
			}
			log.Println(cycleUser)
			indexInFriendList := findID(cycleUser.Friends, targetId)

			friendsListLen := len(cycleUser.Friends)
			cycleUser.Friends[indexInFriendList] = cycleUser.Friends[friendsListLen-1]
			cycleUser.Friends = cycleUser.Friends[:friendsListLen-1]
			log.Println(cycleUser)
			s.Store.Collection("users").UpdateOne(context.TODO(), bson.M{"id": cycleUser.ID}, bson.M{"$set": bson.M{"friends": cycleUser.Friends}})
		}
	}

	log.Println(user.Name)
	ans := "Имя удалённого пользователя " + user.Name

	log.Println(ans)
	s.Store.Collection("users").DeleteOne(context.TODO(), bson.M{"id": user.ID})
	return ans, nil
}

func findID(a []int, aa int) int {
	lenA := len(a)
	counter := 0

	for i := 0; i < lenA; i++ {
		if a[i] == aa {
			break
		} else {
			counter++
		}

	}
	return counter
}

func (s *MongoStorage) GetU(u string) ([]string, error) {
	userID, err := strconv.Atoi(u)
	if err != nil {
		return nil, errors.New("id пользователя введен некорректно")
	}

	res := s.Store.Collection("users").FindOne(context.TODO(), bson.M{"id": userID})
	if res.Err() != nil {
		return nil, res.Err()
	}

	user := entities.MongoUser{}
	if err := res.Decode(&user); err != nil {
		return nil, err
	}

	friendsID := user.Friends
	friendsList := []string{}
	if len(friendsID) != 0 {
		for i := 0; i < len(friendsID); i++ {
			cycleRes := s.Store.Collection("users").FindOne(context.TODO(), bson.M{"id": friendsID[i]})
			cycleUser := entities.MongoUser{}
			if err := cycleRes.Decode(&cycleUser); err != nil {
				return nil, err
			}
			friendsList = append(friendsList, cycleUser.Name)
		}
	}
	return friendsList, nil
}

func (s *MongoStorage) NewAgeU(newAge int, idU string) (string, error) {
	abc, err := strconv.Atoi(idU)
	if err != nil {
		return "", errors.New("id пользователя введен некорректно")
	}

	res := s.Store.Collection("users").FindOne(context.TODO(), bson.M{"id": abc})
	if res.Err() != nil {
		return "", res.Err()
	}

	user := entities.MongoUser{}
	if err := res.Decode(&user); err != nil {
		return "", err
	}

	user.Age = newAge
	s.Store.Collection("users").UpdateOne(context.TODO(), bson.M{"id": abc}, bson.M{"$set": bson.M{"age": user.Age}})
	a := "возраст пользователя успешно обновлён"
	return a, nil
}

func (s *MongoStorage) AllUsers() ([]*entities.MongoUser, error) {

	curs, err := s.Store.Collection("users").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	users := []*entities.MongoUser{}

	for curs.Next(context.TODO()) {

		user := entities.MongoUser{}
		if err := curs.Decode(&user); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}
	return users, nil
}
