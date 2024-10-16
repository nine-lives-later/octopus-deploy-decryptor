package projectExport

type Entity interface {
	EntityID() string
	EntityName() string
	EntitySpaceID() string

	AddToEntityMap(m EntityMap)
}

type EntityMap map[string]Entity // Entity ID -> Entity

func (m EntityMap) Projects() []*Project {
	ret := make([]*Project, 0, len(m))

	for _, e := range m {
		if ee, ok := e.(*Project); ok && ee != nil {
			ret = append(ret, ee)
		}
	}

	return ret
}

func (m EntityMap) VariableSets() []*VariableSet {
	ret := make([]*VariableSet, 0, len(m))

	for _, e := range m {
		if ee, ok := e.(*VariableSet); ok && ee != nil {
			ret = append(ret, ee)
		}
	}

	return ret
}

func (m EntityMap) VariableSetByOwner(ownerID string) *VariableSet {
	for _, e := range m.VariableSets() {
		if e.OwnerId == ownerID {
			return e
		}
	}

	return nil
}

func ReadEntity(filename string) (Entity, error) {
	switch {
	case IsEnvironmentFilename(filename):
		return ReadEnvironment(filename)
	case IsLibraryVariableSetFilename(filename):
		return ReadLibraryVariableSet(filename)
	case IsProjectFilename(filename):
		return ReadProject(filename)
	case IsTenantFilename(filename):
		return ReadTenant(filename)
	case IsTenantVariableFilename(filename):
		return ReadTenantVariable(filename)
	case IsVariableSetFilename(filename):
		return ReadVariableSet(filename)
	}

	return nil, nil // ignore
}
