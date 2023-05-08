package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	mockdb "github.com/JekaTka/services-api/db/mock"
	db "github.com/JekaTka/services-api/db/sqlc"
	"github.com/JekaTka/services-api/util"
)

func randomService(name, description string) db.Service {
	return db.Service{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
	}
}

func randomServiceVersion(serviceID uuid.UUID) db.ServiceVersion {
	return db.ServiceVersion{
		ID:        uuid.New(),
		Changelog: util.RandomString(10),
		Version:   fmt.Sprint(util.RandomInt(0, 100)),
		ServiceID: serviceID,
	}
}

func TestGetServicesAPI(t *testing.T) {
	services := []db.Service{
		randomService("api gateway", "awesome api"),
		randomService("s3", "super storage"),
	}
	testCases := []struct {
		name          string
		params        map[string]string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			params: map[string]string{
				"sort_by": "name",
				"limit":   "5",
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListServicesParams{
					Limit:   5,
					OrderBy: "name",
				}

				store.EXPECT().
					ListServices(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(services, nil)

				store.EXPECT().
					GetServiceVersionsCount(gomock.Any(), gomock.Any()).
					Times(2).
					Return(int64(1), nil)

				store.EXPECT().
					GetServicesCount(gomock.Any()).
					Times(1).
					Return(int64(2), nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotServicesResponse getServicesResponse
				err = json.Unmarshal(data, &gotServicesResponse)
				require.NoError(t, err)

				require.Equal(t, http.StatusOK, recorder.Code)
				require.Equal(t, int64(2), gotServicesResponse.Metadata.TotalItems)
				require.Len(t, gotServicesResponse.Services, 2)
				require.Equal(t, services[0].ID, gotServicesResponse.Services[0].ID)
				require.Equal(t, services[0].Name, gotServicesResponse.Services[0].Name)
				require.Equal(t, services[0].Description, gotServicesResponse.Services[0].Description)
				require.Equal(t, int64(1), gotServicesResponse.Services[0].Versions)
			},
		},
		// TODO add checks for 5xx and 4xx errors
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/v1/services", nil)
			require.NoError(t, err)

			q := request.URL.Query()
			for k, v := range tc.params {
				q.Add(k, v)
			}
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetServiceVersions(t *testing.T) {
	service := randomService("ai", "ai tool")
	serviceVersion := randomServiceVersion(service.ID)

	testCases := []struct {
		name          string
		id            uuid.UUID
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			id:   service.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetVersionsByServiceID(gomock.Any(), gomock.Eq(service.ID)).
					Times(1).
					Return([]db.ServiceVersion{serviceVersion}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotServiceVersionsResponse []db.ServiceVersion
				err = json.Unmarshal(data, &gotServiceVersionsResponse)
				require.NoError(t, err)

				require.Equal(t, http.StatusOK, recorder.Code)
				require.Len(t, gotServiceVersionsResponse, 1)
				require.Equal(t, serviceVersion.ID, gotServiceVersionsResponse[0].ID)
				require.Equal(t, serviceVersion.Changelog, gotServiceVersionsResponse[0].Changelog)
				require.Equal(t, serviceVersion.Version, gotServiceVersionsResponse[0].Version)
			},
		},
		// TODO add checks for 5xx and 4xx errors
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/v1/services/%s/versions", tc.id.String())
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
