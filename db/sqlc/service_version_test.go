package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/JekaTka/services-api/util"
)

func createRandomServiceVersion(t *testing.T, serviceID uuid.UUID) ServiceVersion {
	arg := CreateServiceVersionParams{
		Changelog: util.RandomString(10),
		Version:   fmt.Sprintf("%d.%d.%d", util.RandomInt(0, 10), util.RandomInt(0, 10), util.RandomInt(0, 10)),
		ServiceID: serviceID,
	}

	serviceVersion, err := testQueries.CreateServiceVersion(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, serviceVersion)

	require.Equal(t, arg.Changelog, serviceVersion.Changelog)
	require.Equal(t, arg.Version, serviceVersion.Version)
	require.Equal(t, arg.ServiceID, serviceVersion.ServiceID)

	require.NotZero(t, serviceVersion.ID)
	require.NotZero(t, serviceVersion.CreatedAt)
	require.NotZero(t, serviceVersion.UpdatedAt)

	return serviceVersion
}

func TestCreateServiceVersion(t *testing.T) {
	service := createRandomService(t)
	createRandomServiceVersion(t, service.ID)

	t.Cleanup(func() {
		err := testQueries.DeleteAllServiceVersions(context.Background())
		require.NoError(t, err)
		err = testQueries.DeleteAllServices(context.Background())
		require.NoError(t, err)
	})
}

func TestGetVersionsByServiceID(t *testing.T) {
	service := createRandomService(t)
	expectedServiceVersion := createRandomServiceVersion(t, service.ID)

	serviceVersion, err := testQueries.GetVersionsByServiceID(context.Background(), service.ID)
	require.NoError(t, err)
	require.NotEmpty(t, service)
	require.Len(t, serviceVersion, 1)
	require.Equal(t, []ServiceVersion{expectedServiceVersion}, serviceVersion)

	t.Cleanup(func() {
		err := testQueries.DeleteAllServiceVersions(context.Background())
		require.NoError(t, err)
		err = testQueries.DeleteAllServices(context.Background())
		require.NoError(t, err)
	})
}

func TestGetServiceVersionsCount(t *testing.T) {
	service := createRandomService(t)
	createRandomServiceVersion(t, service.ID)
	createRandomServiceVersion(t, service.ID)
	createRandomServiceVersion(t, service.ID)

	count, err := testQueries.GetServiceVersionsCount(context.Background(), service.ID)
	require.NoError(t, err)
	require.Equal(t, int64(3), count)

	t.Cleanup(func() {
		err := testQueries.DeleteAllServiceVersions(context.Background())
		require.NoError(t, err)
		err = testQueries.DeleteAllServices(context.Background())
		require.NoError(t, err)
	})
}
