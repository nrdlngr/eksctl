package delete

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/vpc"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/printers"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

var (
	clusterConfigFile = ""
)

func deleteClusterCmd(g *cmdutils.Grouping) *cobra.Command {
	p := &api.ProviderConfig{}
	cfg := api.NewClusterConfig()

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Delete a cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doDeleteCluster(p, cfg, cmdutils.GetNameArg(args), cmd); err != nil {
				logger.Critical("%s\n", err.Error())
				os.Exit(1)
			}
		},
	}

	group := g.New(cmd)

	group.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&cfg.Metadata.Name, "name", "n", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, p)
		cmdutils.AddWaitFlag(&wait, fs, "deletion of all resources")
		cmdutils.AddConfigFileFlag(&clusterConfigFile, fs)
	})

	cmdutils.AddCommonFlagsForAWS(group, p, true)

	group.AddTo(cmd)
	return cmd
}

func doDeleteCluster(p *api.ProviderConfig, cfg *api.ClusterConfig, nameArg string, cmd *cobra.Command) error {
	printer := printers.NewJSONPrinter()

	if err := cmdutils.NewMetadataLoader(p, cfg, clusterConfigFile, nameArg, cmd).Load(); err != nil {
		return err
	}

	ctl := eks.New(p, cfg)
	meta := cfg.Metadata

	if !ctl.IsSupportedRegion() {
		return cmdutils.ErrUnsupportedRegion(p)
	}
	logger.Info("using region %s", meta.Region)

	if err := ctl.CheckAuth(); err != nil {
		return err
	}

	logger.Info("deleting EKS cluster %q", meta.Name)
	if err := printer.LogObj(logger.Debug, "cfg.json = \\\n%s\n", cfg); err != nil {
		return err
	}

	stackManager := ctl.NewStackManager(cfg)

	handleErrors := func(errs []error, subject string) error {
		logger.Info("%d error(s) occurred while deleting %s", len(errs), subject)
		for _, err := range errs {
			logger.Critical("%s\n", err.Error())
		}
		return fmt.Errorf("failed to delete %s", subject)
	}

	{
		deprecatedStacks, err := stackManager.ListDeprecatedStacksForDeletion()
		if stackCount := len(deprecatedStacks); stackCount > 0 {
			errs := stackManager.SquentialWaitDeleteStacks(deprecatedStacks)
			if len(errs) > 0 {
				return handleErrors(errs, "deprecated stacks")
			}
			logger.Success("deleted all %s stacks", stackCount)
			return nil
		}
	}

	{
		tryDeleteAllNodeGroups := func(force bool) error {
			errs := stackManager.WaitDeleteAllNodeGroups(force)
			if len(errs) > 0 {
				return handleErrors(errs, "nodegroup(s)")
			}
			return nil
		}
		if err := tryDeleteAllNodeGroups(false); err != nil {
			logger.Info("will retry deleting nodegroup")
			logger.Info("trying to cleanup dangling network interfaces")
			if err := ctl.GetClusterVPC(cfg); err != nil {
				return errors.Wrapf(err, "getting VPC configuration for cluster %q", cfg.Metadata.Name)
			}
			if err := vpc.CleanupNetworkInterfaces(ctl.Provider, cfg); err != nil {
				return err
			}
			if err := tryDeleteAllNodeGroups(true); err != nil {
				return err
			}
		}
		logger.Debug("all nodegroups were deleted")
	}

	var clusterErr bool
	if wait {
		clusterErr = handleIfError(stackManager.WaitDeleteCluster(true), "cluster")
	} else {
		clusterErr = handleIfError(stackManager.DeleteCluster(true), "cluster")
	}

	if clusterErr {
		if handleIfError(ctl.DeprecatedDeleteControlPlane(meta), "control plane") {
			// handleIfError(stackManager.DeprecatedDeleteStackControlPlane(wait), "stack control plane (deprecated)")
		}
	}

	// handleIfError(stackManager.DeprecatedDeleteStackServiceRole(wait), "service group (deprecated)")
	// handleIfError(stackManager.DeprecatedDeleteStackVPC(wait), "stack VPC (deprecated)")
	// handleIfError(stackManager.DeprecatedDeleteStackDefaultNodeGroup(wait), "default nodegroup (deprecated)")

	ctl.MaybeDeletePublicSSHKey(meta.Name)

	kubeconfig.MaybeDeleteConfig(meta)

	if len(deletedResources) == 0 {
		logger.Warning("no EKS cluster resources were found for %q", meta.Name)
	} else {
		logger.Success("the following EKS cluster resource(s) for %q will be deleted: %s. If in doubt, check CloudFormation console", meta.Name, strings.Join(deletedResources, ", "))
	}

	return nil
}
