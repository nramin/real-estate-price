package structs

type RealEstatePrices struct {
	Properties  []PropertyDetails `json:"properties,omitempty" bson:"properties,omitempty"`
	TotalEquity float64           `json:"totalEquity,omitempty" bson:"totalEquity,omitempty"`
	TotalValue  float64           `json:"totalValue,omitempty" bson:"totalValue,omitempty"`
	Success     *bool             `json:"success,omitempty" bson:"success,omitempty"`
	Error       string            `json:"error,omitempty" bson:"error,omitempty"`
}

type PropertyDetails struct {
	Address string  `json:"address,omitempty"`
	Price   float64 `json:"price,omitempty"`
	Equity  float64 `json:"equity,omitempty"`
	Balance float64 `json:"balance,omitempty"`
}

type YamlConfig struct {
	Properties []Property `yaml:"properties"`
}

type Property struct {
	Address string  `yaml:"address"`
	Zillow  string  `yaml:"zillow"`
	Balance float64 `yaml:"balance"`
}
