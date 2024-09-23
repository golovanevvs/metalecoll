package agent

import (
	"testing"

	"github.com/golovanevvs/metalecoll/internal/agent/storage/mapstorage"
	"github.com/stretchr/testify/assert"
)

func TestRegrtmet(t *testing.T) {
	store := mapstorage.NewStorage()

	ag = NewAgent(store, 2, 10)

	RegisterMetrics()

	_, err := ag.store.GetMetricsMap()

	assert.NoError(t, err)
}
