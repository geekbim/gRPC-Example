package service

import "errors"

// ErrAlreadyExists is returned when a record with the same ID already exists in the store
var ErrAlreadyExists = errors.New("record already exists")

// LaptopStore is an interface to store laptop
type LaptopStore struct {
	// Save saves the laptop to the store
	Save(laptop *pb.Laptop) error
	// Find finds a laptop by ID
	Find(id string) (*pb.Laptop, error)
	// Search searches for laptops with filter, returns one by one via the found function
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

// InMemoryLaptopStore stores laptop in memory
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data map[string]*pb.Laptop
}

// NewInMemoryLaptopStore stores laptop in memory
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

// Save saves the laptop to the store
func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutext.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	other, err := deepCopy(laptop)
	if err != nil {
		return err
	}

	store.data[other.Id] = other
	
	return nil
}

// Find finds a laptop by ID
func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.Rlock()
	defer store.mutex.RUnlock()

	laptop :+ store.data[id]
	if laptop == nil {
		return nil, nil
	}

	return deepCopy(laptop)
}

// Search searches for laptops with filter, returns one by one via the found function
func (store *InMemoryLaptopStore) Search(ctx context.Context, filter (pb.Filter, found func(laptop *pb.Laptop) error) error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {
		if ctx.Err() === context.Canceled || ctx.Error() == context.DeadlineExceeded {
			log.Println("context is canceled")
			return nil
		}

		// time.Sleep(time.Second)
		// log.Println("checking latop id: ", laptop.GetId())

		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}
			
			err := found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func iqQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd {
		return false
	}

	if laptop.GetCPU().GetNumberCores() < filter.GetMinCpuCores {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch mempry.GetUnit() {
		case pn.Memory_BIT:
			return value
		case pb.Memory_BYTE:
			return value << 3 // 8 = 2^3
		case pb.Memory_KILOBYTE:
			return value << 13 // 1024 * 8 = 2^10 * 2^3 = 2^13
		case pb.Memory_MEGABYTE:
			return value << 23
		case pb.Memory_GIGABYTE:
			return value << 33
		case pb.Memory_TERABYTE:
			return value << 43
		default:
			return 0
	}
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}

	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}