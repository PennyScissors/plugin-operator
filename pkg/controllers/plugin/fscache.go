package plugin

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	plugincontroller "github.com/pennyscissors/plugin-operator/pkg/generated/controllers/catalog.cattle.io/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	FSCacheRootDir = "plugincache"

	// Cache states used by custom resources
	Cached   = "cached"
	Disabled = "disabled"
	Pending  = "pending"
)

var (
	FsCache = FSCache{}
)

type FSCache struct{}

func (c *FSCache) Sync(cache plugincontroller.UIPluginCache) error {
	// Sync filesystem cache with plugins in the plugin controller's cache
	cachedPlugins, err := cache.List(UIPluginNamesapce, labels.Everything())
	if err != nil {
		return err
	}
	for _, p := range cachedPlugins {
		plugin := p.Spec.Plugin
		if plugin.NoCache {
			logrus.Debugf("skipped caching plugin [Name: %s Version: %s] cache is disabled [noCache: %v]", plugin.Name, plugin.Version, plugin.NoCache)
			continue
		}
		if isCached, err := c.isCached(plugin.Name, plugin.Version); err != nil {
			return err
		} else if isCached {
			logrus.Debugf("skipped caching plugin [Name: %s Version: %s] is already cached", plugin.Name, plugin.Version)
			continue
		}
		err = c.Save(plugin.Name, plugin.Version)
		if err != nil {
			return err
		}
		urlFilesTxt := fmt.Sprintf("%s/%s", plugin.Endpoint, FilesTxtFilename)
		files, err := getPluginFiles(urlFilesTxt)
		if err != nil {
			return err
		}
		for _, file := range files {
			err = fetchFile(plugin.Endpoint, plugin.Name, plugin.Version, file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *FSCache) Save(name, version string) error {
	err := os.MkdirAll(filepath.Join(FSCacheRootDir, name, version, "plugin"), os.ModePerm)
	if err != nil {
		logrus.Debugf("failed to cache plugin [Name: %s Version: %s] in filesystem", name, version)
		return err
	}

	return nil
}

func (c *FSCache) isCached(name, version string) (bool, error) {
	_, err := os.Stat(filepath.Join(FSCacheRootDir, name, version, "plugin"))
	if !errors.Is(err, os.ErrNotExist) {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return false, nil
}

// Might need this later if we decide to clear the cache on exit?
// func (c *FSCache) Clear() {
// 	if _, err := os.Stat(FSCacheRootDir); err == nil {
// 		os.RemoveAll(FSCacheRootDir)
// 	} else {
// 		logrus.Error(err)
// 	}
// }

func getPluginFiles(urlFilesTxt string) ([]string, error) {
	var files []string
	logrus.Debugf("fetching file [%s]", urlFilesTxt)
	resp, err := http.Get(urlFilesTxt)
	if err != nil {
		return files, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return files, err
	}
	files = strings.Split(string(b), "\n")

	return files, nil
}

func fetchFile(endpoint, name, version, file string) error {
	url := endpoint + "/" + file
	logrus.Debugf("fetching file [%s]", url)
	resp, err := http.Get(url)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer resp.Body.Close()
	filepath := filepath.Join(FSCacheRootDir, name, version, file)
	logrus.Debugf("creating file [%s]", url)
	out, err := os.Create(filepath)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	return nil
}
