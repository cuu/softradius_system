package libs

import (
	"fmt"
	"syscall"
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func (this *DiskStatus) GetGB(key string) string{
	if key == "All" {
		return fmt.Sprintf("%.2f GB",float64(this.All)/float64(GB))
	}

	if key == "Used" {
		return fmt.Sprintf("%.2f GB",float64(this.Used)/float64(GB))
	}

	if key == "Free" {
		return fmt.Sprintf("%.2f GB",float64(this.Free)/float64(GB))
	}

	return ""
}


/*
func main() {
	disk := DiskUsage("/")
	fmt.Printf("All: %.2f GB\n", float64(disk.All)/float64(GB))
	fmt.Printf("Used: %.2f GB\n", float64(disk.Used)/float64(GB))
	fmt.Printf("Free: %.2f GB\n", float64(disk.Free)/float64(GB))
}
*/
