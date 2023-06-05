package main

import (
	"context"
	"log"
	"net"

	connectDB "github.com/UserFinder/connect"
	"github.com/UserFinder/initializers"
	"github.com/UserFinder/models"
	pb "github.com/UserFinder/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) CreateUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	// Hash the password using bcrypt
	uuid := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password+uuid), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error hashing password")
	}

	post := models.User{
		ID:       uuid,
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword), // Store the hashed password in the database
	}

	var check struct {
		Result bool
	}

	err = connectDB.DB.Raw("SELECT EXISTS(SELECT 1 FROM `users` WHERE `email` = ?) as result", req.Email).Scan(&check).Error

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error checking email")
	}

	if check.Result {
		return nil, status.Errorf(codes.InvalidArgument, "User already exists")
	}

	// use the `result` variable
	result := connectDB.DB.Create(&post)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "Error creating user")
	}

	return &pb.UserResponse{
		Id:    post.ID,
		Name:  req.Name,
		Email: req.Email,
	}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	// Get the user ID from the request
	email := req.Email
	password := req.Password

	user, err := authenticateUser(email, password)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error authenticating user")
	}

	// Create and return a user response
	response := &pb.UserResponse{
		Id:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
	return response, nil
}

func authenticateUser(email string, password string) (models.User, error) {
	var user models.User
	result := connectDB.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	// Compare the hashed password with the input password

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+user.ID)); err != nil {

		return user, err
	}

	return user, nil
}

func init() {
	initializers.LoadEnv()
	connectDB.InitConnector()
}
func main() {
	// Start a gRPC server
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	server := &server{}
	pb.RegisterUserServiceServer(s,server)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Connect to the database

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
}
