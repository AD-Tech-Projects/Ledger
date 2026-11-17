package repositories

import (
	"fmt"
	"reflect"
)

// Association is a belongs-to association of the model M whose field holds a *T.
// It is declared once next to the model with MustAssociation, so Join only accepts associations of the repository model.
type Association[M any, T any] struct {
	path string
}

// MustAssociation declares the association field of the model M that holds a *T.
// NOTE - it panics when M has no field with this name of type *T, so a wrong declaration fails when the models package
// is loaded (on startup and on any test that imports it) instead of at query time
func MustAssociation[M any, T any](field string) Association[M, T] {
	modelType := reflect.TypeFor[M]()
	targetType := reflect.TypeFor[*T]()

	isStruct := modelType.Kind() == reflect.Struct
	if isStruct {
		structField, found := modelType.FieldByName(field)
		isAssociationField := found && structField.Type == targetType
		if isAssociationField {
			return Association[M, T]{path: field}
		}
	}

	panic(fmt.Sprintf("repositories: %s has no association field %q of type %s", modelType, field, targetType))
}

// Nested is the association next reached through the association parent (e.g. the owner of a transaction's source wallet).
// It is undeclared (zero value) when either association is undeclared.
func Nested[M any, T any, U any](parent Association[M, T], next Association[T, U]) Association[M, U] {
	isDeclared := parent.path != "" && next.path != ""
	if !isDeclared {
		return Association[M, U]{}
	}
	return Association[M, U]{path: parent.path + "." + next.path}
}
