package db

import (
	"context"
	"testing"
	"time"

	"github.com/JekaTka/services-api/util"

	"github.com/stretchr/testify/require"
)

func createRandomService(t *testing.T) Service {
	arg := CreateServiceParams{
		Name:        util.RandomString(6),
		Description: util.RandomString(10),
	}

	service, err := testQueries.CreateService(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, service)

	require.Equal(t, arg.Name, service.Name)
	require.Equal(t, arg.Description, service.Description)

	require.NotZero(t, service.ID)
	require.NotZero(t, service.CreatedAt)
	require.NotZero(t, service.UpdatedAt)

	return service
}

func TestCreateService(t *testing.T) {
	createRandomService(t)

	t.Cleanup(func() {
		err := testQueries.DeleteAllServices(context.Background())
		require.NoError(t, err)
	})
}

func TestGetServiceByID(t *testing.T) {
	expectedService := createRandomService(t)
	service, err := testQueries.GetServiceByID(context.Background(), expectedService.ID)

	require.NoError(t, err)
	require.NotEmpty(t, service)
	require.Equal(t, expectedService.ID, service.ID)
	require.Equal(t, expectedService.Name, service.Name)
	require.Equal(t, expectedService.Description, service.Description)
	require.WithinDuration(t, expectedService.CreatedAt.Time, service.CreatedAt.Time, time.Second)

	t.Cleanup(func() {
		err := testQueries.DeleteAllServices(context.Background())
		require.NoError(t, err)
	})
}

func TestListServices(t *testing.T) {
	servicesParams := []CreateServiceParams{
		{
			Name:        "Availability service",
			Description: "This is highly available service",
		},
		{
			Name:        "Reliability",
			Description: "This is highly reliable service",
		},
		{
			Name:        "Services API",
			Description: "All services",
		},
		{
			Name:        "ALB on premise",
			Description: "",
		},
		{
			Name:        "AI service",
			Description: "Kind of ChatGPT",
		},
	}

	expectedList := make([]Service, 0, 5)
	for _, param := range servicesParams {
		service, err := testQueries.CreateService(context.Background(), param)
		require.NoError(t, err)
		require.NotEmpty(t, service)

		expectedList = append(expectedList, service)
	}

	// get 5 elem ordered by name
	arg := ListServicesParams{
		Limit:   5,
		Offset:  0,
		OrderBy: "name",
	}
	services, err := testQueries.ListServices(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, services, 5)

	require.Equal(t, "AI service", services[0].Name)
	require.Equal(t, "ALB on premise", services[1].Name)
	require.Equal(t, "Availability service", services[2].Name)
	require.Equal(t, "Reliability", services[3].Name)
	require.Equal(t, "Services API", services[4].Name)

	// get 2 elem ordered by description
	arg = ListServicesParams{
		Limit:   2,
		Offset:  0,
		OrderBy: "description",
	}
	services, err = testQueries.ListServices(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, services, 2)

	require.Equal(t, "ALB on premise", services[0].Name)
	require.Equal(t, "Services API", services[1].Name)
	arg = ListServicesParams{
		Limit:   2,
		Offset:  1,
		OrderBy: "name",
		Search:  "ser",
	}
	services, err = testQueries.ListServices(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, services, 1)
	require.Equal(t, "Availability service", services[0].Name)

	t.Cleanup(func() {
		err := testQueries.DeleteAllServices(context.Background())
		require.NoError(t, err)
	})
}

func TestGetServicesCount(t *testing.T) {
	createRandomService(t)
	createRandomService(t)
	createRandomService(t)

	count, err := testQueries.GetServicesCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(3), count)

	t.Cleanup(func() {
		err := testQueries.DeleteAllServices(context.Background())
		require.NoError(t, err)
	})
}
