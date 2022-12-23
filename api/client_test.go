package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_GetTracks(t *testing.T) {
	c := NewClient()
	tracks, err := c.GetTracks("300204")
	require.NoError(t, err)
	fmt.Println(tracks)
}
