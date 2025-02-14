package epochs

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/iotaledger/hive.go/node"
	"github.com/labstack/echo"

	"github.com/iotaledger/goshimmer/packages/clock"
	"github.com/iotaledger/goshimmer/plugins/messagelayer"
	"github.com/iotaledger/goshimmer/plugins/webapi"
	"github.com/iotaledger/goshimmer/plugins/webapi/jsonmodels"
)

var (
	// plugin is the plugin instance of the web API epochs endpoint plugin.
	plugin *node.Plugin
	once   sync.Once
)

// Plugin returns the plugin as a singleton.
func Plugin() *node.Plugin {
	once.Do(func() {
		plugin = node.NewPlugin("WebAPI epochs Endpoint", node.Enabled, func(*node.Plugin) {
			webapi.Server().GET("epochs/:epochID", GetEpoch)
			webapi.Server().GET("epochs/current", GetEpochID)
			webapi.Server().GET("epochs/time/:unixTime", GetEpochID)
			webapi.Server().GET("epochs/oracle/current", GetOracleEpochID)
			webapi.Server().GET("epochs/oracle/time/:unixTime", GetOracleEpochID)
		})
	})

	return plugin
}

// GetEpoch returns the weights and total weight of active nodes for the given epoch ID.
//{
//   "weights": [
//       {
//           "shortNodeID": "4AeXyZ26e4G",
//           "nodeID": "2GtxMQD94KvDH1SJPJV7icxofkyV1njuUZKtsqKmtux5",
//           "mana": 1000000000000000
//       }
//   ],
//   "totalWeight": 1000000000000000
//}
func GetEpoch(c echo.Context) error {
	epochID, err := strconv.ParseUint(c.Param("epochID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, jsonmodels.NewErrorResponse(err))
	}
	start := messagelayer.Tangle().WeightProvider.EpochIDToStartTime(epochID)
	end := messagelayer.Tangle().WeightProvider.EpochIDToEndTime(epochID)
	weights, totalWeight := messagelayer.Tangle().WeightProvider.WeightsOfRelevantSupporters(epochID)
	return c.JSON(http.StatusOK, jsonmodels.NewEpoch(epochID, start, end, weights, totalWeight))
}

// GetOracleEpochID gets the oracle epoch ID for the given time or the current epoch if no time is specified.
//{
//    "epochId": 1484
//}
func GetOracleEpochID(c echo.Context) error {
	var unixTime int64
	if c.Param("unixTime") != "" {
		var err error
		unixTime, err = strconv.ParseInt(c.Param("unixTime"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, jsonmodels.NewErrorResponse(err))
		}
	} else {
		unixTime = clock.SyncedTime().Unix()
	}

	epochID := messagelayer.Tangle().WeightProvider.OracleEpoch(time.Unix(unixTime, 0))
	return c.JSON(http.StatusOK, jsonmodels.EpochID{EpochID: epochID})
}

// GetEpochID gets the epoch ID for the given time or the current epoch if no time is specified.
//{
//    "epochId": 1484
//}
func GetEpochID(c echo.Context) error {
	var unixTime int64
	if c.Param("unixTime") != "" {
		var err error
		unixTime, err = strconv.ParseInt(c.Param("unixTime"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, jsonmodels.NewErrorResponse(err))
		}
	} else {
		unixTime = clock.SyncedTime().Unix()
	}

	epochID := messagelayer.Tangle().WeightProvider.Epoch(time.Unix(unixTime, 0))
	return c.JSON(http.StatusOK, jsonmodels.EpochID{EpochID: epochID})
}
