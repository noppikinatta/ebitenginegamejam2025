package system

type ResourceManager struct {
	resources map[string]int
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: map[string]int{
			"Gold":  100,
			"Iron":  50,
			"Wood":  75,
			"Grain": 60,
			"Mana":  25,
		},
	}
}

func (rm *ResourceManager) GetResource(resourceType string) int {
	return rm.resources[resourceType]
}

func (rm *ResourceManager) AddResource(resourceType string, amount int) {
	rm.resources[resourceType] += amount
}

func (rm *ResourceManager) ConsumeResources(costs map[string]int) bool {
	// Check if we have enough resources
	for resourceType, cost := range costs {
		if rm.resources[resourceType] < cost {
			return false
		}
	}

	// Consume resources
	for resourceType, cost := range costs {
		rm.resources[resourceType] -= cost
	}

	return true
}

func (rm *ResourceManager) GetAllResources() map[string]int {
	result := make(map[string]int)
	for k, v := range rm.resources {
		result[k] = v
	}
	return result
}

func (rm *ResourceManager) SetResource(resourceType string, amount int) {
	rm.resources[resourceType] = amount
}
