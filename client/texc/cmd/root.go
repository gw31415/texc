/*
Copyright © 2020 Amadeus_vn
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/

package cmd

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/gw31415/texc/proto"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	address    string
	block_size int
	verbose    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "texc [main TeX file]",
	Short: "Remote TeX compiler",
	Long:  `Send tex files and receive deliverables.`,
	Args:  cobra.MinimumNArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return err
		}
		defer conn.Close()
		client := proto.NewTexcServiceClient(conn)
		stream, err := client.Sync(context.Background())
		if err != nil {
			return err
		}
		tar_data := bytes.NewBuffer([]byte{})
		tar_w := tar.NewWriter(tar_data)
		paths, err := dirwalk(".")
		if err != nil {
			return err
		}
		var total_size int64 = 0
		for _, path := range paths {
			f, err := os.OpenFile(path, os.O_RDONLY, 0755)
			if err != nil {
				return err
			}
			stat, err := f.Stat()
			if err != nil {
				return err
			}
			tar_w.WriteHeader(&tar.Header{
				Name:    path,
				Mode:    int64(stat.Mode()),
				ModTime: stat.ModTime(),
				Size:    stat.Size(),
			})
			if err != nil {
				return err
			}
			io.Copy(tar_w, f)
			total_size += stat.Size()
			cmd.Printf("Add: %s\n", path)
		}
		tar_w.Close()
		in_pb := new(proto.Input)
		in_pb.Data = make([]byte, block_size)
		var sent_size int = 0
		send_status := "Send: %d%%\n"
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			send_status = "\rSend: %d%%"
		}
		for {
			i, err := tar_data.Read(in_pb.Data)
			sent_size += i
			if err == io.EOF {
				break
			}
			stream.Send(in_pb)
			percent := int64(sent_size) * 100 / total_size
			if percent > 99 {
				percent = 99
			}
			cmd.Printf(send_status, percent)
		}
		cmd.Printf(send_status, 100)
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			cmd.Println()
		}
		stream.Send(&proto.Input{
			Exec: []string{"latexmk", args[0]},
		})
		stream.Send(&proto.Input{
			Dl: getFileNameWithoutExt(args[0]) + ".pdf",
		})
		stream.CloseSend()
		dl := false
		b := new(bytes.Buffer)
		for {
			out, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if out.Stdout != nil {
				if verbose {

				os.Stdout.Write(out.Stdout)
				} else {
					fmt.Print(".")
				}
			}
			if out.Stderr != nil {
				os.Stderr.Write(out.Stderr)
			}
			if out.Data != nil {
				dl = true
				b.Write(out.Data)
			}
		}
		cmd.Println()
		if dl {
			tr := tar.NewReader(b)
			for {
				h, err := tr.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
				wd, _ := os.Getwd()
				path := fmt.Sprintf("%s/%s", wd, h.Name)
				dir := filepath.Dir(path)
				if !h.FileInfo().IsDir() {
					os.MkdirAll(dir, 0755)
					file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0744)
					if err != nil {
						return err
					}
					io.Copy(file, tr)
					cmd.Printf("Download: %s\n", h.Name)
				}
			}
		}
		return nil
	},
}

func dirwalk(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, file := range files {
		if file.IsDir() {
			f, err := dirwalk(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			paths = append(paths, f...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths, nil
}
func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.texc.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "help message for toggle")
	rootCmd.Flags().StringVar(&address, "address", "texc.amas.dev:3475", `address of Texc server`)
	rootCmd.Flags().IntVar(&block_size, "block-size", 0xffff, `block size of sending file`)
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "view output in detail")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".texc" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".texc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
