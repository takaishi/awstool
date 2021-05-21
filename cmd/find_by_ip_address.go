/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/cobra"
)

// findByIpAddressCmd represents the findByIpAddress command
var findByIpAddressCmd = &cobra.Command{
	Use:   "find-by-ip-address",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return err
		}
		var nextToken *string
		ecsClient := ecs.NewFromConfig(cfg)
		for {
			lco, err := ecsClient.ListClusters(context.TODO(), &ecs.ListClustersInput{NextToken: nextToken})
			if err != nil {
				return err
			}
			for _, clusterArn := range lco.ClusterArns {
				arns, err := listTasks(ecsClient, clusterArn)
				if err != nil {
					return err
				}
				size := 100
				var j int
				for i := 0; i < len(arns); i += size {
					j += size
					if j > len(arns) {
						j = len(arns)
					}
					dto, err := ecsClient.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{Tasks: arns[i:j], Cluster: aws.String(clusterArn)})
					if err != nil {
						return err
					}
					for _, task := range dto.Tasks {
						for _, container := range task.Containers {
							for _, ni := range container.NetworkInterfaces {
								if *ni.PrivateIpv4Address == ipAddress {
									fmt.Printf("%s\n", *task.TaskArn)
									return nil
								}
							}

						}
					}
				}
			}
		}
		return nil
	},
}

func listTasks(ecsClient *ecs.Client, clusterArn string) (taskArns []string, err error) {
	var nextToken *string
	for {
		lto, err := ecsClient.ListTasks(context.TODO(), &ecs.ListTasksInput{Cluster: aws.String(clusterArn), NextToken: nextToken})
		if err != nil {
			return nil, err
		}
		taskArns = append(taskArns, lto.TaskArns...)
		if lto.NextToken == nil {
			break
		} else {
			nextToken = lto.NextToken
		}
	}
	return
}

var ipAddress string

func init() {
	rootCmd.AddCommand(findByIpAddressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findByIpAddressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findByIpAddressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	findByIpAddressCmd.Flags().StringVar(&ipAddress, "ip-address", "t", "Help message for toggle")
}
