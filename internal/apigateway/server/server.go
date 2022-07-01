package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/omekov/superapp-backend/internal/auth/delivery/grpc/v1/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// подключиться к auth server
// swagger gateway
// run server
func Run(port, configPath string) error {
	conn, err := grpc.Dial("localhost:4040", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	client := proto.NewAuthClient(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost == r.Method {
			type auth struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}

			cred := auth{}
			err := json.NewDecoder(r.Body).Decode(&cred)
			if err != nil {
				log.Println(err)
			}

			token, err := client.Login(r.Context(), &proto.AuthRequest{
				Username: cred.Username,
				Password: cred.Password,
			})
			if err != nil {
				log.Println(err)
			}

			byteToken, err := json.Marshal(token)
			if err != nil {
				log.Println(err)
			}

			w.WriteHeader(http.StatusOK)
			w.Write(byteToken)
		}
	})
	http.ListenAndServe(":5051", mux)
	return nil
}
