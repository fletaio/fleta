package formulator

import (
	"bytes"
	"encoding/json"

	"github.com/fletaio/fleta/common"
	"github.com/fletaio/fleta/core/types"
	"github.com/fletaio/fleta/encoding"
	"github.com/fletaio/fleta/process/admin"
)

// UpdateAlphaPolicy is used to update reward policy
type UpdateAlphaPolicy struct {
	Timestamp_ uint64
	From_      common.Address
	Policy     *AlphaPolicy
}

// Timestamp returns the timestamp of the transaction
func (tx *UpdateAlphaPolicy) Timestamp() uint64 {
	return tx.Timestamp_
}

// From returns the from address of the transaction
func (tx *UpdateAlphaPolicy) From() common.Address {
	return tx.From_
}

// Validate validates signatures of the transaction
func (tx *UpdateAlphaPolicy) Validate(p types.Process, loader types.LoaderWrapper, signers []common.PublicHash) error {
	sp := p.(*Formulator)

	if tx.From() != sp.admin.AdminAddress(loader, p.Name()) {
		return admin.ErrUnauthorizedTransaction
	}
	if tx.Policy == nil {
		return ErrInvalidAlphaPolicy
	}

	fromAcc, err := loader.Account(tx.From())
	if err != nil {
		return err
	}
	if err := fromAcc.Validate(loader, signers); err != nil {
		return err
	}
	return nil
}

// Execute updates the context by the transaction
func (tx *UpdateAlphaPolicy) Execute(p types.Process, ctw *types.ContextWrapper, index uint16) error {
	if bs, err := encoding.Marshal(tx.Policy); err != nil {
		return err
	} else {
		ctw.SetProcessData(tagAlphaPolicy, bs)
	}
	return nil
}

// MarshalJSON is a marshaler function
func (tx *UpdateAlphaPolicy) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString(`{`)
	buffer.WriteString(`"timestamp":`)
	if bs, err := json.Marshal(tx.Timestamp_); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"from":`)
	if bs, err := tx.From_.MarshalJSON(); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`,`)
	buffer.WriteString(`"policy":`)
	if bs, err := tx.Policy.MarshalJSON(); err != nil {
		return nil, err
	} else {
		buffer.Write(bs)
	}
	buffer.WriteString(`}`)
	return buffer.Bytes(), nil
}
