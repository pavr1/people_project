package auth

import (
	"io"
	"net/http"

	"github.com/pavr1/people_project/people/config"
	log "github.com/sirupsen/logrus"
)

type Auth struct {
	log    *log.Logger
	config *config.Config
}

func NewAuth(log *log.Logger, config *config.Config) *Auth {
	return &Auth{
		log:    log,
		config: config,
	}
}

func (a *Auth) IsValidToken(token string) (int, string, error) {
	req, err := http.NewRequest(http.MethodGet, a.config.Auth.Path, nil)
	if err != nil {
		a.log.WithError(err).Error("Failed to create request")

		return -1, "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Host", a.config.Auth.Host)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.log.WithField("Path", a.config.Auth.Path).WithError(err).Error("Failed to send request")

		return -1, "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.log.WithError(err).Error("Failed to read response")

		return -1, "", err
	}

	return resp.StatusCode, string(body), nil
}
