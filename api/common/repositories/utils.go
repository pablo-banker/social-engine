package repositories

import (
	"errors"
	"fmt"
	"reflect"
	"social-engine/common/repositories/entities"
	"strings"
)

type ContextType string

const (
	GormTestContext ContextType = "gormTestContext"
)

func ParseEntity(result entities.IEntity, entity entities.IEntity) error {
	if result == nil {
		return errors.New("result cannot be nil")
	}
	if entity == nil {
		return errors.New("entity cannot be nil")
	}

	resultVal := reflect.ValueOf(result)
	entityVal := reflect.ValueOf(entity)

	if resultVal.Kind() != reflect.Ptr || resultVal.Elem().Kind() != reflect.Struct {
		return errors.New("result must be a pointer to a struct")
	}
	if entityVal.Kind() != reflect.Ptr || entityVal.Elem().Kind() != reflect.Struct {
		return errors.New("entity must be a pointer to a struct")
	}

	if !resultVal.Type().AssignableTo(entityVal.Type()) &&
		!resultVal.Type().ConvertibleTo(entityVal.Type()) {
		return fmt.Errorf("cannot assign %T to %T", result, entity)
	}

	entityVal.Elem().Set(resultVal.Elem())
	return nil
}

func isEmptyID(id any) bool {
	if id == nil {
		return true
	}

	rv := reflect.ValueOf(id)

	// Se for ponteiro, pega o valor apontado
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return true
		}
		rv = rv.Elem()
	}

	// Se não for válido (tipo inválido, nil interno, etc)
	if !rv.IsValid() {
		return true
	}

	// Compara com zero do tipo
	zero := reflect.Zero(rv.Type())
	return reflect.DeepEqual(rv.Interface(), zero.Interface())
}

func Normalize(s string) string {
	return strings.TrimSpace(strings.ReplaceAll(s, "\t", ""))
}
