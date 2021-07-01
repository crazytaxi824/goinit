// 纯 js 项目用, 不包含 react lint

package js

import (
	_ "embed" // for go:embed file use
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"local/src/util"
)

var createFolders = []string{".vscode", "src"}

var (
	//go:embed cfgfiles/launch.json
	launchJSON []byte

	//go:embed cfgfiles/settings_template.txt
	settingTemplate []byte

	//go:embed cfgfiles/gitignore
	gitignore []byte

	//go:embed cfgfiles/package.json
	packageJSON []byte

	// for unit test 'jest' use
	//go:embed cfgfiles/example.test.js
	exampleTestJS []byte

	//go:embed eslint/eslintrc-js.json
	eslintrcJSON []byte
)

// file content
var mainJS = []byte(`main();

function main() {
  console.log('hello world');
}
`)

// filesAndContent JS project files
var filesAndContent = []util.FileContent{
	{Path: util.LaunchJSONPath, Content: launchJSON},
	{Path: util.GitignorePath, Content: gitignore},
	{Path: "package.json", Content: packageJSON},
	{Path: "src/main.js", Content: mainJS},
}

// for jest use only

const testFolder = "test"

var jestFileContent = util.FileContent{
	Path:    testFolder + "/example.test.js",
	Content: exampleTestJS,
}

func InitProject(tsjsSet *flag.FlagSet, jestflag, eslint, eslintLocal *bool) (suggs []*util.Suggestion, err error) {
	// parse arges first
	// nolint // flag.ExitOnError will do the os.Exit(2)
	tsjsSet.Parse(os.Args[3:])

	ff := util.InitFoldersAndFiles(createFolders, filesAndContent)

	if *jestflag {
		// add jest example test file
		ff.AddFolders(testFolder)
		ff.AddFiles(jestFileContent)
	}

	if *eslint && *eslintLocal {
		// 如果两个选项都有，则报错
		return nil, errors.New("can not setup eslint globally and locally at same time")
	} else if *eslint && !*eslintLocal {
		// 设置 global eslint
		err = initGlobalEslint(ff)
	} else if !*eslint && *eslintLocal {
		// 设置 local eslint
		err = initLocalEslint(ff)
	} else {
		// 不设置 eslint, 只需要设置 settings.json 文件
		err = initWithoutEslint(ff)
	}

	if err != nil {
		return nil, err
	}

	// 写入所需文件
	fmt.Println("init TypeScript project")
	if err := ff.WriteAllFiles(); err != nil {
		return nil, err
	}

	// 安装所有缺失的依赖
	if err := ff.InstallMissingDependencies(); err != nil {
		return nil, err
	}

	return ff.Suggestions(), nil
}

// // 不设置 ESLint, 写入 <project>/.vscode/settings.json 文件.
func initWithoutEslint(ff *util.FoldersAndFiles) error {
	// 直接写 settings.json 文件
	err := addSettingJSON(ff)
	if err != nil {
		return err
	}
	return nil
}

// 设置 local ESLint:
//  - 写入 <project>/eslint/eslintrc-js.json 本地配置文件.
//  - 写入 <project>/.vscode/settings.json 文件.
//  - 安装 ESLint 缺失的本地依赖.
func initLocalEslint(ff *util.FoldersAndFiles) error {
	// 检查 npm 是否安装，把 suggestion 当 error 返回，因为必须要安装依赖
	if sugg := util.CheckCMDInstall("npm"); sugg != nil {
		return errors.New(sugg.String())
	}

	// 获取项目的绝对地址
	projectPath, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	// 添加 <project>/eslint 文件夹，添加 eslintrc-js.json 文件
	// ff.addEslintJSONAndEspath(projectPath + eslintFilePath)
	ff.AddLintConfigAndLintPath(projectPath+eslintFilePath, eslintrcJSON)

	// 设置 settings.json 文件, 将 config 设置为 eslint 配置文件地址
	err = addSettingJSON(ff)
	if err != nil {
		return err
	}

	// 添加 ESLint 缺失的本地依赖
	return ff.AddMissingDependencies(eslintDependencies, "package.json", "")
}

// 设置 global ESLint:
//  - 写入 ~/.vsc/eslint/eslintrc-js.json 全局配置文件.
//  - 写入 ~/.vsc/vsc-config.json 全局配置文件.
//  - 写入 <project>/.vscode/settings.json 文件.
//  - 安装 ESLint 缺失的全局依赖.
func initGlobalEslint(ff *util.FoldersAndFiles) error {
	// 检查 npm 是否安装，把 suggestion 当 error 返回，因为必须要安装依赖
	if sugg := util.CheckCMDInstall("npm"); sugg != nil {
		return errors.New(sugg.String())
	}

	// 获取 ~/.vsc 文件夹地址
	vscDir, err := util.GetVscConfigDir()
	if err != nil {
		return err
	}

	// 通过 vsc-config.json 获取 eslint.JS 配置文件地址.
	err = readEslintPathFromVscCfgJSON(ff, vscDir)
	if err != nil {
		return err
	}

	// 设置 settings.json 文件, 将 configFile 设置为 eslint 配置文件地址
	err = addSettingJSON(ff)
	if err != nil {
		return err
	}

	// 添加 ESLint 缺失的全局依赖
	eslintFolder := vscDir + eslintDirector
	pkgFilePath := eslintFolder + "/package.json"
	return ff.AddMissingDependencies(eslintDependencies, pkgFilePath, eslintFolder)
}
