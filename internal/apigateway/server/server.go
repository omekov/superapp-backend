package server

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/omekov/superapp-backend/doc/statik"
	gw "github.com/omekov/superapp-backend/internal/apigateway/delivery/grpc/v1"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// подключиться к auth server
// grpc gateway
// run http server
func Run(port, configPath string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterAuthHandlerFromEndpoint(ctx, grpcMux, "localhost:4040", opts)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		return err
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	return http.ListenAndServe(":8081", mux)

	// conn, err := grpc.Dial(, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	return err
	// }

	// client := proto.NewAuthClient(conn)

	// mux := http.NewServeMux()
	// mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	// 	if http.MethodPost == r.Method {
	// 		type auth struct {
	// 			Username string `json:"username"`
	// 			Password string `json:"password"`
	// 		}

	// 		cred := auth{}
	// 		err := json.NewDecoder(r.Body).Decode(&cred)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}

	// 		token, err := client.Login(r.Context(), &proto.AuthRequest{
	// 			Username: cred.Username,
	// 			Password: cred.Password,
	// 		})
	// 		if err != nil {
	// 			log.Println(err)
	// 		}

	// 		byteToken, err := json.Marshal(token)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}

	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write(byteToken)
	// 	}
	// })
	// http.ListenAndServe(":4041", mux)
	// return nil
}
