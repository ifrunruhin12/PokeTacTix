package battle

import (
	"context"
	"errors"
	"testing"

	"pokemon-cli/internal/pokemon"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type evolutionTestService struct {
	targets  map[int]*pokemon.Pokemon
	failures map[int]error
	calls    []int
}

func (s *evolutionTestService) GetByID(context.Context, int) (*pokemon.Pokemon, error) {
	return nil, nil
}

func (s *evolutionTestService) GetByName(context.Context, string) (*pokemon.Pokemon, error) {
	return nil, nil
}

func (s *evolutionTestService) EnsureEvolutionChain(context.Context, int, string) (int, error) {
	return 0, nil
}

func (s *evolutionTestService) GetRandomCard(context.Context, bool) (pokemon.Card, error) {
	return pokemon.Card{}, nil
}

func (s *evolutionTestService) GetEvolutionForLevel(_ context.Context, pokemonID, _ int) (*pokemon.Pokemon, error) {
	s.calls = append(s.calls, pokemonID)
	if err := s.failures[pokemonID]; err != nil {
		return nil, err
	}
	return s.targets[pokemonID], nil
}

type evolutionExecCall struct {
	arguments []any
}

type evolutionTestExecutor struct {
	calls []evolutionExecCall
}

func (e *evolutionTestExecutor) Exec(_ context.Context, _ string, arguments ...any) (pgconn.CommandTag, error) {
	e.calls = append(e.calls, evolutionExecCall{arguments: arguments})
	return pgconn.NewCommandTag("UPDATE 1"), nil
}

func TestEvolutionChainAppliesEveryEligibleStage(t *testing.T) {
	charmeleon := &pokemon.Pokemon{ID: 5, Name: "charmeleon", Types: []string{"fire"}}
	charizard := &pokemon.Pokemon{ID: 6, Name: "charizard", Types: []string{"fire", "flying"}}
	service := &evolutionTestService{targets: map[int]*pokemon.Pokemon{
		4: charmeleon,
		5: charizard,
	}}

	targets, err := resolveEvolutionTargets(context.Background(), service, 4, 36)
	require.NoError(t, err)
	require.Equal(t, []*pokemon.Pokemon{charmeleon, charizard}, targets)
	assert.Equal(t, []int{4, 5, 6}, service.calls)

	executor := &evolutionTestExecutor{}
	finalTarget, err := applyEvolution(context.Background(), executor, 10, 20, targets)
	require.NoError(t, err)
	assert.Same(t, charizard, finalTarget)
	require.Len(t, executor.calls, 2)
	assert.Equal(t, 5, executor.calls[0].arguments[1])
	assert.Equal(t, 6, executor.calls[1].arguments[1])
}

func TestResolveEvolutionTargetsDiscardsPartialChainOnError(t *testing.T) {
	service := &evolutionTestService{
		targets: map[int]*pokemon.Pokemon{
			4: {ID: 5, Name: "charmeleon"},
		},
		failures: map[int]error{5: errors.New("cold refresh failed")},
	}

	targets, err := resolveEvolutionTargets(context.Background(), service, 4, 36)
	require.Error(t, err)
	assert.Nil(t, targets)
	assert.Equal(t, []int{4, 5}, service.calls)
}
