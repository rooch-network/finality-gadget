package types

import (
	"encoding/json"
	"fmt"
	"time"
)

/**
 *  ######################################
 *  ### DO NOT EDIT THIS FILE DIRECTLY ###
 *  ######################################
 *
 * This file is generated from the OpenRPC specification
 */

//type Block struct {
//	BlockHash      string `json:"block_hash" description:"block hash"`
//	BlockHeight    uint64 `json:"block_height" description:"block height"`
//	BlockTimestamp uint64 `json:"block_time" description:"block timestamp"`
//}

type BlockView struct {
	BlockHash   string `json:"block_hash" description:"block hash"`
	BlockHeight string `json:"block_height" description:"block height"`
	BlockTime   string `json:"block_time" description:"block timestamp"`
}

// PaginatedResponse is a generic type for paginated responses
type PaginatedResponse[T any] struct {
	Data        []T         `json:"data"`
	HasNextPage bool        `json:"has_next_page"`
	NextCursor  interface{} `json:"next_cursor,omitempty"`
}

// Define specific paginated response types
type PaginatedBlockViews = PaginatedResponse[BlockView]

// QueryOptions represents options for queries
type QueryOptions struct {
	Decode      *bool `json:"decode,omitempty"`
	Descending  *bool `json:"descending,omitempty"`
	FilterOut   *bool `json:"filterOut,omitempty"`
	ShowDisplay *bool `json:"showDisplay,omitempty"`
}

// StateOptions represents options for state queries
type StateOptions struct {
	Decode      *bool   `json:"decode,omitempty"`
	ShowDisplay *bool   `json:"showDisplay,omitempty"`
	StateRoot   *string `json:"stateRoot,omitempty"`
}

// Status represents the overall system status
type Status struct {
	//BitcoinStatus BitcoinStatus `json:"bitcoin_status"`
	RoochStatus   RoochStatus   `json:"rooch_status"`
	ServiceStatus ServiceStatus `json:"service_status"`
}

// ServiceStatus represents different service states
type ServiceStatus string

const (
	ServiceStatusActive         ServiceStatus = "active"
	ServiceStatusMaintenance    ServiceStatus = "maintenance"
	ServiceStatusReadOnlyMode   ServiceStatus = "read-only-mode"
	ServiceStatusDateImportMode ServiceStatus = "date-import-mode"
)

// String returns the string representation of ServiceStatus
func (s ServiceStatus) String() string {
	return string(s)
}

// Helper methods for PaginatedResponse
func NewPaginatedResponse[T any](data []T, hasNextPage bool, nextCursor interface{}) PaginatedResponse[T] {
	return PaginatedResponse[T]{
		Data:        data,
		HasNextPage: hasNextPage,
		NextCursor:  nextCursor,
	}
}

// Helper method for creating new StateOptions
func NewStateOptions(decode, showDisplay bool, stateRoot string) StateOptions {
	return StateOptions{
		Decode:      &decode,
		ShowDisplay: &showDisplay,
		StateRoot:   &stateRoot,
	}
}

// Helper method for creating new QueryOptions
func NewQueryOptions(decode, descending, filterOut, showDisplay bool) QueryOptions {
	return QueryOptions{
		Decode:      &decode,
		Descending:  &descending,
		FilterOut:   &filterOut,
		ShowDisplay: &showDisplay,
	}
}

// Helper method for VMErrorInfo
func NewVMErrorInfo(errorMessage string, executionState []string) VMErrorInfo {
	return VMErrorInfo{
		ErrorMessage:   errorMessage,
		ExecutionState: executionState,
	}
}

// IsActive checks if the service status is active
func (s ServiceStatus) IsActive() bool {
	return s == ServiceStatusActive
}

// IsMaintenance checks if the service status is in maintenance
func (s ServiceStatus) IsMaintenance() bool {
	return s == ServiceStatusMaintenance
}

// Helper method to check if a transaction is pending
func (t TransactionStatusView) IsPending() bool {
	return t == TransactionStatusPending
}

// Helper method to check if a transaction is executed
func (t TransactionStatusView) IsExecuted() bool {
	return t == TransactionStatusExecuted
}

// Helper method to check if a transaction failed
func (t TransactionStatusView) IsFailed() bool {
	return t == TransactionStatusFailed
}

// VMStatusView represents the status of a VM operation
type VMStatusView interface{}

// VMErrorInfo represents VM error information
type VMErrorInfo struct {
	ErrorMessage   string   `json:"error_message"`
	ExecutionState []string `json:"execution_state"`
}

// RoochStatus represents Rooch system status
type RoochStatus struct {
	ChainID           string `json:"chain_id"`
	LatestBlockHash   string `json:"latest_block_hash"`
	LatestBlockHeight string `json:"latest_block_height"`
}

// ScriptCallView represents a script call
type ScriptCallView struct {
	Args   []string `json:"args"`
	Script string   `json:"script"`
	TyArgs []string `json:"ty_args"`
}

// StateKVView represents a state key-value pair
type StateKVView struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// TransactionExecutionInfoView represents transaction execution information
type TransactionExecutionInfoView struct {
	Authentication []TransactionAuthenticatorView `json:"authentication"`
	Hash           string                         `json:"hash"`
	Status         TransactionStatusView          `json:"status"`
}

// TransactionAuthenticatorView represents a transaction authenticator
type TransactionAuthenticatorView struct {
	AuthenticatorType string               `json:"authenticator_type"`
	BitcoinAuthData   *BitcoinAuthDataView `json:"bitcoin_auth_data,omitempty"`
	PublicKey         *string              `json:"public_key,omitempty"`
	Signature         *string              `json:"signature,omitempty"`
}

// BitcoinAuthDataView represents Bitcoin authentication data
type BitcoinAuthDataView struct {
	Address   string `json:"address"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

// TransactionSequenceInfoView represents transaction sequence information
type TransactionSequenceInfoView struct {
	TxHash  string `json:"tx_hash"`
	TxOrder string `json:"tx_order"`
}

// TransactionStatusView represents transaction status
type TransactionStatusView string

const (
	TransactionStatusPending  TransactionStatusView = "pending"
	TransactionStatusExecuted TransactionStatusView = "executed"
	TransactionStatusFailed   TransactionStatusView = "failed"
)

// String returns the string representation of TransactionStatusView
func (t TransactionStatusView) String() string {
	return string(t)
}

// MarshalJSON implements json.Marshaler
func (t TransactionStatusView) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

// UnmarshalJSON implements json.Unmarshaler
func (t *TransactionStatusView) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case string(TransactionStatusPending),
		string(TransactionStatusExecuted),
		string(TransactionStatusFailed):
		*t = TransactionStatusView(s)
		return nil
	default:
		return fmt.Errorf("invalid TransactionStatusView value: %s", s)
	}
}

// MoveFunctionView represents a Move function
type MoveFunctionView struct {
	Name              string   `json:"name"`
	Visibility        string   `json:"visibility"`
	IsEntry           bool     `json:"is_entry"`
	GenericTypeParams []string `json:"generic_type_params"`
	Params            []string `json:"params"`
	Returns           []string `json:"returns"`
}

// MoveStructView represents a Move struct
type MoveStructView struct {
	Name              string          `json:"name"`
	IsNative          bool            `json:"is_native"`
	Abilities         []string        `json:"abilities"`
	GenericTypeParams []string        `json:"generic_type_params"`
	Fields            []MoveFieldView `json:"fields"`
}

// MoveFieldView represents a Move field
type MoveFieldView struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Helper methods for error handling
type RoochError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *RoochError) Error() string {
	return fmt.Sprintf("RoochError: [%d] %s", e.Code, e.Message)
}

// Helper method for creating new RoochError
func NewRoochError(code int, message string) *RoochError {
	return &RoochError{
		Code:    code,
		Message: message,
	}
}

// Helper methods for validation
func ValidateAddress(address string) error {
	if len(address) == 0 {
		return NewRoochError(400, "address cannot be empty")
	}
	// Add more validation logic as needed
	return nil
}

func ValidateHash(hash string) error {
	if len(hash) == 0 {
		return NewRoochError(400, "hash cannot be empty")
	}
	// Add more validation logic as needed
	return nil
}

// Helper method for creating pagination parameters
type PaginationParams struct {
	Limit  int         `json:"limit"`
	Cursor interface{} `json:"cursor,omitempty"`
}

func NewPaginationParams(limit int) PaginationParams {
	return PaginationParams{
		Limit: limit,
	}
}

func (p *PaginationParams) SetCursor(cursor interface{}) {
	p.Cursor = cursor
}

// Helper method for timestamp handling
func ParseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}

func FormatTimestamp(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// Helper methods for MoveStructView
func NewMoveStructView(name string, isNative bool, abilities []string) MoveStructView {
	return MoveStructView{
		Name:      name,
		IsNative:  isNative,
		Abilities: abilities,
	}
}

// Helper method to add fields to MoveStructView
func (m *MoveStructView) AddField(name, fieldType string) {
	m.Fields = append(m.Fields, MoveFieldView{
		Name: name,
		Type: fieldType,
	})
}

// Helper method to add generic type parameters
func (m *MoveStructView) AddGenericTypeParams(params ...string) {
	m.GenericTypeParams = append(m.GenericTypeParams, params...)
}

// Helper methods for ServiceStatus
func (s ServiceStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *ServiceStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	switch str {
	case string(ServiceStatusActive),
		string(ServiceStatusMaintenance),
		string(ServiceStatusReadOnlyMode),
		string(ServiceStatusDateImportMode):
		*s = ServiceStatus(str)
		return nil
	default:
		return fmt.Errorf("invalid ServiceStatus value: %s", str)
	}
}

// Helper methods for error handling with context
type ErrorContext struct {
	Operation string
	Details   map[string]interface{}
}

func NewErrorWithContext(code int, message string, context ErrorContext) *RoochError {
	return &RoochError{
		Code:    code,
		Message: fmt.Sprintf("%s: %s (Operation: %s)", message, formatErrorContext(context)),
	}
}

func formatErrorContext(context ErrorContext) string {
	details := ""
	for k, v := range context.Details {
		details += fmt.Sprintf("%s=%v, ", k, v)
	}
	if len(details) > 0 {
		details = details[:len(details)-2] // Remove trailing comma and space
	}
	return fmt.Sprintf("Operation: %s, Details: {%s}", context.Operation, details)
}

// Helper methods for pagination
func (p *PaginatedResponse[T]) HasMore() bool {
	return p.HasNextPage
}

func (p *PaginatedResponse[T]) GetNextCursor() interface{} {
	return p.NextCursor
}

func (p *PaginatedResponse[T]) GetItems() []T {
	return p.Data
}

// Helper method for creating empty responses
func EmptyPaginatedResponse[T any]() PaginatedResponse[T] {
	return PaginatedResponse[T]{
		Data:        make([]T, 0),
		HasNextPage: false,
	}
}
