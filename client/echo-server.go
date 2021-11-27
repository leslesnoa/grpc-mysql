package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	pb "github.com/leslesnoa/grpc-mysql/pb"
	"google.golang.org/grpc"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.BodyDump(bodyDumpHandler))
	e.GET("/users", getUsers)
	e.Logger.Fatal(e.Start(":1323"))
}

// Echoのリクエスト/レスポンスボディのプラグイン
func bodyDumpHandler(c echo.Context, reqBody, resBody []byte) {
	fmt.Printf("Request Body: %v\n", string(reqBody))
	fmt.Printf("Response Body: %v\n", string(resBody))
}

// Interceptorの定義
func unaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("before call: %s, request: %+v", method, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("after call: %s, response: %+v", method, reply)
	return err
}

func getUsers(c echo.Context) error {
	fmt.Println("starting at getUsersFunc.")
	u := new(User)
	c.Bind(u)
	addr := "localhost:9090"
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(unaryInterceptor))
	// conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	ctx := context.Background()
	// r, err := client.CreateUser(ctx, &pb.CreateUserRequest{Name: u.Name, Email: u.Email})
	// r, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{Id: 2, Name: u.Name, Email: u.Email})
	r, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
	// r, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: 5})
	// r, err := client.SayHello(ctx, &pb.HelloRequest{Name: u.Name, Email: u.Email, Test: "test"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	// log.Printf("Name: %s", r.Name)
	// log.Printf("Id: %v", r.Id)
	// log.Printf("Email: %s", r.Email)
	return c.JSON(http.StatusOK, r)
	// return r.JSON(http.StatusOK, u)
}
