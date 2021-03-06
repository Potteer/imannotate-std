package memorygc

import (
	"time"

	"github.com/smileinnovation/imannotate/api/providers"
)

type MemoryGC struct {
	ImageProvider providers.ImageProvider
	mem           map[string]string
	ttl           time.Duration
}

func NewMemoryGC() *MemoryGC {
	mm := &MemoryGC{
		mem: make(map[string]string),
		ttl: 1 * time.Minute,
	}

	return mm
}

func (gc *MemoryGC) Collect(name string, data string) {

	gc.mem[name] = data
	go func(name, data string) {
		time.Sleep(gc.ttl)
		if _, ok := gc.mem[name]; ok {
			gc.ImageProvider.AddImage(name, data)
		}
	}(name, data)
}

func (gc *MemoryGC) Delete(name string) error {
	delete(gc.mem, name)
	return nil
}

func (gc *MemoryGC) SetImageProvider(p providers.ImageProvider) {
	gc.ImageProvider = p
}

func (gc *MemoryGC) SetTTL(ttl time.Duration) {
	gc.ttl = ttl
}
