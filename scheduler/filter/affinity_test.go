package filter

import (
	"testing"

	"github.com/docker/swarm/cluster"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
)

func TestAffinityFilter(t *testing.T) {
	//TODO: add test for images

	var (
		f     = AffinityFilter{}
		nodes = []*cluster.Node{
			cluster.NewNode("node-0"),
			cluster.NewNode("node-1"),
			cluster.NewNode("node-2"),
		}
		result []*cluster.Node
		err    error
	)

	nodes[0].ID = "node-0-id"
	nodes[0].Name = "node-0-name"
	nodes[0].AddContainer(&cluster.Container{
		Container: dockerclient.Container{
			Id:    "container-0-id",
			Names: []string{"container-0-name"},
		},
	})

	nodes[1].ID = "node-1-id"
	nodes[1].Name = "node-1-name"
	nodes[1].AddContainer(&cluster.Container{
		Container: dockerclient.Container{
			Id:    "container-1-id",
			Names: []string{"container-1-name"},
		},
	})

	nodes[2].ID = "node-2-id"
	nodes[2].Name = "node-2-name"

	// Without constraints we should get the unfiltered list of nodes back.
	result, err = f.Filter(&dockerclient.ContainerConfig{}, nodes)
	assert.NoError(t, err)
	assert.Equal(t, result, nodes)

	// Set a constraint that cannot be fullfilled and expect an error back.
	result, err = f.Filter(&dockerclient.ContainerConfig{
		Env: []string{"affinity:container=does_not_exsits"},
	}, nodes)
	assert.Error(t, err)

	// Set a contraint that can only be filled by a single node.
	result, err = f.Filter(&dockerclient.ContainerConfig{
		Env: []string{"affinity:container=container-0*"},
	}, nodes)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, result[0], nodes[0])

	// This constraint can only be fullfilled by a subset of nodes.
	result, err = f.Filter(&dockerclient.ContainerConfig{
		Env: []string{"affinity:container=container-*"},
	}, nodes)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.NotContains(t, result, nodes[2])

	// Validate node pinning by id.
	result, err = f.Filter(&dockerclient.ContainerConfig{
		Env: []string{"affinity:container=container-0-id"},
	}, nodes)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, result[0], nodes[0])

	// Validate node pinning by name.
	result, err = f.Filter(&dockerclient.ContainerConfig{
		Env: []string{"affinity:container=container-1-name"},
	}, nodes)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, result[0], nodes[1])
}
