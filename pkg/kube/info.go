package kube

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// pod info utils

func Labels() map[string]string {
	fileContent, err := readFiles("/etc/pod/metadata/labels")
	if err != nil {
		return nil
	}
	lines := strings.Split(fileContent[0], "\n")
	labels := make(map[string]string)
	for _, line := range lines {
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			// labels has the format : app.kubernetes.io/name="agent"
			// should delete qutation marks
			labels[parts[0]] = parts[1][1 : len(parts[1])-1]
		}
	}
	return labels
}

// Annotations stores the meta data of the pod and other important information
// "leaseLockName": the name of the lease lock resource
// "leaseLockNamespace": the namespace of the lease lock resource
// "region": the region of the pod
func Annotations() map[string]string {
	fileContent, err := readFiles("/etc/pod/metadata/annotations")
	if err != nil {
		return nil
	}
	lines := strings.Split(fileContent[0], "\n")
	annotations := make(map[string]string)
	for _, line := range lines {
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			annotations[parts[0]] = parts[1]
		}
	}
	return annotations
}

func GetPodNamespace() string {
	fileContent, err := readFiles("/etc/pod/metadata/namespace")
	if err != nil {
		return ""
	}
	lines := strings.Split(fileContent[0], "\n")
	for _, line := range lines {
		return line
	}
	return ""
}

func GetPodName() string {
	fileContent, err := readFiles("/etc/pod/metadata/name")
	if err != nil {
		return ""
	}
	lines := strings.Split(fileContent[0], "\n")
	for _, line := range lines {
		return line
	}
	return ""
}

func GetPodID() string {
	fileContent, err := readFiles("/etc/pod/metadata/id")
	if err != nil {
		return ""
	}
	lines := strings.Split(fileContent[0], "\n")
	for _, line := range lines {
		return line
	}
	return ""
}

func GetNodeName() string {
	po, err := APIClient.GetPod(GetPodNamespace(), GetPodName())
	if err != nil {
		return ""
	}
	return po.Spec.NodeName
}

func Runtime() map[string]string {
	r := map[string]string{
		"os":            runtime.GOOS,
		"arch":          runtime.GOARCH,
		"version":       runtime.Version(),
		"max_procs":     strconv.FormatInt(int64(runtime.GOMAXPROCS(0)), 10),
		"num_goroutine": strconv.FormatInt(int64(runtime.NumGoroutine()), 10),
		"num_cpu":       strconv.FormatInt(int64(runtime.NumCPU()), 10),
	}
	r["hostname"], _ = os.Hostname()
	return r
}

func LocalIP() (ip string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				// check if the ip is ipv6
				if strings.Contains(ip, ":") {
					// add '[' ']' to the ip
					if !strings.HasPrefix(ip, "[") && !strings.HasSuffix(ip, "]") {
						ip = "[" + ip + "]"
					}
				}
				return ip
			}
		}
	}
	return ""
}

func readFiles(dir string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		files = append(files, path)

		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Reading from %v failed", dir)
	}
	list := make([]string, 0)
	for _, path := range files {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "Reading %v failed", path)
		}
		content := string(data)
		duplicate := false
		for _, p := range list {
			if p == content {
				duplicate = true
				break
			}
		}
		if duplicate {
			continue
		}
		list = append(list, content)
	}
	return list, nil
}
