package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/leslesnoa/grpc-mysql/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":9090"
)

type server struct{}

// User is struct
type User struct {
	ID    int32
	Name  string
	Email string
}

func (s *server) GetUser(ctx context.Context, r *pb.GetUserRequest) (*pb.User, error) {
	// var u := User
	log.Printf("Recieved GetUserRequest : %s", r)
	// sql connect
	// 完全なDNS：username:password@protocol(address)/dbname?param=value
	log.Println("starting mysql connect")
	db, err := sql.Open("mysql", "docker:docker@tcp(127.0.0.1:3306)/my_testdb")
	if err != nil {
		log.Panicln("mysql connecting Error!", err)
		panic(err)
	}
	// debug
	if db != nil {
		log.Println("Success SQL connection.")
	}
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead og panic in your app
	}
	log.Println("ping ok.")
	// See "Important settings" section.
	/* 接続のタイムアウト時間を設定 */
	db.SetConnMaxLifetime(time.Minute * 3)
	// アプリが使用する接続の数を制限
	db.SetMaxOpenConns(10)
	// SetMaxOpenConnsと同じ数値を設定推奨
	db.SetMaxIdleConns(10)

	// 1件のレコードのみを取得する場合はQueryRowを利用
	var user User
	user = User{ID: r.GetId()}
	err = db.QueryRow("SELECT * FROM users WHERE id = ?", &user.ID).Scan(&user.Name, &user.Email)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("レコードが存在しません")
	case err != nil:
		panic(err.Error())
	default:
		fmt.Println(user.ID, user.Name)
	}

	defer db.Close()

	return &pb.User{Id: user.ID, Name: user.Name, Email: user.Email}, nil
}
func (s *server) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.User, error) {
	// var u := User
	log.Printf("Recieved CreateUserRequest : %s", r)

	// sql connect
	// 完全なDNS：username:password@protocol(address)/dbname?param=value
	log.Println("starting mysql connect")
	db, err := sql.Open("mysql", "docker:docker@tcp(127.0.0.1:3306)/my_testdb")
	if err != nil {
		log.Panicln("mysql connecting Error!", err)
		panic(err)
	}
	// debug
	if db != nil {
		log.Println("Success SQL connection.")
	}
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead og panic in your app
	}
	log.Println("ping ok.")
	defer db.Close()

	// See "Important settings" section.
	/* 接続のタイムアウト時間を設定 */
	db.SetConnMaxLifetime(time.Minute * 3)
	// アプリが使用する接続の数を制限
	db.SetMaxOpenConns(10)
	// SetMaxOpenConnsと同じ数値を設定推奨
	db.SetMaxIdleConns(10)

	// Insertの場合
	stmtInsert, err := db.Prepare("INSERT INTO users(name) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtInsert.Close()

	user := User{Name: r.GetName()}
	result, err := stmtInsert.Exec(user.Name)
	if err != nil {
		panic(err.Error())
	}
	log.Println(result)

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(lastInsertID)

	return &pb.User{Id: int32(lastInsertID), Name: r.Name, Email: r.Email}, nil
}
func (s *server) UpdateUser(ctx context.Context, r *pb.UpdateUserRequest) (*pb.User, error) {
	// var u := User
	log.Printf("Recieved UpdateUserRequest : %s", r)

	// Update
	// sql connect
	// 完全なDNS：username:password@protocol(address)/dbname?param=value
	log.Println("starting mysql connect")
	db, err := sql.Open("mysql", "docker:docker@tcp(127.0.0.1:3306)/my_testdb")
	if err != nil {
		log.Panicln("mysql connecting Error!", err)
		panic(err)
	}
	// debug
	if db != nil {
		log.Println("Success SQL connection.")
	}
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead og panic in your app
	}
	log.Println("ping ok.")
	defer db.Close()

	// See "Important settings" section.
	/* 接続のタイムアウト時間を設定 */
	db.SetConnMaxLifetime(time.Minute * 3)
	// アプリが使用する接続の数を制限
	db.SetMaxOpenConns(10)
	// SetMaxOpenConnsと同じ数値を設定推奨
	db.SetMaxIdleConns(10)

	user := User{ID: r.GetId(), Name: r.GetName(), Email: r.GetEmail()}
	stmtUpdate, err := db.Prepare("UPDATE users SET name=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtUpdate.Close()
	log.Println(user)

	result, err := stmtUpdate.Exec(user.Name, user.ID)
	if err != nil {
		panic(err.Error())
	}

	rowsAffect, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rowsAffect)

	return &pb.User{Id: user.ID, Name: user.Name, Email: user.Email}, nil
}
func (s *server) DeleteUser(ctx context.Context, r *pb.DeleteUserRequest) (*pb.Empty, error) {
	// var u := User
	log.Printf("Recieved DeleteUserRequest : %s", r)

	// Delete
	// sql connect
	// 完全なDNS：username:password@protocol(address)/dbname?param=value
	log.Println("starting mysql connect")
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/my_testdb")
	if err != nil {
		log.Panicln("mysql connecting Error!", err)
		panic(err)
	}
	// debug
	if db != nil {
		log.Println("Success SQL connection.")
	}
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead og panic in your app
	}
	log.Println("ping ok.")
	defer db.Close()

	// See "Important settings" section.
	/* 接続のタイムアウト時間を設定 */
	db.SetConnMaxLifetime(time.Minute * 3)
	// アプリが使用する接続の数を制限
	db.SetMaxOpenConns(10)
	// SetMaxOpenConnsと同じ数値を設定推奨
	db.SetMaxIdleConns(10)

	user := User{ID: r.GetId()}
	stmtDelete, err := db.Prepare("DELETE FROM users WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtDelete.Close()

	result, err := stmtDelete.Exec(user.ID)
	if err != nil {
		panic(err.Error())
	}

	rowsAffect, err := result.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rowsAffect)

	return &pb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %c", err)
	}
}
