// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package btf

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/cilium/ebpf/btf"
	"github.com/cilium/tetragon/pkg/defaults"
	"github.com/cilium/tetragon/pkg/logger"

	"golang.org/x/sys/unix"
)

var (
	btfObj  *btf.Spec
	btfFile string
)

func btfFileExists(file string) error {
	_, err := os.Stat(file)
	return err
}

func observerFindBTF(ctx context.Context, lib, btf string) (string, error) {
	if btf == "" {
		// Alternative to auto-discovery and/or command line argument we
		// can also set via environment variable.
		tetragonBtfEnv := os.Getenv("TETRAGON_BTF")
		if tetragonBtfEnv != "" {
			if _, err := os.Stat(tetragonBtfEnv); err != nil {
				return btf, err
			}
			return tetragonBtfEnv, nil
		}

		var uname unix.Utsname
		err := unix.Uname(&uname)
		if err != nil {
			return btf, fmt.Errorf("Kernel version lookup (uname -r) failing. Use '--kernel' to set manually: %w", err)
		}
		kernelVersion := unix.ByteSliceToString(uname.Release[:])

		// Preference of BTF files, first search for kernel exposed BTF, then
		// check for vmlinux- hubble metadata, and finally if all those are missing
		// search the lib directory for a btf file.
		if _, err := os.Stat(defaults.DefaultBTFFile); err == nil {
			logger.GetLogger().WithField("btf-file", defaults.DefaultBTFFile).Info("BTF discovery: default kernel btf file found")
			return defaults.DefaultBTFFile, nil
		}
		logger.GetLogger().WithField("btf-file", defaults.DefaultBTFFile).Info("BTF discovery: default kernel btf file does not exist")

		runFile := path.Join(lib, "metadata", "vmlinux-"+kernelVersion)
		if _, err := os.Stat(runFile); err == nil {
			logger.GetLogger().WithField("btf-file", runFile).Info("BTF discovery: candidate btf file found")
			return runFile, nil
		}
		logger.GetLogger().WithField("btf-file", runFile).Info("BTF discovery: candidate btf file does not exist")

		runFile = path.Join(lib, "btf")
		if _, err := os.Stat(runFile); err == nil {
			logger.GetLogger().WithField("btf-file", runFile).Info("BTF discovery: candidate btf file found")
			return runFile, nil
		}
		logger.GetLogger().WithField("btf-file", runFile).Info("BTF discovery: candidate btf file does not exist")

		return btf, fmt.Errorf("Kernel version %q BTF search failed kernel is not included in supported list. Please check Tetragon requirements documentation, then use --btf option to specify BTF path and/or '--kernel' to specify kernel version", kernelVersion)
	}
	if err := btfFileExists(btf); err != nil {
		return btf, fmt.Errorf("User specified BTF does not exist: %w", err)
	}
	logger.GetLogger().WithField("btf-file", btf).Info("BTF file: user specified btf file found")
	return btf, nil
}

func NewBTF() (*btf.Spec, error) {
	spec, err := btf.LoadSpec(btfFile)
	if err != nil {
		return nil, err
	}

	overlay := "/sys/kernel/btf/overlay"
	f, err := os.Open(overlay)
	if err != nil {
		logger.GetLogger().WithError(err).Warnf("failed to open /sys/kernel/btf/overlay")
		return spec, nil
	}
	defer f.Close()

	splitSpec, err := btf.LoadSplitSpecFromReader(f, spec)
	if err != nil {
		logger.GetLogger().WithError(err).Warnf("failed to load /sys/kernel/btf/overlay")
		return spec, nil
	}

	numSymbols := 0
	iter := splitSpec.Iterate()
	for {
		if !iter.Next() {
			break
		}

		t := iter.Type
		id, err := spec.Add(t)
		if err != nil {
			name := t.TypeName()
			if name == "" {
				name = "ANON"
			}
			logger.GetLogger().WithField("id", id).WithField("name", t.TypeName()).WithError(err).Warn("failed to add type")
		} else {
			numSymbols++
		}
	}
	logger.GetLogger().Warnf("Added %d symbols from /sys/kernel/btf/overlay", numSymbols)

	return spec, nil
}

func InitCachedBTF(ctx context.Context, lib, btf string) error {
	var err error

	// Find BTF metdaata and populate btf opaqu object
	btfFile, err = observerFindBTF(ctx, lib, btf)
	if err != nil {
		return fmt.Errorf("tetragon, aborting kernel autodiscovery failed: %w", err)
	}
	btfObj, err = NewBTF()
	return err
}

func GetCachedBTFFile() string {
	return btfFile
}

func GetCachedBTF() *btf.Spec {
	return btfObj
}
