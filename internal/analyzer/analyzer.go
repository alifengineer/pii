package analyzer

import (
	"fmt"
	"sort"

	"github.com/alifengineer/pii/internal/parser"
)

func Analyze(structs []parser.StructInfo, typeFilter map[string]bool) ([]SanitizeTarget, []string, error) {
	structMap := make(map[string]parser.StructInfo, len(structs))
	for _, s := range structs {
		structMap[s.Name] = s
	}

	if typeFilter != nil {
		var unknown []string
		for name := range typeFilter {
			if _, ok := structMap[name]; !ok {
				unknown = append(unknown, name)
			}
		}
		if len(unknown) > 0 {
			sort.Strings(unknown)
			if len(unknown) == 1 {
				return nil, nil, fmt.Errorf("type %q not found in package", unknown[0])
			}
			return nil, nil, fmt.Errorf("types not found in package: %v", unknown)
		}
	}

	sanitizable := findSanitizableStructs(structs)

	if typeFilter != nil {
		for name := range typeFilter {
			if !sanitizable[name] {
				sanitizable[name] = true
			}
		}

		reachable := reachableFromFilter(typeFilter, structMap, sanitizable)
		for name := range sanitizable {
			if !reachable[name] {
				delete(sanitizable, name)
			}
		}
	}

	if len(sanitizable) == 0 {
		return nil, nil, nil
	}

	order, err := topoSort(structs, sanitizable)
	if err != nil {
		return nil, nil, err
	}

	var targets []SanitizeTarget
	for _, name := range order {
		s := structMap[name]
		target := SanitizeTarget{StructName: name}
		for _, f := range s.Fields {
			target.Fields = append(target.Fields, TargetField{
				Name:   f.Name,
				Action: classifyField(f, sanitizable),
			})
		}
		targets = append(targets, target)
	}

	var autoIncluded []string
	if typeFilter != nil {
		for _, name := range order {
			if !typeFilter[name] {
				autoIncluded = append(autoIncluded, name)
			}
		}
	}

	return targets, autoIncluded, nil
}

func reachableFromFilter(
	filter map[string]bool, structMap map[string]parser.StructInfo, sanitizable map[string]bool,
) map[string]bool {
	reachable := make(map[string]bool)
	var walk func(name string)
	walk = func(name string) {
		if reachable[name] {
			return
		}
		reachable[name] = true
		s, ok := structMap[name]
		if !ok {
			return
		}
		for _, f := range s.Fields {
			if f.IsStruct && sanitizable[f.ElemType] {
				walk(f.ElemType)
			}
		}
	}
	for name := range filter {
		if sanitizable[name] {
			walk(name)
		}
	}
	return reachable
}

func classifyField(f parser.FieldInfo, sanitizable map[string]bool) FieldAction {
	if f.HasSanitize && !f.IsSlice && !f.IsPtr {
		return ActionSanitize
	}
	if f.HasSanitize && f.IsSlice {
		return ActionSanitizeSlice
	}
	if f.IsStruct && !f.IsPtr && sanitizable[f.ElemType] {
		return ActionSanitize
	}
	if f.IsStruct && f.IsPtr && sanitizable[f.ElemType] {
		return ActionSanitizePtr
	}
	return ActionCopy
}

func findSanitizableStructs(structs []parser.StructInfo) map[string]bool {
	sanitizable := make(map[string]bool)

	changed := true
	for changed {
		changed = false
		for _, s := range structs {
			if sanitizable[s.Name] {
				continue
			}
			for _, f := range s.Fields {
				if f.HasSanitize {
					sanitizable[s.Name] = true
					changed = true
					break
				}
				if f.IsStruct && sanitizable[f.ElemType] {
					sanitizable[s.Name] = true
					changed = true
					break
				}
			}
		}
	}
	return sanitizable
}

type color int

const (
	white color = iota
	gray
	black
)

func topoSort(structs []parser.StructInfo, sanitizable map[string]bool) ([]string, error) {
	structMap := make(map[string]parser.StructInfo, len(structs))
	for _, s := range structs {
		structMap[s.Name] = s
	}

	colors := make(map[string]color)
	var order []string

	var visit func(name string) error
	visit = func(name string) error {
		switch colors[name] {
		case gray:
			return fmt.Errorf("circular struct reference involving %s", name)
		case black:
			return nil
		}
		colors[name] = gray

		s, ok := structMap[name]
		if ok {
			for _, f := range s.Fields {
				if f.IsStruct && sanitizable[f.ElemType] && f.ElemType != name {
					if err := visit(f.ElemType); err != nil {
						return err
					}
				}
			}
		}

		colors[name] = black
		order = append(order, name)
		return nil
	}

	for _, s := range structs {
		if !sanitizable[s.Name] {
			continue
		}
		if err := visit(s.Name); err != nil {
			return nil, err
		}
	}
	return order, nil
}
