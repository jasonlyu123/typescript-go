package project

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

var (
	_ LanguageExtendabilityHost = (*LSPLanguageExtensionHost)(nil)
	_ LanguageExtendabilityHost = (*NoopLanguageExtensionHost)(nil)
)

type LSPLanguageExtensionHost struct {
	extraFileExtensions []tsoptions.FileExtensionInfo
	client              Client
	files               map[tspath.Path]*virtualDiskFile
	context             context.Context
	toPath              func(fileName string) tspath.Path
}

func NewLSPLanguageExtensionHost(
	extraFileExtensions []tsoptions.FileExtensionInfo,
	client Client,
	context context.Context,
	toPath func(fileName string) tspath.Path,
) *LSPLanguageExtensionHost {
	return &LSPLanguageExtensionHost{
		extraFileExtensions: extraFileExtensions,
		client:              client,
		files:               make(map[tspath.Path]*virtualDiskFile),
		context:             context,
		// files: make(map[tspath.Path]*VirtualDiskFile,
	}
}

func (h *LSPLanguageExtensionHost) GetFiles() map[tspath.Path]*virtualDiskFile {
	return h.files
}

func (h *LSPLanguageExtensionHost) CanHandleFile(fileName string) bool {
	ext := strings.TrimPrefix(filepath.Ext(fileName), ".")
	for _, info := range h.extraFileExtensions {
		if info.Extension == ext {
			return true
		}
	}
	return false
}

func (h *LSPLanguageExtensionHost) GetScriptKindFromFileName(fileName string) core.ScriptKind {
	ext := filepath.Ext(fileName)
	for _, info := range h.extraFileExtensions {
		if info.Extension == ext {
			if info.ScriptKind == core.ScriptKindDeferred {
				return core.GetScriptKindFromFileName(fileName)
			}
			return info.ScriptKind
		}
	}
	return core.GetScriptKindFromFileName(fileName)
}

func (h *LSPLanguageExtensionHost) LoadFile(path string) (*virtualDiskFile, error) {
	params := &lsproto.LanguageExtensionLoadFileParams{
		URI: lsconv.FileNameToDocumentURI(path),
	}
	result, err := h.client.LanguageExtensionLoadFile(h.context, params)
	if err != nil || result == nil {
		return nil, err
	}
	file := newVirtualDiskFile(path, result.Content, result.ScriptKind)
	return file, nil
}

type NoopLanguageExtensionHost struct{}

func (h *NoopLanguageExtensionHost) CanHandleFile(fileName string) bool {
	return false
}

func (h *NoopLanguageExtensionHost) GetScriptKindFromFileName(fileName string) core.ScriptKind {
	return core.GetScriptKindFromFileName(fileName)
}

func (h *NoopLanguageExtensionHost) LoadFile(path string) (*virtualDiskFile, error) {
	return nil, nil
}

func (h *NoopLanguageExtensionHost) GetFiles() map[tspath.Path]*virtualDiskFile {
	return nil
}
