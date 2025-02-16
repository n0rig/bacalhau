package util

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/bacalhau-project/bacalhau/pkg/config"
	"github.com/bacalhau-project/bacalhau/pkg/publicapi/apimodels"
	"github.com/bacalhau-project/bacalhau/pkg/publicapi/client"
	clientv2 "github.com/bacalhau-project/bacalhau/pkg/publicapi/client/v2"
	"github.com/bacalhau-project/bacalhau/pkg/version"
)

func GetAPIClient(ctx context.Context) *client.APIClient {
	legacyTLS := client.LegacyTLSSupport(config.ClientTLSConfig())
	apiClient := client.NewAPIClient(legacyTLS, config.ClientAPIHost(), config.ClientAPIPort())

	if token, err := getAuthToken(config.ClientAPIBase()); err != nil {
		log.Warn().Err(err).Msg("Failed to read access tokens – API calls will be without authorization")
	} else if token != nil {
		apiClient.DefaultHeaders["Authorization"] = token.String()
	}

	return apiClient
}

func GetAPIClientV2() *clientv2.Client {
	base := config.ClientAPIBase()
	tlsConfig := config.ClientTLSConfig()

	bv := version.Get()
	headers := map[string][]string{
		apimodels.HTTPHeaderBacalhauGitVersion: {bv.GitVersion},
		apimodels.HTTPHeaderBacalhauGitCommit:  {bv.GitCommit},
		apimodels.HTTPHeaderBacalhauBuildDate:  {bv.BuildDate.UTC().String()},
		apimodels.HTTPHeaderBacalhauBuildOS:    {bv.GOOS},
		apimodels.HTTPHeaderBacalhauArch:       {bv.GOARCH},
	}

	opts := []clientv2.OptionFn{
		clientv2.WithCACertificate(tlsConfig.CACert),
		clientv2.WithInsecureTLS(tlsConfig.Insecure),
		clientv2.WithTLS(tlsConfig.UseTLS),
		clientv2.WithHeaders(headers),
	}

	if token, err := getAuthToken(base); err != nil {
		log.Warn().Err(err).Msg("Failed to read access tokens – API calls will be without authorization")
	} else if token != nil {
		opts = append(opts, clientv2.WithHTTPAuth(token))
	}

	return clientv2.New(base, opts...)
}

func getAuthToken(baseURL string) (*apimodels.HTTPCredential, error) {
	token, err := ReadToken(baseURL)
	if token != "" {
		return &apimodels.HTTPCredential{
			Scheme: "Bearer",
			Value:  token,
		}, err
	}

	return nil, err
}
