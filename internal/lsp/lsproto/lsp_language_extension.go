package lsproto

import "github.com/microsoft/typescript-go/internal/core"

const (
	MethodLanguageExtensionLoadFile = "$/extensibility/language/loadFile"
)

type ProtocolCustomExtensionInfo struct {
	Extension      string          `json:"extension"`
	IsMixedContent bool            `json:"isMixedContent"`
	ScriptKind     core.ScriptKind `json:"scriptKind"`
}

type LanguageExtensionLoadFileParams struct {
	URI DocumentUri `json:"uri"`
}

type LanguageExtensionLoadFileResult struct {
	Content    string          `json:"content"`
	ScriptKind core.ScriptKind `json:"scriptKind"`
}

var LanguageExtensionLoadFileInfo = RequestInfo[*LanguageExtensionLoadFileParams, *LanguageExtensionLoadFileResult]{Method: MethodLanguageExtensionLoadFile}
