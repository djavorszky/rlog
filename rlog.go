package rlog

import (
	"log"

	"golang.org/x/net/context"
)

// Server implements implements the methods of the rlog service
type Server struct{}

// Register informs the server that a new service has joined
func (s *Server) Register(ctx context.Context, msg *RegisterRequest) (*RegisterResponse, error) {

	log.Printf("Register: %v", msg)

	return &RegisterResponse{}, nil
}

// Debug will log the received message on debug level
func (s *Server) Debug(ctx context.Context, msg *LogMessage) (*LogResponse, error) {

	log.Printf("Debug: %v", msg)

	return &LogResponse{}, nil
}

// Fatal will log the received message on fatal level
func (s *Server) Fatal(ctx context.Context, msg *LogMessage) (*LogResponse, error) {

	log.Printf("Fatal: %v", msg)

	return &LogResponse{}, nil
}

// Error will log the received message on error level
func (s *Server) Error(ctx context.Context, msg *LogMessage) (*LogResponse, error) {

	log.Printf("Error: %v", msg)

	return &LogResponse{}, nil
}

// Warn will log the received message on warn level
func (s *Server) Warn(ctx context.Context, msg *LogMessage) (*LogResponse, error) {

	log.Printf("Warn: %v", msg)

	return &LogResponse{}, nil
}

// Info will log the received message on info level
func (s *Server) Info(ctx context.Context, msg *LogMessage) (*LogResponse, error) {

	log.Printf("Info: %v", msg)

	return &LogResponse{}, nil
}
