package treelite

// #include <stdlib.h>
// #include "include/treelite/c_api.h"
// #include "include/treelite/c_api_common.h"
// #include "include/treelite/c_api_runtime.h"
import "C"
import (
	"errors"
	"io"
	"unsafe"
)

var (
	errLoadModel = errors.New("load model error")
	errFreeModel = errors.New("free model error")
)

// Model is a tree model
type Model struct {
	pointer C.ModelHandle
}

// Close frees internally allocated model
func (m *Model) Close() error {
	ret := C.TreeliteFreeModel(
		m.pointer,
	)
	if ret == -1 {
		return errFreeModel
	}
	return nil
}

// LoadLightGBMModel loads a model file generated by LightGBM (Microsoft/LightGBM).
// The model file must contain a decision tree ensemble.
func LoadLightGBMModel(modelPath string) (*Model, error) {
	var handle C.ModelHandle
	ret := C.TreeliteLoadLightGBMModel(
		C.CString(modelPath),
		&handle,
	)
	if ret == -1 {
		return nil, errLoadModel
	}
	return &Model{
		pointer: handle,
	}, nil
}

// LoadXGBoostModel loads a model file generated by XGBoost (dmlc/xgboost).
// The model file must contain a decision tree ensemble.
func LoadXGBoostModel(modelPath string) (*Model, error) {
	var handle C.ModelHandle
	ret := C.TreeliteLoadXGBoostModel(
		C.CString(modelPath),
		&handle,
	)
	if ret == -1 {
		return nil, errLoadModel
	}
	return &Model{
		pointer: handle,
	}, nil
}

// LoadXGBoostJSON loads a json model file generated by XGBoost (dmlc/xgboost).
// The model file must contain a decision tree ensemble.
func LoadXGBoostJSON(modelPath string) (*Model, error) {
	var handle C.ModelHandle
	ret := C.TreeliteLoadXGBoostJSON(
		C.CString(modelPath),
		&handle,
	)
	if ret == -1 {
		return nil, errLoadModel
	}
	return &Model{
		pointer: handle,
	}, nil
}

// LoadXGBoostJSONString loads a model stored as JSON stringby XGBoost (dmlc/xgboost).
// The model json must contain a decision tree ensemble.
func LoadXGBoostJSONString(modelJSON string) (*Model, error) {
	var handle C.ModelHandle
	ret := C.TreeliteLoadXGBoostJSONString(
		C.CString(modelJSON),
		C.ulong(len(modelJSON)),
		&handle,
	)
	if ret == -1 {
		return nil, errLoadModel
	}
	return &Model{
		pointer: handle,
	}, nil
}

// LoadXGBoostModelFromMemoryBuffer loads XGBoost model from the given reader
// this method once loads all model to memory
func LoadXGBoostModelFromMemoryBuffer(reader io.Reader) (*Model, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var handle C.ModelHandle
	ret := C.TreeliteLoadXGBoostModelFromMemoryBuffer(
		unsafe.Pointer(&buf[0]),
		C.ulong(len(buf)),
		&handle,
	)
	if ret == -1 {
		return nil, errLoadModel
	}
	return &Model{
		pointer: handle,
	}, nil
}

// LoadLightGBMModelFromString Loads a LightGBM model from a string.
// The string should be created with the model_to_string() method in LightGBM.
func LoadLightGBMModelFromString(modelString string) (*Model, error) {
	var handle C.ModelHandle
	ret := C.TreeliteLoadLightGBMModelFromString(
		C.CString(modelString),
		&handle,
	)
	if ret == -1 {
		return nil, errLoadModel
	}
	return &Model{
		pointer: handle,
	}, nil
}
