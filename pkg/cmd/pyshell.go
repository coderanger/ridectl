/*
Copyright 2019 Ridecell, Inc.

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
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/Ridecell/ridectl/pkg/exec"
	"github.com/Ridecell/ridectl/pkg/kubernetes"
)

func init() {
	rootCmd.AddCommand(pyShellCmd)
}

var pyShellCmd = &cobra.Command{
	Use:   "pyshell [flags] <cluster_name>",
	Short: "Open a Python shell on a Summon instance",
	Long:  `Open an interactive Python terminal on a Summon instance running on Kubernetes`,
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Cluster name argument is required")
		}
		if len(args) > 1 {
			return fmt.Errorf("Too many arguments")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		clientset, err := kubernetes.GetClient(kubeconfigFlag)
		if err != nil {
			return errors.Wrap(err, "unable to load Kubernetes configuration")
		}
		pod, err := kubernetes.FindSummonPod(clientset, args[0], fmt.Sprintf("app.kubernetes.io/instance=%s-web", args[0]))
		if err != nil {
			return errors.Wrap(err, "unable to find pod")
		}
		fmt.Printf("Connecting to %s/%s\n", pod.Namespace, pod.Name)

		// Spawn kubectl exec.
		kubectlArgs := []string{"kubectl", "exec", "-it", "-n", pod.Namespace, pod.Name, "--", "bash", "-l", "-c", "python manage.py shell"}
		return exec.Exec(kubectlArgs)
	},
}
