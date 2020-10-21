package prometheus_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/require"

	"github.com/zzpu/ums/metrics/prometheus"

	"github.com/zzpu/ums/internal"
	"github.com/zzpu/ums/x"
)

func TestHandler(t *testing.T) {
	_, reg := internal.NewFastRegistryWithMocks(t)
	router := x.NewRouterAdmin()
	reg.MetricsHandler().SetRoutes(router.Router)
	ts := httptest.NewServer(router)
	defer ts.Close()

	c := http.DefaultClient

	response, err := c.Get(ts.URL + prometheus.MetricsPrometheusPath)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusOK, response.StatusCode)

	textParser := expfmt.TextParser{}
	text, err := textParser.TextToMetricFamilies(response.Body)
	require.NoError(t, err)
	require.EqualValues(t, "go_info", *text["go_info"].Name)
}
