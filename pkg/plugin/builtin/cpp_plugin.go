package builtin

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
	"path"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&CppPlugin{})
}

// CppPlugin implements Plugin for the built-in protoc C++ plugin.
type CppPlugin struct{}

// Name implements part of the Plugin interface.
func (p *CppPlugin) Name() string {
	return "builtin:cpp"
}

// cppGeneratedFileName is a utility function that returns a function that
// computes the name of a predicted generated file having the given
// extension(s) relative to the given dir.
func cppGeneratedFileName(reldir string) func(f *protoc.File) []string {
	return func(f *protoc.File) []string {
		name := f.Name
		if reldir != "" {
			name = path.Join(reldir, name)
		}
		return []string{name + ".pb.cc", name + ".pb.h"}
	}
}

// Configure implements part of the Plugin interface.
func (p *CppPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	cppFiles := protoc.FlatMapFiles(
		cppGeneratedFileName(ctx.Rel),
		protoc.Always,
		ctx.ProtoLibrary.Files()...,
	)

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/builtin", "cpp"),
		Outputs: cppFiles,
		Options: ctx.PluginConfig.GetOptions(),
	}

}
