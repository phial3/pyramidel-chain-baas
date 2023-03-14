package configtxgen

//
//type TopLevel struct {
//	Organizations []Organizations `yaml:"Organizations"`
//	Capabilities  Capabilities    `yaml:"Capabilities"`
//	Application   Application     `yaml:"Application"`
//	Orderer       Orderer         `yaml:"Orderer"`
//	Channel       Channel         `yaml:"Channel"`
//	Profiles      Profiles        `yaml:"Profiles"`
//}
//type Readers struct {
//	Type string `yaml:"Type"`
//	Rule string `yaml:"Rule"`
//}
//type Writers struct {
//	Type string `yaml:"Type"`
//	Rule string `yaml:"Rule"`
//}
//type Admins struct {
//	Type string `yaml:"Type"`
//	Rule string `yaml:"Rule"`
//}
//type Endorsement struct {
//	Type string `yaml:"Type"`
//	Rule string `yaml:"Rule"`
//}
//type Policies struct {
//	Readers     Readers     `yaml:"Readers"`
//	Writers     Writers     `yaml:"Writers"`
//	Admins      Admins      `yaml:"Admins"`
//	Endorsement Endorsement `yaml:"Endorsement"`
//}
//type AnchorPeers struct {
//	Host string `yaml:"Host"`
//	Port int    `yaml:"Port"`
//}
//type Organizations struct {
//	Name             string        `yaml:"Name"`
//	ID               string        `yaml:"ID"`
//	MSPDir           string        `yaml:"MSPDir"`
//	Policies         Policies      `yaml:"Policies"`
//	AnchorPeers      []AnchorPeers `yaml:"AnchorPeers"`
//	OrdererEndpoints []string      `yaml:"OrdererEndpoints"`
//}
//type Channel struct {
//	V20 bool `yaml:"V2_0"`
//}
//type Orderer struct {
//	V20 bool `yaml:"V2_0"`
//}
//type Application struct {
//	V20 bool `yaml:"V2_0"`
//}
//type Capabilities struct {
//	Channel     Channel     `yaml:"Channel"`
//	Orderer     Orderer     `yaml:"Orderer"`
//	Application Application `yaml:"Application"`
//}
//type LifecycleEndorsement struct {
//	Type string `yaml:"Type"`
//	Rule string `yaml:"Rule"`
//}
//type Policies struct {
//	Readers              Readers              `yaml:"Readers"`
//	Writers              Writers              `yaml:"Writers"`
//	Admins               Admins               `yaml:"Admins"`
//	LifecycleEndorsement LifecycleEndorsement `yaml:"LifecycleEndorsement"`
//	Endorsement          Endorsement          `yaml:"Endorsement"`
//}
//type Capabilities struct {
//	V20 bool `yaml:"V2_0"`
//}
//type Application struct {
//	Organizations interface{}  `yaml:"Organizations"`
//	Policies      Policies     `yaml:"Policies"`
//	Capabilities  Capabilities `yaml:"Capabilities"`
//}
//type Consenters struct {
//	Host          string `yaml:"Host"`
//	Port          int    `yaml:"Port"`
//	ClientTLSCert string `yaml:"ClientTLSCert"`
//	ServerTLSCert string `yaml:"ServerTLSCert"`
//}
//type EtcdRaft struct {
//	Consenters []Consenters `yaml:"Consenters"`
//}
//type BatchSize struct {
//	MaxMessageCount   int    `yaml:"MaxMessageCount"`
//	AbsoluteMaxBytes  string `yaml:"AbsoluteMaxBytes"`
//	PreferredMaxBytes string `yaml:"PreferredMaxBytes"`
//}
//type BlockValidation struct {
//	Type string `yaml:"Type"`
//	Rule string `yaml:"Rule"`
//}
//type Policies struct {
//	Readers         Readers         `yaml:"Readers"`
//	Writers         Writers         `yaml:"Writers"`
//	Admins          Admins          `yaml:"Admins"`
//	BlockValidation BlockValidation `yaml:"BlockValidation"`
//}
//type Orderer struct {
//	OrdererType   string      `yaml:"OrdererType"`
//	Addresses     []string    `yaml:"Addresses"`
//	EtcdRaft      EtcdRaft    `yaml:"EtcdRaft"`
//	BatchTimeout  string      `yaml:"BatchTimeout"`
//	BatchSize     BatchSize   `yaml:"BatchSize"`
//	Organizations interface{} `yaml:"Organizations"`
//	Policies      Policies    `yaml:"Policies"`
//}
//type Policies struct {
//	Readers Readers `yaml:"Readers"`
//	Writers Writers `yaml:"Writers"`
//	Admins  Admins  `yaml:"Admins"`
//}
//type Channel struct {
//	Policies     Policies     `yaml:"Policies"`
//	Capabilities Capabilities `yaml:"Capabilities"`
//}
//type Policies struct {
//	Readers         Readers         `yaml:"Readers"`
//	Writers         Writers         `yaml:"Writers"`
//	Admins          Admins          `yaml:"Admins"`
//	BlockValidation BlockValidation `yaml:"BlockValidation"`
//}
//type Orderer struct {
//	OrdererType   string          `yaml:"OrdererType"`
//	Addresses     []string        `yaml:"Addresses"`
//	EtcdRaft      EtcdRaft        `yaml:"EtcdRaft"`
//	BatchTimeout  string          `yaml:"BatchTimeout"`
//	BatchSize     BatchSize       `yaml:"BatchSize"`
//	Organizations []Organizations `yaml:"Organizations"`
//	Policies      Policies        `yaml:"Policies"`
//	Capabilities  Capabilities    `yaml:"Capabilities"`
//}
//type SampleConsortium struct {
//	Organizations []Organizations `yaml:"Organizations"`
//}
//type Consortiums struct {
//	SampleConsortium SampleConsortium `yaml:"SampleConsortium"`
//}
//type ThreeOrgsOrdererGenesis struct {
//	Policies     Policies     `yaml:"Policies"`
//	Capabilities Capabilities `yaml:"Capabilities"`
//	Orderer      Orderer      `yaml:"Orderer"`
//	Consortiums  Consortiums  `yaml:"Consortiums"`
//}
//type Policies struct {
//	Readers Readers `yaml:"Readers"`
//	Writers Writers `yaml:"Writers"`
//	Admins  Admins  `yaml:"Admins"`
//}
//type ThreeOrgsChannel struct {
//	Consortium   string       `yaml:"Consortium"`
//	Policies     Policies     `yaml:"Policies"`
//	Capabilities Capabilities `yaml:"Capabilities"`
//	Application  Application  `yaml:"Application"`
//}
//type Profiles struct {
//	ThreeOrgsOrdererGenesis ThreeOrgsOrdererGenesis `yaml:"ThreeOrgsOrdererGenesis"`
//	ThreeOrgsChannel        ThreeOrgsChannel        `yaml:"ThreeOrgsChannel"`
//}
