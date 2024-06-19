package mi_server_cmd

import (
	"context"
	"log"
	"os"
	"os/signal"

	mi "github.com/mt3hr/mi/src/app/mi/mi"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := serverCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.MousetrapHelpText = "" // Windowsでマウスから起動しても怒られないようにする
	serverCmd.PersistentFlags().StringVarP(&mi.ConfigFileName, "config_file", "c", "", "使用するコンフィグファイル")
	serverCmd.PersistentFlags().StringVarP(&mi.TagStructFile, "tag_struct_file", "t", "", "使用するタグ構造ファイル")
}

var (
	serverCmd = &cobra.Command{
		Use:              "mi_server",
		PersistentPreRun: mi.PersistentPreRun,
		Run: func(_ *cobra.Command, _ []string) {
			err := mi.LoadRepositories()
			if err != nil {
				log.Fatal(err)
			}
			defer mi.LoadedRepositories.Close()
			interceptCh := make(chan os.Signal)
			signal.Notify(interceptCh, os.Interrupt)
			go func() {
				<-interceptCh
				mi.LoadedRepositories.Close()
				os.Exit(0)
			}()
			mi.LoadedRepositories, err = mi.WrapT(mi.LoadedRepositories)
			if err != nil {
				log.Fatal(err)
			}

			mi.LoadedRepositories.UpdateCache(context.Background())
			err = mi.LaunchServer()
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)
