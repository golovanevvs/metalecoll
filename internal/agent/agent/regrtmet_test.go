package agent

import (
	"testing"

	"github.com/golovanevvs/metalecoll/internal/agent/autil"
	"github.com/golovanevvs/metalecoll/internal/agent/storage/amapstorage"
	"github.com/stretchr/testify/assert"
)

func TestRegrtmet(t *testing.T) {
	store := amapstorage.NewStorage()

	ag = NewAgent(store, 2, 10)

	RegisterMetrics()

	_, err := autil.GMM(ag.store)

	assert.NoError(t, err)
}
