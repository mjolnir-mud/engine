package yaml_data_source

type pluign struct{}

func (p pluign) Name() string {
	return "yaml_data_source"
}

func (p pluign) Start() error {
	return nil
}

func (p pluign) Stop() error {
	return nil
}
