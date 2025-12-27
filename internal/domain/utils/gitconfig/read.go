package gitconfig

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/go-git/go-git/v5/config"
	cc "github.com/go-git/go-git/v5/plumbing/format/config"
)

func GetGitConfigUserPath() string {
	userPath, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(userPath, ".gitconfig")
}

func readConfigFile(path string) (*config.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return config.ReadConfig(file)
}

func writeConfigFile(path string, cfg *config.Config) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := cfg.Marshal()
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}

func AddNewIncludeIfNotExists(newInclude string, filePath string) error {
	cfg, err := readConfigFile(filePath)
	if err != nil {
		return err
	}
	sectionInclude := cfg.Raw.Section("include")
	//init section if not exists
	if sectionInclude == nil {
		sectionInclude = &cc.Section{
			Name:    "include",
			Options: cc.Options{},
		}
		cfg.Raw.Sections = append(cfg.Raw.Sections, sectionInclude)
	}

	if slices.ContainsFunc(sectionInclude.Options, func(opt *cc.Option) bool {
		return opt.Key == "path" && opt.Value == newInclude
	}) {
		//already exists
		return nil
	}

	sectionInclude.Options = append(sectionInclude.Options, &cc.Option{
		Key:   "path",
		Value: newInclude,
	})

	return writeConfigFile(filePath, cfg)
}
