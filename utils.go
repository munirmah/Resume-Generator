package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
)

func (s *resume) sanitizeResume() error {
	v := reflect.ValueOf(s).Elem()
	walkStruct(v)
	return nil
}

func sanitize(in string) (string, error) {
	replacements := map[string]string{
		"&":    "\\&",
		"%":    "\\%",
		"-":    "\\textendash{}",
		"$":    "\\$",
		"#":    "\\#",
		"<":    "\\textless{}",
		">":    "\\textgreater{}",
		"^":    "\\^",
		"\xA0": "~", // Non-breaking space
		"~":    "\\textasciitilde{}",
	}

	if strings.Contains(in, "\\write18") {
		return "", fmt.Errorf(`security risk: \\write18 found in input`)
	}

	re := regexp.MustCompile(`([&%$#\-<>^\xA0~])`)
	out := re.ReplaceAllStringFunc(in, func(match string) string {
		return replacements[match]
	})

	return out, nil
}

func listify(items []string, delimiter string) string {
	return strings.Join(items, delimiter+" ")
}

func openFile(file string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", file)
	case "darwin":
		cmd = exec.Command("open", file)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", file)
	default:
		return fmt.Errorf("unsupported platform")
	}
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	return nil
}

func walkStruct(v reflect.Value) {
	switch v.Kind() {
	case reflect.String:
		s := v.String()
		if s != "" {
			sanitized, err := sanitize(s)
			if err != nil {
				log.Fatalf("Error sanitizing string: %v", err)
			}
			v.SetString(sanitized)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).PkgPath != "" {
				continue
			}
			walkStruct(v.Field(i))
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			walkStruct(v.Index(i))
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			walkStruct(v.MapIndex(k))
		}
	default:
		return
	}
}

func getFilename(path string) string {
	f := filepath.Base(path)
	return strings.TrimSuffix(f, filepath.Ext(f))
}

func overwriteStruct(target, source interface{}) error {
	tv := reflect.ValueOf(target).Elem()
	sv := reflect.ValueOf(source).Elem()
	if tv.Kind() != reflect.Struct || sv.Kind() != reflect.Struct {
		return fmt.Errorf("target and source must be structs")
	}
	if tv.Type() != sv.Type() {
		return fmt.Errorf("target and source must be of the same type")
	}
	for i := 0; i < sv.NumField(); i++ {
		if tv.Field(i).CanSet() && sv.Field(i).CanSet() {
			tf := tv.Field(i)
			sf := sv.Field(i)
			if tf.Kind() == reflect.Struct && sf.Kind() == reflect.Struct {
				err := overwriteStruct(tf.Addr().Interface(), sf.Addr().Interface())
				if err != nil {
					return err
				}
			} else if !sf.IsZero() {
				tf.Set(sf)
			}
		}
	}
	return nil
}

func checkDependencies(dependencies []string) error {
	log.Info("Checking dependencies")
	for _, dep := range dependencies {
		_, err := exec.LookPath(dep)
		if err != nil {
			return fmt.Errorf("%s not found in PATH", dep)
		}
	}
	return nil
}
