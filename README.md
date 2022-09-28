# Custom Quality Gate sonarqube using Golang

This tools for custom quality gate when scanning sast in sonarqube. this tools can running on jenkins , bamboo , or other ci/cd

# Requirement
- java-11-openjdk
- sonar-scanner (included )


# Installation
```bash
mkdir /opt/go-sonar
echo "export PATH=PATH:/opt/go-sonar" > ~/.bashrc
which go-sonar-scanner
```

# Execute in bash 

## sample1 using flag

```bash
go-sonar-scanner --baseurl="http://0.0.0.0:9001" --auth=flag --username=<USERNAME> --password=<PASSWORD> --projectKey=testing123 --qualityGate="sample-qg-1" --args="-Dsonar.login=<TOKEN> -Dsonar.projectKey=farm-app"
```

## sample2 using env
```bash
export sonaruser=<USERNAME>
export sonarpass=<PASSWORD>
go-sonar-scanner --baseurl="http://0.0.0.0:9001" --auth=env --projectKey=farm-app --qualityGate="sample-qg-1" --args="-Dsonar.login=<TOKEN> -Dsonar.projectKey=farm-app"
```

## sample3 running on jenkins

```groovy
pipeline {
	agent any 
	stages {
		stage("Testing") {
			steps {
				script {
                    go-sonar-scanner --baseurl="http://0.0.0.0:9001" --auth=flag --username=<USERNAME> --password=<PASSWORD> --projectKey=farm-app --qualityGate="sample-qg-1" --args="-Dsonar.login=<TOKEN> -Dsonar.projectKey=farm-app"
                }
            }
        }
    }
}
```

# FRIENDS :D
@agambewe
