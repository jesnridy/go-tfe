package tfe

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHKeysList(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	orgTest, orgTestCleanup := createOrganization(t, client)
	defer orgTestCleanup()

	kTest1, _ := createSSHKey(t, client, orgTest)
	kTest2, _ := createSSHKey(t, client, orgTest)

	t.Run("without list options", func(t *testing.T) {
		ks, err := client.SSHKeys.List(ctx, orgTest.Name, SSHKeyListOptions{})
		require.NoError(t, err)
		assert.Contains(t, ks, kTest1)
		assert.Contains(t, ks, kTest2)
	})

	t.Run("with list options", func(t *testing.T) {
		t.Skip("paging not supported yet in API")
		// Request a page number which is out of range. The result should
		// be successful, but return no results if the paging options are
		// properly passed along.
		ks, err := client.SSHKeys.List(ctx, orgTest.Name, SSHKeyListOptions{
			ListOptions: ListOptions{
				PageNumber: 999,
				PageSize:   100,
			},
		})
		require.NoError(t, err)
		assert.Empty(t, ks)
	})

	t.Run("without a valid organization", func(t *testing.T) {
		ks, err := client.SSHKeys.List(ctx, badIdentifier, SSHKeyListOptions{})
		assert.Nil(t, ks)
		assert.EqualError(t, err, "Invalid value for organization")
	})
}

func TestSSHKeysCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	orgTest, orgTestCleanup := createOrganization(t, client)
	defer orgTestCleanup()

	t.Run("with valid options", func(t *testing.T) {
		options := SSHKeyCreateOptions{
			Name:  String(randomString(t)),
			Value: String(randomString(t)),
		}

		k, err := client.SSHKeys.Create(ctx, orgTest.Name, options)
		require.NoError(t, err)

		// Get a refreshed view from the API.
		refreshed, err := client.SSHKeys.Read(ctx, k.ID)
		require.NoError(t, err)

		for _, item := range []*SSHKey{
			k,
			refreshed,
		} {
			assert.NotEmpty(t, item.ID)
			assert.Equal(t, *options.Name, item.Name)
		}
	})

	t.Run("when options is missing name", func(t *testing.T) {
		k, err := client.SSHKeys.Create(ctx, "foo", SSHKeyCreateOptions{
			Value: String(randomString(t)),
		})
		assert.Nil(t, k)
		assert.EqualError(t, err, "Name is required")
	})

	t.Run("when options is missing value", func(t *testing.T) {
		k, err := client.SSHKeys.Create(ctx, "foo", SSHKeyCreateOptions{
			Name: String(randomString(t)),
		})
		assert.Nil(t, k)
		assert.EqualError(t, err, "Value is required")
	})

	t.Run("when options has an invalid organization", func(t *testing.T) {
		k, err := client.SSHKeys.Create(ctx, badIdentifier, SSHKeyCreateOptions{
			Name: String("foo"),
		})
		assert.Nil(t, k)
		assert.EqualError(t, err, "Invalid value for organization")
	})
}

func TestSSHKeysRead(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	orgTest, orgTestCleanup := createOrganization(t, client)
	defer orgTestCleanup()

	kTest, _ := createSSHKey(t, client, orgTest)

	t.Run("when the SSH key exists", func(t *testing.T) {
		k, err := client.SSHKeys.Read(ctx, kTest.ID)
		require.NoError(t, err)
		assert.Equal(t, kTest, k)
	})

	t.Run("when the SSH key does not exist", func(t *testing.T) {
		k, err := client.SSHKeys.Read(ctx, "nonexisting")
		assert.Nil(t, k)
		assert.EqualError(t, err, "Error: not found")
	})

	t.Run("without a valid SSH key ID", func(t *testing.T) {
		k, err := client.SSHKeys.Read(ctx, badIdentifier)
		assert.Nil(t, k)
		assert.EqualError(t, err, "Invalid value for SSH key ID")
	})
}

func TestSSHKeysUpdate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	orgTest, orgTestCleanup := createOrganization(t, client)
	defer orgTestCleanup()

	t.Run("with valid options", func(t *testing.T) {
		kBefore, kTestCleanup := createSSHKey(t, client, orgTest)
		defer kTestCleanup()

		kAfter, err := client.SSHKeys.Update(ctx, kBefore.ID, SSHKeyUpdateOptions{
			Name:  String(randomString(t)),
			Value: String(randomString(t)),
		})
		require.NoError(t, err)

		assert.Equal(t, kBefore.ID, kAfter.ID)
		assert.NotEqual(t, kBefore.Name, kAfter.Name)
	})

	t.Run("when updating the name", func(t *testing.T) {
		kBefore, kTestCleanup := createSSHKey(t, client, orgTest)
		defer kTestCleanup()

		kAfter, err := client.SSHKeys.Update(ctx, kBefore.ID, SSHKeyUpdateOptions{
			Name: String("updated-key-name"),
		})
		require.NoError(t, err)

		assert.Equal(t, kBefore.ID, kAfter.ID)
		assert.Equal(t, "updated-key-name", kAfter.Name)
	})

	t.Run("when updating the value", func(t *testing.T) {
		kBefore, kTestCleanup := createSSHKey(t, client, orgTest)
		defer kTestCleanup()

		kAfter, err := client.SSHKeys.Update(ctx, kBefore.ID, SSHKeyUpdateOptions{
			Value: String("updated-key-value"),
		})
		require.NoError(t, err)

		assert.Equal(t, kBefore.ID, kAfter.ID)
		assert.Equal(t, kBefore.Name, kAfter.Name)
	})

	t.Run("without a valid SSH key ID", func(t *testing.T) {
		w, err := client.SSHKeys.Update(ctx, badIdentifier, SSHKeyUpdateOptions{})
		assert.Nil(t, w)
		assert.EqualError(t, err, "Invalid value for SSH key ID")
	})
}

func TestSSHKeysDelete(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	orgTest, orgTestCleanup := createOrganization(t, client)
	defer orgTestCleanup()

	kTest, _ := createSSHKey(t, client, orgTest)

	t.Run("with valid options", func(t *testing.T) {
		err := client.SSHKeys.Delete(ctx, kTest.ID)
		require.NoError(t, err)

		// Try loading the workspace - it should fail.
		_, err = client.SSHKeys.Read(ctx, kTest.ID)
		assert.EqualError(t, err, "Error: not found")
	})

	t.Run("when the SSH key does not exist", func(t *testing.T) {
		err := client.SSHKeys.Delete(ctx, kTest.ID)
		assert.EqualError(t, err, "Error: not found")
	})

	t.Run("when the SSH key ID is invalid", func(t *testing.T) {
		err := client.SSHKeys.Delete(ctx, badIdentifier)
		assert.EqualError(t, err, "Invalid value for SSH key ID")
	})
}
