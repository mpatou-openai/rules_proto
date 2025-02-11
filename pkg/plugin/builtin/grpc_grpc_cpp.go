package builtin

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
	"path"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&GrpcGrpcCppPlugin{})
}

// GrpcGrpcCppPlugin implements Plugin for the built-in protoc python plugin.
type GrpcGrpcCppPlugin struct{}

// Name implements part of the Plugin interface.
func (p *GrpcGrpcCppPlugin) Name() string {
	return "grpc:grpc:cpp"
}

// cppGRPCGeneratedFileName is a utility function that returns a function that
// computes the name of a predicted generated file having the given
// extension(s) relative to the given dir.
func cppGRPCGeneratedFileName(reldir string) func(f *protoc.File) []string {
	return func(f *protoc.File) []string {
		name := f.Name
		if reldir != "" {
			name = path.Join(reldir, name)
		}
		return []string{name + ".grpc.pb.cc", name + ".grpc.pb.h"}
	}
}

// Configure implements part of the Plugin interface.
func (p *GrpcGrpcCppPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !protoc.HasServices(ctx.ProtoLibrary.Files()...) {
		return nil
	}
	cppFiles := protoc.FlatMapFiles(
		cppGRPCGeneratedFileName(ctx.Rel),
		protoc.HasService,
		ctx.ProtoLibrary.Files()...,
	)

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/grpc/grpc", "protoc-gen-grpc-cpp"),
		Outputs: cppFiles,
		Options: ctx.PluginConfig.GetOptions(),
	}
}
