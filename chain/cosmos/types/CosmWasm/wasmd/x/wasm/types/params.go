package types

// Modified for crosschain to remove various parts not needed for compatibility

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var AllAccessTypes = []AccessType{
	AccessTypeNobody,
	AccessTypeOnlyAddress,
	AccessTypeAnyOfAddresses,
	AccessTypeEverybody,
}

func (a AccessType) With(addrs ...sdk.AccAddress) AccessConfig {
	switch a {
	case AccessTypeNobody:
		return AllowNobody
	case AccessTypeOnlyAddress:
		if n := len(addrs); n != 1 {
			panic(fmt.Sprintf("expected exactly 1 address but got %d", n))
		}
		if err := sdk.VerifyAddressFormat(addrs[0]); err != nil {
			panic(err)
		}
		return AccessConfig{Permission: AccessTypeOnlyAddress, Address: addrs[0].String()}
	case AccessTypeEverybody:
		return AllowEverybody
	case AccessTypeAnyOfAddresses:
		bech32Addrs := make([]string, len(addrs))
		for i, v := range addrs {
			bech32Addrs[i] = v.String()
		}
		return AccessConfig{Permission: AccessTypeAnyOfAddresses, Addresses: bech32Addrs}
	}
	panic("unsupported access type")
}

func (a AccessType) String() string {
	switch a {
	case AccessTypeNobody:
		return "Nobody"
	case AccessTypeOnlyAddress:
		return "OnlyAddress"
	case AccessTypeEverybody:
		return "Everybody"
	case AccessTypeAnyOfAddresses:
		return "AnyOfAddresses"
	}
	return "Unspecified"
}

func (a *AccessType) UnmarshalText(text []byte) error {
	for _, v := range AllAccessTypes {
		if v.String() == string(text) {
			*a = v
			return nil
		}
	}
	*a = AccessTypeUnspecified
	return nil
}

func (a AccessType) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *AccessType) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(a)
}

func (a *AccessType) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, data []byte) error {
	return json.Unmarshal(data, a)
}

func (a AccessConfig) Equals(o AccessConfig) bool {
	return a.Permission == o.Permission && a.Address == o.Address
}

var (
	DefaultUploadAccess = AllowEverybody
	AllowEverybody      = AccessConfig{Permission: AccessTypeEverybody}
	AllowNobody         = AccessConfig{Permission: AccessTypeNobody}
)

// DefaultParams returns default wasm parameters
func DefaultParams() Params {
	return Params{
		CodeUploadAccess:             AllowEverybody,
		InstantiateDefaultPermission: AccessTypeEverybody,
	}
}

func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// ValidateBasic performs basic validation on wasm parameters
func (p Params) ValidateBasic() error {
	if err := validateAccessType(p.InstantiateDefaultPermission); err != nil {
		return errors.Wrap(err, "instantiate default permission")
	}
	if err := validateAccessConfig(p.CodeUploadAccess); err != nil {
		return errors.Wrap(err, "upload access")
	}
	return nil
}

func validateAccessConfig(i interface{}) error {
	v, ok := i.(AccessConfig)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return v.ValidateBasic()
}

func validateAccessType(i interface{}) error {
	return nil
}

// ValidateBasic performs basic validation
func (a AccessConfig) ValidateBasic() error {
	return nil
}

// Allowed returns if permission includes the actor.
// Actor address must be valid and not nil
func (a AccessConfig) Allowed(actor sdk.AccAddress) bool {
	switch a.Permission {
	case AccessTypeNobody:
		return false
	case AccessTypeEverybody:
		return true
	case AccessTypeOnlyAddress:
		return a.Address == actor.String()
	case AccessTypeAnyOfAddresses:
		for _, v := range a.Addresses {
			if v == actor.String() {
				return true
			}
		}
		return false
	default:
		panic("unknown type")
	}
}