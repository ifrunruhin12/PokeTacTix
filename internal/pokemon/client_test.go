package pokemon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPokeAPIClientEscapesPokemonIdentifiers(t *testing.T) {
	var paths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.EscapedPath())
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewPokeAPIClient(server.URL, time.Second)
	_, err := client.FetchPokemonRaw(context.Background(), "Mr/Mime")
	require.NoError(t, err)
	_, err = client.FetchSpeciesRaw(context.Background(), "Mr/Mime")
	require.NoError(t, err)

	assert.Equal(t, []string{"/pokemon/mr%2Fmime", "/pokemon-species/mr%2Fmime"}, paths)
}
