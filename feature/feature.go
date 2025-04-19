package feature

import (
	"flag"
	"qbittorrent_exporter/lib/log"
	"sync"
)

type FeatureFlag int

const (
	TRANSIENT_STATE FeatureFlag = iota
)

func (f FeatureFlag) String() string {
	var featureName = map[FeatureFlag]string{
		TRANSIENT_STATE: "transient-state",
	}
	return featureName[f]
}

var (
	flags map[FeatureFlag]bool = map[FeatureFlag]bool{}
	lock  sync.Mutex
)

func Use(features map[FeatureFlag]bool) func() {
	featureValues := make(map[FeatureFlag]*bool)
	for ff, val := range features {
		featureValues[ff] = &val
		flag.BoolVar(featureValues[ff], "ff-"+ff.String(), val, "[FeatureFlag]["+ff.String()+"]")
	}

	return func() {
		for f, value := range featureValues {
			Set(f, *value)
		}
	}
}

func Get(flag FeatureFlag) bool {
	lock.Lock()
	defer lock.Unlock()
	val, ok := flags[flag]
	if !ok {
		log.Warn("[FEATURE_FLAG][" + flag.String() + "] does not exist")
		return false
	}
	return val
}

func Set(flag FeatureFlag, value bool) {
	lock.Lock()
	defer lock.Unlock()
	flags[flag] = value
}
