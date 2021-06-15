package main

import (
	"context"
	"fmt"
	"go-learn/ent/ent"
	"go-learn/ent/ent/car"
	"go-learn/ent/ent/group"
	"go-learn/ent/ent/user"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

func main() {
	client, err := ent.Open("sqlite3", "file:ent/sqlite/test.db?_fk=1", ent.Debug())
	if nil != err {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	if err := client.Schema.Create(context.Background()); nil != err {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	/*
		// ariel, err := CreateUser(context.Background(), client, 30, "Ariel")
		// if nil != err {
		// 	log.Fatalf("failed creating user: %v", err)
		// }

		// neta, err := CreateUser(context.Background(), client, 28, "Neta")
		// if nil != err {
		// 	log.Fatalf("failed creating user: %v", err)
		// }

		ariel, err := QueryUser(context.Background(), client, "Ariel")
		if nil != err {
			log.Fatalf("failed query user: %v", err)
		}

		neta, err := QueryUser(context.Background(), client, "Neta")
		if nil != err {
			log.Fatalf("failed query user: %v", err)
		}

		ford, err := QueryCar(context.Background(), client, "Ford")
		if nil != err {
			log.Fatalf("failed query car: %v", err)
		}

		tesla, err := QueryCar(context.Background(), client, "Tesla")
		if nil != err {
			log.Fatalf("failed query car: %v", err)
		}

		mazda, err := QueryCar(context.Background(), client, "Mazda")
		if nil != err {
			log.Fatalf("failed query car: %v", err)
		}

		// mazda, err := CreateCar(context.Background(), client, "Mazda")
		// if nil != err {
		// 	log.Fatalf("failed creating car: %v", err)
		// }

		err = UpdateCarOwner(context.Background(), ariel, tesla, mazda)
		if nil != err {
			log.Fatalf("failed update user: %v", err)
		}

		err = UpdateCarOwner(context.Background(), neta, ford)
		if nil != err {
			log.Fatalf("failed update user: %v", err)
		}

		_, err = CreateGroup(context.Background(), client, "GitLab", ariel, neta)
		if nil != err {
			log.Fatalf("failed creating group: %v", err)
		}

		_, err = CreateGroup(context.Background(), client, "GitHub", ariel)
		if nil != err {
			log.Fatalf("failed creating group: %v", err)
		}
	*/

	// u, err := QueryUser(context.Background(), client)
	// if nil != err {
	// 	log.Fatalf("failed query user: %v", err)
	// }

	// _, err = CreateCars(context.Background(), client, u)
	// if nil != err {
	// 	log.Fatalf("failed create user: %v", err)
	// }

	// err = QueryCars(context.Background(), u)
	// if nil != err {
	// 	log.Fatalf("failed query car/s: %v", err)
	// }

	// err = UpdateGroup(context.Background(), client)
	// if nil != err {
	// 	log.Fatalf("failed update group: %v", err)
	// }

	err = QueryGithub(context.Background(), client)
	if nil != err {
		log.Fatalf("failed query car/s: %v", err)
	}

	err = QueryArielCars(context.Background(), client)
	if nil != err {
		log.Fatalf("failed query car/s: %v", err)
	}

	err = QueryGroupWithUsers(context.Background(), client)
	if nil != err {
		log.Fatalf("failed query group/s: %v", err)
	}
}

func QueryGroupWithUsers(ctx context.Context, client *ent.Client) error {
	groups, err := client.Group.
		Query().
		Where(group.HasUsers()).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting groups: %w", err)
	}
	log.Println("groups returned:", groups) // Output: (Group(Name=GitHub), Group(Name=GitLab),)
	return nil
}

func QueryArielCars(ctx context.Context, client *ent.Client) error {
	// Get "Ariel" from previous steps.
	a8m := client.User.
		Query().
		Where(
			user.HasCars(),
			user.Name("Ariel"),
		).
		OnlyX(ctx)
		// Get the groups, that a8m is connected to:
	cars, err := a8m.
		QueryGroups(). // (Group(Name=GitHub), Group(Name=GitLab),)
		QueryUsers().  // (User(Name=Ariel, Age=30), User(Name=Neta, Age=28),)
		QueryCars().   //
		Where(car.Not( //  Get Neta and Ariel cars, but filter out
			car.ModelEQ("Mazda"), //  those who named "Mazda"
		)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
	}
	log.Println("cars returned:", cars) // Output: (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Ford, RegisteredAt=<Time>),)
	return nil
}

func QueryGithub(ctx context.Context, client *ent.Client) error {
	cars, err := client.Group.
		Query().
		Where(group.NameEQ("GitHub")). // (Group(Name=GitHub),)
		QueryUsers().                  // (User(Name=Ariel, Age=30),)
		QueryCars().                   // (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Mazda, RegisteredAt=<Time>),)
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %w", err)
	}
	log.Println("cars returned:", cars) // Output: (Car(Model=Tesla, RegisteredAt=<Time>), Car(Model=Mazda, RegisteredAt=<Time>),)
	return nil
}

func UpdateGroup(ctx context.Context, client *ent.Client) error {
	err := client.Group.Update().SetName("GitHub").Where(group.IDEQ(2)).Exec(ctx)
	if nil != err {
		return errors.Wrap(err, "failed updating group")
	}
	log.Println("group was updated")

	return nil
}

func CreateGroup(ctx context.Context, client *ent.Client, name string, us ...*ent.User) (*ent.Group, error) {
	g, err := client.Group.Create().SetName(name).AddUsers(us...).Save(ctx)
	if nil != err {
		return nil, errors.Wrap(err, "failed creating group")
	}
	log.Println("group was created: ", g)

	return g, nil
}

func CreateUser(ctx context.Context, client *ent.Client, age int, name string) (*ent.User, error) {
	u, err := client.User.Create().SetAge(age).SetName(name).Save(ctx)
	if nil != err {
		return nil, errors.Wrap(err, "failed creating user")
	}
	log.Println("user was created: ", u)

	return u, nil
}

func QueryUser(ctx context.Context, client *ent.Client, name string) (*ent.User, error) {
	u, err := client.User.Query().Where(user.NameEQ(name)).First(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed querying user")
	}
	log.Println("user returned: ", u)

	return u, nil
}

func CreateCars(ctx context.Context, client *ent.Client, userNoCar *ent.User) (*ent.User, error) {
	tesla, err := client.Car.Create().SetModel("Tesla").SetRegisteredAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating car")
	}
	log.Println("car was created: ", tesla)

	ford, err := client.Car.Create().SetModel("Ford").SetRegisteredAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating car")
	}
	log.Println("car was created: ", ford)

	u, err := userNoCar.Update().AddCars(tesla, ford).Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed update user")
	}
	log.Println("user was updated: ", u)

	return u, nil
}

func QueryCars(ctx context.Context, userHasCars *ent.User) error {
	cars, err := userHasCars.QueryCars().All(ctx)
	if err != nil {
		return errors.Wrap(err, "failed querying user cars")
	}
	log.Println("returned cars:", cars)

	for _, car := range cars {
		owner, err := car.QueryOwner().Only(ctx)
		if err != nil {
			return errors.Wrapf(err, "failed querying car %q owner", car.Model)
		}
		log.Printf("car %q owner: %q\n", car.Model, owner.Name)
	}

	ford, err := userHasCars.QueryCars().Where(car.ModelEQ("Ford")).Only(ctx)
	if err != nil {
		return errors.Wrap(err, "failed querying user car")
	}
	log.Println("returned ford car:", ford)

	return nil
}

func CreateCar(ctx context.Context, client *ent.Client, mode string) (*ent.Car, error) {
	car, err := client.Car.Create().SetModel(mode).SetRegisteredAt(time.Now()).Save(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating car")
	}
	log.Println("car was created: ", car)

	return car, nil
}

func QueryCar(ctx context.Context, client *ent.Client, mode string) (*ent.Car, error) {
	car, err := client.Car.Query().Where(car.ModelEQ(mode)).Only(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed querying user car")
	}
	log.Println("returned car:", car)

	return car, nil
}

func UpdateCarOwner(ctx context.Context, u *ent.User, cs ...*ent.Car) error {
	for _, c := range cs {
		_ = c.Update().ClearOwner().Exec(ctx)
	}
	_, err := u.Update().AddCars(cs...).Save(ctx)
	if err != nil {
		return errors.Wrap(err, "failed update user")
	}
	log.Println("user was updated: ", u)

	return nil
}
