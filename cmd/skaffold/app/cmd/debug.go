/*
Copyright 2019 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"io"

	"github.com/GoogleContainerTools/skaffold/cmd/skaffold/app/cmd/commands"
	debugging "github.com/GoogleContainerTools/skaffold/pkg/skaffold/debug"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/deploy"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewCmdDebug describes the CLI command to run a pipeline in debug mode.
func NewCmdDebug(out io.Writer) *cobra.Command {
	cmdUse := "debug"
	return commands.
		New(out).
		WithLongDescription(cmdUse, "Runs a pipeline file in debug mode", "Similar to `dev`, but configures the pipeline for debugging.").
		WithFlags(func(f *pflag.FlagSet) {
			AddFlags(f, cmdUse)
		}).
		NoArgs(cancelWithCtrlC(context.Background(), doDebug))
}

func doDebug(ctx context.Context, out io.Writer) error {
	// HACK: disable watcher to prevent redeploying changed containers during debugging
	// TODO: enable file-sync but avoid redeploys of artifacts being debugged
	if len(opts.TargetImages) == 0 {
		opts.TargetImages = []string{"none"}
	}

	deploy.AddManifestTransform(debugging.ApplyDebuggingTransforms)

	return doDev(ctx, out)
}
