const superb = require('superb')
const { exec } = require('child_process')
const fs = require('fs')

// https://sao.js.org/#/create?id=config-file
module.exports = {
  enforceNewFolder: true,
  showTip: false,
  gitInit: false,
  installDependencies: false,
  prompts: {
    namespace: {
      message: 'What is your vsc namespace?',
      default: 'dwarvesf',
    },
    domain: {
      message: 'What is your vsc domain?',
      default: 'github.com',
    },
    description: {
      message: 'How would you describe the new project?',
      default: `my ${superb()} Go project`,
    },
  },
  data({ domain, namespace }) {
    return {
      domainDir: `${domain}/${namespace}/`,
    }
  },
  move: {
    gitignore: '.gitignore',
  },
  post(
    {
      answers: { domain, namespace },
      folderName,
      folderPath,
      log,
      chalk,
    },
    stream
  ) {
    // check for GOPATH env
    if (!process.env.GOPATH) {
      log.error(
        `${chalk.magenta(
          '$GOPATH'
        )} is not set, it is mandatory for Go projects`
      )
      process.exit(1)
    }
    // check if same project src already exist
    const srcPath = `${process.env.GOPATH}/src/${domain + '/' + namespace}`
    const projectPath = `${srcPath}/${folderName}`
    if (fs.existsSync(projectPath)) {
      log.error(
        `${chalk.magenta(
          projectPath
        )} already exist, please remove it or use a different project name!`
      )
      process.exit(1)
    }
    // move src to srcPath, because of how GOPATH works. In future releases
    // where Go module is more stable we wouldn't have to do this.
    exec(
      `mkdir -p ${srcPath} && mv -n ${folderPath} ${srcPath}/${folderName}`,
      (err, stdout, stderr) => {
        if (err) {
          log.error(err.message)
          process.exit(1)
        }
        if (stderr) {
          log.error(stderr)
          process.exit(1)
        }
        // tips
        log.success('Done, let the hacking begin!')
        log.info(
          `Type ${chalk.magenta(
            'cd ' + srcPath + '/' + folderName
          )} to get started!`
        )
      }
    )
  },
}
