package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	v1 "github.com/pennyscissors/plugin-operator/pkg/apis/catalog.cattle.io/v1"
	plugincontroller "github.com/pennyscissors/plugin-operator/pkg/generated/controllers/catalog.cattle.io/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
)

var (
	Index = SafeIndex{}
)

type SafeIndex struct {
	mu        sync.RWMutex
	Generated time.Time
	Entries   map[string]Entry
}
type Entry struct {
	*v1.UIPluginEntry
	Generated time.Time
}

func (s *SafeIndex) Generate(namespace string, cache plugincontroller.UIPluginCache) error {
	logrus.Debug("generating index from plugin controller's cache")
	s.mu.Lock()
	defer s.mu.Unlock()
	cachedPlugins, err := cache.List(namespace, labels.Everything())
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	s.Entries = make(map[string]Entry, len(cachedPlugins))
	for _, plugin := range cachedPlugins {
		entry := Entry{
			UIPluginEntry: &plugin.Spec.Plugin,
		}
		entry.Generated = time.Now()
		logrus.Debugf("adding plugin to index: %+v", *entry.UIPluginEntry)
		s.Entries[entry.Name] = entry
	}

	return nil
}

func (s *SafeIndex) SyncWithFsCache() error {
	files, err := filepath.Glob(FSCacheRootDir + "/*/*")
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	for _, file := range files {
		// Sync up entries in the filesystem cache with the index
		// Entries in the filesystem cache that aren't in the index will be deleted
		logrus.Debugf("syncing index with filesystem cache")
		// splits /root/pluginName/pluginVersion/* from a fs cache path
		s := strings.Split(file, "/")
		name := s[1]
		version := s[2]
		_, ok := Index.Entries[name]
		if ok && Index.Entries[name].Version == version {
			continue
		} else {
			// Delete plugin entry from filesystem cache
			err = os.RemoveAll(filepath.Join(FSCacheRootDir, name))
			if err != nil {
				logrus.Errorf("failed to delete entry [Name: %s Version: %s] from filesystem cache: %s", name, version, err.Error())
				return err
			}
			logrus.Debugf("deleted plugin entry from cache [Name: %s Version: %s]", name, version)
		}
	}

	return nil
}
