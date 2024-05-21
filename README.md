# tokenclone

![Alt text](assets/Octocatopher.png "Octocatopher")

A small Golang utility to clone a GitHub repository using Github app credentials.

## Build

- Clone repo
- Build the binary
    - Defaults to an osx build
    - OS X build option ```make compile-osx```    
    - Linux arm build option ```make compile-linux```   
    - Windows build option ```make compile-windows```
- Your binary to run is in bin/

## Requirements

You need to have a [Github App](https://docs.github.com/en/apps/creating-github-apps) built and deployed to the repos you want to clone. Do that first.

## Usage

Linux/OS X:<br>
```./tokenclone --app_id <id> --pem_path <path to your pem> --repo_url <git repo, https clone link>, and --clone_dir <directory to clone to>```

Windows:<br>
```tokenclone.exe --app_id <id> --pem_path <path to your pem> --repo_url <git repo, https clone link>, and --clone_dir <directory to clone to>```

## Help

```./tokenclone --help```

## Where is this helpful?

You can use this in a bunch of ways. One example of a use for this is Ansible Tower which currently is limited to using PATs for its Github integration. PATs and things like service account users both overscope access as well as allow long aged tokens, by using this utility you can avoid that.